package firestore

import (
	"context"
	"fmt"
	"go-chat-app/service/user"
	identifier "go-chat-app/util"

	"google.golang.org/api/iterator"
)

const (
	_userCollection = "user"
)

var _ user.UserService = (*UserService)(nil)

// UserService represents a service for managing users
type UserService struct {
	*FirestoreClient
}

// NewUserService returns a new instance of UserService
func NewUserService(client *FirestoreClient) *UserService {
	return &UserService{
		client,
	}
}

// FindUserByID retrieves a user by ID
func (s *UserService) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	userDocs := s.client.Collection(_userCollection).Where("id", "==", id)
	iter := userDocs.Documents(ctx)

	users := make([]*user.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user.User

		doc.DataTo(&user)

		users = append(users, &user)
	}

	return users[0], nil
}

// FindUserByID retrieves a user by ID
func (s *UserService) FindUsers(ctx context.Context) ([]*user.User, error) {
	userDocs := s.client.Collection(_userCollection)
	iter := userDocs.Documents(ctx)

	users := make([]*user.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user.User

		doc.DataTo(&user)

		users = append(users, &user)
	}

	return users, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *user.User) error {
	// TODO: user could possible provide a new uuid - probably need to account for this and make sure they can't
	uuid := identifier.NewUUID()
	user.ID = uuid

	ref, res, err := s.client.Collection(_userCollection).Add(ctx, &user)
	if err != nil {
		return err
	}

	fmt.Println(ref, res)

	return nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id int, upd user.UserUpdate) (*user.User, error) {
	panic("unimplemented")
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	panic("unimplemented")
}
