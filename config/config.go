package config

import "github.com/orvice/utils/env"

var (
	IPNotifyFormat = "[%s] ip changed, old IP: %s new IP: %s"
)

func Init() {
	inf := env.Get("IP_CHANGE_FORMAT")
	if len(inf) != 0 {
		IPNotifyFormat = inf
	}
}
