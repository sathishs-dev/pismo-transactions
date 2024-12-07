package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sathishs-dev/pismo-transactions/internal/meta/writer"
)

func withRecovery(next http.Handler, tracesCh chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := debug.Stack()

				writer.WriteJSON(
					w,
					http.StatusInternalServerError,
					writer.ErrorResponse{
						Code:  "recover_error",
						Title: "recover",
						Trace: fmt.Sprintf("%v, %s", err, stackTrace),
					},
				)
				if tracesCh != nil {
					tracesCh <- string(stackTrace)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}
}
