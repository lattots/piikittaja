package user

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

func (u *User) SendMessage(b *bot.Bot, msg string) error {
	ctx := context.Background()
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.ID,
		Text:   msg,
	})
	if err != nil {
		return fmt.Errorf("error sending message: %v\n", err)
	}
	return nil
}
