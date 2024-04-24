package notify

import (
	"context"
	"log/slog"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
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

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramNotifier{
		bot:    bot,
		token:  token,
		chatID: chatID,
	}, nil
}

func (t *TelegramNotifier) Send(ctx context.Context, s string) error {
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
