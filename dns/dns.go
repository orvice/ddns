package dns

import "context"

type DNS interface {
	GetIP(ctx context.Context, domain string) (string, error)
	UpdateIP(ctx context.Context, domain, ip string) error
}
