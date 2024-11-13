package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bykclk/twocli/internal/crypto"
)

const dataFile = "data/accounts.db"

// Account represents an account with a name and encrypted secret.
type Account struct {
	Name            string `json:"name"`
	EncryptedSecret []byte `json:"encrypted_secret"`
}

// LoadAccounts loads and decrypts the accounts from the data file.
func LoadAccounts(masterPassword string) ([]Account, error) {
	// Check if data file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return []Account{}, nil // No accounts yet
	}

	// Read the encrypted data file
	encryptedData, err := os.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	// Decrypt the data
	jsonData, err := crypto.DecryptData(encryptedData, masterPassword)
	if err != nil {
		return nil, errors.New("incorrect master password")
	}

	// Unmarshal JSON data
	var accounts []Account
	if err = json.Unmarshal(jsonData, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// saveAccounts encrypts and saves the accounts to the data file.
func saveAccounts(accounts []Account, masterPassword string) error {
	// Marshal accounts to JSON
	jsonData, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	// Encrypt the data
	encryptedData, err := crypto.EncryptData(jsonData, masterPassword)
	if err != nil {
		return err
	}

	// Ensure data directory exists
	if err = os.MkdirAll(filepath.Dir(dataFile), 0700); err != nil {
		return err
	}

	// Write encrypted data to file
	err = os.WriteFile(dataFile, encryptedData, 0600)
	if err != nil {
		return err
	}

	return nil
}

// AddAccount adds a new account to the storage.
func AddAccount(name, secret, masterPassword string) error {
	accounts, err := LoadAccounts(masterPassword)
	if err != nil {
		return err
	}

	// Check for duplicate account name
	for _, acc := range accounts {
		if strings.EqualFold(acc.Name, name) {
			return errors.New("account with this name already exists")
		}
	}

	// Encrypt the secret
	encryptedSecret, err := crypto.EncryptData([]byte(secret), masterPassword)
	if err != nil {
		return err
	}

	// Add the new account
	accounts = append(accounts, Account{
		Name:            name,
		EncryptedSecret: encryptedSecret,
	})

	// Save accounts
	if err = saveAccounts(accounts, masterPassword); err != nil {
		return err
	}

	return nil
}

// GetAccountSecret retrieves and decrypts the secret for a given account name.
func GetAccountSecret(name, masterPassword string) (string, error) {
	accounts, err := LoadAccounts(masterPassword)
	if err != nil {
		return "", err
	}

	// Find the account
	for _, acc := range accounts {
		if strings.EqualFold(acc.Name, name) {
			// Decrypt the secret
			secretData, err := crypto.DecryptData(acc.EncryptedSecret, masterPassword)
			if err != nil {
				return "", err
			}
			return string(secretData), nil
		}
	}

	return "", errors.New("account not found")
}

func DeleteAccount(name, masterPassword string) error {
	accounts, err := LoadAccounts(masterPassword)
	if err != nil {
		return err
	}

	// Find and delete the account
	index := -1
	for i, acc := range accounts {
		if strings.EqualFold(acc.Name, name) {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("account not found")
	}

	// Remove the account from the slice
	accounts = append(accounts[:index], accounts[index+1:]...)

	// Save the updated accounts
	if err = saveAccounts(accounts, masterPassword); err != nil {
		return err
	}

	return nil
}

// UpdateAccount updates the secret of an existing account.
func UpdateAccount(name, newSecret, masterPassword string) error {
	accounts, err := LoadAccounts(masterPassword)
	if err != nil {
		return err
	}

	// Find the account
	found := false
	for i, acc := range accounts {
		if strings.EqualFold(acc.Name, name) {
			// Encrypt the new secret
			encryptedSecret, err := crypto.EncryptData([]byte(newSecret), masterPassword)
			if err != nil {
				return err
			}

			// Update the secret
			accounts[i].EncryptedSecret = encryptedSecret
			found = true
			break
		}
	}

	if !found {
		return errors.New("account not found")
	}

	// Save the updated accounts
	if err = saveAccounts(accounts, masterPassword); err != nil {
		return err
	}

	return nil
}
