package util

import "unicode/utf8"

const maxlen = 180

func Truncate(s string) string {
	if utf8.RuneCountInString(s) <= maxlen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxlen]) + "..."
}
