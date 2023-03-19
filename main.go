package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

type User struct {
	Id        int     `xlsx:"id"`
	FirstName string  `xlsx:"first_name"`
	LastName  string  `xlsx:"last_name"`
	Email     string  `xlsx:"email"`
	Gender    bool    `xlsx:"gender"`
	Balance   float32 `xlsx:"balance"`
}

func main() {

	f, _ := excelize.OpenFile("Book1.xlsx")
	rows, _ := f.GetRows("Sheet1")

	var dataExcel []User

	firstRow := map[string]int{}

	for i, row := range rows[0] {
		firstRow[row] = i
	}

	dataExcel = make([]User, 0, len(rows)-1)

	for _, row := range rows[1:3] {
		u := User{}
		v := reflect.ValueOf(&u)
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

		dataExcel = append(dataExcel, u)
	}

	fmt.Println(dataExcel)
}

func convertType(objType string, value string) interface{} {
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
