package dns

import (
	"context"
	"os"
	"testing"

	"github.com/libdns/cloudflare"
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
