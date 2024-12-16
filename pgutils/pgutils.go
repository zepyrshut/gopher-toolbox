package pgutils

import (
	"log/slog"
	"math"
	"math/big"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func NumericToFloat64(n pgtype.Numeric) float64 {
	val, err := n.Value()
	if err != nil {
		slog.Error("error getting numeric value", "error", err)
		return 0
	}

	strValue, ok := val.(string)
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

func NumericToInt64(n pgtype.Numeric) int64 {
	return n.Int.Int64() * int64(math.Pow(10, float64(n.Exp)))
}

func FloatToNumeric(number float64, precision int) (value pgtype.Numeric) {
	parse := strconv.FormatFloat(number, 'f', precision, 64)
	slog.Info("parse", "parse", parse)

	if err := value.Scan(parse); err != nil {
		slog.Error("error scanning numeric", "error", err)
	}
	return value
}

func AddNumeric(a, b pgtype.Numeric) pgtype.Numeric {
	minExp := a.Exp
	if b.Exp < minExp {
		minExp = b.Exp
	}

	aInt := new(big.Int).Set(a.Int)
	bInt := new(big.Int).Set(b.Int)

	for a.Exp > minExp {
		aInt.Mul(aInt, big.NewInt(10))
		a.Exp--
	}
	for b.Exp > minExp {
		bInt.Mul(bInt, big.NewInt(10))
		b.Exp--
	}

	resultado := new(big.Int).Add(aInt, bInt)

	return pgtype.Numeric{
		Int:   resultado,
		Exp:   minExp,
		Valid: true,
	}
}

func SubtractNumeric(a, b pgtype.Numeric) pgtype.Numeric {
	minExp := a.Exp
	if b.Exp < minExp {
		minExp = b.Exp
	}

	aInt := new(big.Int).Set(a.Int)
	bInt := new(big.Int).Set(b.Int)

	for a.Exp > minExp {
		aInt.Mul(aInt, big.NewInt(10))
		a.Exp--
	}
	for b.Exp > minExp {
		bInt.Mul(bInt, big.NewInt(10))
		b.Exp--
	}

	resultado := new(big.Int).Sub(aInt, bInt)

	return pgtype.Numeric{
		Int:   resultado,
		Exp:   minExp,
		Valid: true,
	}
}
