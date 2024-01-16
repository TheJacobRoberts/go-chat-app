package model

import "go-chat-app/service/user_service"

type User struct {
	documentID string
	User       *user_service.User
}

func NewUser(documentID string, user *user_service.User) *User {
	return &User{
		documentID: documentID,
		User:       user,
	}
}

func (u *User) GetDocumentID() string {
	return u.documentID
}
