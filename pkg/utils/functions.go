package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CapitalizeWords(word string) string {
	word = strings.TrimSpace(word)
	caser := cases.Title(language.Spanish)
	return caser.String(word)
}

func GetFirstNameAndFirstLastName(name string) string {
	name = strings.TrimSpace(name)
	nameArray := strings.Split(name, " ")
	if len(nameArray) <= 2 {
		return name
	}

	return nameArray[0] + " " + nameArray[2]
}
