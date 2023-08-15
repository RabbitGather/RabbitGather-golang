package graceful_shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/meowalien/go-meowalien-lib/errs"
)

type GracefulShutdown interface {
	// Add adds a stop mission to the stop stack.
	Add(name string, f func(ctx context.Context) error)
	// SetOnSystemStopContext replace the function that returns the context and cancel function used when the system stops.
	// The default function returns a background context and a empty cancel function.
	SetOnSystemStopContext(fc func() (context.Context, context.CancelFunc))
	// WaitAllStop waits for all stop missions to finish or the context is done.
	WaitAllStop(ctx context.Context) error
	// StopNow run all stop missions immediately and wait for all stop missions to finish or the context is done.
	StopNow(ctx context.Context) error
}

// NewGracefulShutdown creates a new GracefulShutdown, witch listens to SIGINT and SIGTERM.
// When the program receives SIGINT or SIGTERM, it will run all stop missions in the stop stack in the reverse order of adding.
func NewGracefulShutdown() GracefulShutdown {
	g := &gracefulShutdown{
		stopStack: make([]stopMission, 0),
		onSystemStopContext: func() (context.Context, context.CancelFunc) {
			return context.Background(), func() {}
		},
		c:            make(chan os.Signal, 1),
		afterAllStop: make(chan struct{}),
	}
	signal.Notify(g.c, syscall.SIGINT, syscall.SIGTERM)
	go g.listen()
	return g
}

type gracefulShutdown struct {
	afterAllStop        chan struct{}
	c                   chan os.Signal
	onSystemStopContext func() (context.Context, context.CancelFunc)
	stopStack           []stopMission
}

func (g *gracefulShutdown) listen() {
	defer close(g.afterAllStop)
	theSignal := <-g.c
	log.Println("GracefulShutdown: received signal: ", theSignal)
	signal.Stop(g.c)
	err := g.runStopStack()
	if err != nil {
		log.Fatalln("GracefulShutdown: finished with error: ", err)
	}
}

func (g *gracefulShutdown) Add(name string, f func(ctx context.Context) error) {
	g.stopStack = append(g.stopStack, stopMission{
		f:    f,
		name: name,
	})
}

func (g *gracefulShutdown) SetOnSystemStopContext(fc func() (context.Context, context.CancelFunc)) {
	g.onSystemStopContext = fc
}

func (g *gracefulShutdown) runStopStack() error {
	ctx, cancel := g.onSystemStopContext()
	defer cancel()
	var err error
	for i := len(g.stopStack) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic when running stopStack:%s , %+v\n", g.stopStack[i].name, r)
				}
			}()
			err1 := g.stopStack[i].f(ctx)
			err = errs.New(err, err1)

		}()
	}
	return err
}

func (g *gracefulShutdown) WaitAllStop(ctx context.Context) error {
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
		return g.WaitAllStop(ctx)

	default:
		// already notified, wait for stop
		return g.WaitAllStop(ctx)
	}
}

type stopMission struct {
	f    func(ctx context.Context) error
	name string
}
