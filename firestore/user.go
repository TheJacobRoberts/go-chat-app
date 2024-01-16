package firestore

import (
	"context"
	"errors"
	"time"

	"go-chat-app/firestore/types"
	"go-chat-app/model"
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
func (s *UserService) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	userDoc, err := findUserByID(ctx, s.client, id)
	if err != nil {
		return nil, err
	}

	return userDoc.User, nil
}

// FindUserByID retrieves a user by ID
func (s *UserService) FindUsers(ctx context.Context, filter model.UserFilter) ([]*model.User, error) {
	users := make([]*model.User, 0)

	userDocs, _, err := findUsers(ctx, s.client, filter)
	if err != nil {
		return nil, err
	}
	for _, userDoc := range userDocs {
		users = append(users, userDoc.User)
	}

	return users, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	if err := createUser(ctx, s.client, user); err != nil {
		return err
	}
	return nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id string, update model.UserUpdate) (*model.User, error) {
	user, err := updateUser(ctx, s.client, id, update)
	if err != nil {
		return nil, err
	}
	return user, err
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	if err := deleteUser(ctx, s.client, id); err != nil {
		return err
	}
	return nil
}

func applyFilters(collection *firestore.CollectionRef, filter model.UserFilter) firestore.Query {
	query := collection.Query

	// Apply filters
	if v := filter.ID; v != nil {
		query = collection.Where("id", "==", filter.ID)
	}
	if v := filter.Email; v != nil {
		query = collection.Where("email", "==", filter.Email)
	}
	if v := filter.Limit; v != 0 {
		query = collection.Limit(filter.Limit)
	}
	if v := filter.Offset; v != 0 {
		query = query.Offset(filter.Offset)
	}

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
func findUsers(ctx context.Context, client *firestore.Client, filter model.UserFilter) (_ []*types.UserDocument, n int, err error) {
	query := applyFilters(client.Collection(_userCollection), filter)

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
		var user model.User

		doc.DataTo(&user)

		users = append(users, types.NewUserDocument(doc.Ref.ID, &user))
	}

	return users, len(users), err
}

// findUserByID fetches a user by their ID
func findUserByID(ctx context.Context, client *firestore.Client, id string) (*types.UserDocument, error) {
	userDocs, _, err := findUsers(ctx, client, model.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(userDocs) == 0 {
		return nil, errors.New("user not found")
	}
	return userDocs[0], nil
}

// findUserByID fetches a user by their email
func findUserByEmail(ctx context.Context, client *firestore.Client, email string) (*types.UserDocument, error) {
	userDoc, _, err := findUsers(ctx, client, model.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(userDoc) == 0 {
		return nil, errors.New("user not found")
	}
	return userDoc[0], nil
}

// createUser creates a new user
func createUser(ctx context.Context, client *firestore.Client, user *model.User) error {
	// Validate
	if err := user.Validate(); err != nil {
		return err
	}

	_, _, err := client.Collection(_userCollection).Add(ctx, &user_service.User{
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
func updateUser(ctx context.Context, client *firestore.Client, id string, upd model.UserUpdate) (*model.User, error) {
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

	_, err = client.Collection(_userCollection).Doc(userDoc.ID).Update(ctx, []firestore.Update{
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
	userDoc, err := findUserByID(ctx, client, id)
	if err != nil {
		return err
	}

	if _, err := client.Collection(_userCollection).Doc(userDoc.ID).Delete(ctx); err != nil {
		return err
	}

	return nil
}
