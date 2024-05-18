package adapter

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/j03hanafi/halo-suster/common/configs"
)

var (
	jwtCache     *cache.Cache
	jwtCacheOnce sync.Once
)

func GetJWTCache() *cache.Cache {
	jwtCacheOnce.Do(func() {
		exp := (time.Duration(configs.Get().JWT.Expire) * time.Second) - time.Minute
		jwtCache = cache.New(exp, exp)
	})

	if jwtCache == nil {
		panic("jwt cache is not initialized")
	}
	return jwtCache
}
