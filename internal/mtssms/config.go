package mtssms

type Config struct {
	BaseURL   string `env:"MTS_URL, required"`
	Login     string `env:"MTS_LOGIN, required"`
	Password  string `env:"MTS_PASSWORD, required"`
	Signuture string `env:"MTS_SIGN, required"`
}
