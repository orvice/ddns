package ip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIP(t *testing.T) {
	ip, err := NewIfconfigCo().GetIP()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)
}
