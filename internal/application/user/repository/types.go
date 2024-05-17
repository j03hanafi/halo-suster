package repository

import (
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var userPool = sync.Pool{
	New: func() any {
		return new(user)
	},
}

func userAcquire() *user {
	return userPool.Get().(*user)
}

func userRelease(t *user) {
	*t = user{}
	userPool.Put(t)
}

type user struct {
	ID        ulid.ULID `db:"id"`
	NIP       string    `db:"nip"`
	Name      string    `db:"name"`
	Password  string    `db:"password"`
	IsIT      bool      `db:"is_it"`
	CreatedAt time.Time `db:"created_at"`
}
