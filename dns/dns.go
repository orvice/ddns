package dns

type DNS interface {
	GetIP(domain string) (string, error)
	UpdateIP(domain, ip string) error
}
