// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/fdully/msms/internal/mtssms"
	"github.com/fdully/msms/internal/serverenv"
	"github.com/fdully/msms/pkg/logging"
	"github.com/sethvargo/go-envconfig"
)

type MTSConfigProvider interface {
	MTSConfig() *mtssms.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config interface{}) (*serverenv.ServerEnv, error) {
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections, resolving secrets, and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config interface{}, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	// Build a list of mutators. This list will grow as we initialize more of the
	// configuration, such as the secret manager.
	var mutatorFuncs []envconfig.MutatorFunc

	// Build a list of options to pass to the server env.
	var serverEnvOpts []serverenv.Option

	// Process first round of environment variables.
	if err := envconfig.ProcessWith(ctx, config, l, mutatorFuncs...); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	if provider, ok := config.(MTSConfigProvider); ok {
		logger.Info("configuring MTS provider")

		mtsConfig := provider.MTSConfig()
		m, err := mtssms.New(mtsConfig.BaseURL, mtsConfig.Login, mtsConfig.Password, mtsConfig.Signuture)
		if err != nil {
			return nil, fmt.Errorf("unable to create MTS client: %w", err)
		}

		// Update serverEnv setup.
		serverEnvOpts = append(serverEnvOpts, serverenv.WithMTSProvider(m))
	}

	return serverenv.New(ctx, serverEnvOpts...), nil
}
