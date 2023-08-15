// service is the interface that connects the business logic and dependencies.
package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/graceful_shutdown"
)

type Servers struct {
	stopFunctions map[string]func(ctx context.Context) error

	waitingCount int
	cond         *sync.Cond

	grpcServer                   *GRPCServer
	manuallyGracefulStopChan     chan context.Context
	manuallyGracefulStopChanOnce sync.Once
	GracefulShutdown             graceful_shutdown.GracefulShutdown
}

func (g *Servers) AddService(svr *Service) {
	switch {
	case svr.GRPCService != nil:
		g.addGRPCService(svr.GRPCService)
		//fallthrough
		//case svr.HTTPService != nil:
		//	g.addHTTPService(svr.HTTPService)
	}

}

// StartAll starts all servers and block until any server stopped, or GracefulStop() is called.
func (g *Servers) StartAll() (err error) {
	grpcServerErrChan := g.startGrpcServer()

	// if any service is added, the listen and serve above will
	// start goroutine to listen and serve, so the g.waitingCount will not be 0
	if g.waitingCount == 0 {
		return errs.New("no server to start")
	}
	var stopCtx context.Context
	select {
	case err = <-grpcServerErrChan:
	case stopCtx = <-g.manuallyGracefulStopChan:
		//	the error is nil
	}
	if stopCtx != nil {
		// stop all servers if any server stopped with error
		g.stopAllServers(stopCtx)
		return
	}

	// when the program reaches here, it means that the GracefulStop() is called, do nothing
	return
}

// GracefulStop will stop all servers gracefully and wait for them to stop.
// Return error if context is canceled or timeout.
func (g *Servers) GracefulStop(ctx context.Context) (err error) {
	g.manuallyGracefulStopChanOnce.Do(func() {
		close(g.manuallyGracefulStopChan)
		err = g.WaitAllServerStop(ctx)
	})
	return
}

// WaitAllServerStop will block until all servers are stopped, no matter they stopped gracefully or not.
func (g *Servers) WaitAllServerStop(ctx context.Context) error {
	if g.cond == nil {
		g.cond = sync.NewCond(&sync.Mutex{})
	}
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	for g.waitingCount > 0 {
		g.cond.Wait()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
	return nil
}

func (g *Servers) addGoRoutine(delta int) {
	if g.cond == nil {
		g.cond = sync.NewCond(&sync.Mutex{})
	}
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.waitingCount += delta
}

func (g *Servers) doneGoRoutine() {
	if g.cond == nil {
		g.cond = sync.NewCond(&sync.Mutex{})
	}
	g.cond.L.Lock()
	defer g.cond.L.Unlock()
	g.waitingCount--
	g.cond.Broadcast()

}

func (g *Servers) stopAllServers(ctx context.Context) {
	for name, stopFunc := range g.stopFunctions {
		fmt.Println("stopping server: ", name)
		// since all stop functions are protected by sync.Once, it is safe to call them all
		stopFunc(ctx)
	}
}

func (g *Servers) addStopFunc(name string, stopFunc func(ctx context.Context) (err error)) {
	once := sync.Once{}
	// make sure stopFunc is only called once
	stopFunc = func(ctx context.Context) (err1 error) {
		once.Do(func() {
			err1 = stopFunc(ctx)
		})
		return
	}
	g.stopFunctions[name] = stopFunc
	g.GracefulShutdown.Add(name, stopFunc)
}

func (g *Servers) startGrpcServer() (errChan <-chan error) {
	ec := make(chan error, 1)
	errChan = ec
	if g.grpcServer == nil {
		// if grpcServer is not set, return nil to block the select statement infinitely
		ec = nil
		return
	}

	g.addGoRoutine(1)
	go func() {
		defer g.doneGoRoutine()
		defer close(ec)

		g.addStopFunc("grpcServer", g.grpcServer.GracefulStop)
		ec <- g.grpcServer.ListenAndServe()
	}()
	return
}

func (g *Servers) addGRPCService(rpcServers GRPCService) {
	if g.grpcServer == nil {
		g.grpcServer = &GRPCServer{
			servers: g,
		}
	}

	g.grpcServer.AddService(rpcServers)
}
