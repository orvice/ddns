package notify

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orvice/ddns/internal/config"
)

type Notifier interface {
	Send(ctx context.Context, s string) error
}

// tgSender is the minimal surface TelegramNotifier needs from the bot API.
// Implemented by *tgbotapi.BotAPI; trivially fakable in tests.
type tgSender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

var _ Notifier = (*TelegramNotifier)(nil)

type TelegramNotifier struct {
	bot    tgSender
	chatID int64
	logger *slog.Logger
}

func NewTelegramNotifier(conf *config.Config) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		return nil, err
	}
	return &TelegramNotifier{
		bot:    bot,
		chatID: conf.TelegramChatID,
		logger: slog.Default().With("notifier", "telegram"),
	}, nil
}

// NewTelegramNotifierWithSender builds a TelegramNotifier with an injected
// sender. Intended for tests.
func NewTelegramNotifierWithSender(bot tgSender, chatID int64, logger *slog.Logger) *TelegramNotifier {
	if logger == nil {
		logger = slog.Default()
	}
	return &TelegramNotifier{bot: bot, chatID: chatID, logger: logger}
}

func (t *TelegramNotifier) Send(_ context.Context, s string) error {
	msg := tgbotapi.NewMessage(t.chatID, s)
	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error("send message failed", "chat_id", t.chatID, "error", err)
		return err
	}
	t.logger.Info("message sent", "chat_id", t.chatID)
	return nil
}

// MultiNotifier fans out one Send call to multiple underlying notifiers.
// Errors are logged but not aggregated, matching prior best-effort semantics.
type MultiNotifier struct {
	notifiers []Notifier
}

func NewMultiNotifier(ns ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: ns}
}

func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

func (m *MultiNotifier) Send(ctx context.Context, s string) error {
	for _, n := range m.notifiers {
		_ = n.Send(ctx, s)
	}
	return nil
}
