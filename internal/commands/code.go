package commands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/bykclk/twocli/internal/storage"
	"github.com/bykclk/twocli/internal/totp"
)

type CodeCommand struct{}

func NewCodeCommand() *CodeCommand {
	return &CodeCommand{}
}

func (c *CodeCommand) Name() string {
	return "code"
}

func (c *CodeCommand) Description() string {
	return "Generate TOTP code for an account"
}

func (c *CodeCommand) Run(args []string) error {
	fs := flag.NewFlagSet("code", flag.ContinueOnError)
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

	secret, err := storage.GetAccountSecret(*name, masterPassword)
	if err != nil {
		return err
	}

	if err = totp.ValidateSecret(secret); err != nil {
		return fmt.Errorf("invalid secret key for account '%s': %v", *name, err)
	}

	code, err := totp.GenerateCode(secret)
	if err != nil {
		return err
	}

	fmt.Printf("Your TOTP code for '%s' is: %06d\n", *name, code)
	return nil
}
