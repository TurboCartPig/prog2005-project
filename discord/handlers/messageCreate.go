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

	// Check if the bot was mentioned in the message
	mentioned := false
	for _, user := range msg.Mentions {
		if user.ID == s.State.User.ID {
			mentioned = true
			break
		}
	}

	// We only respond to messages were we are mentioned
	if !mentioned {
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
			url := command[2:]
			chReg := types.ChannelRegistration{
				ChannelID:  msg.ChannelID,
				RepoWebURL: url,
			}
			firestore.SaveChannelRegistration(&chReg)
		} else if strings.HasPrefix(command, "unsub") {
			// Unsubscribe from repo
		}
	}

	// If the bot was mentioned, send a message back
	// if mentioned {
	// 	_, err := s.ChannelMessageSend(msg.ChannelID, "I'm a very friendly and nice bot, I would neeeeveer insults you, uwu")
	// 	if err != nil {
	// 		log.Println("Failed to send message: ", err)
	// 	}
	// }
}
