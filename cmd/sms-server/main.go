package main

import (
	"context"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	"github.com/fdully/msms/internal/setup"
	"github.com/fdully/msms/internal/smsgw"
	"github.com/fdully/msms/pkg/logging"
	"github.com/fdully/msms/pkg/server"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var cfg smsgw.Config
	env, err := setup.Setup(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	smsgwServer, err := smsgw.NewServer(&cfg, env)
	if err != nil {
		return fmt.Errorf("smsgw.NewServer: %w", err)
	}

	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("error listening: %s", err.Error())
	}
	// Close the listener when the application closes.
	defer l.Close()

	srv, err := server.NewFromListener(l)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}
	logger.Infow("server listening", "address", srv.Addr())
	return srv.ServeHTTPHandler(ctx, smsgwServer.Routes(ctx))
}
