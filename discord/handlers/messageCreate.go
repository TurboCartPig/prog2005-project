package handlers

import (
	"fmt"
	"log"
	"strings"

	"developer-bot/endpoints/firestore"
	"developer-bot/endpoints/types"

	"github.com/bwmarrin/discordgo"
)

// MessageCreate handles messages being sent in any channel the bot has access to.
func MessageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if msg.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(msg.Content, "!") {
		command := msg.Content[1:]

		if strings.HasPrefix(command, "help") {
			// Print help message
			helpMessage := discordgo.MessageSend{
				Embed: &discordgo.MessageEmbed{
					Title: "Developer-bot help instructions",
					Description: "How to do stuff with developer-bot!\n" +
						"1. Register the bot with GitLab according to readme.\n" +
						"2. Subscribe to a GitLab project from the Discord channel where you want to receive notifications.\n" +
						"3. ???\n" +
						"4. Huge profit!",
					Color: 1752220,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "!help",
							Value:  "Print this help message",
							Inline: false,
						},
						{
							Name: "!sub <gitlab repo url>",
							Value: "Subscribe to a GitLab project and receive notifications when deadlines are posted" +
								". The notifications will appear only in the Discord channel that subscribed to them." +
								"\n\n" +
								"Example: `!sub https://git.gvk.idi.ntnu.no/course/prog2005`\n",
							Inline: false,
						},
						{
							Name: "!unsub <gitlab repo url>",
							Value: "Unsubscribe from a GitLab project." +
								"\n\n" +
								"Example: `!unsub https://git.gvk.idi.ntnu.no/course/prog2005`\n",
							Inline: false,
						},
					},
				},
			}
			_, err := s.ChannelMessageSendComplex(msg.ChannelID, &helpMessage)
			if err != nil {
				log.Println("Failed to send message: ", err)
			}
		} else if strings.HasPrefix(command, "sub ") {
			// Subscribe to gitlab repo notifications
			url := command[4:]
			chReg := types.ChannelRegistration{
				ChannelID:  msg.ChannelID,
				RepoWebURL: url,
			}
			firestore.SaveChannelRegistration(&chReg)

			log.Printf("subscribing from a channel at %s", url)
			_, _ = s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Subscribing to %s", url))
		} else if strings.HasPrefix(command, "unsub ") {
			url := command[6:]
			// Unsubscribe from repo
			log.Printf("unsubscribing from a channel at %s", url)
			_, _ = s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Unsubscribing from %s", url))
		}
	}
}
