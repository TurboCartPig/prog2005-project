package discord

import (
	"developer-bot/firestore"
	"developer-bot/types"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	// Commands lists all the slash commands to register
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Print help message",
		},
		{
			Name:        "sub",
			Description: "Subscribe to a GitLab repository",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "repo",
					Description: "URL of a GitLab repository",
					Required:    true,
				},
			},
		},
		{
			Name:        "unsub",
			Description: "Unsubscribe to a GitLab repository",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "repo",
					Description: "URL of a GitLab repository",
					Required:    true,
				},
			},
		},
		{
			Name:        "deadlines",
			Description: "Get all deadlines from subscribed GitLab repo",
		},
	}
	// CommandHandlers defines what functions to call when slash commands are used
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help":      commandHandlerHelp,
		"sub":       commandHandlerSub,
		"unsub":     commandHandlerUnsub,
		"deadlines": commandHandlerDeadlines,
	}
)

// Repspond with a help message
func commandHandlerHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title: "Developer-bot help instructions",
				Description: "How to do stuff with developer-bot!\n" +
					"1. Register the bot with GitLab according to readme.\n" +
					"2. Subscribe to a GitLab project from the Discord channel where you want to receive notifications.\n" +
					"3. ???\n" +
					"4. Huge profit!",
				Color: 1752220,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "/help",
						Value:  "Print this help message",
						Inline: false,
					},
					{
						Name: "/sub <gitlab repo url>",
						Value: "Subscribe to a GitLab project and receive notifications when deadlines are posted" +
							". The notifications will appear only in the Discord channel that subscribed to them." +
							"\n\n" +
							"Example: `/sub https://git.gvk.idi.ntnu.no/course/prog2005`\n",
						Inline: false,
					},
					{
						Name: "/unsub <gitlab repo url>",
						Value: "Unsubscribe from a GitLab project." +
							"\n\n" +
							"Example: `/unsub https://git.gvk.idi.ntnu.no/course/prog2005`\n",
						Inline: false,
					},
					{
						Name:   "/deadlines",
						Value:  "Get all the deadlines from the subscribed repo",
						Inline: false,
					},
				},
			}},
		},
	})
}

// Subscribe to gitlab repo notifications
func commandHandlerSub(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url := i.Data.Options[0].StringValue()

	chReg := types.ChannelRegistration{
		ChannelID:  i.ChannelID,
		RepoWebURL: url,
	}
	firestore.SaveChannelRegistration(&chReg)

	log.Printf("subscribing from a channel at %s", url)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("Subscribing to %s", url),
		},
	})
	if err != nil {
		log.Println("Failed to send subscription confimation")
	}
}

// Unsubscribe from repo
func commandHandlerUnsub(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url := i.Data.Options[0].StringValue()

	err := firestore.DeleteChannelRegistations(i.ChannelID)
	if err != nil {
		log.Printf("failed while unsubscribing from a channel at %s. %s", url, err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("Failed while unsubscribing from %s, try again later...", url),
			},
		})
	}

	log.Printf("unsubscribing from a channel at %s", url)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("Subscribing to %s", url),
		},
	})
	if err != nil {
		log.Println("Failed to send subscription confimation")
	}
}

// Respond with all deadlines for the repos registered for this channel
func commandHandlerDeadlines(s *discordgo.Session, i *discordgo.InteractionCreate) {
	repoURL, err := firestore.GetRepoURLByChannelID(i.ChannelID)
	if err != nil {
		log.Println("Failed to get repoURL")
		return
	}

	// Get deadlines from firestore
	deadlines := firestore.GetDeadlinesByRepoURL(repoURL)

	var fields []*discordgo.MessageEmbedField
	for _, elem := range deadlines {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   elem.Title,
			Value:  elem.Description + "\n\nDue: " + elem.DueDate,
			Inline: false,
		})
	}

	log.Println("Posting deadlines")
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Deadlines",
				Description: "Hei",
				Color:       10181046,
				Fields:      fields,
			}},
		},
	})
	if err != nil {
		log.Println("Failed to post deadlines in response to deadlines command")
	}
}
