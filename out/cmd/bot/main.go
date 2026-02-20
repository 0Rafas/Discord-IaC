package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"TitanBot/internal/bot"
	"TitanBot/internal/cache"
	"TitanBot/internal/database"
	"TitanBot/internal/metrics"
)

func main() {
	log.Println("Starting TitanBot...")

	// Initialize discord client
	client, err := bot.NewClient()
	if err != nil {
		log.Fatalf("Error creating discord session: %v", err)
	}

	// Initialize Database
	err = database.Init()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Initialize Redis Cache
	err = cache.Init()
	if err != nil {
		log.Fatalf("Cache initialization failed: %v", err)
	}

	// Start Prometheus Metrics Server
	go metrics.StartServer(":2112")

	// Open connection to discord
	err = client.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down cleanly...")
	client.Close()
}
