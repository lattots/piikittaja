package telegramutil

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func SendMessage(ctx context.Context, b *bot.Bot, id int64, msg string) error {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: id,
		Text:   msg,
	})
	return err
}

func SendMessageToAll(ctx context.Context, b *bot.Bot, ids []int64, msg string) error {
	var err error
	for _, id := range ids {
		err = SendMessage(ctx, b, id, msg)
		if err != nil {
			return fmt.Errorf("Error sending message to user with id \"%d\": %w", id, err)
		}
	}

	return nil
}
