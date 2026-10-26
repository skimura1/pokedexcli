package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type config struct {
	commandRegistry map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: \n")
	for commandName, commandMap := range config.commandRegistry {
		fmt.Printf("%s: %s\n", commandName, commandMap.description)
	}
	return nil
}

func cleanInput(text string) []string {
	lowercasedString := strings.ToLower(text)
	return strings.Fields(lowercasedString)
}

func repl(config *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		if scanner.Scan() {
			inputArray := cleanInput(scanner.Text())
			if command, ok := config.commandRegistry[inputArray[0]]; ok {
				if err := command.callback(config); err != nil {
					fmt.Println("Error ocurred: ", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
