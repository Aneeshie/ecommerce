package validator

import (
	"errors"
	"unicode"
)

var (
	ErrPasswordTooShort           = errors.New("Password is too short")
	ErrPasswordTooLong            = errors.New("Password is too long")
	ErrPasswordMissingUppercase   = errors.New("Password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase   = errors.New("Password must contain at least one lowercase letter")
	ErrPasswordMissingDigit       = errors.New("Password must contain at least one digit")
	ErrPasswordMissingSpecialChar = errors.New("Password must contain at least one special character")
)

func ValidatePassword(password string) error {

	if len(password) < 6 {
		return ErrPasswordTooShort
	}
	if len(password) > 20 {
		return ErrPasswordTooLong
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecialChar := false

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		default:
			hasSpecialChar = true
		}

	}

	if !hasUpper {
		return ErrPasswordMissingUppercase
	}
	if !hasLower {
		return ErrPasswordMissingLowercase
	}
	if !hasDigit {
		return ErrPasswordMissingDigit
	}
	if !hasSpecialChar {
		return ErrPasswordMissingSpecialChar
	}
	return nil
}
