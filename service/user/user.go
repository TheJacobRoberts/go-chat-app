package user

import (
	"context"
	"errors"
	"time"
)

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

type UserService interface {
	FindUserByID(ctx context.Context, id string) (*User, error)
}


