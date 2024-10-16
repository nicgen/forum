package lib

import (
	"regexp"
	"unicode"
)

// ? Function to check is the username is in the correct format
func IsValidPassword(password string) bool {
	// Check if the username is at least 8 characters long
	if len(password) <= 8 {
		return false
	}

	// Regular expression pattern for special characters
	specialChars := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`)

	hasSpecialChar := specialChars.MatchString(password)
	hasNumber := false

	// Check for at least one number
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}

	// Return true if all conditions are met
	return hasSpecialChar && hasNumber
}

// ? Function to check is the email is in the correct format
func IsValidEmail(email string) bool {
	// Regular expression pattern
	// ^[^@\s]+: Start with one or more characters that are not @ or whitespace
	// @: Must contain an @ symbol
	// [^@\s]+: Followed by one or more characters that are not @ or whitespace
	// \.: Must contain a dot
	// [^@\s]{2,}: Ends with at least two characters that are not @ or whitespace
	pattern := `^[^@\s]+@[^@\s]+\.[^@\s]{2,}$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Test the email against the pattern
	return regex.MatchString(email)
}
