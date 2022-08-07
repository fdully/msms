package smsgw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fdully/msms/internal/mtssms"
	"github.com/fdully/msms/pkg/logging"
)

func (s *Server) handleSend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx).Named("handleSend")
		logger.Debugw("starting")
		defer logger.Debugw("finishing")

		tel := r.FormValue("telephones")
		msg := r.FormValue("message")

		tels := strings.Split(tel, ",")

		if err := mtssms.CheckMessage(msg); err != nil {
			logger.Info(err)
			s.h.RenderJSON(w, http.StatusBadRequest, err)
			return
		}

		if s.CheckBadPhones(tel) {
			err := fmt.Errorf("containes bad phones %s", tel)
			logger.Info(err)
			s.h.RenderJSON(w, http.StatusBadRequest, err)

			return
		}

		var checkedTelephones []string
		for _, v := range tels {
			vv, err := mtssms.CheckSingleTelephone(v)
			if err != nil {
				logger.Info(err)
				continue
			}

			checkedTelephones = append(checkedTelephones, vv)
		}

		if len(checkedTelephones) < 1 {
			err := fmt.Errorf("bad telephones %s", tel)
			logger.Info(err)
			s.h.RenderJSON(w, http.StatusBadRequest, err)
			return
		}

		mts := s.env.MTS()

		mtsReq := mtssms.CreateRequestWithSingleMessage(checkedTelephones, msg, s.config.MTS.Signuture)
		res, err := mts.Send(ctx, mtsReq)
		if err != nil {
			s.h.RenderJSON(w, http.StatusBadGateway, err)
			return
		}

		logger.Infow(r.RemoteAddr, tel, msg)

		s.h.RenderJSON(w, http.StatusOK, res)
	})
}
