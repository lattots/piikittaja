package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/lattots/piikittaja/pkg/auth"
)

func main() {
	err := godotenv.Load("./assets/.env")
	if err != nil {
		log.Fatalln("error loading environment variables: ", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Admin Manager - Manage Access Rights to Piikki Web App")

	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		log.Fatalln("error getting database URL from environment variables")
	}

	// Database handle is created.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatalln("error connecting to the database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("error pinging database: %w", err)
	}

	adminDB := auth.NewAdminDB(db)

	for {
		fmt.Print("Enter command:\n")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1) // Trailing newline is removed

		switch text {
		case "exit":
			fmt.Println("Terminating...")
			os.Exit(0)
		case "check":
			email := getEmail(reader)
			if adminDB.IsAdmin(email) {
				fmt.Println("User is already an admin")
			} else {
				fmt.Println("User is not an admin")
			}
		case "add":
			email := getEmail(reader)
			err = adminDB.AddAdmin(email)
			if err != nil {
				log.Fatalln("error adding admin to database: %w", err)
			}
		case "remove":
			email := getEmail(reader)
			err = adminDB.RemoveAdmin(email)
			if err != nil {
				log.Fatalln("error removing admin from database: %w", err)
			}
		default:
			fmt.Println("Invalid command")
		}
	}
}

func getEmail(reader *bufio.Reader) string {
	fmt.Print("Enter email address:\n")
	email, _ := reader.ReadString('\n')
	email = strings.Replace(email, "\n", "", -1) // Trailing newline is removed

	return email
}
