package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/go-telegram/bot"
	tgmodels "github.com/go-telegram/bot/models"

	amountparser "github.com/lattots/piikittaja/pkg/amount_parser"
	"github.com/lattots/piikittaja/pkg/auth"
	"github.com/lattots/piikittaja/pkg/models"
	telegramutil "github.com/lattots/piikittaja/pkg/telegram"
	"github.com/lattots/piikittaja/pkg/transaction"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

type handler struct {
	usrStore   userstore.UserStore
	traHandler transaction.TransactionHandler
	adminStore auth.AdminStore
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

	adminStore, err := auth.NewAdminDB(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating admin store: %s\n", err)
	}

	return &handler{usrStore: usrStore, traHandler: traHandler, adminStore: adminStore}, nil
}

func (h *handler) defaultHandler(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	sender := update.Message.From
	receivedMessage := update.Message.Text

	amount, err := amountparser.ParseToCents(receivedMessage)
	// if function errors, the message is not an amount, and it should be handled as unknown command
	// if function doesn't error, amount exists, and it should be handled as new tab

	if err != nil || !isValidAmount(amount) {
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

	_, err = h.traHandler.Withdraw(u, amount)
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

	err = sendPaymentConfirmation(ctx, b, u, amount)
	if err != nil {
		log.Printf("error while sending payment confirmation to %s: %s\n", sender.Username, err)
		handleInternalError(ctx, b, sender)
	}

	u.Username = sender.Username
	u.FirstName = sender.FirstName
	u.LastName = sender.LastName
	err = h.usrStore.Update(u)
	if err != nil {
		log.Printf("error updating user %s: %s\n", u.Username, err)
		return
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

	u := models.NewUser(int(sender.ID), sender.Username, sender.FirstName, sender.LastName)
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

	msg := fmt.Sprintf("Saldosi on nyt: %s", amountparser.String(u.Balance))
	err = telegramutil.SendMessage(context.TODO(), b, int64(sender.ID), msg)
	if err != nil {
		log.Printf("error sending error message to user %s: %s", sender.Username, err)
	}

	u.Username = sender.Username
	u.FirstName = sender.FirstName
	u.LastName = sender.LastName
	err = h.usrStore.Update(u)
	if err != nil {
		log.Printf("error updating user %s: %s\n", u.Username, err)
		return
	}
}

func (h *handler) handleNewTelegramAdmin(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	fullText := update.Message.Text
	parts := strings.Fields(fullText)

	// Expected format: /lisaatg username
	if len(parts) != 2 {
		msg := "Käyttö: /lisaatg <käyttäjänimi>"
		err := telegramutil.SendMessage(ctx, b, update.Message.From.ID, msg)
		if err != nil {
			handleInternalError(ctx, b, update.Message.From)
			log.Printf("error sending invalid command message: %s\n", err)
		}
		return
	}

	username := parts[1]

	from := update.Message.From
	sender, err := h.usrStore.GetByID(int(from.ID))
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error getting sender from store: %s\n", err)
		return
	}

	if !sender.IsAdmin {
		msg := "Valitettavasti sinulla ei ole oikeutta luoda uusia pääkäyttäjiä. Womp womp:)"
		err = telegramutil.SendMessage(ctx, b, from.ID, msg)
		if err != nil {
			handleInternalError(ctx, b, from)
			log.Printf("error sending unauthorised message: %s\n", err)
		}
		return
	}

	newAdmin, err := h.usrStore.GetByUsername(username)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error getting new admin user from store: %s\n", err)
		return
	}

	newAdmin.IsAdmin = true

	err = h.usrStore.Update(newAdmin)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error updating new admin user in store: %s\n", err)
		return
	}

	msg := fmt.Sprintf("Lisätty! Käyttäjä \"%s\" on nyt Telegram-botin ylläpitäjä.", newAdmin.Username)
	err = telegramutil.SendMessage(ctx, b, from.ID, msg)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error sending success message to new admin issuer: %s\n", err)
	}
}

func (h *handler) handleNewWebAdmin(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	fullText := update.Message.Text
	parts := strings.Fields(fullText)

	// Expected format: /lisaanetti john@doe.com
	if len(parts) != 2 {
		msg := "Käyttö: /lisaanetti <sähköposti>"
		err := telegramutil.SendMessage(ctx, b, update.Message.From.ID, msg)
		if err != nil {
			handleInternalError(ctx, b, update.Message.From)
			log.Printf("error sending invalid command message: %s\n", err)
		}
		return
	}

	email := parts[1]

	from := update.Message.From
	sender, err := h.usrStore.GetByID(int(from.ID))
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error getting sender from store: %s\n", err)
		return
	}

	if !sender.IsAdmin {
		msg := "Valitettavasti sinulla ei ole oikeutta luoda uusia pääkäyttäjiä. Womp womp:)"
		err = telegramutil.SendMessage(ctx, b, from.ID, msg)
		if err != nil {
			handleInternalError(ctx, b, from)
			log.Printf("error sending unauthorised message: %s\n", err)
		}
		return
	}

	err = h.adminStore.AddAdmin(email)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error adding new web admin: %s\n", err)
		return
	}

	msg := fmt.Sprintf("Lisätty! Käyttäjä \"%s\" on nyt nettisivun ylläpitäjä.", email)
	err = telegramutil.SendMessage(ctx, b, from.ID, msg)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error sending success message to new admin issuer: %s\n", err)
	}
}

func (h *handler) handlePaymentReminder(ctx context.Context, b *bot.Bot, update *tgmodels.Update) {
	from := update.Message.From
	sender, err := h.usrStore.GetByID(int(from.ID))
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error getting sender from store: %s\n", err)
		return
	}

	if !sender.IsAdmin {
		msg := "Sinulla ei ole tuollaisia lupia. Kannattaa ensi vaaleissa harkita IEyttä, jotta pääset tähän klubiin:)"
		err = telegramutil.SendMessage(ctx, b, int64(sender.ID), msg)
		if err != nil {
			handleInternalError(ctx, b, from)
			log.Printf("error sending message to user \"%s\": %s", sender.Username, err)
			return
		}
	}

	users, err := h.usrStore.GetUsers()
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error getting users from store: %s", err)
		return
	}

	msgCounter := 0
	for _, usr := range users {
		// Only users with negative balance are reminded
		if usr.Balance >= 0 {
			continue
		}

		msg := fmt.Sprintf(
			"Hei, %s! On käynyt ilmi, että sitä on rilluteltu, jonka seurauksena saldo on päässyt pakkaselle. "+
				"Maksathan velkasi ensitilassa IE:lle.\n\nNykyinen saldosi: %s\n\n"+
				"Paatuneelta piikittäjältä maksaminen sujuu varmasti jo kuin tanssi, mutta muiden kohdalla suosittelen "+
				"kääntymään ohjeistuksen puoleen komennolla /maksaminen. "+
				"Saldoa on myös mahdollista kerryttää etukäteen, jos luulet, että lähitulevaisuudessa korkki taas aukeaa...",
			usr.Username,
			amountparser.String(usr.Balance),
		)

		err = telegramutil.SendMessage(context.TODO(), b, int64(usr.ID), msg)
		if errors.Is(bot.ErrorForbidden, err) {
			log.Printf("User: %s has probably blocked PiikkiBotti...\nError: %s\n", usr.Username, err.Error())
		} else if err != nil {
			log.Printf("fatal error while sending message to %s: %s\n", usr.Username, err.Error())
		}
		msgCounter++
	}

	msg := fmt.Sprintf("Maksumuistutukset lähetetty %d käyttäjälle:) Fyrkkaa tulossa!", msgCounter)
	err = telegramutil.SendMessage(ctx, b, int64(sender.ID), msg)
	if err != nil {
		handleInternalError(ctx, b, from)
		log.Printf("error sending message to user \"%s\": %s", sender.Username, err)
		return
	}
}

const assetDirectory = "./assets/telegram_bot/"

func sendPaymentConfirmation(ctx context.Context, b *bot.Bot, u *models.User, amount int) error {
	userID := int64(u.ID)
	photo := fmt.Sprintf("%s%d.webp", assetDirectory, amount)
	msg := fmt.Sprintf("Onnistui! Nauti herkuistasi:)\nNykyinen piikkisi on %s", amountparser.String(u.Balance))

	params, err := telegramutil.GetSendPhotoParams(userID, photo, msg)
	if err != nil {
		return err
	}

	_, err = b.SendPhoto(ctx, params)
	return err
}

func isValidAmount(amount int) bool {
	validAmounts := []int{
		1_00, 1_50, 2_00, 3_00, 4_00, 10_00,
	}
	return slices.Contains(validAmounts, amount)
}
