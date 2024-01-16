package firestore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-chat-app/firestore/types"
	"go-chat-app/service/user_service"
	identifier "go-chat-app/util"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const (
	_userCollection = "user"
)

var _ user_service.UserService = (*UserService)(nil)

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
func (s *UserService) FindUserByID(ctx context.Context, id string) (*user_service.User, error) {
	userDocs := s.client.Collection(_userCollection).Where("id", "==", id)
	iter := userDocs.Documents(ctx)

	users := make([]*user_service.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user_service.User

		doc.DataTo(&user)

		users = append(users, &user)
	}

	return users[0], nil
}

// FindUserByID retrieves a user by ID
func (s *UserService) FindUsers(ctx context.Context) ([]*user_service.User, error) {
	userDocs := s.client.Collection(_userCollection)
	iter := userDocs.Documents(ctx)

	users := make([]*user_service.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user_service.User

		doc.DataTo(&user)

		users = append(users, &user)
	}

	return users, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *user_service.User) error {
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
func (s *UserService) UpdateUser(ctx context.Context, id int, upd user_service.UserUpdate) (*user_service.User, error) {
	userDocs := s.client.Collection(_userCollection).Where("id", "==", id)
	iter := userDocs.Documents(ctx)

	users := make([]*user_service.User, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user_service.User

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

func applyFilters(collection *firestore.CollectionRef, filter user_service.UserFilter) (query firestore.Query) {
	query = collection.Where("id", "==", filter.ID)
	query = collection.Where("email", "==", filter.Email)
	query.Limit(filter.Limit)
	query.Offset(filter.Offset)

	return query
}

// CREATE
// - createUser

// READ
// - findUserByID
// - findUserByEmail
// - findUsers

// UPDATE
// - updateUser

// DELETE
// - deleteUser

// findUsers returns a list of users matching a filter, as well as a count of total matching users
func findUsers(ctx context.Context, client *firestore.Client, filter user_service.UserFilter) (_ []*types.UserDocument, n int, err error) {
	collection := client.Collection(_userCollection)

	query := applyFilters(collection, filter)

	iter := query.Documents(ctx)

	users := make([]*types.UserDocument, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var user types.User

		doc.DataTo(&user)

		users = append(users, types.NewUserDocument(doc.Ref.ID, &user))
	}

	return users, len(users), err
}

// findUserByID fetches a user by their ID
func findUserByID(ctx context.Context, client *firestore.Client, id string) (*types.UserDocument, error) {
	users, _, err := findUsers(ctx, client, user_service.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, errors.New("User not found")
	}
	return users[0], nil
}

// findUserByID fetches a user by their email
func findUserByEmail(ctx context.Context, client *firestore.Client, email string) (*types.UserDocument, error) {
	users, _, err := findUsers(ctx, client, user_service.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, errors.New("User not found")
	}
	return users[0], nil
}

// createUser creates a new user
func createUser(ctx context.Context, client *firestore.Client, user *types.User) error {
	// Validate
	if err := user.Validate(); err != nil {
		return err
	}

	collection := client.Collection(_userCollection)

	_, _, err := collection.Add(ctx, &user_service.User{
		ID:        identifier.NewUUID(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		return err
	}

	return nil
}

// updateUser updates fields on a user object
func updateUser(ctx context.Context, client *firestore.Client, id string, upd user_service.UserUpdate) (*types.User, error) {
	userDoc, err := findUserByID(ctx, client, id)
	if err != nil {
		return userDoc.User, err
	}

	// Update fields
	if v := upd.Name; v != nil {
		userDoc.User.Name = *v
	}
	if v := upd.Email; v != nil {
		userDoc.User.Email = *v
	}

	// Set UpdatedAt field
	userDoc.User.UpdatedAt = time.Now().UTC()

	// Validate
	if err := userDoc.User.Validate(); err != nil {
		return userDoc.User, err
	}

	collection := client.Collection(_userCollection)

	doc := collection.Doc(userDoc.ID)

	doc.Update(ctx, []firestore.Update{
		{
			Path:  "name",
			Value: userDoc.User.Name,
		},
		{
			Path:  "email",
			Value: userDoc.User.Email,
		},
		{
			Path:  "updatedAt",
			Value: time.Now().UTC(),
		},
	})

	if err != nil {
		return userDoc.User, err
	}

	return userDoc.User, nil
}

// findUserByID fetches a user by their ID
func deleteUser(ctx context.Context, client *firestore.Client, id string) error {
	user, err := findUserByID(ctx, client, id)
	if err != nil {
		return err
	}

	if _, err := client.Collection(_userCollection).Doc(user.ID).Delete(ctx); err != nil {
		return err
	}

	return nil
}
