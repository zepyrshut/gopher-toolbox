package gorender

import (
	"strings"

	spanish "github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	esTranslations "github.com/go-playground/validator/v10/translations/es"
)

type FormData struct {
	HasErrors bool
	Errors    map[string]string
	Values    map[string]string
}

func NewForm() FormData {
	return FormData{
		HasErrors: false,
		Errors:    map[string]string{},
		Values:    map[string]string{},
	}
}

// AddError añade errores a la estructura FormData, es un mapa cuya clave es una
// cadena de carecteres. Hay que tener en cuenta que cuando se hace una
// validación, se llama a esta función cuya clave es el nombre del campo con lo
// cual si hay más de un error de validación se sobreescriben el anterior y sólo
// se muestra el último error.
func (fd *FormData) AddError(field, message string) {
	fd.HasErrors = true
	fd.Errors[field] = message
}

func (fd *FormData) AddValue(field, value string) {
	fd.Values[field] = value
}

type ValidationError struct {
	Field  string
	Reason string
}

func (fd *FormData) ValidateStruct(s interface{}) (map[string]string, error) {
	spanishTranslator := spanish.New()
	uni := ut.New(spanishTranslator, spanishTranslator)
	trans, _ := uni.GetTranslator("es")
	validate := validator.New()
	_ = esTranslations.RegisterDefaultTranslations(validate, trans)
	errors := make(map[string]string)
	var validationErrors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fd.AddError("form-error", "Error de validación de datos.")
			return errors, err
		}

		for _, err := range err.(validator.ValidationErrors) {
			fieldName, _ := trans.T(err.Field())
			message := strings.Replace(err.Translate(trans), err.Field(), fieldName, -1)

			validationErrors = append(validationErrors, ValidationError{
				Field:  strings.ToLower(err.Field()),
				Reason: correctMessage(message),
			})
		}

		for _, err := range validationErrors {
			errors[err.Field] = err.Reason
		}

		if len(errors) > 0 {
			fd.Errors = errors
			fd.HasErrors = true
		}

		return errors, err
	}

	return errors, nil
}

func correctMessage(s string) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	if runes[len(runes)-1] != '.' {
		runes = append(runes, '.')
	}

	return string(runes)
}
