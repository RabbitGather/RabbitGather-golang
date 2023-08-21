package server

import "context"

// Launcher is the interface that defines the methods of a server lifecycle.
type Launcher interface {
	// Name returns the name of the server.
	Name() string
	// GracefulStop stops the server gracefully.
	GracefulStop(ctx context.Context) (err error)
	// ListenAndServe starts the server and blocks until the server is stopped.
	// Returns error if the server fails to start or stops unexpectedly.
	ListenAndServe() (err error)
}
