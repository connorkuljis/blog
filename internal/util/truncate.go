package util

import "unicode/utf8"

func Truncate(s string) string {
	if utf8.RuneCountInString(s) <= 80 {
		return s
	}
	runes := []rune(s)
	return string(runes[:80]) + "..."
}
