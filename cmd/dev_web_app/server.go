package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lattots/piikittaja/pkg/handler"
)

func main() {
	router := http.NewServeMux()

	host := os.Getenv("HOST_URL")
	if host == "" {
		log.Fatalln("HOST_URL not provided in environment variables")
	}

	h, err := handler.NewHandler(host)
	if err != nil {
		log.Fatalln("error creating new handler:", err)
	}

	// API routes
	router.HandleFunc("GET /users/{userId}", h.GetUserByID)
	router.HandleFunc("GET /users", h.GetUsers)
	router.HandleFunc("GET /users/{userId}/transactions", h.GetUserTransactions)
	router.HandleFunc("POST /users/{userId}/transactions", h.NewTransaction)
	router.HandleFunc("GET /transactions", h.GetTransactions)

	const port = ":8080"
	fmt.Printf("Server started on port %s\n", port)

	if err = http.ListenAndServe(port, router); err != nil {
		log.Fatalln("unexpected error: ", err)
	}
}
