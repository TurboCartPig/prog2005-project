package endpoints

import (
	"developer-bot/discord"
	"developer-bot/firestore"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"developer-bot/types"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Version of the API
const Version = "v1"

// getPort from environment variable.
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	return port
}

// Serve the REST API over HTTP.
func Serve() {
	port := getPort()
	router := routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(handler2 http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	// If the walker failed. 
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging error: %s\n", err.Error())
	}

	// Create a new firestore client in the firestore package
	firestore.NewFirestoreClient()

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// routes sets up routes.
func routes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Post("/developer", developer)

	return router
}

// developer handles POST requests from GitLab.
func developer(w http.ResponseWriter, r *http.Request) {
	var newWebhook types.WebhookData
	err := json.NewDecoder(r.Body).Decode(&newWebhook)

	// If the body is not a issue body, simply ignore it.
	if err != nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	processWebhook(&newWebhook)

	// Respond that we have accepted the request and we are processing it further
	w.WriteHeader(http.StatusAccepted)
}

// Figures out whether or not a webhook is a 'deadline' or a 'vote', and makes appropriate actions thereafter
func processWebhook(webhook *types.WebhookData) {
	// If the webhook is a deadline
	if isLabel(webhook, "deadline") {
		deadline := types.Deadline{
			RepoWebURL:  webhook.Project.WebURL,
			Title:       webhook.ObjectAttributes.Title,
			Description: webhook.ObjectAttributes.Description,
			DueDate:     webhook.ObjectAttributes.DueDate,
			IssueWebURL: webhook.ObjectAttributes.URL,
		}
		firestore.SaveDeadlineToFirestore(&deadline)
		sendDeadlineToDiscord(&deadline)
	}
	// if the webhook is a vote
	if isLabel(webhook, "vote") {
		// Divide into different options
		options := strings.Split(webhook.ObjectAttributes.Description, "+==")
		var opt []types.Option
		for i, elem := range options {
			content := strings.Split(elem, "+--") // Split between title and description
			opt = append(opt, types.Option{
				Title:       content[0],
				Description: content[1],
				EmojiCode:   types.VotingEmojis[i],
			})
		}
		vote := types.Vote{
			RepoWebURL:  webhook.Project.WebURL,
			Title:       webhook.ObjectAttributes.Title,
			Options:     opt,
			IssueWebURL: webhook.ObjectAttributes.URL,
		}
		discord.SendVoteToDiscord(vote)
	}
}

// Sends a 'deadline' to Discord
func sendDeadlineToDiscord(deadline *types.Deadline) {
	discordMessage := discordgo.MessageSend{
		Content: "New deadline posted:",
		Embed: &discordgo.MessageEmbed{
			URL:         deadline.IssueWebURL,
			Title:       deadline.Title,
			Description: deadline.Description,
			Color:       15158332,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "DUE DATE",
					Value:  deadline.DueDate,
					Inline: false,
				},
			},
		},
	}
	channelID := firestore.GetChannelIDByRepoURL(deadline.RepoWebURL)
	for _, elem := range channelID {
		discord.SendComplexMessage(elem, &discordMessage)
	}
}

// Checks if a webhook is of a specific label type or not
func isLabel(webhook *types.WebhookData, labelIdentifier string) bool {
	for _, label := range webhook.Labels {
		// If the correct label has been found in the webhook
		if label.Title == labelIdentifier {
			return true
		}
	}
	return false
}
