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

	Addr      string   `env:"SMSGWAPP_ADDRESS, default=localhost:8090"`
	BadPhones []string `env:"SMSGWAPP_BAD_PHONES"`
}

func (c *Config) MTSConfig() *mtssms.Config {
	return &c.MTS
}
