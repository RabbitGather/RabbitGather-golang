package server

import "context"

type Server interface {
	// GracefulStop stops the server gracefully.
	GracefulStop(ctx context.Context) (err error)
	// ListenAndServe starts the server and blocks until the server is stopped.
	// Returns error if the server fails to start or stops unexpectedly.
	ListenAndServe() (err error)
}
