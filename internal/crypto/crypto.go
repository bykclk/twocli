package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateKey derives a key from the password using PBKDF2.
func GenerateKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

// EncryptData encrypts data using AES-256-GCM with the given password.
func EncryptData(data []byte, password string) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	// Derive key from password
	key := GenerateKey(password, salt)

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt data
	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	// Return salt + nonce + ciphertext
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptData decrypts data using AES-256-GCM with the given password.
func DecryptData(data []byte, password string) ([]byte, error) {
	if len(data) < 28 {
		return nil, errors.New("invalid data")
	}

	// Extract salt, nonce, and ciphertext
	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	// Derive key from password
	key := GenerateKey(password, salt)

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("incorrect password or corrupted data")
	}

	return plaintext, nil
}
