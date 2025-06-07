package main

import (
	"strings"
	"os"
	"bufio"
	"fmt"
)

type cliCommand struct {
	name			string
	description		string
	callback		func() error
}

var commandsRegistry map[string]cliCommand

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil 
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for _, value := range commandsRegistry {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}


func initializeCommandRegistry() {
	commandsRegistry = make(map[string]cliCommand)

	commandsRegistry["exit"] = cliCommand{
			name:			"exit",
			description:	"Exit the Pokedex",
			callback:		commandExit,
	}
	commandsRegistry["help"] = cliCommand{
			name:			"help",
			description:	"Displays a help message",
			callback:		commandHelp,
	}
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
		command, exists := commandsRegistry[commandName]
		if exists {
			command.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}
