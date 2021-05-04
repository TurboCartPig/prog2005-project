package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"developer-bot/endpoints/firestore"
)

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
		}
	}

	// If the bot was mentioned, send a message back
	if mentioned {
		_, err := s.ChannelMessageSend(msg.ChannelID, "I'm a very friendly and nice bot, I would neeeeveer insults you, uwu")
		if err != nil {
			log.Println("Failed to send message: ", err)
		}
	}
}
