package smsgw

import (
	"net/http"

	"github.com/fdully/msms/pkg/logging"
)

func (s *Server) handleRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx).Named("handleRoot")
		logger.Debugw("starting")
		defer logger.Debugw("finishing")

		s.h.RenderJSON(w, http.StatusOK, nil)
	})
}
