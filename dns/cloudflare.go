package dns

import (
	"github.com/libdns/cloudflare"
	"github.com/orvice/ddns/internal/config"
)

func NewCloudFlare() LibDNS {
	provider := cloudflare.Provider{APIToken: config.GetConfig().CFToken}
	return &provider
}
