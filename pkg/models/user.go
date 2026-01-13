package models

type User struct {
	ID       int
	Username string
	Balance  int // Balance in cents
	IsAdmin  bool
}

func NewUser(id int, username string) *User {
	return &User{ID: id, Username: username, Balance: 0}
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Balance:  u.Balance,
	}
}
