package domain

import (
	"net/http"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

const (
	RoleIT        = "IT"
	RoleNurse     = "Nurse"
	UserFromToken = "loggedInUser"
)

var UserPool = sync.Pool{
	New: func() any {
		return new(User)
	},
}

func UserAcquire() *User {
	return UserPool.Get().(*User)
}

func UserRelease(t *User) {
	*t = User{}
	UserPool.Put(t)
}

type User struct {
	ID        ulid.ULID
	NIP       string
	Name      string
	Password  string
	Role      string
	ImgURL    string
	CreatedAt time.Time
}

type ErrDuplicateNIP struct{}

func (e ErrDuplicateNIP) Error() string {
	return "NIP already registered"
}

func (e ErrDuplicateNIP) Status() int {
	return http.StatusConflict
}

type ErrUserNotFound struct{}

func (e ErrUserNotFound) Error() string {
	return "User not found"
}

func (e ErrUserNotFound) Status() int {
	return http.StatusNotFound
}

type ErrInvalidNIP struct{}

func (e ErrInvalidNIP) Error() string {
	return "Invalid NIP"
}

func (e ErrInvalidNIP) Status() int {
	return http.StatusNotFound
}

type ErrInvalidPassword struct{}

func (e ErrInvalidPassword) Error() string {
	return "Invalid password"
}

func (e ErrInvalidPassword) Status() int {
	return http.StatusBadRequest
}

type ErrAccessNotAllowed struct{}

func (e ErrAccessNotAllowed) Error() string {
	return "Access not allowed"
}

func (e ErrAccessNotAllowed) Status() int {
	return http.StatusBadRequest
}

type ErrNotFoundOrNotNurse struct{}

func (e ErrNotFoundOrNotNurse) Error() string {
	return "User not found or is not a nurse"
}

func (e ErrNotFoundOrNotNurse) Status() int {
	return http.StatusBadRequest
}