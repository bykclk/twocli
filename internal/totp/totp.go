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
	if _, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret); err != nil {
		return err
	}
	return nil
}

// GenerateCode generates a TOTP code based on the provided secret key.
// The secret should be a base32 encoded string.
func GenerateCode(secret string) (uint32, error) {
	if err := ValidateSecret(secret); err != nil {
		return 0, fmt.Errorf("invalid secret key: %v", err)
	}

	// Decode the base32 encoded secret key
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return 0, errors.New("failed to decode secret key")
	}

	// Calculate the time step (number of 30-second intervals since Unix epoch)
	epochSeconds := timeNow().Unix()
	timeStep := uint64(epochSeconds / 30)

	// Convert time step to byte array (big-endian)
	var timeBytes [8]byte
	binary.BigEndian.PutUint64(timeBytes[:], timeStep)

	// Compute HMAC-SHA1 hash
	hash := hmac.New(sha1.New, key)
	hash.Write(timeBytes[:])
	hmacHash := hash.Sum(nil)

	// Perform dynamic truncation to extract a 4-byte string
	offset := hmacHash[len(hmacHash)-1] & 0x0F
	if offset+4 > byte(len(hmacHash)) {
		return 0, errors.New("invalid HMAC hash length")
	}
	truncatedHash := hmacHash[offset : offset+4]

	// Convert truncated hash to uint32
	code := binary.BigEndian.Uint32(truncatedHash) & 0x7FFFFFFF

	// Modulo operation to get the final TOTP code (6 digits)
	otp := code % 1000000

	return otp, nil
}
