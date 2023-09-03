package graceful_shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GracefulShutdown interface {
	// Add adds a stop mission to the stop stack.
	Add(name string, f func(ctx context.Context))
	// StopNow run all stop missions immediately and wait for all stop missions to finish or the context is done.
	StopNow(ctx context.Context) error
}

type Config func(ctx context.Context) context.Context

func ShutdownDeadLine(timeout time.Duration) Config {
	return func(ctx context.Context) (x context.Context) {
		x, _ = context.WithTimeout(ctx, timeout)
		return x
	}
}

// NewGracefulShutdown creates a new GracefulShutdown, witch listens to SIGINT and SIGTERM.
// When the program receives SIGINT or SIGTERM, it will run all stop missions in the stop stack in the reverse order of adding.
// onSystemStopContext is a function that returns a context and a cancel function, which will be called when the program
// receives SIGINT or SIGTERM. and the context will be used to run all stop missions.
// set onSystemStopContext to nil to use the default context.Background() and func(){}.
func NewGracefulShutdown(cf ...Config) GracefulShutdown {
	onSystemStopContext := func() context.Context {
		ctx := context.Background()
		for _, f := range cf {
			ctx = f(ctx)
		}
		return ctx
	}

	g := &gracefulShutdown{
		stopStack:           make([]stopMission, 0),
		onSystemStopContext: onSystemStopContext,
		c:                   make(chan os.Signal, 1),
		afterAllStop:        make(chan struct{}),
	}
	signal.Notify(g.c, syscall.SIGINT, syscall.SIGTERM)
	go g.listen()
	return g
}

type gracefulShutdown struct {
	afterAllStop        chan struct{}
	c                   chan os.Signal
	onSystemStopContext func() context.Context
	stopStack           []stopMission
}

func (g *gracefulShutdown) listen() {
	defer close(g.afterAllStop)
	theSignal := <-g.c
	log.Println("GracefulShutdown: received signal: ", theSignal)
	signal.Stop(g.c)
	g.runStopStack()
}

func (g *gracefulShutdown) Add(name string, f func(ctx context.Context)) {
	g.stopStack = append(g.stopStack, stopMission{
		f:    f,
		name: name,
	})
}

func (g *gracefulShutdown) runStopStack() {
	ctx := g.onSystemStopContext()
	for i := len(g.stopStack) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic when running stopStack:%s , %+v\n", g.stopStack[i].name, r)
				}
			}()
			g.stopStack[i].f(ctx)
		}()
	}
	return
}

func (g *gracefulShutdown) waitAllStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-g.afterAllStop:
		return nil
	}
}

func (g *gracefulShutdown) StopNow(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// abort
		return ctx.Err()
	case <-g.afterAllStop:
		// already stopped
		return nil
	case g.c <- syscall.SIGTERM:
		// successfully sent
		return g.waitAllStop(ctx)

	default:
		// already notified, wait for stop
		return g.waitAllStop(ctx)
	}
}

type stopMission struct {
	f    func(ctx context.Context)
	name string
}
