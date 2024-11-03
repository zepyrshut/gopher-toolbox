package pgutils

import (
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToFloat64(n pgtype.Numeric) float64 {
	value, err := n.Value()
	if err != nil {
		slog.Error("error getting numeric value", "error", err)
		return 0
	}

	strValue, ok := value.(string)
	if !ok {
		slog.Error("error converting numeric value to string")
		return 0
	}

	floatValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		slog.Error("error converting string to float", "error", err)
		return 0
	}

	return floatValue
}

func FloatToNumeric(number float64) (value pgtype.Numeric) {
	parse := strconv.FormatFloat(number, 'f', -1, 64)
	if err := value.Scan(parse); err != nil {
		slog.Error("error scanning float to numeric", "error", err)
	}
	return value
}
