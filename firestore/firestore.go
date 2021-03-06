package firestore

import (
	"context"
	"developer-bot/types"
	"log"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	firebase "firebase.google.com/go/v4"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

const (
	ChannelRegistrationsCollection = "channel-registrations"
	DeadlinesCollection            = "deadlines"
)

// client is the local firestore client.
var client *firestore.Client

// NewFirestoreClient creates and initializes a new firestore client.
func NewFirestoreClient() {
	// Use GOOGLE_APPLICATION_CREDENTIALS env var to find the service account key
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	// If a new firebase app could not be created
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	// If the firestore client could not be created.
	if err != nil {
		log.Fatalln(err)
	}
}

// ShutdownClient the firebase client.
func ShutdownClient() {
	client.Close()
}

// SaveDeadlineToFirestore registers a deadline in firestore.
func SaveDeadlineToFirestore(deadline *types.Deadline) {
	ctx := context.Background()
	_, _, err := client.Collection(DeadlinesCollection).Add(ctx, *deadline)
	// If the deadline could not be added to firestore.
	if err != nil {
		log.Println(err)
	}
}

// SaveChannelRegistration registers a channelID and repo url pair in firebase.
func SaveChannelRegistration(channelRegistration *types.ChannelRegistration) {
	ctx := context.Background()
	_, err := client.Collection(ChannelRegistrationsCollection).Doc(channelRegistration.ChannelID).Set(ctx, *channelRegistration)
	// If the channelRegistration could not be added to firestore.
	if err != nil {
		log.Println(err)
	}
}

// GetChannelIDByRepoURL gets the channelIDs associated with a given repo url.
func GetChannelIDByRepoURL(repoURL string) []string {
	ctx := context.Background()
	iter := client.Collection(ChannelRegistrationsCollection).Where("RepoWebURL", "==", repoURL).Documents(ctx)

	var channelIDs []string
	var cr types.ChannelRegistration
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		err = doc.DataTo(&cr)
		if err != nil {
			break
		}
		channelIDs = append(channelIDs, cr.ChannelID)
	}
	return channelIDs
}

// GetRepoURLByChannelID gets the repo url for a given channelID.
func GetRepoURLByChannelID(channelID string) (string, error) {
	ctx := context.Background()
	docref := client.Collection(ChannelRegistrationsCollection).Doc(channelID)
	docsnap, err := docref.Get(ctx)
	// If the document could not get fetched 
	if err != nil {
		log.Println("Failed to get document from firebase:", err)
		return "", err
	}

	var cr types.ChannelRegistration
	err = docsnap.DataTo(&cr)
	// If the document could not be converted to a ChannelRegistration struct
	if err != nil {
		log.Println("Failed to parse the document into ChannelRegistration")
		return "", err
	}

	return cr.RepoWebURL, nil
}

// GetDeadlinesByRepoURL gets all the deadlines for a given repo.
func GetDeadlinesByRepoURL(repoURL string) []types.Deadline {
	ctx := context.Background()
	iter := client.Collection(DeadlinesCollection).Where("RepoWebURL", "==", repoURL).Documents(ctx)

	var deadlines []types.Deadline
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			log.Println("Failed to get document from firestore: ", err)
			continue
		}

		var deadline types.Deadline
		err = doc.DataTo(&deadline)
		// If the document could not be converted to a Deadline struct
		if err != nil {
			log.Println("Failed to parse document into Deadline", err)
			continue
		}

		// If the duedate of the deadline has expired, delete it from firestore and don't add them to the slice.
		if deadline.DueDate < time.Now().String()[:10] {
			_, err = client.Collection(DeadlinesCollection).Doc(doc.Ref.ID).Delete(ctx)
			if err != nil {
				log.Printf("Could not delete document with ID %s", doc.Ref.ID)
			}
		} else {
			deadlines = append(deadlines, deadline)
		}
	}

	return deadlines
}

// DeleteChannelRegistrations deletes all registered repos for a given channelID.
func DeleteChannelRegistrations(channelID string) error {
	ctx := context.Background()
	_, err := client.Collection(ChannelRegistrationsCollection).Doc(channelID).Delete(ctx)
	return err
}

// GetBotToken gets the discord bot token from Google cloud's secret manager.
func GetBotToken() (string, error) {
	// The name/path to the secret stored in Google cloud's secret manager.
	name := "projects/515783087290/secrets/discord-bot-token/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	// If the secretmanager client could not be created
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
