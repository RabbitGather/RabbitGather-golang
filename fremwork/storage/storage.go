package storage

type Storage[T any] interface {
	Store(T) error
}
