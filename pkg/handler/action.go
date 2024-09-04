package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/lattots/piikittaja/pkg/user"
)

func (h *Handler) HandleUserAction(w http.ResponseWriter, r *http.Request) {
	actionForm := h.tmpl.Lookup("action-form")
	if actionForm == nil {
		log.Println("error finding action-form block in template")
		http.Error(w, "action-form block not found in template", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form")
	}

	var actionStatus string

	username := r.Form.Get("username")
	action := r.Form.Get("action-type")
	amount, err := strconv.Atoi(r.Form.Get("amount"))

	if err != nil {
		log.Println("error converting amount to integer")
		actionStatus = "Jokin meni pieleen..."
	} else {
		actionStatus = "Onnistui!"
	}

	usr, err := user.GetUserByUsername(h.DB, username)
	if err != nil {
		log.Println("error fetching user from database", err)
		http.Error(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	if action == "borrow" {
		_, err = usr.AddToTab(amount)
	} else if action == "pay" {
		_, err = usr.PayBackTab(amount)
	} else {
		log.Println("unknown action", err)
		http.Error(w, "unknown action", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Println("error adding to users tab", err)
		http.Error(w, "error adding to users tab", http.StatusInternalServerError)
		return
	}

	// Create a map to hold the status message
	data := map[string]interface{}{
		"ActionStatus": actionStatus,
	}

	// HX-trigger in the header triggers an action called "newUserAction".
	// This is used to refresh other parts of the UI.
	w.Header().Set("HX-Trigger", "newUserAction")
	if err := actionForm.Execute(w, data); err != nil {
		log.Println("error executing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
