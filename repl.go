package main

import (
	"strings"
	"os"
	"bufio"
	"fmt"

	"github.com/skimura1/pokedexcli/internal/pokedexapi"
)

type config struct {
	pokeapiClient pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
}


func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan();
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		commandName := input[0]

		command, exists := getCommands()[commandName]
		if exists {
			err := command.callback(cfg); if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name			string
	description		string
	callback		func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"help": {
			name:		 "help",
			description: "Displays a help message",
			callback:	 commandHelp,
		},
		"map": {
			name:		 "map",
			description: "Display next 20 locations",
			callback:	 commandMapf,
		},
		"mapb": {
			name:		"mapb",
			description: "Diplays previous 20 locations",
			callback:	 commandMapb,
		},
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex",
			callback:	 commandExit,
		},
	}
}
