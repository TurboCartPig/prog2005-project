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
			Description: "Unsubscribe from a GitLab repository",
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
		{
			Name:        "endvote",
			Description: "End the ongoing vote, and cast the result",
		},
	}
	// CommandHandlers defines what functions to call when slash commands are used
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help":      commandHandlerHelp,
		"sub":       commandHandlerSub,
		"unsub":     commandHandlerUnsub,
		"deadlines": commandHandlerDeadlines,
		"endvote":   commandHandlerEndvote,
	}
)

// Respond with a help message
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
	// Extract the url from the discord interaction
	url := i.Data.Options[0].StringValue()

	// Build a registration and save it to firestore
	chReg := types.ChannelRegistration{
		ChannelID:  i.ChannelID,
		RepoWebURL: url,
	}
	firestore.SaveChannelRegistration(&chReg)

	// Respond to the command with a confimation
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
	// Extract the url from the discord interaction
	url := i.Data.Options[0].StringValue()

	// Delete the registration from firestore
	err := firestore.DeleteChannelRegistrations(i.ChannelID)
	if err != nil {
		log.Printf("failed while unsubscribing from a channel at %s. %s", url, err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("Failed while unsubscribing from %s, try again later...", url),
			},
		})
		return
	}

	// Respond to the command with a confimation
	log.Printf("unsubscribing from a channel at %s", url)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf("Unsubscribing from %s", url),
		},
	})
	if err != nil {
		log.Println("Failed to send subscription confimation")
	}
}

// Respond with all deadlines for the repos registered for this channel
func commandHandlerDeadlines(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Get the repo associated with this channel
	repoURL, err := firestore.GetRepoURLByChannelID(i.ChannelID)
	if err != nil {
		log.Println("Failed to get repoURL")
		return
	}

	// Get deadlines from firestore
	deadlines := firestore.GetDeadlinesByRepoURL(repoURL)

	// Build all the feilds for an embed
	var fields []*discordgo.MessageEmbedField
	for _, elem := range deadlines {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   elem.Title,
			Value:  elem.Description + "\n\nDue: " + elem.DueDate,
			Inline: false,
		})
	}

	// Respond to the command with all the deadlines
	log.Println("Posting deadlines")
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Deadlines",
				Description: "",
				Color:       10181046,
				Fields:      fields,
			}},
		},
	})
	if err != nil {
		log.Println("Failed to post deadlines in response to deadlines command")
	}
}

func commandHandlerEndvote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Print("command " + i.ChannelID)
	// Send endvote signal to the appropriate vote tracker
	if c, ok := votingChannels[i.ChannelID]; ok {
		c <- 1
	}

	// Acknowledge the command, the processing of the vote is handled elsewere
	log.Println("Ending vote")
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Processing results",
		},
	})
	if err != nil {
		log.Println("Failed to post acknowledgement of endvote")
	}
}
