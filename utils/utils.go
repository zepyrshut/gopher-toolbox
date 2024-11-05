package utils

import (
	"fmt"
	"log/slog"
	"time"
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
