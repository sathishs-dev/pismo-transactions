package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/sathishs-dev/pismo-transactions/internal/meta/server"
	"github.com/sathishs-dev/pismo-transactions/pkg/handler"
)

func initWebServer(l zerolog.Logger, h handler.Handler) server.HTTPServer {
	web := chi.NewMux()
	web.Use(
		loggerMiddleware(l),
	)

	web.Route("/accounts", func(r chi.Router) {
		r.Post("/", h.CreateAccount())
		r.Get("/{accountId}", h.GetAccount())
	})

	web.Post("/transactions", h.CreateTransaction())

	return server.New(web)
}
