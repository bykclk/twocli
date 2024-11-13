package storage

import (
	"os"
	"testing"
)

func cleanup() {
	if err := os.Remove(dataFile); err != nil && !os.IsNotExist(err) {
		println("Warning: Failed to clean up data file:", err)
	}
}

func TestAddDeleteAccount(t *testing.T) {
	defer cleanup()

	masterPassword := "testpassword"
	accountName := "TestAccount"
	secret := "JBSWY3DPEHPK3PXP"

	// Add account
	err := AddAccount(accountName, secret, masterPassword)
	if err != nil {
		t.Fatalf("Failed to add account: %v", err)
	}

	// Delete account
	err = DeleteAccount(accountName, masterPassword)
	if err != nil {
		t.Fatalf("Failed to delete account: %v", err)
	}

	// Try to get the deleted account
	_, err = GetAccountSecret(accountName, masterPassword)
	if err == nil {
		t.Fatalf("Expected error when getting deleted account")
	}
}

func TestUpdateAccount(t *testing.T) {
	defer cleanup()

	masterPassword := "testpassword"
	accountName := "TestAccount"
	secret := "OLDSECRET"
	newSecret := "NEWSECRET"

	// Add account
	err := AddAccount(accountName, secret, masterPassword)
	if err != nil {
		t.Fatalf("Failed to add account: %v", err)
	}

	// Update account
	err = UpdateAccount(accountName, newSecret, masterPassword)
	if err != nil {
		t.Fatalf("Failed to update account: %v", err)
	}

	// Get the updated secret
	retrievedSecret, err := GetAccountSecret(accountName, masterPassword)
	if err != nil {
		t.Fatalf("Failed to get account secret: %v", err)
	}

	if retrievedSecret != newSecret {
		t.Fatalf("Expected secret '%s', got '%s'", newSecret, retrievedSecret)
	}
}
