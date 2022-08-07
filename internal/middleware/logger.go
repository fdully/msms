package middleware

import (
	"net/http"

	"github.com/fdully/msms/pkg/logging"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// PopulateLogger populates the logger onto the context.
func PopulateLogger(originalLogger *zap.SugaredLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logger := originalLogger

			// Only override the logger if it's the default logger. This is only used
			// for testing and is intentionally a strict object equality check because
			// the default logger is a global default in the logger package.
			if existing := logging.FromContext(ctx); existing == logging.DefaultLogger() {
				logger = existing
			}

			// If there's a request ID, set that on the logger.
			if id := RequestIDFromContext(ctx); id != "" {
				logger = logger.With("request_id", id)
			}

			ctx = logging.WithLogger(ctx, logger)
			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
