package handlers

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"developer-bot/endpoints/firestore"
	"developer-bot/endpoints/types"
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
		} else if strings.HasPrefix(command, "sub") {
			// Subscribe to gitlab repo notifications
			url := command[4:]
			chReg := types.ChannelRegistration{
				ChannelID:  msg.ChannelID,
				RepoWebURL: url,
			}
			firestore.SaveChannelRegistration(&chReg)
		} else if strings.HasPrefix(command, "unsub") {
			// Unsubscribe from repo
		}
	}
}
