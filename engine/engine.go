package engine

import "github.com/anurag-327/neuron-core/runner"

// ExecutionService is the core execution service which is responsible for executing code in a sandboxed environment
type ExecutionService struct {
	runner runner.Runner
}

// NewExecutionService creates a new execution service
func NewExecutionService(runner runner.Runner) *ExecutionService {
	return &ExecutionService{runner: runner}
}
