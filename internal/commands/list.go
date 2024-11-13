package commands

import (
	"fmt"
)

type ListCommand struct{}

func NewListCommand() *ListCommand {
	return &ListCommand{}
}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Description() string {
	return "List all saved accounts"
}

func (c *ListCommand) Run(_ []string) error {
	accounts, _, err := loadAccountsWithAttempts()
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		fmt.Println("No accounts found.")
		return nil
	}

	fmt.Println("Saved accounts:")
	for _, acc := range accounts {
		fmt.Printf("- %s\n", acc.Name)
	}

	return nil
}
