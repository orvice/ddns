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
		return NewCloudFlare()
	case "aliyun":
		return NewAliyun()
	}
	return nil
}

// cloudflare
func NewCloudFlare() LibDNS {
	provider := cloudflare.Provider{APIToken: config.GetConfig().CFToken}
	return &provider
}

// aliyun
func NewAliyun() LibDNS {
	provider := alidns.Provider{
		AccKeyID:     config.GetConfig().AliyunAccessKeyID,
		AccKeySecret: config.GetConfig().AliyunAccessKeySecret,
	}
	return &provider
}
