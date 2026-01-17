package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/anurag-327/neuron-core/conn"
	"github.com/anurag-327/neuron-core/pkg/api"
	"github.com/anurag-327/neuron-core/pkg/logger"
	"github.com/anurag-327/neuron-core/runner/docker/pool"
	"github.com/anurag-327/neuron-core/runtime"
	"github.com/anurag-327/neuron-core/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
)

type ExecOptions struct {
	Cmd         []string
	ExecTimeout time.Duration

	AttachStdout bool
	AttachStderr bool
	AttachStdin  bool
}

type ExecResult struct {
	Stdout         string
	Stderr         string
	ExitCode       int
	ErrType        api.SandboxError
	ErrMsg         string
	ContainerDirty bool
	Duration       time.Duration
}

func ExecuteDockerExec(
	ctx context.Context,
	dockerClient *client.Client,
	containerID string,
	opts ExecOptions,
	phase string,
) ExecResult {

	appLogger := logger.GetGlobalLogger()
	result := ExecResult{}
	result.ContainerDirty = false
	result.ExitCode = 1
	startTime := time.Now()

	// Create exec (no timeout here)
	execResp, err := dockerClient.ContainerExecCreate(
		ctx,
		containerID,
		container.ExecOptions{
			Cmd:          opts.Cmd,
			AttachStdout: opts.AttachStdout,
			AttachStderr: opts.AttachStderr,
			AttachStdin:  opts.AttachStdin,
		},
	)
	if err != nil {
		appLogger.Error(time.Now(), "Exec create failed", map[string]interface{}{
			"container_id": containerID,
			"phase":        phase,
			"error":        err.Error(),
		})
		result.ErrType = api.ErrSandboxError
		result.ErrMsg = "Exec create failed"
		return result
	}

	// Go-side timeout (protect runner)
	execCtx, cancel := context.WithTimeout(ctx, opts.ExecTimeout)
	defer cancel()

	// Attach BEFORE start
	attach, err := dockerClient.ContainerExecAttach(
		execCtx,
		execResp.ID,
		container.ExecStartOptions{},
	)
	appLogger.Info(time.Now(), "Exec attach", map[string]interface{}{
		"container_id": containerID,
		"phase":        phase,
	})
	if err != nil {
		appLogger.Error(time.Now(), "Exec attach failed", map[string]interface{}{
			"container_id": containerID,
			"phase":        phase,
			"error":        err.Error(),
		})
		result.ErrType = api.ErrSandboxError
		result.ErrMsg = "Exec attach failed"
		result.ContainerDirty = true
		return result
	}
	defer attach.Close()

	// Start exec
	if err := dockerClient.ContainerExecStart(
		context.Background(),
		execResp.ID,
		container.ExecStartOptions{},
	); err != nil {
		appLogger.Error(time.Now(), "Exec start failed", map[string]interface{}{
			"container_id": containerID,
			"phase":        phase,
			"error":        err.Error(),
		})
		result.ErrType = api.ErrSandboxError
		result.ErrMsg = "Exec start failed"
		result.ContainerDirty = true
		return result
	}

	// Read stdout / stderr asynchronously
	var stdoutBuf, stderrBuf bytes.Buffer
	done := make(chan error, 1)

	go func() {
		_, err := stdcopy.StdCopy(&stdoutBuf, &stderrBuf, attach.Reader)
		done <- err
	}()

	// Wait for completion OR Go timeout
	select {

	case <-execCtx.Done():
		appLogger.Error(time.Now(), "Execution timeout (Go context)", map[string]interface{}{
			"container_id": containerID,
			"timeout_sec":  opts.ExecTimeout,
			"phase":        phase,
		})
		result.ErrType = api.ErrTLE
		result.ErrMsg = api.MsgTLE
		result.ContainerDirty = true
		return result

	case err := <-done:
		if err != nil {
			appLogger.Error(time.Now(), "Output read failed", map[string]interface{}{
				"container_id": containerID,
				"timeout_sec":  opts.ExecTimeout,
				"phase":        phase,
				"error":        err.Error(),
			})
			result.ErrType = api.ErrSandboxError
			result.ErrMsg = "Output read failed"
			result.ContainerDirty = true
			return result
		}
	}

	// Inspect exit code
	inspect, err := dockerClient.ContainerExecInspect(
		context.Background(),
		execResp.ID,
	)
	if err != nil {
		appLogger.Error(time.Now(), "Exec inspect failed", map[string]interface{}{
			"container_id": containerID,
			"timeout_sec":  opts.ExecTimeout,
			"phase":        phase,
			"error":        err.Error(),
		})
		result.ErrType = api.ErrSandboxError
		result.ErrMsg = "Exec inspect failed"
		result.ContainerDirty = true
		result.Duration = time.Since(startTime)
		return result
	}

	result.Stdout = stdoutBuf.String()
	result.Stderr = stderrBuf.String()
	result.ExitCode = inspect.ExitCode

	// Classify exit code
	switch inspect.ExitCode {
	case 0:
		// success

	case 124, 137:
		result.ErrType = api.ErrTLE
		result.ErrMsg = api.MsgTLE

	default:
		result.ErrType = api.ErrRuntimeError
		result.ErrMsg = "Runtime error"
	}

	result.Duration = time.Since(startTime)
	return result
}

func acquireContainer(ctx context.Context, language string) (*pool.ContainerPool, string, error) {
	p := pool.Manager.GetPool(language)
	if p == nil {
		return nil, "", errors.New("invalid language: " + language)
	}
	containerID, err := p.Get(ctx)
	if err != nil {
		return nil, "", err
	}
	return p, containerID, nil
}

func prepareWorkspace(spec api.RunConfig) (string, string, runtime.FileNames, error) {
	runID := uuid.New().String()
	relativeBasePath := fmt.Sprintf("/tmp/runner/job_%s", runID)
	rootPath, _ := os.Getwd()
	basePath := filepath.Join(rootPath, relativeBasePath)

	if err := os.MkdirAll(basePath, 0777); err != nil {
		return "", "", runtime.FileNames{}, fmt.Errorf("failed to make directory: %w", err)
	}

	// Force 0777 permissions
	if err := os.Chmod(basePath, 0777); err != nil {
		utils.DeleteFolder(basePath)
		return "", "", runtime.FileNames{}, fmt.Errorf("failed to update permissions: %w", err)
	}

	languageConfig, ok := runtime.LanguageRegistry[spec.Language]
	if !ok {
		utils.DeleteFolder(basePath)
		return "", "", runtime.FileNames{}, fmt.Errorf("unsupported language: %s", spec.Language)
	}

	names := BuildFileNames(BuildFilesParam{
		BasePath:  basePath,
		FileName:  languageConfig.EntryFile.FileName,
		Extension: languageConfig.EntryFile.Extension,
	})

	if err := utils.WriteContentToFile(names.PathFull, []byte(spec.Code), 0777); err != nil {
		utils.DeleteFolder(basePath)
		return "", "", runtime.FileNames{}, fmt.Errorf("write code failed: %w", err)
	}

	inputPath := filepath.Join(basePath, "input.txt")
	if err := utils.WriteContentToFile(inputPath, []byte(spec.Input), 0777); err != nil {
		utils.DeleteFolder(basePath)
		return "", "", runtime.FileNames{}, fmt.Errorf("write input failed: %w", err)
	}

	return runID, basePath, names, nil
}

func executeCommand(
	ctx context.Context,
	cli *conn.DockerClient,
	containerID string,
	basePath string,
	cmdStr string,
	timeout time.Duration,
	phase string,
) ExecResult {
	containerJobPath := filepath.Join("/sandbox", filepath.Base(basePath))

	// We need to construct the sh -c command for timeout
	execCmd := []string{
		"sh", "-c",
		fmt.Sprintf(
			"cd %s && timeout -s KILL %.2fs sh -c '%s'",
			containerJobPath,
			timeout.Seconds(),
			cmdStr,
		),
	}

	// Go-side timeout is slightly larger to allow `timeout` command to kill first
	goCtxTimeout := timeout + 2*time.Second

	return ExecuteDockerExec(ctx, cli.Client, containerID, ExecOptions{
		Cmd:          execCmd,
		ExecTimeout:  goCtxTimeout,
		AttachStdout: true,
		AttachStderr: true,
	}, phase)
}
