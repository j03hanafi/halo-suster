package handler

import "sync"

var baseResponsePool = sync.Pool{
	New: func() any {
		return new(baseResponse)
	},
}

func baseResponseAcquire() *baseResponse {
	return baseResponsePool.Get().(*baseResponse)
}

func baseResponseRelease(t *baseResponse) {
	*t = baseResponse{}
	baseResponsePool.Put(t)
}

type baseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

var versionPool = sync.Pool{
	New: func() any {
		return new(version)
	},
}

func versionAcquire() *version {
	return versionPool.Get().(*version)
}

func versionRelease(t *version) {
	*t = version{}
	versionPool.Put(t)
}

type version struct {
	Version string `json:"version"`
}

var healthPool = sync.Pool{
	New: func() any {
		return new(health)
	},
}

func healthAcquire() *health {
	return healthPool.Get().(*health)
}

func healthRelease(t *health) {
	*t = health{}
	healthPool.Put(t)
}

type health struct {
	Status     string `json:"status"`
	IdleConns  int32  `json:"idle_conns"`
	TotalConns int32  `json:"total_conns"`
	MaxConns   int32  `json:"max_conns"`
}
