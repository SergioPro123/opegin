package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CapitalizeWords(word string) string {
	caser := cases.Title(language.Spanish)
	return caser.String(word)
}
