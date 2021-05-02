package main

import (
	"developer-bot/endpoints"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go endpoints.Serve()

	// Wait for signal from the os before exiting
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
}
