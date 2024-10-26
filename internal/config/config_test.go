package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DOMAIN", "example.com")
	os.Setenv("TELEGRAM_CHATID", "1234567890")
	conf, err := New()
	assert.Nil(t, err)

	t.Log(conf)
	assert.Equal(t, "example.com", conf.Domain)
	assert.Equal(t, int64(1234567890), conf.TelegramChatID)
}
