package ip

import "testing"

func TestIfconfigCo(t *testing.T) {
	cli := NewIfconfigCo()
	ip, err := cli.GetIP()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(ip)
}
