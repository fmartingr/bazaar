package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func ExtractPrice(raw string) string {
	re := regexp.MustCompile("[^0-9,.]+")
	return strings.Replace(re.ReplaceAllString(raw, ""), ",", ".", 1)
}

func StripLastCharacter(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
