package runner

import (
	"context"

	"github.com/anurag-327/neuron-core/conn"
	"github.com/anurag-327/neuron-core/pkg/api"
	"github.com/anurag-327/neuron-core/runner/docker"
)

// Runner is the interface for running code in a sandboxed environment
// It provides a common interface for different runner implementations

type Runner interface {
	Run(ctx context.Context, spec api.RunConfig) api.RunResult
	Health(ctx context.Context) error
}

func NewRunner(client *conn.DockerClient) Runner {
	return docker.NewDockerRunner(client)
}
