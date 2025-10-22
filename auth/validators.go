package auth

import (
	"regexp"
	"unicode"
)

func IsUsernameValid(username string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`).MatchString(username)
}

func IsPasswordValid(password string) bool {
	re := regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)

	if !re.MatchString(password) {
		return false
	}

	var hasDigit bool = false
	var hasLetter bool = false

	for _, value := range password {
		if unicode.IsLetter(value) {
			hasLetter = true
			break
		}
	}

	for _, value := range password {
		if unicode.IsDigit(value) {
			hasDigit = true
			break
		}
	}

	return hasDigit && hasLetter
}
