package totp

import (
	"testing"
	"time"
)

// Mock time for testing
type mockTime struct {
	currentTime time.Time
}

func (m *mockTime) Now() time.Time {
	return m.currentTime
}

func TestValidateSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "Valid secret",
			secret:  "JBSWY3DPEHPK3PXP",
			wantErr: false,
		},
		{
			name:    "Valid secret with spaces",
			secret:  "JBSW Y3DP EHPK 3PXP",
			wantErr: false,
		},
		{
			name:    "Invalid characters",
			secret:  "INVALID!@#$",
			wantErr: true,
		},
		{
			name:    "Empty secret",
			secret:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSecret(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCode(t *testing.T) {
	// Set up mock time
	mockTime := &mockTime{
		currentTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	originalTimeNow := timeNow
	timeNow = mockTime.Now
	defer func() { timeNow = originalTimeNow }()

	tests := []struct {
		name          string
		secret        string
		wantCode      uint32
		wantRemaining int64
		wantErr       bool
	}{
		{
			name:          "Valid secret at start of period",
			secret:        "JBSWY3DPEHPK3PXP",
			wantRemaining: 30,
			wantErr:       false,
		},
		{
			name:          "Valid secret mid period",
			secret:        "JBSWY3DPEHPK3PXP",
			wantRemaining: 15,
			wantErr:       false,
		},
		{
			name:          "Invalid secret",
			secret:        "INVALID!@#$%^&*",
			wantRemaining: 0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid secret mid period" {
				mockTime.currentTime = mockTime.currentTime.Add(15 * time.Second)
			}

			got, err := GenerateCode(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check remaining time calculation
				if got.RemainingSeconds != tt.wantRemaining {
					t.Errorf("GenerateCode() remaining = %v, want %v", got.RemainingSeconds, tt.wantRemaining)
				}

				// Verify code is 6 digits
				if got.Code > 999999 {
					t.Errorf("GenerateCode() code = %v, want 6 digits", got.Code)
				}

				// Verify code changes every 30 seconds
				mockTime.currentTime = mockTime.currentTime.Add(30 * time.Second)
				newCode, _ := GenerateCode(tt.secret)
				if got.Code == newCode.Code {
					t.Error("GenerateCode() code should change after 30 seconds")
				}
			}
		})
	}
}

func TestRemainingTimeCalculation(t *testing.T) {
	// Set up mock time
	mockTime := &mockTime{
		currentTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	originalTimeNow := timeNow
	timeNow = mockTime.Now
	defer func() { timeNow = originalTimeNow }()

	secret := "JBSWY3DPEHPK3PXP"
	testTimes := []struct {
		secondsFromStart int
		wantRemaining    int64
	}{
		{0, 30},  // Start of period
		{15, 15}, // Mid period
		{29, 1},  // End of period
		{30, 30}, // Start of next period
	}

	for _, tt := range testTimes {
		t.Run(string(rune(tt.secondsFromStart)), func(t *testing.T) {
			mockTime.currentTime = time.Date(2024, 1, 1, 0, 0, tt.secondsFromStart, 0, time.UTC)
			got, err := GenerateCode(secret)
			if err != nil {
				t.Errorf("GenerateCode() unexpected error = %v", err)
				return
			}

			if got.RemainingSeconds != tt.wantRemaining {
				t.Errorf("Remaining time = %v, want %v at %v seconds from start",
					got.RemainingSeconds, tt.wantRemaining, tt.secondsFromStart)
			}
		})
	}
}
