package helper

import "regexp"

func ClearString(str string) string {
	var nonAlphanumericRegex = regexp.MustCompile(`[^0-9\-]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
