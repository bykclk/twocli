package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bykclk/twocli/internal/storage"
	"github.com/bykclk/twocli/internal/totp"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorCyan   = "\033[36m"
)

// Progress bar characters
const (
	leftBracket  = "│"
	rightBracket = "│"
	fullBlock    = "■"
	emptyBlock   = "·"
)

func getProgressColor(remaining int64) string {
	switch {
	case remaining > 15:
		return colorGreen
	case remaining > 5:
		return colorYellow
	default:
		return colorRed
	}
}

func generateProgressBar(remaining, total int64) string {
	width := 20 // Progress bar width
	filled := int(float64(remaining) / float64(total) * float64(width))
	color := getProgressColor(remaining)

	// Build progress bar
	progress := leftBracket
	for i := 0; i < width; i++ {
		if i < filled {
			progress += color + fullBlock + colorReset
		} else {
			progress += emptyBlock
		}
	}
	progress += rightBracket

	return progress
}

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

func generateAndDisplayCode(name, secret string, quit chan struct{}) error {
	totpInfo, err := totp.GenerateCode(secret)
	if err != nil {
		return err
	}

	// Clear the line and move cursor to beginning
	fmt.Print("\033[2K\r")
	fmt.Printf("%sYour TOTP code for '%s' is:%s %s%06d%s",
		colorCyan, name, colorReset,
		colorGreen, totpInfo.Code, colorReset)

	// Display countdown
	remaining := totpInfo.RemainingSeconds
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for remaining > 0 {
		select {
		case <-quit:
			fmt.Println("\nExiting...")
			return nil
		case <-ticker.C:
			progressBar := generateProgressBar(remaining, 30)
			timeColor := getProgressColor(remaining)

			// Clear the line and move cursor to beginning
			fmt.Printf("\033[2K\r")
			fmt.Printf("%sYour TOTP code for '%s' is:%s %s%06d%s %s %s%ds%s",
				colorCyan, name, colorReset,
				colorGreen, totpInfo.Code, colorReset,
				progressBar,
				timeColor, remaining, colorReset)

			remaining--
		}
	}
	return nil
}

func (c *CodeCommand) Run(args []string) error {
	fs := flag.NewFlagSet("code", flag.ContinueOnError)
	name := fs.String("name", "", "Account name")
	autoRefresh := fs.Bool("auto", false, "Automatically generate new codes")

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

	// Setup signal handling for graceful exit
	quit := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		close(quit)
	}()

	fmt.Println("Press Ctrl+C to exit")

	for {
		if err := generateAndDisplayCode(*name, secret, quit); err != nil {
			return err
		}

		select {
		case <-quit:
			return nil
		default:
			if !*autoRefresh {
				return nil
			}
		}
	}
}
