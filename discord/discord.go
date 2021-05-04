package discord

import (
	"log"
	"sync"

	"developer-bot/discord/handlers"
	"developer-bot/endpoints/firestore"

	"github.com/bwmarrin/discordgo"
)

type Message interface{}

type MessageSend struct {
	ChannelID string
	Content   string
}

type Shutdown struct{}

var messages = make(chan Message)

func SendMessage(channelid, content string) {
	messages <- MessageSend{ChannelID: channelid, Content: content}
}

func SendShutdown() {
	messages <- Shutdown{}
}

// RunBot runs the discord bot until a signal or interrupt from the os signal that it should quit.
func RunBot(wg *sync.WaitGroup) {
	defer log.Println("Discord bot shut down")
	defer wg.Done() // Decrement wg AFTER sessing is closed

	// Get the bot token and open a discord session with it
	token, err := firestore.GetBotToken()
	if err != nil {
		return
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Failed to create a discord session\n", err)
		return
	}

	// Add event handler
	session.AddHandler(handlers.MessageCreate)

	// Specify minimal gateway intents
	// See https://discord.com/developers/docs/topics/gateway#gateway-intents
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	// Open the websocket that facilitates communication with discord
	err = session.Open()
	if err != nil {
		log.Println("Failed to open a websocket for communicating with discord\n", err)
		return
	}
	defer session.Close()

	for {
		input := <-messages
		switch t := input.(type) {
		case MessageSend:
			_, err := session.ChannelMessageSend(t.ChannelID, t.Content)
			if err != nil {
				log.Println("Failed to send message: ", err)
			}
		case Shutdown:
			return
		}
	}
}
