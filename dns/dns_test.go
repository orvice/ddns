package dns

import (
	"context"
	"os"
	"testing"

	"github.com/libdns/cloudflare"
	"github.com/orvice/ddns/internal/config"
)

func TestCloudFlare(t *testing.T) {
	provider := cloudflare.Provider{APIToken: os.Getenv("CF_TOKEN")}
	records, err := provider.GetRecords(context.Background(), os.Getenv("CF_ZONE"))
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(records)
}

func TestNew_DefaultsToCloudflare(t *testing.T) {
	d, err := New(&config.Config{})
	if err != nil || d == nil {
		t.Fatalf("expected default cloudflare provider, got err=%v d=%v", err, d)
	}
}

func TestNew_Aliyun(t *testing.T) {
	d, err := New(&config.Config{DNSProvider: "aliyun"})
	if err != nil || d == nil {
		t.Fatalf("expected aliyun provider, got err=%v d=%v", err, d)
	}
}

func TestNew_UnknownProvider(t *testing.T) {
	if _, err := New(&config.Config{DNSProvider: "nope"}); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}
