package bot

import (
	"fmt"
	"os"

	"TitanBot/internal/handlers"
	"github.com/bwmarrin/discordgo"
)

func NewClient() (*discordgo.Session, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN environment variable is not set")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Register Intents
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// Sharding is enabled
	// In a real production bot, this would be dynamically allocated by a master broker
	// Here we set up local sharding configuration.
	session.ShardCount = 1
	session.ShardID = 0

	// Initialize the Event Router (Worker Pool)
	router := handlers.NewRouter(session, 100)
	router.StartWorkers()

	// Register event handlers
	session.AddHandler(router.OnReady)
	session.AddHandler(router.OnMessageCreate)

	return session, nil
}
