package docker

import (
	"context"
	"time"

	"github.com/anurag-327/neuron-core/conn"
	"github.com/anurag-327/neuron-core/pkg/api"
	"github.com/anurag-327/neuron-core/pkg/logger"
	"github.com/anurag-327/neuron-core/runtime"
	"github.com/anurag-327/neuron-core/utils"
)

// DockerRunner implements the Runner interface for running code in a Docker container
type DockerRunner struct {
	client *conn.DockerClient
}

// NewDockerRunner creates a new DockerRunner instance
func NewDockerRunner(client *conn.DockerClient) *DockerRunner {
	return &DockerRunner{client: client}
}

// Run implements the Runner interface for running code in a Docker container
// 1. Acquire Container
// 2. Prepare Workspace
// 3. Compile Phase (if needed)
// 4. Run Phase
// 5. Process Result
func (d *DockerRunner) Run(ctx context.Context, spec api.RunConfig) api.RunResult {

	appLogger := logger.GetGlobalLogger()
	result := api.RunResult{ExitCode: 1}

	totalStart := time.Now()
	var compileDuration time.Duration

	// 1. Acquire Container
	p, containerID, err := acquireContainer(ctx, spec.Language)
	if err != nil {
		appLogger.Error(time.Now(), "Failed to acquire container", map[string]interface{}{
			"language": spec.Language,
			"error":    err.Error(),
		})
		result.ErrType = api.ErrInternalError
		result.ErrMsg = err.Error()
		return result
	}

	defer func() {
		if result.ContainerDirty {
			p.ReplaceContainer(containerID)
		} else {
			p.Put(containerID)
		}
	}()

	runID, basePath, names, err := prepareWorkspace(spec)
	if err != nil {
		appLogger.Error(time.Now(), "Failed to prepare workspace", map[string]interface{}{
			"containerID": containerID,
			"language":    spec.Language,
			"error":       err.Error(),
		})
		result.ErrType = api.ErrInternalError
		result.ErrMsg = "Failed to prepare workspace"
		return result
	}

	defer func() {
		utils.DeleteFolder(basePath)
	}()

	languageConfig := runtime.LanguageRegistry[spec.Language]

	// 3. Compile Phase (if needed)
	if languageConfig.CompileCmd != nil {
		compileCmdStr := languageConfig.CompileCmd(names)
		compileTimeout := languageConfig.ResourceLimits.TimeMs + 2*time.Second

		cStart := time.Now()
		res := executeCommand(ctx, d.client, containerID, basePath, compileCmdStr, compileTimeout, "compile")
		compileDuration = time.Since(cStart)

		if res.ExitCode != 0 {
			final := ProcessResult(spec.Language, res.ExitCode, res.Stdout, res.Stderr, runID)
			final.ContainerDirty = res.ContainerDirty
			switch res.ErrType {
			case api.ErrInternalError, api.ErrSandboxError:
				final.ErrType = api.ErrCompilationError
				final.ErrMsg = res.ErrMsg
			default:
				final.ErrType = api.ErrCompilationError
				final.ErrMsg = "Compilation Failed"
			}
			// Copy this back to main result so defer knows if dirty
			result = final
			return result
		}
	}

	// 4. Run Phase
	runCmd := languageConfig.RunCmd(names)
	runTimeout := languageConfig.ResourceLimits.TimeMs

	rStart := time.Now()
	res := executeCommand(ctx, d.client, containerID, basePath, runCmd, runTimeout, "run")
	runDuration := time.Since(rStart)

	// 5. Process Result
	final := ProcessResult(spec.Language, res.ExitCode, res.Stdout, res.Stderr, runID)
	final.ContainerDirty = res.ContainerDirty

	if res.ErrType != "" {
		final.ErrType = res.ErrType
		final.ErrMsg = res.ErrMsg
	}

	totalDuration := time.Since(totalStart)

	final.Metrics = api.TimeMetrics{
		Total:   totalDuration.Milliseconds(),
		Compile: compileDuration.Milliseconds(),
		Run:     runDuration.Milliseconds(),
	}

	result = final
	return final
}

func (d *DockerRunner) Health(ctx context.Context) error {
	_, err := d.client.Ping(ctx)
	return err
}
