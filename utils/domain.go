package utils

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func GetDomainSuffix(domain string) string {
	s, _ := publicsuffix.Domain(domain)
	return s
}
