package id

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

var entropyPool = sync.Pool{
	New: func() any {
		return rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	},
}

var timePool = sync.Pool{
	New: func() any {
		return ulid.Timestamp(time.Now())
	},
}

func New() ulid.ULID {
	entropy := entropyPool.Get().(io.Reader)
	defer entropyPool.Put(entropy)

	now := timePool.Get().(uint64)
	defer timePool.Put(now)

	newID, err := ulid.New(now, entropy)
	if err != nil {
		return ulid.ULID{}
	}
	return newID
}
