package smsgw

import (
	"github.com/fdully/msms/internal/mtssms"
	"github.com/fdully/msms/internal/setup"
)

// Compile-time check to assert this config matches requirements.
var (
	_ setup.MTSConfigProvider = (*Config)(nil)
)

type Config struct {
	MTS mtssms.Config

	Port      string   `env:"PORT, default=8080"`
	BadPhones []string `env:"BAD_PHONES"`
}

func (c *Config) MTSConfig() *mtssms.Config {
	return &c.MTS
}
