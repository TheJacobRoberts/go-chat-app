package main

import (
	"context"
	"fmt"
	"go-chat-app/firestore"
)

func main() {
	ctx := context.Background()

	fsClient, err := firestore.NewClient(ctx, "project-id")
	if err != nil {
		panic(err)
	}
	defer fsClient.Close()

	userService := firestore.NewUserService(fsClient)

	user, err := userService.FindUserByID(ctx, "user-id")
	if err != nil {
		panic(err)
	}

	fmt.Println(user.Name)
}
