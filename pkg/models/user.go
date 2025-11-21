package models

type User struct {
	ID       int
	Username string
	Balance  int
	IsAdmin  bool
}

func NewUser(id int, username string) *User {
	return &User{ID: id, Username: username, Balance: 0}
}
