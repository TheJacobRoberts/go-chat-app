package types

import (
	"go-chat-app/model"
)

type UserDocument struct {
	ID   string
	User *model.User
}

func NewUserDocument(id string, user *model.User) *UserDocument {
	return &UserDocument{
		ID:   id,
		User: user,
	}
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
