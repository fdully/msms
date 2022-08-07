package middleware

import (
	"net/http"

	"github.com/fdully/msms/pkg/logging"
	"github.com/gorilla/mux"
)

// Recovery recovers from panics and other fatal errors. It keeps the server and
// service running, returning 500 to the caller while also logging the error in
// a structured format.
func Recovery() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger := logging.FromContext(ctx).Named("middleware.Recover")

			defer func() {
				if p := recover(); p != nil {
					logger.Errorw("http handler panic", "panic", p)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
