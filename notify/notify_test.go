package notify

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fakeBot struct {
	sent []tgbotapi.Chattable
	err  error
}

func (f *fakeBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.sent = append(f.sent, c)
	return tgbotapi.Message{}, f.err
}

func TestTelegramNotifier_Send(t *testing.T) {
	bot := &fakeBot{}
	n := NewTelegramNotifierWithSender(bot, 42, nil)

	if err := n.Send(context.Background(), "hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bot.sent) != 1 {
		t.Fatalf("expected 1 send, got %d", len(bot.sent))
	}
	msg, ok := bot.sent[0].(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", bot.sent[0])
	}
	if msg.ChatID != 42 || msg.Text != "hello" {
		t.Errorf("unexpected message: chatID=%d text=%q", msg.ChatID, msg.Text)
	}
}

func TestTelegramNotifier_SendError(t *testing.T) {
	n := NewTelegramNotifierWithSender(&fakeBot{err: errors.New("boom")}, 1, nil)
	if err := n.Send(context.Background(), "x"); err == nil {
		t.Fatal("expected error")
	}
}

type captureNotifier struct{ msgs []string }

func (c *captureNotifier) Send(_ context.Context, s string) error {
	c.msgs = append(c.msgs, s)
	return nil
}

func TestMultiNotifier_FansOut(t *testing.T) {
	a, b := &captureNotifier{}, &captureNotifier{}
	m := NewMultiNotifier(a, b)
	_ = m.Send(context.Background(), "hi")
	if len(a.msgs) != 1 || len(b.msgs) != 1 {
		t.Errorf("expected each notifier to receive 1 msg, got a=%v b=%v", a.msgs, b.msgs)
	}
}
