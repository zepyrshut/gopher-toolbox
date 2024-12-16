package pgutils

import (
	"fmt"
	"log/slog"
	"math/big"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func Test_NumericToFloat(t *testing.T) {
	tests := []struct {
		given    pgtype.Numeric
		expected float64
	}{
		{given: pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true}, expected: 0},
		{given: pgtype.Numeric{Int: big.NewInt(5), Exp: 0, Valid: true}, expected: 5},
		{given: pgtype.Numeric{Int: big.NewInt(10), Exp: 0, Valid: true}, expected: 10},
		{given: pgtype.Numeric{Int: big.NewInt(1000), Exp: -2, Valid: true}, expected: 10.00},
		{given: pgtype.Numeric{Int: big.NewInt(1000), Exp: -3, Valid: true}, expected: 1.000},
		{given: pgtype.Numeric{Int: big.NewInt(1000), Exp: -4, Valid: true}, expected: 0.1000},
		{given: pgtype.Numeric{Int: big.NewInt(2555), Exp: -2, Valid: true}, expected: 25.55},
		{given: pgtype.Numeric{Int: big.NewInt(-1), Exp: -2, Valid: true}, expected: -0.01},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("given %v, expected %v", test.given, test.expected), func(t *testing.T) {
			actual := NumericToFloat64(test.given)
			assert.Equal(t, test.expected, actual)
		})
	}

}

func Test_FloatToNumeric(t *testing.T) {
	tests := []struct {
		given     float64
		expected  pgtype.Numeric
		precision int
	}{
		{given: 0.0, expected: pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true}, precision: 0},
		{given: 25.50, expected: pgtype.Numeric{Int: big.NewInt(2550), Exp: -2, Valid: true}, precision: 2},
		{given: 100.0, expected: pgtype.Numeric{Int: big.NewInt(1), Exp: 2, Valid: true}, precision: 0},
		{given: 0.0001, expected: pgtype.Numeric{Int: big.NewInt(0), Exp: -2, Valid: true}, precision: 2},
		{given: 0.0001, expected: pgtype.Numeric{Int: big.NewInt(1), Exp: -4, Valid: true}, precision: 4},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("given %v, expected %v", test.given, test.expected), func(t *testing.T) {
			actual := FloatToNumeric(test.given, test.precision)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_AddNumeric(t *testing.T) {
	valueA := pgtype.Numeric{Int: big.NewInt(1), Exp: -3, Valid: true}
	valueB := pgtype.Numeric{Int: big.NewInt(2), Exp: -2, Valid: true}

	slog.Info("valueA", "valueA", valueA)
	slog.Info("valueB", "valueB", valueB)

	actual := AddNumeric(valueA, valueB)

	slog.Info("actual", "actual", actual)
	assert.Equal(t, pgtype.Numeric{Int: big.NewInt(21), Exp: -3, Valid: true}, actual)
}

func Test_SubtractNumeric(t *testing.T) {
	valueA := pgtype.Numeric{Int: big.NewInt(1), Exp: -3, Valid: true}
	valueB := pgtype.Numeric{Int: big.NewInt(2), Exp: -2, Valid: true}

	actual := SubtractNumeric(valueA, valueB)

	slog.Info("actual", "actual", actual)
	assert.Equal(t, pgtype.Numeric{Int: big.NewInt(-19), Exp: -3, Valid: true}, actual)
}
