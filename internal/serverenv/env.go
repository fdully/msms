// Package serverenv defines common parameters for the sever environment.
package serverenv

import (
	"context"

	"github.com/fdully/msms/internal/mtssms"
)

// ServerEnv represents latent environment configuration for servers in this application.
type ServerEnv struct {
	mts mtssms.Provider
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(*ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(ctx context.Context, opts ...Option) *ServerEnv {
	env := &ServerEnv{}

	for _, f := range opts {
		env = f(env)
	}

	return env
}

func WithMTSProvider(m mtssms.Provider) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.mts = m
		return s
	}
}

func (env *ServerEnv) MTS() mtssms.Provider {
	return env.mts
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	return nil
}
