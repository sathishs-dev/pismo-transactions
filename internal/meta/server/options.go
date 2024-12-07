package server

import (
	"log"
	"time"
)

type options struct {
	port           int
	disableRecover bool
	readTimeout    time.Duration
	writeTimeout   time.Duration
	errorLogger    *log.Logger
}

// Option - Used to extend the server's functionalities during initialization
type Option func(*options)

// WithPort - This will serve the server on the specified port
func WithPort(port int) Option {
	return func(o *options) {
		o.port = port
	}
}

// WithReadTimeout - Will set the read timeout on the server
func WithReadTimeout(t time.Duration) Option {
	return func(o *options) {
		o.readTimeout = t
	}
}

// WithWriteTimeout - Will set the write timeout on the server
func WithWriteTimeout(t time.Duration) Option {
	return func(o *options) {
		o.writeTimeout = t
	}
}

// WithErrorLogger - Will set a custom error logger on the server
func WithErrorLogger(l *log.Logger) Option {
	return func(o *options) {
		o.errorLogger = l
	}
}

// WithoutRecover - Will disable the default recover middleware
func WithoutRecover() Option {
	return func(o *options) {
		o.disableRecover = true
	}
}
