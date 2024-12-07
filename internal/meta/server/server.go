package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

const DefaultWebPort = 8080

type (
	httpServer struct {
		server   *http.Server
		tracesCh chan string
		errorLog *log.Logger
	}

	HTTPServer interface {
		Start() error
		Stop(context.Context) error
	}
)

func New(h http.Handler, opts ...Option) HTTPServer {
	setup := options{
		port: DefaultWebPort,
	}

	for _, o := range opts {
		o(&setup)
	}

	tracesCh := make(chan string, runtime.GOMAXPROCS(0))
	if !setup.disableRecover {
		h = withRecovery(h, tracesCh)
	}

	return &httpServer{
		tracesCh: tracesCh,
		errorLog: setup.errorLogger,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", setup.port),
			Handler:      h,
			ErrorLog:     setup.errorLogger,
			ReadTimeout:  setup.readTimeout,
			WriteTimeout: setup.writeTimeout,
		},
	}
}

func (h *httpServer) Start() error {
	go h.outputStackTraces()

	return h.server.ListenAndServe()
}

func (h *httpServer) Stop(ctx context.Context) error {
	close(h.tracesCh)

	return h.server.Shutdown(ctx)
}

func (h *httpServer) outputStackTraces() {
	for trace := range h.tracesCh {
		if h.errorLog != nil {
			h.errorLog.Printf("recovered from panic: \n %s", trace)
		}
	}
}
