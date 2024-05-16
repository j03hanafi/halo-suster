package id

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropyPool = sync.Pool{
	New: func() any {
		return rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404
	},
}

var zero ulid.ULID

func New() ulid.ULID {
	entropy := entropyPool.Get().(io.Reader)
	defer entropyPool.Put(entropy)

	newID, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return ulid.ULID{}
	}
	return newID
}

func IsZero(id ulid.ULID) bool {
	return id.Compare(zero) == 0
}
