package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
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

func LoadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value)
	}
	return scanner.Err()
}

func Slugify(s string) string {
	s = strings.ToLower(s)

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
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

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
