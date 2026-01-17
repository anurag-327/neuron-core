package docker

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anurag-327/neuron-core/pkg/api"
	"github.com/anurag-327/neuron-core/runtime"
	"github.com/anurag-327/neuron-core/utils"
)

func DetectLanguageError(language, stdout, stderr string) (api.SandboxError, string) {

	s := stderr
	c := stdout + "\n" + stderr // Some runtime errors print to stdout

	// C++
	if language == "cpp" {
		// Compiler errors (from g++)
		if strings.Contains(s, "error:") ||
			strings.Contains(s, "fatal error:") ||
			strings.Contains(s, "undefined reference") {
			return api.ErrCompilationError, api.MsgCompilationError
		}

		// Runtime crash detection
		if strings.Contains(s, "Segmentation fault") ||
			strings.Contains(s, "core dumped") ||
			strings.Contains(s, "abort") ||
			strings.Contains(s, "floating point exception") {
			return api.ErrRuntimeError, api.MsgRuntimeError
		}

	}

	// Go
	if language == "go" {
		if strings.Contains(s, "undefined:") ||
			strings.Contains(s, "cannot use") ||
			strings.Contains(s, "no required module") {
			return api.ErrCompilationError, api.MsgCompilationError
		}

		if strings.Contains(c, "panic:") ||
			strings.Contains(c, "runtime error:") {
			return api.ErrRuntimeError, api.MsgRuntimeError
		}
	}

	// Python
	if language == "python" {
		if strings.Contains(s, "SyntaxError") ||
			strings.Contains(s, "IndentationError") {
			return api.ErrCompilationError, api.MsgCompilationError
		}

		if strings.Contains(c, "Traceback (most recent call last):") {
			return api.ErrRuntimeError, api.MsgRuntimeError
		}
	}

	// Java
	if language == "java" {
		if strings.Contains(s, "error:") ||
			strings.Contains(s, "cannot find symbol") ||
			strings.Contains(s, "symbol not found") {
			return api.ErrCompilationError, api.MsgCompilationError
		}

		if strings.Contains(c, "Exception in thread") {
			return api.ErrRuntimeError, api.MsgRuntimeError
		}
	}

	// JavaScript (Node.js)
	if language == "javascript" {
		if strings.Contains(s, "SyntaxError:") {
			return api.ErrCompilationError, api.MsgCompilationError
		}

		if strings.Contains(c, "TypeError:") ||
			strings.Contains(c, "ReferenceError:") ||
			strings.Contains(c, "UnhandledPromiseRejectionWarning") {
			return api.ErrRuntimeError, api.MsgRuntimeError
		}
	}

	// LAST RESORT CHECK — strict runtime detection
	// Only treat stderr as runtime error if it contains *real* crash signals.

	if isMeaningfulRuntimeErrorGeneric(stderr) {
		return api.ErrRuntimeError, api.MsgRuntimeError
	}

	// Everything OK
	return api.ErrRuntimeError, api.MsgRuntimeError
}

func isMeaningfulRuntimeErrorGeneric(stderr string) bool {
	s := strings.ToLower(stderr)

	// ignore if nothing meaningful
	if strings.TrimSpace(s) == "" {
		return false
	}

	// ignore logs like [info], [debug], etc.
	if strings.Contains(s, "[info]") ||
		strings.Contains(s, "[debug]") ||
		strings.Contains(s, "note:") {
		return false
	}

	// ignore warnings
	if strings.Contains(s, "warning") {
		return false
	}

	// real runtime errors
	crashPatterns := []string{
		"segmentation fault",
		"core dumped",
		"panic:",
		"runtime error",
		"traceback (most recent call last):",
		"exception in thread",
		"nullpointerexception",
		"typeerror:",
		"referenceerror:",
		"indexerror:",
		"valueerror:",
		"abort",
		"illegal instruction",
		"floating point exception",
	}

	for _, pat := range crashPatterns {
		if strings.Contains(s, pat) {
			return true
		}
	}

	return false
}

type BuildFilesParam struct {
	FileName  string
	Extension string
	BasePath  string
}

func BuildFileNames(param BuildFilesParam) runtime.FileNames {
	full := param.FileName + "." + param.Extension
	return runtime.FileNames{
		FileName: param.FileName,
		FullName: full,
		PathBase: filepath.Join(param.BasePath, param.FileName),
		PathFull: filepath.Join(param.BasePath, full),
	}
}

func ProcessResult(lang string, status int, stdout, stderr string, jobID string) api.RunResult {
	// Truncate first to prevent massive strings from hitting regex or downstream DB
	const MaxOutputSize = 256 * 1024 // 256KB
	stdout = utils.TruncateOutput(stdout, MaxOutputSize)
	stderr = utils.TruncateOutput(stderr, MaxOutputSize)

	cleanStderr := utils.SanitizeOutput(stderr)
	cleanStdout := utils.SanitizeOutput(stdout)

	result := api.RunResult{
		Stdout:   cleanStdout,
		Stderr:   cleanStderr,
		ExitCode: int64(status),
	}

	switch status {
	case 0: // Success
		return result

	case 1: // Runtime Error
		errorType, message := DetectLanguageError(lang, cleanStdout, cleanStderr)
		result.ErrType = errorType
		result.ErrMsg = message
		return result

	case 124, 143, 137: // SIGTERM / OOM / Timeout
		// For TLE/OOM, output is often huge/infinite or irrelevant. Discard it.
		result.ErrType = api.ErrTLE
		result.ErrMsg = api.MsgTLE
		result.Stdout = ""
		result.Stderr = ""
		return result

	case 139, 136, 134: // Segmentation Fault (SIGSEGV)
		result.ErrType = api.ErrRuntimeError
		result.ErrMsg = api.MsgRuntimeError
		return result

	default:
		fmt.Println("Unknown exit code:", status)
		errorType, message := DetectLanguageError(lang, cleanStdout, cleanStderr)
		result.ErrType = errorType
		result.ErrMsg = message
		return result
	}
}
