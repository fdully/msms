package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// contextKeyRequestID is the unique key in the context where the request ID is
// stored.
const contextKeyRequestID = contextKey("request_id")

// PopulateRequestID populates the request context with a random UUID if one
// does not already exist.
func PopulateRequestID() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if existing := RequestIDFromContext(ctx); existing == "" {
				u, err := uuid.NewRandom()
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				ctx = withRequestID(ctx, u.String())
				r = r.Clone(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDFromContext pulls the request ID from the context, if one was set.
// If one was not set, it returns the empty string.
func RequestIDFromContext(ctx context.Context) string {
	v := ctx.Value(contextKeyRequestID)
	if v == nil {
		return ""
	}

	t, ok := v.(string)
	if !ok {
		return ""
	}
	return t
}

// withRequestID sets the request ID on the provided context, returning a new
// context.
func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, id)
}
