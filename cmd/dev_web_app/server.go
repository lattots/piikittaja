package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/lattots/piikittaja/pkg/handler"
)

func main() {
	router := http.NewServeMux()

	// Serve static files from the 'assets/web_app' directory
	staticFileHandler := http.StripPrefix("/assets/web_app/", http.FileServer(http.Dir("./assets/web_app")))
	router.Handle("/assets/web_app/", staticFileHandler)

	h, err := handler.NewHandler()
	if err != nil {
		log.Fatalln("error creating new handler:", err)
	}

	// Views
	router.HandleFunc("/", h.HandleIndex)
	router.HandleFunc("GET /user-view", h.HandleUserView)
	router.HandleFunc("GET /login", h.HandleLogin)

	// Actions
	router.HandleFunc("POST /action", h.HandleUserAction)

	const port = ":3000"
	fmt.Printf("Server started on port %s\n", port)

	if err = http.ListenAndServe(port, router); err != nil {
		log.Fatalln("unexpected error: ", err)
	}
}
