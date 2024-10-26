package ip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := NewIfGetter().GetIP()
	t.Log(ip, err)
}
