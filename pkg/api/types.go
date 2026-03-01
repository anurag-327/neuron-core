package api

type RunConfig struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Limit    Limit  `json:"limit"`
	Input    string `json:"input"`
}

type RunResult struct {
	Stdout         string       `json:"stdout"`
	Stderr         string       `json:"stderr"`
	ErrType        SandboxError `json:"err_type"`
	ErrMsg         string       `json:"err_msg"`
	ExitCode       int64        `json:"exit_code"`
	ContainerDirty bool         `json:"-"`
	Metrics        TimeMetrics  `json:"metrics"`
}

type TimeMetrics struct {
	Total   int64 `json:"total_ms"`
	Compile int64 `json:"compile_ms"`
	Run     int64 `json:"run_ms"`
}
