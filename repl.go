package main

import "strings"

func cleanInput(text string) []string {
	lowercasedString := strings.ToLower(text)
	return strings.Fields(lowercasedString)
}
