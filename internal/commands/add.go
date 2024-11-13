package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/bykclk/twocli/internal/storage"
	"github.com/bykclk/twocli/internal/totp"
)

type AddCommand struct{}

func NewAddCommand() *AddCommand {
	return &AddCommand{}
}

func (c *AddCommand) Name() string {
	return "add"
}

func (c *AddCommand) Description() string {
	return "Add a new account"
}

func (c *AddCommand) Run(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	name := fs.String("name", "", "Account name")
	secret := fs.String("secret", "", "Account secret key (base32 encoded)")

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

	masterPassword, err := promptForMasterPassword()
	if err != nil {
		return err
	}

	if err = storage.AddAccount(*name, *secret, masterPassword); err != nil {
		if err.Error() == "incorrect master password" {
			fmt.Println("Incorrect master password.")
		}
		return err
	}

	fmt.Println("Account added successfully.")
	return nil
}
