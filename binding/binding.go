package binding

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

type FormBinding struct{}

func (FormBinding) Bind(r *http.Request, obj any) error {
	if r == nil {
		return errors.New("request is nil")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	if r.Form == nil {
		return errors.New("form is nil")
	}

	return mapForm(obj, r.Form)
}

func mapForm(obj any, form map[string][]string) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("obj must be a non-nil pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("obj must be a pointer to a struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		formTag := fieldType.Tag.Get("form")

		if formTag == "" {
			formTag = fieldType.Name
		}

		if values, ok := form[formTag]; ok && len(values) > 0 {
			value := values[0]
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(intValue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				uintValue, err := strconv.ParseUint(value, 10, 64)
				if err != nil {
					return err
				}
				field.SetUint(uintValue)
			case reflect.Float32, reflect.Float64:
				floatValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				field.SetFloat(floatValue)
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return err
				}
				field.SetBool(boolValue)
			default:
				return errors.New("unsupported field type")
			}
		}
	}
	return nil
}
