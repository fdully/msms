package smsgw

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/fdully/msms/internal/mtssms"
	"github.com/fdully/msms/pkg/logging"
)

func (s *Server) handleISE() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := logging.FromContext(ctx).Named("handleISE")
		logger.Debugw("starting")
		defer logger.Debugw("finishing")

		tel := r.FormValue("telephones")
		msg := r.FormValue("message")

		ttel, err := mtssms.CheckSingleTelephone(tel)
		if err != nil {
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

		passcode := getPasscode(msg)
		if passcode == "" {
			err := fmt.Errorf("doesn't contain passcode %s", msg)
			logger.Info(err)
			s.h.RenderJSON(w, http.StatusBadRequest, err)
			return
		}

		mts := s.env.MTS()

		mtsReq := mtssms.CreateRequestWithSingleMessage([]string{ttel}, "passcode: "+passcode, s.config.MTS.Signuture)
		res, err := mts.Send(ctx, mtsReq)
		if err != nil {
			s.h.RenderJSON(w, http.StatusBadGateway, err)
			return
		}

		logger.Infow(r.RemoteAddr, tel, msg)

		s.h.RenderJSON(w, http.StatusOK, res)
	})
}

func getPasscode(msg string) string {
	m := regexp.MustCompile(`Имя пользователя: (\d{6})`)
	mm := m.FindStringSubmatch(msg)

	if len(mm) > 1 {
		return mm[1]
	}

	// changing pattern
	m = regexp.MustCompile(`Username: (\d{6})`)
	mm = m.FindStringSubmatch(msg)

	if len(mm) > 1 {
		return mm[1]
	}

	return ""
}
