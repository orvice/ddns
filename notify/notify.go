package notify

import (
	"context"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/orvice/ddns/internal/config"
)

var notifiers []Notifier

type Notifier interface {
	Send(ctx context.Context, s string) error
}

var _ Notifier = new(TelegramNotifier)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	token  string
	chatID int64
}

func NewTelegramNotifier(conf *config.Config) (Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		return nil, err
	}
	return &TelegramNotifier{
		bot:    bot,
		token:  conf.TelegramToken,
		chatID: conf.TelegramChatID,
	}, nil
}

func (t *TelegramNotifier) Send(_ context.Context, s string) error {
	msg := tgbotapi.NewMessage(t.chatID, s)
	resp, err := t.bot.Send(msg)
	if err != nil {
		slog.Error("send message error", "error", err.Error())
		return err
	}
	slog.Info("send message success", "resp", resp)
	return nil
}

func Init() {
	notifiers = make([]Notifier, 0)
}

func AddNotifier(n Notifier) {
	notifiers = append(notifiers, n)
}

func Notify(ctx context.Context, s string) {
	for _, n := range notifiers {
		_ = n.Send(ctx, s)
	}
}
