package main

import (
	"developer-bot/discord"
	"developer-bot/endpoints"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	go endpoints.Serve()

	wg.Add(1)
	go discord.RunBot(&wg)

	// Wait for signal from the os before exiting
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	discord.SendShutdown()

	// Wait for all goroutines to finish before exiting
	wg.Wait()
}
