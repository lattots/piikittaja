package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/lattots/gipher"
	telegramutil "github.com/lattots/piikittaja/pkg/telegram"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		log.Fatalln("error getting database URL from environment variables")
	}

	h, err := newHandler(dbURL)
	if err != nil {
		log.Fatalln("error creating telegram bot handler:", err)
	}

	fmt.Println("Creating bot...")
	b, err := bot.New(os.Getenv("TELEGRAM_TOKEN"), bot.WithDefaultHandler(h.defaultHandler))
	if err != nil {
		log.Fatalln("error creating bot:\n", err)
	}

	var commands []models.BotCommand
	commands = append(commands, models.BotCommand{
		Command:     "/piikki",
		Description: "Näe nykyinen piikki.",
	})
	commands = append(commands, models.BotCommand{
		Command:     "/piikkaa",
		Description: "Lisää piikkiin.",
	})
	commands = append(commands, models.BotCommand{
		Command:     "/terve",
		Description: "Tevre!",
	})
	commands = append(commands, models.BotCommand{
		Command:     "/maksaminen",
		Description: "Maksuohjeet.",
	})
	commands = append(commands, models.BotCommand{
		Command:     "/apua",
		Description: "Apua!",
	})
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: commands})
	if err != nil {
		log.Fatalln("error setting commands for bot:\n", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, h.handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/piikkaa", bot.MatchTypeExact, handleGetAmountInput)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/piikki", bot.MatchTypeExact, h.handleGetBalance)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/terve", bot.MatchTypeExact, handleGreet)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/maksaminen", bot.MatchTypeExact, handlePaymentInfo)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/apua", bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/lisaatg", bot.MatchTypePrefix, h.handleNewTelegramAdmin)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/lisaanetti", bot.MatchTypePrefix, h.handleNewWebAdmin)

	fmt.Println("Bot created successfully")

	b.Start(ctx)
}

func handleGreet(ctx context.Context, b *bot.Bot, update *models.Update) {
	senderUsername := update.Message.Chat.Username

	msg := fmt.Sprintf(
		"Terve vaan, %s", senderUsername,
	)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})

	if err != nil {
		log.Fatalln("error sending message:\n", err)
	}
}

func handleGetAmountInput(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := requestKeyboardInput(ctx, b, update)
	if err != nil {
		log.Fatalln("error requesting keyboard input:\n", err)
	}
}

func handlePaymentInfo(ctx context.Context, b *bot.Bot, update *models.Update) {
	params, err := getSendPhotoParams(update)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = b.SendPhoto(ctx, params)
	if err != nil {
		log.Fatalln("error sending message:", err)
	}
}

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := "Olen PiikkiBotti. Autan killan tärkeimpiä vapaaehtoisia kirjanpitotehtävissä.\n\n" +
		"/piikki: Nähdäksesi nykyisen saldosi.\n" +
		"/piikkaa: Käyttääksesi saldoasi.\n" +
		"/terve: Tervehtiäksesi PiikkiBottia.\n" +
		"/maksaminen: Saadaksesi rahan tallettamiseen liittyvät ohjeet.\n" +
		"/apua: Saadaksesi apua."

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})

	if err != nil {
		log.Fatalln("error sending message:\n", err)
	}
}

func requestKeyboardInput(ctx context.Context, b *bot.Bot, update *models.Update) error {
	keyboard := [][]models.KeyboardButton{
		{
			{Text: "1€"},
			{Text: "2€"},
		},
		{
			{Text: "5€"},
			{Text: "10€"},
		},
	}

	keyboardMarkup, err := json.Marshal(models.ReplyKeyboardMarkup{
		Keyboard:        keyboard,
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	})
	if err != nil {
		return fmt.Errorf("error marshalling keyboard markup: %v", err)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Paljonko piikataan?🤑",
		ReplyMarkup: string(keyboardMarkup),
	})
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func createAnimation(amount, transactionId int) error {
	if !isValidAmount(amount) {
		return fmt.Errorf("error creating animation for amount: %d", amount)
	}

	// Ensure that the tmp/ directory exists
	// This directory is used to temporarily store created animation files
	tmpPath := filepath.Join(".", "assets", "telegram_bot", "tmp")
	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating tmp directory: %w", err)
	}

	backgroundFilename := fmt.Sprintf("./assets/telegram_bot/%d€.gif", amount)
	outputFilename := fmt.Sprintf("./assets/telegram_bot/tmp/%d.gif", transactionId)
	fontFilename := "./assets/telegram_bot/Raleway-Black.ttf"

	err = gipher.CreateTimeStampGIF(backgroundFilename, outputFilename, fontFilename)
	return err
}

func getSendAnimationParams(update *models.Update, transactionId, userBalance int) (*bot.SendAnimationParams, error) {
	animationFile, err := os.Open(fmt.Sprintf("./assets/telegram_bot/tmp/%d.gif", transactionId))
	if err != nil {
		return nil, fmt.Errorf("error opening GIF file with ID %d: %s", transactionId, err)
	}

	reader := bufio.NewReader(animationFile)

	animation := &models.InputFileUpload{
		Filename: "rahaa",
		Data:     reader,
	}

	params := &bot.SendAnimationParams{
		ChatID:    update.Message.Chat.ID,
		Width:     100,
		Height:    100,
		Duration:  1,
		Animation: animation,
		Caption:   fmt.Sprintf("Saldosi on nyt %d€", userBalance),
	}

	return params, nil
}

func getSendPhotoParams(update *models.Update) (*bot.SendPhotoParams, error) {
	msg := "Vai että maksun aika lähestyy...\n\n" +
		"Näin se tapahtuu:\n" +
		"1. Saavu kiltahuoneelle rahat mukanasi\n" +
		"2. Etsi kuvien perusteella postilaatikko ja kirjekuori\n" +
		"3. Sujauta rahat kirjekuoreen ja kirjoita Telegram-käyttäjäsi kuoreen\n" +
		"4. Tiputa kirjekuori postilaatikkoon ja kumarra/niiaa kolmesti\n\n" +
		"Kuten arvata saattaa, maksun käsittelyssä menee joitain päiviä. " +
		"Älä siis hätäile, vaikka piikkisi ei välittömästi kuittaudu maksetuksi."

	photoFile, err := os.Open("./assets/telegram_bot/payment.png")
	if err != nil {
		return nil, fmt.Errorf("error opening photo file %w", err)
	}

	reader := bufio.NewReader(photoFile)

	photo := &models.InputFileUpload{
		Filename: "payment",
		Data:     reader,
	}

	params := &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   photo,
		Caption: msg,
	}

	return params, nil
}

func handleInternalError(ctx context.Context, b *bot.Bot, sender *models.User) {
	msg := "Nyt kävi kämmi. Kokeile kohta uudestaan!"
	err := telegramutil.SendMessage(ctx, b, int64(sender.ID), msg)
	if err != nil {
		log.Printf("error sending error message to user %s: %s", sender.Username, err)
	}
}
