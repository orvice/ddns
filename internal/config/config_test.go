package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DOMAIN", "example.com")
	os.Setenv("TELEGRAM_CHATID", "1234567890")
	err := LoadConfig(".")
	assert.Nil(t, err)

	t.Log(GetConfig())
	assert.Equal(t, "example.com", GetConfig().Domain)
	assert.Equal(t, int64(1234567890), GetConfig().TelegramChatID)
}
