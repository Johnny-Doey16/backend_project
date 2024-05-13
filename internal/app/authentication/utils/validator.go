package utils

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	startsWithLetterRegex                   = regexp.MustCompile(`^[A-Za-z]`)
	isValidCharactersRegex                  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	hasConsecutiveUnderscoresOrHyphensRegex = regexp.MustCompile(`(__|-{2,})`)
	endsWithUnderscoreOrHyphenRegex         = regexp.MustCompile(`[_-]$`)
	isValidEmailRegex                       = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	isValidPhoneRegex                       = regexp.MustCompile(`^\+?1?\d{10}$`)
	hasUppercaseRegex                       = regexp.MustCompile(`[A-Z]`)
	hasLowercaseRegex                       = regexp.MustCompile(`[a-z]`)
	hasNumberRegex                          = regexp.MustCompile(`[0-9]`)
	hasSpecialCharRegex                     = regexp.MustCompile(`[!@#$%^&*()]`)
	cleanNumberRegex                        = regexp.MustCompile(`[\s-().]`)
	hasTenDigitsRegex                       = regexp.MustCompile(`^\d{10}$`)
	cityRegex                               = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	postalCodeRegex                         = regexp.MustCompile(`^\d{5}$`)
	stateRegex                              = regexp.MustCompile(`^[a-zA-Z]+$`)
	countryRegex                            = regexp.MustCompile(`^[a-zA-Z]+$`)
)

const (
	usernameMinLen = 4
	usernameMaxLen = 30
	passwordMinLen = 8
	passwordMaxLen = 64
	fullNameMinLen = 3
	fullNameMaxLen = 30
	emailMinLen    = 3
	emailMaxLen    = 254
)

func ValidateString(value string, min int, max int) error {
	if n := len(value); n < min || n > max {
		return fmt.Errorf("length must be between %d and %d characters", min, max)
	}
	return nil
}

// other validation functions...

func ValidateUsername(username string) bool {
	if len(username) < usernameMinLen || len(username) > usernameMaxLen {
		return false
	}
	if !startsWithLetterRegex.MatchString(username) {
		return false
	}
	if !isValidCharactersRegex.MatchString(username) {
		return false
	}
	if hasConsecutiveUnderscoresOrHyphensRegex.MatchString(username) {
		return false
	}
	if endsWithUnderscoreOrHyphenRegex.MatchString(username) {
		return false
	}
	return true
}

func ValidatePhone(phoneNumber string) bool {
	cleanNumber := cleanNumberRegex.ReplaceAllString(phoneNumber, "")
	return hasTenDigitsRegex.MatchString(cleanNumber)
}

func ValidatePassword(password string) bool {
	if len(password) <= passwordMinLen || len(password) >= passwordMaxLen {
		return false
	}
	if !hasUppercaseRegex.MatchString(password) || !hasLowercaseRegex.MatchString(password) || !hasNumberRegex.MatchString(password) || !hasSpecialCharRegex.MatchString(password) {
		return false
	}
	commonPatterns := []string{"123456", "password"}
	for _, pattern := range commonPatterns {
		if password == pattern {
			return false
		}
	}
	return true
}

func ValidateEmail(value string) error {
	// if err := ValidateString(value, 3, 254); err != nil {
	// 	return fmt.Errorf("invalid email length: %s", err)
	// }

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email: %s", err)
	}
	return nil
}

func ValidateAddress(address string) bool {
	return address != ""
}

func ValidateCity(city string) bool {
	return cityRegex.MatchString(city)
}

func ValidatePostalCode(postalCode string) bool {
	return postalCodeRegex.MatchString(postalCode)
}

func ValidateState(state string) bool {
	return stateRegex.MatchString(state)
}

func ValidateCountry(country string) bool {
	return countryRegex.MatchString(country)
}

func IsEmailFormat(value string) bool {
	return isValidEmailRegex.MatchString(value)
}

func IsPhoneFormat(value string) bool {
	return isValidPhoneRegex.MatchString(value)
}

// other validation functions...
