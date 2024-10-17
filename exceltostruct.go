package exceltostruct

import (
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

func Convert[T any](bookPath, sheetName string) (dataExcel []T) {
	f, _ := excelize.OpenFile(bookPath)
	rows, _ := f.GetRows(sheetName)

	firstRow := map[string]int{}

	for i, row := range rows[0] {
		firstRow[row] = i
	}

	t := new(T)
	dataExcel = make([]T, 0, len(rows)-1)

	for _, row := range rows[1:] {
		v := reflect.ValueOf(t)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag.Get("xlsx")
			objType := v.Field(i).Type().String()

			if j, ok := firstRow[tag]; ok {
				field := v.Field(i)
				if len(row) > j {
					d := row[j]
					elementConverted := convertType(objType, d)
					field.Set(reflect.ValueOf(elementConverted))
				}
			}
		}

		dataExcel = append(dataExcel, *t)
	}

	return dataExcel
}

func convertType(objType string, value string) any {
	switch objType {
	case "int":
		valueInt, _ := strconv.Atoi(value)
		return valueInt
	case "bool":
		valueBool, _ := strconv.ParseBool(value)
		return valueBool
	case "float32":
		valueFloat, _ := strconv.ParseFloat(value, 32)
		return float32(valueFloat)
	case "string":
		return value
	}
	return value
}
