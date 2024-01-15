package main

import (
	"context"
	"fmt"
	"time"

	"go-chat-app/firestore"
	"go-chat-app/service/user"
)

func main() {
	ctx := context.Background()

	fsClient, err := firestore.NewFirestoreClient(ctx, "go-chat-app-411202")
	if err != nil {
		panic(err)
	}
	defer fsClient.Close()

	userService := firestore.NewUserService(fsClient)

	u, err := userService.FindUserByID(ctx, "c91f5d35-837f-4939-bfd3-46fc1b120bb3")
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] %s\n", u.ID, u.Name)

	users, err := userService.FindUsers(ctx)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Printf("[%s] %s\n", user.ID, user.Name)
	}

	if err := userService.CreateUser(ctx, &user.User{
		Name:      "Test",
		Email:     "testing@example.com",
		CreatedAt: time.Now().UTC().Add(-time.Minute * 60),
		UpdatedAt: time.Now().UTC().Add(-time.Minute * 59),
	}); err != nil {
		panic(err)
	}
}
