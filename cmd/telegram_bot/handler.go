package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"

	"github.com/lattots/piikittaja/pkg/models"
	telegramutil "github.com/lattots/piikittaja/pkg/telegram"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

type handler struct {
	usrStore   userstore.UserStore
	traHandler transaction.TransactionHandler
}

func newHandler(dbURL string) (*handler, error) {
	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating user store: %w", err)
	}
	traStore, err := transaction.NewMariaDBStore(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction store: %w", err)
	}
	traHandler := transaction.NewTransactionHandler(traStore)

	return &handler{usrStore: usrStore, traHandler: traHandler}, nil
}

func (h *handler) defaultHandler(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	sender := update.Message.From
	receivedMessage := update.Message.Text

	amount, err := getAmount(receivedMessage)
	// if function errors, the message is not an amount, and it should be handled as unknown command
	// if function doesn't error, amount exists, and it should be handled as new tab

	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "En ymmärtänyt tuota. Kirjoita /apua saadaksesi apua.",
		})
		if err != nil {
			log.Fatalln("error sending message:\n", err)
		}
		return
	}

	exists, err := h.usrStore.Exists(int(sender.ID))
	if err != nil {
		log.Printf("error checking if user %s exists: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
		return
	}

	// If user doesn't already exist, new account is created to database.
	// After this the function continues normally
	if !exists {
		h.handleStart(ctx, b, update)
	}

	u, err := h.usrStore.GetByID(int(sender.ID))
	if err != nil {
		log.Printf("error getting user from store: %s", err)
		handleInternalError(ctx, b, sender)
		return
	}

	transactionId, err := h.traHandler.Withdraw(u, amount)
	if errors.Is(err, transaction.ErrNotEnoughBalance) {
		msg := "Tili ammottaa tyhjyyttään :O\n\nMene töihin!"
		log.Println(err)
		err := telegramutil.SendMessage(context.TODO(), b, int64(sender.ID), msg)
		if err != nil {
			log.Printf("error while sending message to %s: %s\n", sender.Username, err)
		}
		return
	} else if err != nil {
		log.Printf("error while sending message to %s: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
		return
	}

	err = createAnimation(amount, transactionId)
	if err != nil {
		log.Fatalln(err)
	}

	params, err := getSendAnimationParams(update, transactionId, u.Balance)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = b.SendAnimation(ctx, params)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Remove(fmt.Sprintf("./assets/telegram_bot/tmp/%d.gif", transactionId))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h *handler) handleStart(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	sender := update.Message.From

	msg := fmt.Sprintf(
		"Hyvää päivää, %s. Olet avannut PiikkiBotin. Onnittelut erinomaisesta valinnasta!\n\n"+
			"Olet sitten kokenut piikittäjä tai portista astuva noviisi, saat apua kirjoittamalla /apua\n\n"+
			"PiikkiBotti toimii kuin henkilökohtainen pankkitili, jolle voit tallettaa rahaa seuraavasti /maksaminen\n\n",
		sender.Username,
	)
	err := telegramutil.SendMessage(context.TODO(), b, int64(sender.ID), msg)
	if err != nil {
		log.Printf("fatal error while sending message to %s: %s\n", sender.Username, err)
	}

	exists, err := h.usrStore.Exists(int(sender.ID))
	if err != nil {
		log.Printf("error checking if user %s exists: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
		return
	}
	// If user already has an account, function returns as no new user needs to be created
	if exists {
		return
	}

	u := models.NewUser(int(sender.ID), sender.Username)
	err = h.usrStore.Insert(u)
	if err != nil {
		log.Printf("error checking if user %s exists: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
		return
	}
}

func (h *handler) handleGetBalance(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	sender := update.Message.From

	exists, err := h.usrStore.Exists(int(sender.ID))
	if err != nil {
		log.Printf("error checking if user %s exists: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
		return
	}

	if !exists {
		h.handleStart(ctx, b, update)
		return
	}

	u, err := h.usrStore.GetByID(int(sender.ID))
	if err != nil {
		log.Printf("error getting user from store: %s", err)
		handleInternalError(ctx, b, sender)
		return
	}

	msg := fmt.Sprintf("Saldosi on nyt: %d€", u.Balance)
	err = telegramutil.SendMessage(context.TODO(), b, int64(sender.ID), msg)
	if err != nil {
		log.Printf("error sending error message to user %s: %s", sender.Username, err)
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

	// Inputted amount is validated.
	if !isValidAmount(amount) {
		// If amount is not valid, function errors.
		return 0, fmt.Errorf("amount is not valid: %d", amount)
	}

	return amount, nil
}

func isValidAmount(amount int) bool {
	validAmounts := []int{
		1, 2, 5, 10,
	}
	isValid := false
	for _, validAmount := range validAmounts {
		if amount == validAmount {
			isValid = true
			break
		}
	}
	return isValid
}
