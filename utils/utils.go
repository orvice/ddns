package utils

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func GetDomainSuffix(domain string) string {
	s, _ := publicsuffix.Domain(domain)
	return s
}

func NewLogger() *slog.Logger {
	w := os.Stderr
	logger := slog.New(
		tint.NewHandler(w, &tint.Options{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		}),
	)
	return logger
}
