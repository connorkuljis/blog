package util

import "unicode/utf8"

const defaultMaxLen = 180

// Truncate shortens s to n runes if n > 0, otherwise uses defaultMaxLen.
func Truncate(s string, n ...int) string {
	max := defaultMaxLen
	if len(n) > 0 && n[0] > 0 {
		max = n[0]
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	runes := []rune(s)
	return string(runes[:max]) + "..."
}
