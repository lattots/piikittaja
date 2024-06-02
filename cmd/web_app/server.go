package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lattots/piikittaja/pkg/user"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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

	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/action", handleUserAction)
	router.HandleFunc("/user-view", handleUserView)

	certFilePath := "/etc/letsencrypt/live/piikki.stadi.ninja/fullchain.pem"
	keyFilePath := "/etc/letsencrypt/live/piikki.stadi.ninja/privkey.pem"

	port := ":3000"
	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServeTLS(port, certFilePath, keyFilePath, router); err != nil {
		log.Fatalln("unexpected error: ", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	indexTempl, err := template.ParseFiles("./assets/web_app/html/index.html")
	if err != nil {
		log.Println("error parsing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Define a template with the name "index" to ensure only that block is executed
	tmpl := indexTempl.Lookup("index")
	if tmpl == nil {
		log.Println("error finding index block in template")
		http.Error(w, "index block not found in template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		log.Println("error executing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Action struct {
	Username string
	Action   string
	Amount   int
}

func handleUserAction(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/web_app/html/index.html")
	if err != nil {
		log.Println("error parsing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actionForm := tmpl.Lookup("action-form")
	if actionForm == nil {
		log.Println("error finding action-form block in template")
		http.Error(w, "action-form block not found in template", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println("error parsing form")
	}

	var actionStatus string

	amount, err := strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		log.Println("error converting amount to integer")
		actionStatus = "Jokin meni pieleen..."
	} else {
		actionStatus = "Onnistui!"
	}
	action := Action{
		Username: r.FormValue("username"),
		Action:   r.FormValue("action-type"),
		Amount:   amount,
	}

	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		log.Println("error getting database URL", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	// Database handle is created for user.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println("error connecting to database: ", err)
		http.Error(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	usr, err := user.GetUserByUsername(db, action.Username)
	if err != nil {
		log.Println("error fetching user from database", err)
		http.Error(w, "error fetching user", http.StatusInternalServerError)
		return
	}

	if action.Action == "borrow" {
		_, err = usr.AddToTab(action.Amount)
	} else if action.Action == "pay" {
		_, err = usr.PayBackTab(action.Amount)
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

	if err := actionForm.Execute(w, data); err != nil {
		log.Println("error executing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUserView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/web_app/html/index.html")
	if err != nil {
		log.Println("error parsing html template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	singleUser := tmpl.Lookup("single-user-view")
	if singleUser == nil {
		log.Println("error finding single-user-view block in template")
		http.Error(w, "single-user-view block not found in template", http.StatusInternalServerError)
		return
	}

	// Database URL is read from environment variables.
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		log.Println("error getting database URL", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	// Database handle is created for user.
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println("error connecting to database: ", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	users, err := user.GetUsers(db)
	if err != nil {
		log.Println("error fetching users from db", err)
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
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
