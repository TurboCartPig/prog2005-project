package firestore

import (
	"context"
	"developer-bot/endpoints/types"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

var client *firestore.Client

// NewFirestoreClient creates and initializes a new firestore client.
func NewFirestoreClient() {
	// Use GOOGLE_APPLICATION_CREDENTIALS env var to find the service account key
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// Saves webhook data from gitlab to firestore.
func SaveWebhookToFirestore (webhook *types.WebhookData) {
	ctx := context.Background()
	_, _, err := client.Collection("deadlines").Add(ctx,*webhook)
	if err != nil {
		fmt.Println(err)
	}
}