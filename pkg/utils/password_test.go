package utils

import (
	"testing"
)

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"empty string", "", true},
		{"too short with letters and digits", "Ab1", true},
		{"exactly 7 chars", "Abcdef1", true},
		{"only digits 8 chars", "12345678", true},
		{"only letters 8 chars", "abcdefgh", true},
		{"only uppercase 8 chars", "ABCDEFGH", true},
		{"only special chars", "!@#$%^&*", true},
		{"digits and specials no letters", "1234!@#$", true},
		{"letters and specials no digits", "abcd!@#$", true},

		{"valid lowercase and digit", "abcdefg1", false},
		{"valid uppercase and digit", "ABCDEFG1", false},
		{"valid mixed case and digit", "Abcdefg1", false},
		{"exactly 8 chars valid", "aaaaaaa1", false},
		{"long valid password", "abcdefghijklmnop123", false},
		{"digit at start", "1abcdefg", false},
		{"unicode letter with digit", "àbcdefg1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr && err == nil {
				t.Errorf("expected error for password %q, got nil", tt.password)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error for password %q, got %v", tt.password, err)
			}
			if tt.wantErr && err != nil && err != ErrWeakPassword {
				t.Errorf("expected ErrWeakPassword, got %v", err)
			}
		})
	}
}

func TestValidatePasswordStrength_ErrorIdentity(t *testing.T) {
	err := ValidatePasswordStrength("short")
	if err != ErrWeakPassword {
		t.Errorf("expected ErrWeakPassword sentinel, got %v", err)
	}
}
