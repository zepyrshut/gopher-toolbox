package binding

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_mapForm(t *testing.T) {
	var someStruct struct {
		StringType  string  `form:"stringtype"`
		IntType     int     `form:"inttype"`
		Int8Type    int8    `form:"int8type"`
		Int16Type   int16   `form:"int16type"`
		Int32Type   int32   `form:"int32type"`
		Int64Type   int64   `form:"int64type"`
		UintType    uint    `form:"uinttype"`
		Uint8Type   uint8   `form:"uint8type"`
		Uint16Type  uint16  `form:"uint16type"`
		Uint32Type  uint32  `form:"uint32type"`
		Uint64Type  uint64  `form:"uint64type"`
		Float32Type float32 `form:"float32type"`
		Float64Type float64 `form:"float64type"`
		BoolType    bool    `form:"booltype"`
	}

	formData := map[string][]string{
		"stringtype": {"stringType"},
		"inttype":    {"-2147483647"},
		"int8type":   {"-127"},
		"int16type":  {"-32767"},
		"int32type":  {"-2147483647"},
		"int64type":  {"-9223372036854775807"},
		"uinttype":   {"4294967295"},
		"uint8type":  {"255"},
		"uint16type": {"65535"},
		"uint32type": {"4294967295"},
		"uint64type": {"18446744073709551615"},
		"float32type": {
			"3.1415927",
		},
		"float64type": {
			"3.141592653589793",
		},
		"booltype": {"true"},
	}

	err := mapForm(&someStruct, formData)
	require.NoError(t, err)
	require.Equal(t, "stringType", someStruct.StringType)
	require.Equal(t, int(-2147483647), someStruct.IntType)
	require.Equal(t, int8(-127), someStruct.Int8Type)
	require.Equal(t, int16(-32767), someStruct.Int16Type)
	require.Equal(t, int32(-2147483647), someStruct.Int32Type)
	require.Equal(t, int64(-9223372036854775807), someStruct.Int64Type)
	require.Equal(t, uint(4294967295), someStruct.UintType)
	require.Equal(t, uint8(255), someStruct.Uint8Type)
	require.Equal(t, uint16(65535), someStruct.Uint16Type)
	require.Equal(t, uint32(4294967295), someStruct.Uint32Type)
	require.Equal(t, uint64(18446744073709551615), someStruct.Uint64Type)
	require.Equal(t, float32(3.1415927), someStruct.Float32Type)
	require.Equal(t, float64(3.141592653589793), someStruct.Float64Type)
	require.Equal(t, true, someStruct.BoolType)
	t.Log(someStruct)
}
