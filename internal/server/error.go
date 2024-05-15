package server

type ErrorHandler interface {
	Error() string
	Status() int
}
