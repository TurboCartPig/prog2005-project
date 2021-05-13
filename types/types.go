package types

import "github.com/bwmarrin/discordgo"

// WebhookData is a reduced version of the payload GitLab posts as part of it's webhook.
type WebhookData struct {
	EventType string `json:"event_type"`
	User      User   `json:"user"`
	Project   struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		WebURL      string `json:"web_url"`
	} `json:"project"`
	ObjectAttributes struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		AuthorID    int      `json:"author_id"`
		DueDate     string   `json:"due_date"`
		ProjectID   int      `json:"project_id"`
		Labels      []Labels `json:"labels"`
		URL         string   `json:"url"`
	} `json:"object_attributes"`
	Labels     []Labels `json:"labels"`
	Repository struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
}

// User is a user, part of WebhookData.
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Labels is a label, part of WebhookData.
type Labels struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	ProjectID   int         `json:"project_id"`
	Description interface{} `json:"description"`
	Type        string      `json:"type"`
	GroupID     interface{} `json:"group_id"`
}

// ChannelRegistration as stored in firebase.
type ChannelRegistration struct {
	ChannelID  string `json:"channel_id"`
	RepoWebURL string `json:"repo_web_url"`
}

// Deadline represents a deadline from GitLab as stored in firebase.
type Deadline struct {
	RepoWebURL  string `json:"repo_web_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IssueWebURL string `json:"issue_web_url"`
	DueDate     string `json:"due_date"`
}

// Vote represnets a vote in discord.
type Vote struct {
	RepoWebURL  string   `json:"repo_web_url"`
	Title       string   `json:"title"`
	Options     []Option `json:"options"`
	IssueWebURL string   `json:"issue_web_url"`
}

// Option is an option as part of a vote.
type Option struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EmojiCode   string `json:"emoji"`
}

// Message is an empty interface for all message signals.
type Message interface{}

// MessageSendComplex sends a complex message in discord.
type MessageSendComplex struct {
	ChannelID string
	Message   *discordgo.MessageSend
}

// MessageSendComplexWithFollowUp sends a message in discord,
// and follows it up later with a followup function.
type MessageSendComplexWithFollowUp struct {
	ChannelID string
	Message   *discordgo.MessageSend
	FollowUp  func(messageID, channelID string, object interface{})
	Object    interface{}
}

// Shutdown signal for discord.
type Shutdown struct{}

// VotingEmojis is an emoji that represents a number.
var VotingEmojis = []string{
	"1️⃣",
	"2️⃣",
	"3️⃣",
	"4️⃣",
	"5️⃣",
	"6️⃣",
	"7️⃣",
	"8️⃣",
	"9️⃣",
}
