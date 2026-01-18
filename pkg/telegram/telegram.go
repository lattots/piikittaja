package telegramutil

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func GetSendPhotoParams(userID int64, photoFilepath, message string) (*bot.SendPhotoParams, error) {
	photoFile, err := os.Open(photoFilepath)
	if err != nil {
		return nil, fmt.Errorf("error opening photo file %w", err)
	}

	_, filename := filepath.Split(photoFilepath)

	reader := bufio.NewReader(photoFile)

	photo := &models.InputFileUpload{
		Filename: filename,
		Data:     reader,
	}

	params := &bot.SendPhotoParams{
		ChatID:  userID,
		Photo:   photo,
		Caption: message,
	}

	return params, nil
}
