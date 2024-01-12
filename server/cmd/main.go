package main

import (
	"go-chat-app/server"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	s := server.NewServer(nil, 8080, logger)

	if err := s.Start(); err != nil {
		panic(err)
	}

}
