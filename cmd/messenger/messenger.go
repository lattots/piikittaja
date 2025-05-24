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
		log.Printf("Sending message to %s...\n", u.Username)
		msg := fmt.Sprintf(
			"MUUTOKSIA PIIKKIBOTIN TOIMINTAAN - LUE TÄMÄ!!\n\n" +
				"Koska joihinkin teistä ronteista ei voi luottaa velkojen maksussa, PiikkiBotti toimii jatkossa " +
				"pankkitilin tavoin. Tämä tarkoittaa siis sitä, että ilmaisia lounaita ei enää ole, ja tilin saldoa " +
				"tulee kerryttää etukäteen. PiikkiBotti ei siis enää myönnä lainaa.\n\n" +
				"PiikkiBottia kannattaa kuitenkin edelleen käyttää, sillä käteisen kantaminen ns. \"sucks ass\".",
		)
		err := u.SendMessage(b, msg)
		if errors.Is(bot.ErrorForbidden, err) {
			log.Printf("User: %s has probably blocked PiikkiBotti...\nError: %s\n", u.Username, err)
		} else if err != nil {
			log.Printf("fatal error while sending message to %s: %s\n", u.Username, err)
		}
	}
	fmt.Println("Payment reminders sent!")
}
