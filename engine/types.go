package engine

import "github.com/anurag-327/neuron-core/pkg/api"

type Limit struct {
	MemoryKB int64
	TimeMs   int64
}

type ExecuteSpec struct {
	Code     string
	Language string
	Limit    Limit
	Input    string
}

type ExecuteResult struct {
	Stdout   string           `json:"stdout"`
	Stderr   string           `json:"stderr"`
	ExitCode int              `json:"exit_code"`
	ErrType  api.SandboxError `json:"err_type"`
	ErrMsg   string           `json:"err_msg"`
	Metrics  api.TimeMetrics  `json:"metrics"`
}
