package lib

type Constructor[T any] interface {
	New() T
}
