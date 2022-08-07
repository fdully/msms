package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Maintainable is an interface that determines if the implementer can supply
// maintenance mode settings.
type Maintainable interface {
	MaintenanceMode() bool
}

func ProcessMaintenance(cfg Maintainable) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.MaintenanceMode() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, `{"error": "please try again later"}`)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
