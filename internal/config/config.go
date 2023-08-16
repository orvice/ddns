package config

import "github.com/orvice/utils/env"

var (
	DOMAIN      string
	UPDATE_TIME int
)

const (
	DNS_MODE_MU = "mu"
)

var (
	DNS_MODE string

	NODE_ID int
	API_URI string

	TELEGRAM_TOKEN  string
	TELEGRAM_CHATID int64
)

func GetConfigFromEnv() {
	DNS_MODE = env.Get("DNS_MODE", "cf")

	DOMAIN = env.Get("DOMAIN")
	UPDATE_TIME = env.GetInt("UPDATE_TIME", 300)

	NODE_ID = env.GetInt("MU_NODE_ID")
	API_URI = env.Get("API_URI")

	TELEGRAM_CHATID = int64(env.GetInt("TELEGRAM_CHATID"))
	TELEGRAM_TOKEN = env.Get("TELEGRAM_TOKEN")
}
