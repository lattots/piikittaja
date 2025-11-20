package userstore_test

import (
	"errors"
	"os"
	"slices"
	"testing"

	"github.com/lattots/piikittaja/pkg/models"
	userstore "github.com/lattots/piikittaja/pkg/user_store"
)

func TestInsert(t *testing.T) {
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	const (
		testUserID   = 1
		testUsername = "test user 1"
	)
	user := models.NewUser(testUserID, testUsername)

	err = usrStore.Insert(user)
	if err != nil {
		t.Fatalf("error inserting user to store: %s", err)
	}

	user, err = usrStore.GetByID(testUserID)
	if err != nil || user == nil {
		t.Fatalf("error getting user from store: %s", err)
	}

	err = usrStore.Insert(user)
	if !errors.Is(err, userstore.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists got %s", err)
	}

	err = usrStore.Remove(testUserID)
	if err != nil {
		t.Fatalf("error removing inserted user: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	const (
		testUserID   = 1
		testUsername = "test user 1"
	)
	user := models.NewUser(testUserID, testUsername)

	err = usrStore.Insert(user)
	if err != nil {
		t.Fatalf("error inserting user to store: %s", err)
	}

	const newUsername = "tester user"
	user.Username = newUsername
	err = usrStore.Update(user)
	if err != nil {
		t.Fatalf("error updating user: %s", err)
	}

	user, err = usrStore.GetByID(testUserID)
	if err != nil {
		t.Fatalf("error getting user from store: %s", err)
	}

	if user.Username != newUsername {
		t.Fatalf("expected updated username to be \"%s\" got \"%s\"", newUsername, user.Username)
	}

	err = usrStore.Remove(testUserID)
	if err != nil {
		t.Fatalf("error removing inserted user: %s", err)
	}
}

func TestGetUserByUsername(t *testing.T) {
	// At this point there should be 2 users in the store from previous tests
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	targetIDs := []int{
		123,
		456,
		789,
	}

	targetUsernames := []string{
		"Kanni Kala",
		"Humalainen Haukka",
		"Vodka Kotka",
	}

	for i := range targetUsernames {
		user, err := usrStore.GetByUsername(targetUsernames[i])
		if err != nil {
			t.Fatalf("error fetching existing user from store: %s", err)
		}

		if user.ID != targetIDs[i] {
			t.Fatalf("found wrong ID for user \"%s\": want %d got %d", targetUsernames[i], targetIDs[i], user.ID)
		}
	}
}

func TestGetUsers(t *testing.T) {
	// This test checks that the user store can correctly fetch multiple users
	// At this point there should be 2 users in the store from previous tests
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	targetIDs := []int{
		123,
		456,
		789,
	}

	users, err := usrStore.GetUsers()
	if err != nil {
		t.Fatalf("error fetching users")
	}

	for _, u := range users {
		if u.Username == "" {
			t.Fatalf("found user with no username. User ID: %d", u.ID)
		}
		if !slices.Contains(targetIDs, u.ID) {
			t.Fatalf("user ID \"%d\" not found in target ID's", u.ID)
		}
	}
}

func TestRemove(t *testing.T) {
	dbURL := os.Getenv("DATABASE_APP")
	if dbURL == "" {
		t.Fatal("DATABASE_APP environment variable is not set")
	}

	usrStore, err := userstore.NewMariaDBStore(dbURL)
	if err != nil {
		t.Fatalf("error creating user store: %s", err)
	}
	defer usrStore.Close()

	const targetID = 123

	err = usrStore.Remove(targetID)
	if err != nil {
		t.Fatalf("error removing user with ID %d: %s", targetID, err)
	}

	users, err := usrStore.GetUsers()
	if err != nil {
		t.Fatalf("error getting users from store: %s", err)
	}

	found := false
	for _, u := range users {
		if u.ID == targetID {
			found = true
			break
		}
	}
	if found {
		t.Fatalf("found removed user from store")
	}
}
