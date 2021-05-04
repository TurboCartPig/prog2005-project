package main

import (
	"developer-bot/discord"
	"developer-bot/endpoints"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	messages := make(chan discord.Message)

	go endpoints.Serve()

	wg.Add(1)
	go discord.RunBot(messages, &wg)

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
