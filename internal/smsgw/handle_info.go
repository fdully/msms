package smsgw

import (
	"net/http"

	"github.com/fdully/msms/internal/mtssms"
	"github.com/fdully/msms/pkg/logging"
	"github.com/gorilla/mux"
)

func (s *Server) handleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx).Named("handleInfo")
		logger.Debugw("starting")
		defer logger.Debugw("finishing")

		vars := mux.Vars(r)
		id := vars["smsID"]

		if id == "" {
			s.h.RenderJSON(w, http.StatusBadRequest, nil)
			return
		}

		mts := s.env.MTS()

		mtsReq := mtssms.CreateInfoRequest(id)
		res, err := mts.Info(ctx, mtsReq)
		if err != nil {
			s.h.RenderJSON(w, http.StatusBadGateway, err)
			return
		}

		s.h.RenderJSON(w, http.StatusOK, res)
	})
}
