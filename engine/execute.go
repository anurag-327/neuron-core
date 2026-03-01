package engine

import (
	"context"

	"github.com/anurag-327/neuron-core/pkg/api"
)

// Execute executes the code in a sandboxed environment
// validates code and language, runs the code in a sandboxed environment and returns the result
func (e *ExecutionService) Execute(ctx context.Context, spec ExecuteSpec) (*api.RunResult, error) {
	result := e.runner.Run(ctx, api.RunConfig{
		Code:     spec.Code,
		Language: spec.Language,
		Limit: api.Limit{
			MemoryKB: spec.Limit.MemoryKB,
			TimeMs:   spec.Limit.TimeMs,
		},
		Input: spec.Input,
	})
	return &result, nil
}

func (e *ExecutionService) Health(ctx context.Context) error {
	return e.runner.Health(ctx)
}
