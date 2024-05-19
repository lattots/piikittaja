package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"github.com/lattots/gipher"
	"github.com/lattots/piikittaja/src/user"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
)

func main() {
	err := godotenv.Load("../../data/.env")
	if err != nil {
		log.Fatalln("error loading .env file: ", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Println("Creating bot...")
	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), bot.WithDefaultHandler(defaultHandler))
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
		Description: "Terve!",
	})
	commands = append(commands, models.BotCommand{
		Command:     "/apua",
		Description: "Apua!",
	})
	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{Commands: commands})
	if err != nil {
		log.Fatalln("error setting commands for bot:\n", err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handleStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/piikkaa", bot.MatchTypeExact, handleGetAmountInput)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/piikki", bot.MatchTypeExact, handleGetTab)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/terve", bot.MatchTypeExact, handleGreet)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/apua", bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handleHelp)

	fmt.Println("Bot created successfully")

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	sender := update.Message.From
	receivedMessage := update.Message.Text

	amount, err := getAmount(receivedMessage)

	var response string
	// if function errors, the message is not an amount, and it should be handled as unknown command
	// if function doesn't error, amount exists, and it should be handled as new tab
	if err == nil {
		usr, err := user.NewUser(int(sender.ID), sender.Username)
		if err != nil {
			log.Fatalln(err)
		}

		transactionId, err := usr.AddToTab(amount)
		if err != nil {
			log.Fatalln(err)
		}

		err = createAnimation(amount, transactionId)
		if err != nil {
			log.Fatalln(err)
		}

		params, err := getSendAnimationParams(update, transactionId)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = b.SendAnimation(ctx, params)
		if err != nil {
			log.Fatalln(err)
		}

		balance, err := usr.GetBalance()
		if err != nil {
			log.Fatalln(err)
		}

		response = fmt.Sprintf("Saldosi on nyt %d€", balance)
	} else {
		response = "En ymmärtänyt tuota. Kirjoita /apua saadaksesi apua."
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   response,
	})
	if err != nil {
		log.Fatalln("error sending message:\n", err)
	}

}

func handleStart(ctx context.Context, b *bot.Bot, update *models.Update) {
	senderUsername := update.Message.Chat.Username

	msg := fmt.Sprintf(
		"Hyvää päivää, %s. Olet avannut PiikkiBotin. Onnittelut erinomaisesta valinnasta!",
		senderUsername,
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

func handleGetTab(ctx context.Context, b *bot.Bot, update *models.Update) {
	sender := update.Message.From
	u, err := user.NewUser(int(sender.ID), sender.Username)
	if err != nil {
		log.Fatalln(err)
	}

	var response string

	tab, err := u.GetBalance()
	if err != nil {
		log.Println(err)
		response = "En löytänyt piikkiäsi..."
	} else {
		response = fmt.Sprintf("Saldosi on nyt: %d€", tab)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   response,
	})

	if err != nil {
		log.Fatalln("error sending message:\n", err)
	}
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

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	msg := "Olen PiikkiBotti. Autan killan tärkeimpiä vapaaehtoisia kirjanpitotehtävissä.\n\n" +
		"/piikki: Nähdäksesi nykyisen piikkisi.\n" +
		"/piikkaa: Lisätäksesi haluamasi summa piikkiin.\n" +
		"/terve: Tervehtiäksesi PiikkiBottia.\n" +
		"/apua: Saadaksesi apua."

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})

	if err != nil {
		log.Fatalln("error sending message:\n", err)
	}
}

func getAmount(s string) (int, error) {
	re := regexp.MustCompile(`^\d+`)
	match := re.FindString(s)

	if match == "" {
		return 0, errors.New("input doesn't contain amount")
	}

	// try to convert string to int
	// if string can't be converted, function errors
	amount, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func addToUserTab(sender *models.User, amount int) (transactionId int, err error) {
	u, err := user.NewUser(int(sender.ID), sender.Username)
	if err != nil {
		return 0, err
	}

	transactionId, err = u.AddToTab(amount)
	if err != nil {
		return 0, err
	}

	return transactionId, err
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
	var backgroundFilename string
	switch amount {
	case 1:
		backgroundFilename = "./static/1€.gif"
	case 2:
		backgroundFilename = "./static/2€.gif"
	case 5:
		backgroundFilename = "./static/5€.gif"
	case 10:
		backgroundFilename = "./static/10€.gif"
	}

	outputFilename := fmt.Sprintf("%d.gif", transactionId)
	const fontFilename = "./static/Raleway-Black.ttf"

	err := gipher.CreateTimeStampGIF(backgroundFilename, outputFilename, fontFilename)
	return err
}

func getSendAnimationParams(update *models.Update, transactionId int) (*bot.SendAnimationParams, error) {
	animationFile, err := os.Open(fmt.Sprintf("%d.gif", transactionId))
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
	}

	return params, nil
}
