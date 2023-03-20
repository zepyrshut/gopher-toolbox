package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strconv"
)

type User struct {
	Id        int     `xlsx:"id"`
	FirstName string  `xlsx:"first_name" validate:"required"`
	LastName  string  `xlsx:"last_name"`
	Email     string  `xlsx:"email"`
	Gender    bool    `xlsx:"gender"`
	Balance   float32 `xlsx:"balance"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()

	data := excelToStruct[User]("Book1.xlsx", "Sheet1")

	ok, errList := someValidation(data)
	if ok {
		fmt.Println("Data is valid")
		fmt.Println(data)
	} else {
		fmt.Println("Data is not valid")
		fmt.Println(errList)
	}
}

func someValidation[T any](data []T) (bool, []error) {
	var errList []error

	for _, v := range data {
		err := validate.Struct(v)
		if err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return false, errList
	} else {
		return true, nil
	}
}

func excelToStruct[T any](bookPath, sheetName string) (dataExcel []T) {
	f, _ := excelize.OpenFile(bookPath)
	rows, _ := f.GetRows(sheetName)

	firstRow := map[string]int{}

	for i, row := range rows[0] {
		firstRow[row] = i
	}

	t := new(T)
	dataExcel = make([]T, 0, len(rows)-1)

	for _, row := range rows[1:3] {
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
