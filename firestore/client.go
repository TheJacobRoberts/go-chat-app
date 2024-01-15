package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// FirestoreClient represents an internal client to interact with firestore
type FirestoreClient struct {
	client *firestore.Client
}

// NewFirestoreClient returns a new instance of FirestoreClient
func NewFirestoreClient(ctx context.Context, projectID string) (*FirestoreClient, error) {
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

	return &FirestoreClient{client}, nil
}

func (c *FirestoreClient) Close() error {
	return c.client.Close()
}
