package commands

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bykclk/twocli/internal/storage"

	"github.com/bykclk/twocli/internal/totp"
)

type UpdateCommand struct{}

func NewUpdateCommand() *UpdateCommand {
	return &UpdateCommand{}
}

func (c *UpdateCommand) Name() string {
	return "update"
}

func (c *UpdateCommand) Description() string {
	return "Update the secret key of an existing account"
}

func (c *UpdateCommand) Run(args []string) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	name := fs.String("name", "", "Account name")
	secret := fs.String("secret", "", "New account secret key (base32 encoded)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *secret == "" {
		fs.Usage()
		return errors.New("both -name and -secret are required")
	}

	if err := totp.ValidateSecret(*secret); err != nil {
		return fmt.Errorf("invalid secret key: %v", err)
	}

	_, masterPassword, err := loadAccountsWithAttempts()
	if err != nil {
		return err
	}

	if err = storage.UpdateAccount(*name, *secret, masterPassword); err != nil {
		return err
	}

	fmt.Printf("Account '%s' updated successfully.\n", *name)
	return nil
}
