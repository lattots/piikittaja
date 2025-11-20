package telegramutil

import (
	"context"

	"github.com/go-telegram/bot"
)

func SendMessage(ctx context.Context, b *bot.Bot, id int64, msg string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: id,
		Text:   msg,
	})
	return err
}
