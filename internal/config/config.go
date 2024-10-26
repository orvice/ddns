package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DNSMode        string `mapstructure:"DNS_MODE"`
	Domain         string `mapstructure:"DOMAIN"`
	TelegramChatID int64  `mapstructure:"TELEGRAM_CHATID"`
	TelegramToken  string `mapstructure:"TELEGRAM_TOKEN"`
}

var (
	config Config
)

func Init() (err error) {
	err = LoadConfig(".")
	if err != nil {
		return
	}
	return
}

func GetConfig() Config {
	return config
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")

	viper.SetConfigType("env")
	viper.AutomaticEnv()
	// viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
