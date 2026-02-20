package main

import (
	"fmt"
	"os"

	"discord-iac/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := commands.InitCommand()
		if err != nil {
			fmt.Printf("Error running init: %v\n", err)
			os.Exit(1)
		}
	case "generate":
		err := commands.GenerateCommand()
		if err != nil {
			fmt.Printf("Error running generate: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Discord Instructure-as-Code CLI")
	fmt.Println("Usage:")
	fmt.Println("  discord-iac init      Creates a sample bot.yaml configuration file.")
	fmt.Println("  discord-iac generate  Generates the Discord bot architecture based on bot.yaml.")
}
