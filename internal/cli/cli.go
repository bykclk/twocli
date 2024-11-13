package cli

import (
	"fmt"
	"os"
)

type Command interface {
	Name() string
	Description() string
	Run(args []string) error
}

func Run(commands []Command) {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		printUsage(commands)
		os.Exit(1)
	}

	cmdName := os.Args[1]

	for _, cmd := range commands {
		if cmd.Name() == cmdName {
			if err := cmd.Run(os.Args[2:]); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	fmt.Printf("Unknown command: %s\n", cmdName)
	printUsage(commands)
	os.Exit(1)
}

func printUsage(commands []Command) {
	fmt.Println("Usage:")
	for _, cmd := range commands {
		fmt.Printf("  %s - %s\n", cmd.Name(), cmd.Description())
	}
}
