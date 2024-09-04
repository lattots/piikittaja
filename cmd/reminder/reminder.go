package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"

	"github.com/lattots/piikittaja/pkg/user"
)

func main() {
	err := godotenv.Load("./assets/.env")
	if err != nil {
		log.Fatalln("error loading .env file: ", err)
	}

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
			fmt.Println("Sending message...")
			msg := fmt.Sprintf(
				"Hei, %s! On käynyt ilmi, että sitä on rilluteltu, jonka seurauksena saldo on päässyt pakkaselle. "+
					"Maksathan velkasi ensitilassa IE:lle.\n\nNykyinen saldosi: %d\n\n"+
					"Velka on helpointa maksaa ottamalla sopiva rahasumma mukaan seuraavan kerran, kun luulet tapaavasi toista IE:tä. "+
					"Saldoa on myös mahdollista kerryttää etukäteen, jos luulet, että lähitulevaisuudessa korkki taas aukeaa...",
				u.Username,
				u.Balance,
			)
			err := u.SendMessage(b, msg)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	fmt.Println("Payment reminders sent!")

}
