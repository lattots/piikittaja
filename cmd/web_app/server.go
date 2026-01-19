package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lattots/piikittaja/pkg/auth"
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
	router.HandleFunc("GET /users/{userId}", auth.RequireAuthentication(h.GetUserByID, h.Auth))
	router.HandleFunc("GET /users", auth.RequireAuthentication(h.GetUsers, h.Auth))
	router.HandleFunc("GET /users/{userId}/transactions", auth.RequireAuthentication(h.GetTransactions, h.Auth))
	router.HandleFunc("POST /transactions", auth.RequireAuthentication(h.NewTransaction, h.Auth))

	// Auth
	router.HandleFunc("GET /auth/{provider}/callback", h.HandleAuthCallback)
	router.HandleFunc("GET /logout/{provider}", h.HandleLogout)
	router.HandleFunc("GET /auth/{provider}", h.HandleProviderLogin)

	const port = ":8080"
	fmt.Printf("Server started on port %s\n", port)

	if err = http.ListenAndServe(port, router); err != nil {
		log.Fatalln("unexpected error: ", err)
	}
}
