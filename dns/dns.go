package dns

import (
	"fmt"

	"github.com/libdns/alidns"
	"github.com/libdns/cloudflare"
	"github.com/libdns/libdns"
	"github.com/orvice/ddns/internal/config"
)

const (
	ProviderCloudflare = "cloudflare"
	ProviderAliyun     = "aliyun"
)

type LibDNS interface {
	libdns.RecordGetter
	libdns.RecordAppender
	libdns.RecordSetter
}

// New returns the DNS provider selected by conf.DNSProvider. An empty value
// defaults to Cloudflare; an unknown value returns an error so misconfiguration
// fails loudly.
func New(conf *config.Config) (LibDNS, error) {
	switch conf.DNSProvider {
	case "", ProviderCloudflare:
		return NewCloudFlare(conf), nil
	case ProviderAliyun:
		return NewAliyun(conf), nil
	default:
		return nil, fmt.Errorf("unknown dns provider: %q", conf.DNSProvider)
	}
}

func NewCloudFlare(conf *config.Config) LibDNS {
	provider := cloudflare.Provider{APIToken: conf.CFToken}
	return &provider
}

func NewAliyun(conf *config.Config) LibDNS {
	provider := alidns.Provider{
		AccKeyID:     conf.AliyunAccessKeyID,
		AccKeySecret: conf.AliyunAccessKeySecret,
	}
	return &provider
}
