package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	if u, err := gothic.CompleteUserAuth(w, r); err == nil {
		log.Println("User already logged in:", u)
		h.HandleIndex(w, r)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (h *Handler) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Calling HandleAuthCallback...")
	provider := r.PathValue("provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	usr, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		_, err = fmt.Fprintln(w, err)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	err = h.Auth.SaveSession(w, r, usr)
	if err != nil {
		log.Println(err)
	}

	for name, values := range w.Header() {
		for _, value := range values {
			log.Printf("Header: %s: %s\n", name, value)
		}
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		log.Fatalln(err)
	}
	err = h.Auth.RemoveSession(w, r)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to remove session: %w", err))
	}

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
