package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram/bot"

	"github.com/lattots/piikittaja/pkg/user"
)

func main() {
	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		log.Fatalln("error getting database URL from environment variables")
	}

	// Database handle is created.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalln("error opening database connection:", err)
	}

	fmt.Println("Creating bot...")
	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatalln("error creating bot:\n", err)
	}

	users, err := user.GetUsers(db)
	if err != nil {
		log.Fatalln("error fetching users:", err)
	}

	for _, u := range users {
		if u.Balance < 0 {
			log.Printf("Sending message to %s...\n", u.Username)
			msg := fmt.Sprintf(
				"Hei, %s! On käynyt ilmi, että sitä on rilluteltu, jonka seurauksena saldo on päässyt pakkaselle. "+
					"Maksathan velkasi ensitilassa IE:lle.\n\nNykyinen saldosi: %d\n\n"+
					"Paatuneelta piikittäjältä maksaminen sujuu varmasti jo kuin tanssi, mutta muiden kohdalla suosittelen "+
					"kääntymään ohjeistuksen puoleen komennolla /maksaminen. "+
					"Saldoa on myös mahdollista kerryttää etukäteen, jos luulet, että lähitulevaisuudessa korkki taas aukeaa...",
				u.Username,
				u.Balance,
			)
			err := u.SendMessage(b, msg)
			if errors.Is(bot.ErrorForbidden, err) {
				log.Printf("User: %s has probably blocked PiikkiBotti...\nError: %s\n", u.Username, err.Error())
			} else if err != nil {
				log.Printf("fatal error while sending message to %s: %s\n", u.Username, err.Error())
			}
		}
	}
	fmt.Println("Payment reminders sent!")
}
