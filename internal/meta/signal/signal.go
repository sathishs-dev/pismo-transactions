package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	terminateSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
	}

	globalSig   signalWrap
	globalSigMu sync.Mutex
)

func init() {
	initGlobalSig()
}

func initGlobalSig() {
	globalSigMu.Lock()
	defer globalSigMu.Unlock()

	if globalSig.sigCh != nil {
		signal.Stop(globalSig.sigCh)
	}

	globalSig = signalWrap{
		sigCh:      make(chan os.Signal, 1),
		shutdownCh: make(chan struct{}),
	}
}

type syncList struct {
	sync.RWMutex
	funcs []func()
}

func (s *syncList) add(fn func()) {
	s.Lock()
	defer s.Unlock()

	s.funcs = append(s.funcs, fn)
}

func (s *syncList) list() []func() {
	s.RLock()
	defer s.RUnlock()

	return s.funcs
}

type signalWrap struct {
	funcs        syncList
	sigCh        chan os.Signal
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
}

func Add(fn func()) {
	globalSig.funcs.add(fn)
}

func New(sigs ...os.Signal) <-chan struct{} {
	ctx := NewWithContext(context.Background(), sigs...)
	return ctx.Done()
}

func NewWithContext(parent context.Context, sigs ...os.Signal) context.Context {
	if sigs == nil {
		sigs = terminateSignals
	}
	signal.Notify(globalSig.sigCh, sigs...)

	ctx, cancel := context.WithCancel(detachedContext{parent: parent})

	go func() {
		select {
		case <-parent.Done():
		case <-globalSig.sigCh:
		case <-globalSig.shutdownCh:
		}

		funcs := globalSig.funcs.list()
		for i := len(funcs) - 1; i >= 0; i-- {
			funcs[i]()
		}

		cancel()
	}()

	return ctx
}

type detachedContext struct{ parent context.Context }

func (d detachedContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (d detachedContext) Done() <-chan struct{}             { return nil }
func (d detachedContext) Err() error                        { return nil }
func (d detachedContext) Value(key interface{}) interface{} { return d.parent.Value(key) }

func Shutdown() {
	globalSigMu.Lock()
	defer globalSigMu.Unlock()

	globalSig.shutdownOnce.Do(func() {
		close(globalSig.shutdownCh)
	})
}
