package esfaker

import (
	"math/rand"
	"strings"
)

const uppercaseAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "0123456789"
const symbols = "!@#$%^&*()_+{}|:<>?~"

var maleNames = []string{
	"Pedro", "Juan", "Pepe", "Francisco", "Luis", "Carlos", "Javier", "José", "Antonio", "Manuel",
}
var femaleNames = []string{
	"María", "Ana", "Isabel", "Laura", "Carmen", "Rosa", "Julia", "Elena", "Sara", "Lucía",
}
var lastNames = []string{
	"García", "Fernández", "González", "Rodríguez", "López", "Martínez", "Sánchez", "Pérez", "Gómez", "Martín",
}

func MaleName() string {
	return maleNames[rand.Intn(len(maleNames))]
}

func FemaleName() string {
	return femaleNames[rand.Intn(len(femaleNames))]
}

func Name() string {
	allNames := append(maleNames, femaleNames...)
	return allNames[rand.Intn(len(allNames))]
}

func LastName() string {
	return lastNames[rand.Intn(len(lastNames))]
}

func Email(beforeAt string) string {
	return beforeAt + "@" + AllChars(5, 10) + ".local"
}

func Int(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func Float(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func Bool() bool {
	return rand.Intn(2) == 0
}

func Chars(min, max int) string {
	var sb strings.Builder
	k := len(lowercaseAlphabet)

	for i := 0; i < rand.Intn(max-min+1)+min; i++ {
		c := lowercaseAlphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func AllChars(min, max int) string {
	allChars := uppercaseAlphabet + lowercaseAlphabet + numbers + symbols
	var sb strings.Builder
	k := len(allChars)

	for i := 0; i < rand.Intn(max-min+1)+min; i++ {
		c := allChars[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func AllCharsOrEmpty(min, max int) string {
	if Bool() {
		return ""
	}
	return AllChars(min, max)
}

func AllCharsOrNil(min, max int) *string {
	if Bool() {
		return nil
	}
	s := AllChars(min, max)
	return &s
}

func NumericString(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(numbers[rand.Intn(len(numbers))])
	}

	return sb.String()
}
