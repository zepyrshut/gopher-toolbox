package gorender

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func or(a, b string) bool {
	if a == "" && b == "" {
		return false
	}
	return true
}

// containsErrors hace una función similar a "{{ with index ... }}" con el
// añadido de que puede pasarle más de un argumento y comprobar si alguno de
// ellos está en el mapa de errores.
//
// Ejemplo:
//
//	{{ if containsErrors .FormData.Errors "name" "email" }}
//	 {{index .FormData.Errors "name" }}
//	 {{index .FormData.Errors "email" }}
//	{{ end }}
func containsErrors(errors map[string]string, names ...string) bool {
	for _, name := range names {
		if _, ok := errors[name]; ok {
			return true
		}
	}
	return false
}

func loadTranslations(language string) map[string]string {
	translations := make(map[string]string)
	filePath := fmt.Sprintf("%s.translate", language)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening translation file:", err)
		return translations
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			translations[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading translation file:", err)
	}

	return translations
}

func translateKey(key string) string {
	translations := loadTranslations("es_ES")
	translated := translations[key]
	if translated != "" {
		return translated
	}
	return key
}
