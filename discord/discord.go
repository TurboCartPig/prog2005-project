package discord

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"developer-bot/discord/handlers"

	"github.com/bwmarrin/discordgo"
)

// RunBot runs the discord bot until a signal or interrupt from the os signal that it should quit.
func RunBot(wg *sync.WaitGroup) {
	defer wg.Done() // Decrement wg AFTER sessing is closed

	// Setup notification from the os to stop the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Get the bot token and open a discord session with it
	token := getToken()
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

	// Wait until the server should stop
	<-stop
}

// getToken from either TOKEN or TOKEN_FILE environment variables.
// There are two options for passing the secret token to the program:
// 1. TOKEN contains the secret token directly and should not be used in production, but it's fine for development.
// 2. TOKEN_FILE contains the full path to a file that contains the secret.
// TOKEN_FILE is preferred since it can be more secure.
func getToken() (token string) {
	if value, set := os.LookupEnv("TOKEN"); set { // Is the TOKEN envvar set
		log.Println("Found TOKEN")
		token = value
	} else if file, set := os.LookupEnv("TOKEN_FILE"); set { // Is the TOKEN_FILE envvar set
		log.Println("Found TOKEN_FILE")
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal("Failed to read token file")
		}
		token = strings.TrimSpace(string(data))
	} else {
		log.Fatal("No token found. Set TOKEN or TOKEN_FILE")
	}

	return token
}
