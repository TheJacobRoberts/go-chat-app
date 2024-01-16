package main

import (
	"context"
	"fmt"

	"go-chat-app/firestore"
	"go-chat-app/model"
)

func main() {
	ctx := context.Background()

	fsClient, err := firestore.NewFirestoreClient(ctx, "go-chat-app-411202")
	if err != nil {
		panic(err)
	}
	defer fsClient.Close()

	userService := firestore.NewUserService(fsClient)

	user, err := userService.UpdateUser(ctx, "cf7f6087-4d59-43aa-914d-7ef6857c1e87", model.UserUpdate{
		Name: &string("Jane Joy Doe"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(user.String())
}
