package types

import (
	"errors"
	"time"
)

type UserDocument struct {
	ID   string
	User *User
}

func NewUserDocument(id string, user *User) *UserDocument {
	return &UserDocument{
		ID:   id,
		User: user,
	}
}

type User struct {
	ID string `json:"id"`

	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate returns an error if the user contains invalid fields
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user name required")
	}

	return nil
}
