package helper

import "strings"

func EmailNotValid(username, password string) bool {
	username, password = strings.ToLower(username), strings.ToLower(password)

	return strings.Contains(password, username)
}
