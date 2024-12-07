package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/rs/zerolog"
)

// getLoglevel receives loglevel string and converts its to zerolog.Level
func getLoglevel(l string) zerolog.Level {
	switch strings.ToLower(l) {
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		log.Printf("invalid loglevel %s passed, configuring INFO as default level", l)
		return zerolog.InfoLevel
	}
}

// responseLogWriter wrapper for http.ResponseWriter
type responseLogWriter struct {
	http.ResponseWriter
	statusCode int
}

// extends the responsewriter
func newResponseLogWriter(w http.ResponseWriter) *responseLogWriter {
	return &responseLogWriter{
		w, 0,
	}
}

// WriteHeader wrapper for response.WriteHeader
func (r *responseLogWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// loggerMiddleware logs the each http request with its necessary fields to the log.Output
func loggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			nrw := newResponseLogWriter(w)
			reqStartedAt := time.Now()

			defer func() {
				reqEndedAt := time.Now()

				var event *zerolog.Event
				if nrw.statusCode >= 400 {
					event = logger.Error()
				} else {
					event = logger.Info()
				}

				event.Str("http_canonical_path", getCanonicalPath(r))
				event.Str("http_request_host", r.Host)
				event.Str("http_request_method", r.Method)
				event.Str("http_request_path", r.URL.Path)
				event.Float64("http_response_time", float64(reqEndedAt.Sub(reqStartedAt).Nanoseconds())/100000.0)
				event.Int("http_status", nrw.statusCode)

				event.Msgf("%s %s %d", r.Method, r.URL.RequestURI(), nrw.statusCode)
			}()
			next.ServeHTTP(nrw, r)
		}
		return http.HandlerFunc(fn)
	}
}

// getCanonicalPath used to get the canonical path from the request
func getCanonicalPath(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		return pattern
	}

	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// No matching pattern, so just return the request path.
		// if no canonical path matches we are simply returning the URL path
		return routePath
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}
