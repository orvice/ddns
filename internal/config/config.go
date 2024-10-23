package config

import "github.com/orvice/utils/env"

var (
	DOMAIN string
)

var (
	DNSMode string

	TelegramToken  string
	TelegramChatID int64
)

func GetConfigFromEnv() {
	DNSMode = env.Get("DNS_MODE", "cf")
	DOMAIN = env.Get("DOMAIN")

	TelegramChatID = int64(env.GetInt("TELEGRAM_CHATID"))
	TelegramToken = env.Get("TELEGRAM_TOKEN")
}
