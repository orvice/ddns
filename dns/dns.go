package dns

import (
	"github.com/libdns/libdns"
)

type LibDNS interface {
	libdns.RecordGetter
	libdns.RecordAppender
	libdns.RecordSetter
}
