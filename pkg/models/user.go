package models

type User struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
	Balance   int // Balance in cents
	IsAdmin   bool
}

func NewUser(id int, username, firstName, lastName string) *User {
	return &User{
		ID:        id,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Balance:   0,
	}
}

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Balance   int    `json:"balance"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Balance:   u.Balance,
	}
}
