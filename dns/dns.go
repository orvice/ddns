package dns

import (
	"github.com/libdns/alidns"
	"github.com/libdns/libdns"
	"github.com/orvice/ddns/internal/config"
)

type LibDNS interface {
	libdns.RecordGetter
	libdns.RecordAppender
	libdns.RecordSetter
}

// aliyun
func NewAliyun() LibDNS {
	provider := alidns.Provider{
		AccKeyID:     config.GetConfig().AliyunAccessKeyID,
		AccKeySecret: config.GetConfig().AliyunAccessKeySecret,
	}
	return &provider
}
