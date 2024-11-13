package totp

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateCode(t *testing.T) {
	// Known test vectors from RFC 6238 Appendix B
	secretBase32 := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	// List of test timestamps and expected TOTP codes (as strings)
	testCases := []struct {
		timestamp int64
		expected  string
	}{
		{59, "287082"},
		{1111111109, "081804"},
		{1111111111, "050471"},
		{1234567890, "005924"},
		{2000000000, "279037"},
		{20000000000, "353130"},
	}

	for _, tc := range testCases {
		// Override timeNow function for testing
		timeNow = func() time.Time {
			return time.Unix(tc.timestamp, 0)
		}

		code, err := GenerateCode(secretBase32)
		if err != nil {
			t.Fatalf("Error generating TOTP code: %v", err)
		}

		// Convert code to string
		codeStr := fmt.Sprintf("%06d", code)

		if codeStr != tc.expected {
			t.Errorf("At time %d, expected code %s, got %06d", tc.timestamp, tc.expected, code)
		}
	}

	// Reset timeNow to default after tests
	timeNow = time.Now
}

func TestValidateSecret(t *testing.T) {
	validSecret := "JBSWY3DPEHPK3PXP"
	invalidSecret := "INVALID_SECRET"

	if err := ValidateSecret(validSecret); err != nil {
		t.Errorf("Valid secret was marked invalid: %v", err)
	}

	if err := ValidateSecret(invalidSecret); err == nil {
		t.Error("Invalid secret was not detected")
	}
}
