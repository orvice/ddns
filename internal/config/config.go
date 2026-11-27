package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DNSProvider    string `mapstructure:"DNS_PROVIDER"`
	Domain         string `mapstructure:"DOMAIN"`
	TelegramChatID int64  `mapstructure:"TELEGRAM_CHATID"`
	TelegramToken  string `mapstructure:"TELEGRAM_TOKEN"`

	CFToken string `mapstructure:"CF_TOKEN"`

	AliyunAccessKeyID     string `mapstructure:"ALIYUN_ACCESS_KEY_ID"`
	AliyunAccessKeySecret string `mapstructure:"ALIYUN_ACCESS_KEY_SECRET"`
}

// New loads configuration from environment variables (and optionally an `app`
// file in the working directory). It uses a fresh viper instance so tests can
// run in parallel without sharing state.
func New() (*Config, error) {
	return Load(".")
}

// Load is like New but lets callers specify the search path. Useful for tests.
func Load(path string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()

	// Viper's AutomaticEnv only resolves keys that have been registered.
	// Bind every struct tag so env vars are picked up reliably.
	for _, key := range []string{
		"DNS_PROVIDER",
		"DOMAIN",
		"TELEGRAM_CHATID",
		"TELEGRAM_TOKEN",
		"CF_TOKEN",
		"ALIYUN_ACCESS_KEY_ID",
		"ALIYUN_ACCESS_KEY_SECRET",
	} {
		_ = v.BindEnv(key)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
