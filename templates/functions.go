package templates

import (
	"errors"
	"strconv"
	"time"
)

func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func FormatDateSpanish(date time.Time) string {
	months := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	days := []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}

	dayName := days[date.Weekday()]
	day := date.Day()
	month := months[date.Month()-1]
	year := date.Year()

	return dayName + ", " + strconv.Itoa(day) + " de " + month + " de " + strconv.Itoa(year)
}
