package dns

import (
	"os"

	"github.com/libdns/cloudflare"
)

func NewCloudFlare() LibDNS {
	provider := cloudflare.Provider{APIToken: os.Getenv("CF_TOKEN")}
	return &provider
}
