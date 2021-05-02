package main

import (
	"developer-bot/discord"
	"developer-bot/endpoints"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	go endpoints.Serve()

	wg.Add(1)
	go discord.RunBot(&wg)

	// Wait for all gorutines to finish before exiting
	wg.Wait()
}
