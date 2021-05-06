package discord

import (
	"developer-bot/types"
	"log"
	"sync"

	"developer-bot/discord/handlers"
	"developer-bot/firestore"

	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session

var messages = make(chan types.Message)

func SendMessage(channelID, content string) {
	messages <- types.MessageSend{ChannelID: channelID, Content: content}
}

func SendComplexMessage(channelID string, message *discordgo.MessageSend) {
	messages <- types.MessageSendComplex{
		ChannelID: channelID,
		Message:   message,
	}
}

func SendComplexMessageWithFollowUp(
	channelID string,
	message *discordgo.MessageSend,
	object interface{},
	followUp func(string, string, interface{}),
) {
	messages <- types.MessageSendComplexWithFollowUp{
		ChannelID: channelID,
		Message:   message,
		FollowUp:  followUp,
		Object:    object,
	}
}

func SendShutdown() {
	messages <- types.Shutdown{}
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

	session, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Failed to create a discord session\n", err)
		return
	}

	// Add handler to notify when the bot is ready
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up and running") })
	// Add handler for new messages being posted
	session.AddHandler(handlers.MessageCreate)
	// Add handler for individual slash commands
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := CommandHandlers[i.Data.Name]; ok {
			handler(s, i)
		}
	})

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
		case types.MessageSend:
			_, err := session.ChannelMessageSend(t.ChannelID, t.Content)
			if err != nil {
				log.Println("Failed to send message: ", err)
			}
		case types.MessageSendComplex:
			_, err := session.ChannelMessageSendComplex(t.ChannelID, t.Message)
			if err != nil {
				log.Println("Failed to send message: ", err)
			}
		case types.MessageSendComplexWithFollowUp:
			message, err := session.ChannelMessageSendComplex(t.ChannelID, t.Message)
			if err != nil {
				log.Println("Failed to send message: ", err)
			}
			t.FollowUp(message.ID, t.ChannelID, t.Object)
		case types.Shutdown:
			return
		}
	}
}

// Register slash commands.
// NOTE: Apparently we only need to do this every time we change the slash commands,
//       not everytime we start the bot
// nolint:deadcode,unused // Will be used in future code
func registerSlashCommands(session *discordgo.Session) {
	for _, command := range Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command)
		if err != nil {
			log.Printf("Failed to create command: %v, error %v", command.Name, err)
			return
		}
	}
}

func HandleVote(messageID, channelID string, object interface{}) {
	if t, ok := object.(*types.Vote); ok {
		for _, elem := range t.Options {
			err := session.MessageReactionAdd(channelID, messageID, elem.EmojiCode)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
