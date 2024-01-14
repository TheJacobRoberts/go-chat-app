package firestore

import (
	"context"
	"fmt"
	"go-chat-app/service/user"

	"google.golang.org/api/iterator"
)

var _ user.UserService = (*UserService)(nil)

type UserService struct {
	client *Client
}

// NewUserService returns a new instance of UserService
func NewUserService(client *Client) *UserService {
	return &UserService{
		client: client,
	}
}

func (s *UserService) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	userDocs := s.client.client.Collection("user").Where("id", "==", id)
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
		// fmt.Println(userFeed)
		fmt.Printf("Document data: %#v\n", user)
		fmt.Println(doc.Data())

		users = append(users, &user)
	}
	// [END fs_get_all_docs]
	return users[0], nil
}
