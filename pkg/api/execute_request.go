package api

type ExecuteRequest struct {
	Code     string `json:"code"`
	Input    string `json:"input"`
	Language string `json:"language"`
	Limit    Limit  `json:"limit"`
}

type Limit struct {
	TimeMs   int64 `json:"timeMs"`
	MemoryKB int64 `json:"memoryKB"`
}
