package utils

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func CorrectTimezone(timeStamp time.Time) time.Time {
	loc, _ := time.LoadLocation("Europe/Madrid")
	return timeStamp.In(loc)
}

func GetBool(value string) bool {
	return value == "true"
}

func LogAndReturnError(err error, message string) error {
	slog.Error(message, "error", err.Error())
	return fmt.Errorf("%s: %w", message, err)
}

func GetBoolFromString(s string) bool {
	if s == "S" || s == "s" {
		return true
	}
	return false
}

func Slugify(s string) string {
	s = strings.ToLower(s)

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	s, _, _ = transform.String(t, s)

	s = strings.ReplaceAll(s, " ", "-")

	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	s = b.String()

	re := regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	return s
}
