package main

import (
	"strings"
	"os"
	"bufio"
	"fmt"
)

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan();
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		commandName := input[0]

		fmt.Println("Your command was:", commandName)
	}
}
