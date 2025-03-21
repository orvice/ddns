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

func New() (*Config, error) {
	path := "."
	viper.AddConfigPath(path)
	viper.SetConfigName("app")

	viper.SetConfigType("env")
	viper.AutomaticEnv()
	// viper.ReadInConfig()
	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}
