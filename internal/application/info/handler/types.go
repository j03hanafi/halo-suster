package handler

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type version struct {
	Version string `json:"version"`
}

type health struct {
	Status     string `json:"status"`
	IdleConns  int32  `json:"idle_conns"`
	TotalConns int32  `json:"total_conns"`
	MaxConns   int32  `json:"max_conns"`
}
