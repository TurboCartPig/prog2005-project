package endpoints

import (
	"developer-bot/discord"
	"developer-bot/endpoints/firestore"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"os"

	"developer-bot/endpoints/types"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Version of the API
const Version = "v1"

// getPort from environment variable
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	return port
}

// Serve the REST API over HTTP
func Serve() {
	port := getPort()
	router := routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(handler2 http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging error: %s\n", err.Error())
	}

	// Create a new firestore client in the firestore package. 
	firestore.NewFirestoreClient()

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// routes sets up routes
func routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Post("/endpoints", developer)

	return router
}

func developer(w http.ResponseWriter, r *http.Request) {
	var newWebhook types.WebhookData
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, "Working")

	processWebhook(&newWebhook)
	w.WriteHeader(http.StatusOK)
}

func processWebhook(webhook *types.WebhookData) {
	if isDeadline(webhook) {
		deadline := types.Deadline{
			RepoWebURL:  webhook.Project.WebURL,
			Title:       webhook.ObjectAttributes.Title,
			Description: webhook.ObjectAttributes.Description,
			DueDate:     webhook.ObjectAttributes.DueDate,
			IssueWebURL: webhook.ObjectAttributes.Url,
		}
		firestore.SaveDeadlineToFirestore(&deadline)
		sendMessageToDiscord(&deadline)
	}
}

func sendMessageToDiscord(deadline *types.Deadline) {
	discordMessage := discordgo.MessageSend{
		Content: "New deadline posted:",
		Embed: &discordgo.MessageEmbed{
			URL:         deadline.IssueWebURL,
			Title:       deadline.Title,
			Description: deadline.Description,
			Color:       15158332,
			Fields: []*discordgo.MessageEmbedField {
				{
					Name:   "DUE DATE",
					Value:  deadline.DueDate,
					Inline: false,
				},
			},
		},
	}
	channelID := firestore.GetChannelIDByRepoURL(deadline.RepoWebURL)
	discord.SendComplexMessage(channelID, &discordMessage)
}

func isDeadline(webhook *types.WebhookData) bool {
	for _, label := range webhook.Labels {
		if label.Title == "deadline" {
			return true
		}
	}
	return false
}
