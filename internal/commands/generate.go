package commands

import (
	"fmt"
	"os"

	"discord-iac/internal/config"
	"discord-iac/internal/generator"
	"gopkg.in/yaml.v3"
)

// GenerateCommand parses bot.yaml and generates the bot codebase
func GenerateCommand() error {
	filename := "bot.yaml"

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w\nHave you run `discord-iac init` yet?", filename, err)
	}

	var cfg config.BotConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	fmt.Printf("Parsed bot.yaml for '%s'. Starting generation...\n", cfg.Name)

	err = generator.GenerateBot(&cfg)
	if err != nil {
		return fmt.Errorf("failed to generate bot architecture: %w", err)
	}

	fmt.Println("Successfully generated Discord bot architecture!")
	fmt.Println("Next steps:")
	fmt.Println("  1. cd out")
	fmt.Println("  2. go mod tidy")
	fmt.Println("  3. go build ./cmd/bot")
	return nil
}
