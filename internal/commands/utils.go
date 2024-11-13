package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bykclk/twocli/internal/storage"
)

const maxPasswordAttempts = 3

func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := readPassword()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}

func readPassword() (string, error) {
	// Disable input echoing
	cmd := exec.Command("stty", "-echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var password string
	// Handle the error from Scanln
	if _, err := fmt.Scanln(&password); err != nil {
		// Re-enable input echoing before returning the error
		cmd = exec.Command("stty", "echo")
		cmd.Stdin = os.Stdin
		_ = cmd.Run() // Ignore any error here as it's non-critical
		return "", err
	}

	// Re-enable input echoing
	cmd = exec.Command("stty", "echo")
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return "", err
	}
	fmt.Println()

	return password, nil
}

func confirmAction(prompt string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	if response == "yes" || response == "y" {
		return true, nil
	}
	return false, nil
}

func promptForMasterPassword() (string, error) {
	return promptPassword("Enter master password: ")
}

func loadAccountsWithAttempts() ([]storage.Account, string, error) {
	var masterPassword string
	var accounts []storage.Account
	var err error

	for attempts := 0; attempts < maxPasswordAttempts; attempts++ {
		masterPassword, err = promptForMasterPassword()
		if err != nil {
			return nil, "", err
		}

		accounts, err = storage.LoadAccounts(masterPassword)
		if err == nil {
			return accounts, masterPassword, nil
		}

		if err.Error() == "incorrect master password" {
			fmt.Println("Incorrect master password. Please try again.")
			continue
		} else {
			return nil, "", err
		}
	}

	return nil, "", errors.New("maximum password attempts exceeded")
}
