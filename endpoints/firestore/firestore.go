package firestore

import (
	"context"
	"developer-bot/endpoints/types"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	firebase "firebase.google.com/go/v4"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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

// ShutdownClient the firebase client.
func ShutdownClient() {
	client.Close()
}

// Saves webhook data from gitlab to firestore.
func SaveWebhookToFirestore(webhook *types.WebhookData) {
	ctx := context.Background()
	_, _, err := client.Collection("deadlines").Add(ctx, *webhook)
	if err != nil {
		fmt.Println(err)
	}
}

func SaveChannelRegistration(channelRegistration *types.ChannelRegistration) {
	ctx := context.Background()
	_, _, err := client.Collection("channel-registrations").Add(ctx, *channelRegistration)
	if err != nil {
		fmt.Println(err)
	}
}

// GetBotToken gets the discord bot token from google cloud's secret manager.
func GetBotToken() (string, error) {
	// The name/path to the secret stored in google cloud's secret manager.
	name := "projects/515783087290/secrets/discord-bot-token/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Printf("failed to create secretmanager client: %v", err)
		return "", err
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Printf("failed to access secret version: %v", err)
		return "", err
	}

	return string(result.Payload.Data), nil
}
