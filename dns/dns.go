package dns

import (
	"github.com/libdns/alidns"
	"github.com/libdns/cloudflare"
	"github.com/libdns/libdns"
	"github.com/orvice/ddns/internal/config"
)

type LibDNS interface {
	libdns.RecordGetter
	libdns.RecordAppender
	libdns.RecordSetter
}

func New(conf *config.Config) LibDNS {
	switch conf.DNSProvider {
	case "cloudflare":
		return NewCloudFlare(conf)
	case "aliyun":
		return NewAliyun(conf)
	}
	return nil
}

// cloudflare
func NewCloudFlare(conf *config.Config) LibDNS {
	provider := cloudflare.Provider{APIToken: conf.CFToken}
	return &provider
}

// aliyun
func NewAliyun(conf *config.Config) LibDNS {
	provider := alidns.Provider{
		AccKeyID:     conf.AliyunAccessKeyID,
		AccKeySecret: conf.AliyunAccessKeySecret,
	}
	return &provider
}
