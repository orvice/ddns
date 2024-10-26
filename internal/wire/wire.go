//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/orvice/ddns/dns"
	"github.com/orvice/ddns/internal/app"
	"github.com/orvice/ddns/internal/config"
	"github.com/orvice/ddns/internal/ip"
	"github.com/orvice/ddns/utils"
)

func NewApp() (*app.App, error) {
	wire.Build(app.New, config.New, dns.New, ip.NewIfGetter, utils.NewLogger)
	return &app.App{}, nil
}
