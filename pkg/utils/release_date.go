package utils

import (
	"time"

	"github.com/goodsign/monday"
)

func ParseReleaseDate(layout, raw string, locale monday.Locale) (*time.Time, error) {
	result, err := monday.Parse(layout, raw, locale)
	return &result, err
}
