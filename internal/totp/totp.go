package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"time"
)

// timeNow is a variable to allow overriding in tests.
var timeNow = time.Now

func ValidateSecret(secret string) error {
	// Remove any whitespace and convert to uppercase
	secret = strings.ToUpper(strings.ReplaceAll(secret, " ", ""))

	// Check for empty secret
	if secret == "" {
		return errors.New("secret cannot be empty")
	}

	// Try to decode the secret
	if _, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret); err != nil {
		return fmt.Errorf("invalid base32 encoding: %v", err)
	}

	return nil
}

// TOTPInfo contains the generated code and its validity information
type TOTPInfo struct {
	Code             uint32
	RemainingSeconds int64
}

// GenerateCode generates a TOTP code and returns it along with remaining validity time
func GenerateCode(secret string) (TOTPInfo, error) {
	// Validate secret first
	if err := ValidateSecret(secret); err != nil {
		return TOTPInfo{}, fmt.Errorf("invalid secret key: %v", err)
	}

	// Decode the base32 encoded secret key
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return TOTPInfo{}, fmt.Errorf("failed to decode secret key: %v", err)
	}

	// Calculate the time step and remaining seconds
	epochSeconds := timeNow().Unix()
	timeStep := uint64(epochSeconds / 30)
	remainingSeconds := 30 - (epochSeconds % 30)

	// Convert time step to byte array
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, timeStep)

	// Calculate HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(msg)
	hash := h.Sum(nil)

	// Get offset
	offset := hash[len(hash)-1] & 0xf

	// Generate 4-byte code
	binary := binary.BigEndian.Uint32(hash[offset:]) & 0x7fffffff
	code := binary % 1000000

	return TOTPInfo{
		Code:             code,
		RemainingSeconds: remainingSeconds,
	}, nil
}
