// service is the interface that connects the business logic and dependencies.
package manager

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/graceful_shutdown"
	"github.com/meowalien/go-meowalien-lib/schedule"

	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/server"
)

type Manager interface {
	server.Launcher
}

type Constructor struct {
	Retryer          schedule.Retryer
	GracefulShutdown graceful_shutdown.GracefulShutdown
	Servers          []server.Launcher
}

func (m Constructor) New() Manager {
	return &manager{
		Constructor:   m,
		stopFunctions: make(map[string]func(ctx context.Context) error),
		cond:          sync.NewCond(&sync.Mutex{}),
	}
}

// manager is the manager of all Manager, it will control all server's lifecycle.
type manager struct {
	Constructor

	isRunning atomic.Bool

	stopFunctions map[string]func(ctx context.Context) error

	waitingCount int
	cond         *sync.Cond

	manuallyGracefulStopChanOnce sync.Once

	error error
}

func (g *manager) Name() string {
	return "Manager"
}

// ListenAndServe starts all Manager and block until any server stopped, or GracefulStop() is called.
func (g *manager) ListenAndServe() (err error) {
	defer func() {
		if err != nil {
			g.isRunning.Store(false)
		}
	}()

	if len(g.Servers) == 0 {
		return errs.New("no server to start")
	}
	for _, server := range g.Servers {
		g.startServer(server)
	}

	if !g.isRunning.CompareAndSwap(false, true) {
		return errs.New("already started")
	}

	_ = g.WaitAllServerStop(context.Background())
	return g.getError()
}

// GracefulStop will stop all Manager gracefully and wait for them to stop.
// Return error if context is canceled or timeout.
func (g *manager) GracefulStop(ctx context.Context) (err error) {
	fmt.Println("Manager GracefulStop")
	g.manuallyGracefulStopChanOnce.Do(func() {
		g.stopAllServers(ctx)
		fmt.Println("waiting for all servers to stop")
		err1 := g.WaitAllServerStop(ctx)
		g.addError(err1)
	})
	return g.getError()
}

// WaitAllServerStop will block until all Manager are stopped, no matter they stopped gracefully or not.
// Return error if context is canceled or timeout.
func (g *manager) WaitAllServerStop(ctx context.Context) error {
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

func (g *manager) addGoRoutine(delta int) {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()

	g.waitingCount += delta
}

func (g *manager) doneGoRoutine() {
	g.cond.L.Lock()
	defer g.cond.L.Unlock()
	g.waitingCount--
	g.cond.Broadcast()

}

func (g *manager) stopAllServers(ctx context.Context) {
	for name, stopFunc := range g.stopFunctions {
		fmt.Println("stopping server: ", name)
		// since all stop functions are protected by sync.Once, it is safe to call them all
		err1 := stopFunc(ctx)
		if err1 != nil {
			g.addError(err1)
		}

		//err = errs.Connect(err, err1)
	}
	return
}

func (g *manager) addStopFunc(name string, stopFunc func(ctx context.Context) (err error)) {
	once := sync.Once{}
	// make sure stopFunc is only called once
	stopFuncWithOnce := func(ctx context.Context) (err1 error) {
		once.Do(func() {
			err1 = stopFunc(ctx)
		})
		return
	}
	g.stopFunctions[name] = stopFuncWithOnce
	g.GracefulShutdown.Add(name, stopFuncWithOnce)
}

func (g *manager) startServer(svr server.Launcher) {
	g.addGoRoutine(1)
	go func(svr server.Launcher) {
		defer g.doneGoRoutine()

		g.addStopFunc(svr.Name(), func(ctx context.Context) (err error) {
			return svr.GracefulStop(ctx)
		})
		err := g.Retryer.Try(context.Background(), func(ctx context.Context) error {
			return svr.ListenAndServe()
		})
		if err != nil {
			g.addError(err)
		}
	}(svr)
}

func (g *manager) addError(err error) {
	g.error = fmt.Errorf("%w\n%s", g.error, err.Error())
}

func (g *manager) getError() error {
	err := g.error
	g.error = nil
	return err
}
