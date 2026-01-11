package handler

import (
	"log"
	"net/http"

	amountparser "github.com/lattots/piikittaja/pkg/amount_parser"
	"github.com/lattots/piikittaja/pkg/models"
)

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	index := h.tmpl.Lookup("index")
	if index == nil {
		log.Println("error finding index block in template")
		http.Error(w, "index block not found in template", http.StatusInternalServerError)
		return
	}
	if err := index.Execute(w, nil); err != nil {
		log.Println("error executing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleUserView(w http.ResponseWriter, r *http.Request) {
	singleUser := h.tmpl.Lookup("single-user-view")
	if singleUser == nil {
		log.Println("error finding single-user-view block in template")
		http.Error(w, "single-user-view block not found in template", http.StatusInternalServerError)
		return
	}

	var users []*models.User
	var err error

	userFilter := r.URL.Query().Get("username")

	if userFilter == "" {
		users, err = h.usrStore.GetUsers()
		if err != nil {
			log.Println("error fetching users from db", err)
			http.Error(w, "error fetching users", http.StatusInternalServerError)
			return
		}
	} else {
		users, err = h.usrStore.SearchUsers(userFilter)
		if err != nil {
			log.Println("error fetching users from db", err)
			http.Error(w, "error fetching users", http.StatusInternalServerError)
			return
		}
	}

	for _, u := range users {
		displayData := struct {
			ID       int
			Username string
			Balance  string
		}{
			ID:       u.ID,
			Username: u.Username,
			Balance:  amountparser.String(u.Balance),
		}
		err := singleUser.Execute(w, displayData)
		if err != nil {
			log.Println("error executing single user template")
			http.Error(w, "could not execute single user template", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	login := h.tmpl.Lookup("login")
	if login == nil {
		log.Println("error finding login block in template")
		http.Error(w, "login block not found in template", http.StatusInternalServerError)
		return
	}

	err := login.Execute(w, nil)
	if err != nil {
		log.Println("error executing html template")
		http.Error(w, "error executing html template", http.StatusInternalServerError)
		return
	}
}
