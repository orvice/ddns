package main

import "github.com/orvice/utils/env"

var (
	DOMAIN      string
	UPDATE_TIME int
)

var (
	CF_API_KEY   string
	CF_API_EMAIL string
)

func GetConfigFromEnv() {
	DOMAIN = env.Get("DOMAIN")
	UPDATE_TIME = env.GetInt("UPDATE_TIME", 300)

	CF_API_KEY = env.Get("CF_API_KEY")
	CF_API_EMAIL = env.Get("CF_API_EMAIL")
}
