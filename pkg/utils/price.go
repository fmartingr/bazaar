package utils

import (
	"regexp"
	"strings"
)

func ExtractPrice(raw string) string {
	re := regexp.MustCompile("[^0-9,.]+")
	return strings.Replace(re.ReplaceAllString(raw, ""), ",", ".", 1)
}
