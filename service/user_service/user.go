package user_service

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
	FindUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// UserFilter represents a filter passed to FindUsers()
type UserFilter struct {
	// Filtering fields.
	ID    *string `json:"id"`
	Email *string `json:"email"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}
