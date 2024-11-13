package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	password := "testpassword"
	data := []byte("secret data")

	encryptedData, err := EncryptData(data, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedData, err := DecryptData(encryptedData, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(data, decryptedData) {
		t.Fatalf("Decrypted data does not match original")
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	password := "testpassword"
	wrongPassword := "wrongpassword"
	data := []byte("secret data")

	encryptedData, err := EncryptData(data, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	_, err = DecryptData(encryptedData, wrongPassword)
	if err == nil {
		t.Fatalf("Decryption should have failed with wrong password")
	}
}

func TestEncryptDecryptConsistency(t *testing.T) {
	password := "testpassword"
	data := []byte(`{"name":"TestAccount","encrypted_secret":"somedata"}`)

	encryptedData, err := EncryptData(data, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decryptedData, err := DecryptData(encryptedData, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(data, decryptedData) {
		t.Fatalf("Decrypted data does not match original")
	}
}
