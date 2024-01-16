package model

import (
	"encoding/json"
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

// String returns user fields as json string
func (u *User) String() string {
	b, _ := json.MarshalIndent(u, "  ", "  ")
	return string(b)
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
