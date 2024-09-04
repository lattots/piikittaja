package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/lattots/piikittaja/pkg/auth"
	"github.com/lattots/piikittaja/pkg/handler"
)

func main() {
	err := godotenv.Load("./assets/.env")
	if err != nil {
		log.Fatalln("error loading .env file: ", err)
	}

	router := http.NewServeMux()

	// Serve static files from the 'assets/web_app' directory
	staticFileHandler := http.StripPrefix("/assets/web_app/", http.FileServer(http.Dir("./assets/web_app")))
	router.Handle("/assets/web_app/", staticFileHandler)

	h, err := handler.NewHandler()
	if err != nil {
		log.Fatalln("error creating new handler:", err)
	}

	// Views
	router.HandleFunc("/", auth.RequireAuthentication(h.HandleIndex, h.Auth))
	router.HandleFunc("GET /user-view", auth.RequireAuthentication(h.HandleUserView, h.Auth))
	router.HandleFunc("GET /login", h.HandleLogin)

	// Actions
	router.HandleFunc("POST /action", auth.RequireAuthentication(h.HandleUserAction, h.Auth))

	// Auth
	router.HandleFunc("GET /auth/{provider}/callback", h.HandleAuthCallback)
	router.HandleFunc("GET /logout/{provider}", h.HandleLogout)
	router.HandleFunc("GET /auth/{provider}", h.HandleProviderLogin)

	certFilePath := "/etc/letsencrypt/live/piikki.stadi.ninja/fullchain.pem"
	keyFilePath := "/etc/letsencrypt/live/piikki.stadi.ninja/privkey.pem"

	port := ":3000"
	fmt.Printf("Server started on port %s\n", port)

	// If server is run in local test environment, it doesn't use tls
	if os.Getenv("ENVIRONMENT") == "local" {
		if err = http.ListenAndServe(port, router); err != nil {
			log.Fatalln("unexpected error: ", err)
		}
	} else {
		if err = http.ListenAndServeTLS(port, certFilePath, keyFilePath, router); err != nil {
			log.Fatalln("unexpected error: ", err)
		}
	}
}
