package firestore

import (
	"context"
	"errors"
	"fmt"
	"go-chat-app/service/user"
	identifier "go-chat-app/util"

	"cloud.google.com/go/firestore"
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
	// TODO: doesn't matter for now as it gets overridden anyway
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

	// TODO: Fix this method
	user := users[0]

	_, err := s.client.Collection(_userCollection).Doc("DC").Update(ctx, []firestore.Update{
		{
			Path:  "name",
			Value: upd.Name,
		},
		{
			Path:  "email",
			Value: upd.Email,
		},
	})
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	panic("unimplemented")
}

func findUserByID(ctx context.Context, client *firestore.Client, id int) (*user.User, error) {
	a, _, err := findUsers(ctx, client, &user.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, errors.New("User not found.")
	}
	return a[0], nil
}

func findUsers(ctx context.Context, client *firestore.Client, filter *user.UserFilter) (_ []*user.User, n int, err error) {
	users := make([]*user.User, 0)
	err = client.RunTransaction(ctx, func(ctx context.Context, t *firestore.Transaction) error {
		colRef := client.Collection(_userCollection)

		query := applyFilters(colRef, filter)

		iter := query.Documents(ctx)

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var user user.User

			doc.DataTo(&user)

			users = append(users, &user)
		}

		return nil
	})

	return users, len(users), err
}

func applyFilters(colRef *firestore.CollectionRef, filter *user.UserFilter) (query firestore.Query) {
	query = colRef.Where("id", "==", filter.ID)
	query = colRef.Where("email", "==", filter.Email)
	query.Limit(filter.Limit)
	query.Offset(filter.Offset)

	return query
}
