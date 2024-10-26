package ip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := NewIfconfigCo().GetIP()
	t.Log(ip, err)
}
