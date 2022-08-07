package smsgw

import (
	"context"
	"strings"

	"github.com/fdully/msms/internal/middleware"
	"github.com/fdully/msms/internal/serverenv"
	"github.com/fdully/msms/pkg/logging"
	"github.com/fdully/msms/pkg/render"
	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	env    *serverenv.ServerEnv
	h      *render.Renderer
}

func NewServer(cfg *Config, env *serverenv.ServerEnv) (*Server, error) {

	return &Server{
		env:    env,
		config: cfg,
		h:      render.NewRenderer(),
	}, nil
}

// Routes returns the router for this server.
func (s *Server) Routes(ctx context.Context) *mux.Router {
	logger := logging.FromContext(ctx).Named("smsgw")

	r := mux.NewRouter()
	r.Use(middleware.Recovery())
	r.Use(middleware.PopulateRequestID())
	r.Use(middleware.PopulateLogger(logger))

	r.Handle("/", s.handleRoot())
	r.Handle("/api/v1/sms/send", s.handleSend()).Methods("GET")
	r.Handle("/api/v1/sms/ise", s.handleISE()).Methods("GET")
	r.Handle("/api/v1/sms/info/{smsID}", s.handleInfo()).Methods("GET")

	return r
}

func (s *Server) CheckBadPhones(tel string) bool {
	for _, v := range s.config.BadPhones {
		if strings.Contains(tel, v) {
			return true
		}
	}

	return false
}
