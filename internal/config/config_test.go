package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DOMAIN", "example.com")
	err := LoadConfig(".")
	assert.Nil(t, err)

	t.Log(GetConfig())
	assert.Equal(t, "example.com", GetConfig().Domain)

}
