package handler

import (
	"log"
	"net/http"

	"github.com/lattots/piikittaja/pkg/user"
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

	users := make([]user.User, 0)

	userFilter := r.URL.Query().Get("username")

	if userFilter == "" {
		var err error
		users, err = user.GetUsers(h.DB)
		if err != nil {
			log.Println("error fetching users from db", err)
			http.Error(w, "error fetching users", http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		users, err = user.SearchUsers(h.DB, userFilter)
		if err != nil {
			log.Println("error fetching users from db", err)
			http.Error(w, "error fetching users", http.StatusInternalServerError)
			return
		}
	}

	for _, u := range users {
		err := singleUser.Execute(w, u)
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
