package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

// GenerateEmailToken creates a secure random token for email verification
func GenerateEmailToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendVerificationEmail - placeholder for email sending
// In production, integrate with services like SendGrid, AWS SES, etc.
func SendVerificationEmail(email, token, displayName string) error {
	// For now, just log the verification link
	verificationURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
	log.Printf("[MOCK EMAIL] Verification email for %s (%s): %s", displayName, email, verificationURL)
	return nil
}

// SendPasswordResetEmail - placeholder for password reset email
func SendPasswordResetEmail(email, token, displayName string) error {
	// For now, just log the reset link
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)
	log.Printf("[MOCK EMAIL] Password reset email for %s (%s): %s", displayName, email, resetURL)
	return nil
}

// IsValidEmail performs basic email validation
func IsValidEmail(email string) bool {
	// Basic email validation - in production use a proper email validation library
	return len(email) > 5 &&
		len(email) < 254 &&
		contains(email, "@") &&
		contains(email, ".")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
