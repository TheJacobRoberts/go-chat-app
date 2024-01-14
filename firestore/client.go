package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type Client struct {
	client *firestore.Client
}

func NewClient(ctx context.Context, projectID string) (*Client, error) {
	conf := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
