package handlers

import (
	"fmt"
	"log"
	"strings"

	"developer-bot/endpoints/firestore"
	"developer-bot/endpoints/types"

	"github.com/bwmarrin/discordgo"
)

const HelpMessage = "Help!!!!"

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
			_, err := s.ChannelMessageSend(msg.ChannelID, HelpMessage)
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

			log.Printf("subscribeing from a channel at %s", url)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Subscribing to %s", url))
		} else if strings.HasPrefix(command, "unsub ") {
			url := command[6:]
			// Unsubscribe from repo
			log.Printf("unsubscribeing from a channel at %s", url)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Unsubscribing from %s", url))
		}
	}
}
