package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/bykclk/twocli/internal/storage"
)

type DeleteCommand struct{}

func NewDeleteCommand() *DeleteCommand {
	return &DeleteCommand{}
}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func (c *DeleteCommand) Description() string {
	return "Delete an existing account"
}

func (c *DeleteCommand) Run(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	name := fs.String("name", "", "Account name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		fs.Usage()
		return errors.New("-name is required")
	}

	_, masterPassword, err := loadAccountsWithAttempts()
	if err != nil {
		return err
	}

	// Confirm deletion
	confirmed, err := confirmAction(fmt.Sprintf("Are you sure you want to delete the account '%s'? (yes/no): ", *name))
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	if err = storage.DeleteAccount(*name, masterPassword); err != nil {
		return err
	}

	fmt.Printf("Account '%s' deleted successfully.\n", *name)
	return nil
}
