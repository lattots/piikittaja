package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type Service struct {
	adminStore AdminStore // Database for storing admin information
}

func NewService(cookieStore *sessions.CookieStore, dbURL string) (*Service, error) {
	// Database handle is created for auth service
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	gothic.Store = cookieStore

	callbackURL := buildCallbackURL("google")

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), callbackURL),
	)

	adminDB, err := NewAdminDB(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating admin store: %w", err)
	}

	return &Service{adminStore: adminDB}, nil
}

// GetSessionUser returns the user in the current session. Returns an error if no user is found.
func (s *Service) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, sessionName)
	if err != nil {
		return goth.User{}, err
	}

	usr := session.Values["user"]
	if usr == nil {
		return goth.User{}, errors.New("no user found")
	}

	return usr.(goth.User), nil
}

// SaveSession saves current user to the session.
func (s *Service) SaveSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := gothic.Store.Get(r, sessionName)

	session.Values["user"] = user
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("error saving user session: %s", err), http.StatusInternalServerError)
		return err
	}

	return nil
}

// RemoveSession removes the current session from the session store. This is used at logout to delete existing user.
func (s *Service) RemoveSession(w http.ResponseWriter, r *http.Request) error {
	session, err := gothic.Store.Get(r, sessionName)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting from request session: %s", err), http.StatusInternalServerError)
		return err
	}

	session.Values["user"] = goth.User{} // Session user is set to nil
	session.Options.MaxAge = -1          // Session expires immediately

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("error saving user session: %s", err), http.StatusInternalServerError)
	}
	return nil
}

// RequireAuthentication returns the handlerFunc after user has signed in.
// It acts as a middleware to require users to be admins in order to access the site.
func RequireAuthentication(handlerFunc http.HandlerFunc, auth *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		isAdmin, err := auth.adminStore.IsAdmin(usr.Email)
		if err != nil {
			log.Printf("error checking if user is admin: %s\n", err)
			http.Redirect(w, r, "/login", http.StatusInternalServerError)
			return
		}
		if !isAdmin {
			log.Println("User is not authorized to access this resource")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Println("User is authorized to access this resource:", usr.Email)
		handlerFunc(w, r)
	}
}

// Automatically builds the OAuth2 callback URL with the specified provider.
func buildCallbackURL(provider string) string {
	var url string
	if os.Getenv("ENVIRONMENT") == "local" {
		url = "http://localhost"
	} else {
		url = os.Getenv("HOST_URL")
	}
	return fmt.Sprintf("%s/api/auth/%s/callback", url, provider)
}
