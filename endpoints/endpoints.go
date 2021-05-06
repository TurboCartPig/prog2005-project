package endpoints

import (
	"developer-bot/discord"
	"developer-bot/endpoints/firestore"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

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

	router.Post("/developer", developer)

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
	if isLabel(webhook, "vote") {
		options := strings.Split(webhook.ObjectAttributes.Description, "+==")
		var opt []types.Option
		for i, elem := range options {
			content := strings.Split(elem, "+--")
			opt = append(opt, types.Option{
				Title:       content[0],
				Description: content[1],
				EmojiCode: types.VotingEmojis[i],
			})
		}
		vote := types.Vote{
			RepoWebURL:  webhook.Project.WebURL,
			Title:       webhook.ObjectAttributes.Title,
			Options:     opt,
			IssueWebURL: webhook.ObjectAttributes.URL,
		}
		sendVoteToDiscord(&vote)
	}
}

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

func sendVoteToDiscord(vote *types.Vote) {
	var fields []*discordgo.MessageEmbedField
	for _, elem := range vote.Options {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   elem.Title,
			Value:  elem.Description,
			Inline: false,
		})
	}
	discordMessage := discordgo.MessageSend{
		Content: "New vote",
		Embed: &discordgo.MessageEmbed{
			Color: 10181046,
			Title:  vote.Title,
			Fields: fields,
		},
	}
	channelID := firestore.GetChannelIDByRepoURL(vote.RepoWebURL)
	for _, elem := range channelID {
		discord.SendComplexMessageWithFollowUp(elem,&discordMessage,vote,handleVote)
	}
}

func handleVote(messageID , channelID string, object interface{}) {
	session := discord.GetDiscordSession()

	if t, ok := object.(*types.Vote); ok {
		for _, elem := range t.Options {
			err := session.MessageReactionAdd(channelID, messageID, elem.EmojiCode)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func isLabel(webhook *types.WebhookData, labelIdentifier string) bool {
	for _, label := range webhook.Labels {
		if label.Title == labelIdentifier {
			return true
		}
	}
	return false
}
