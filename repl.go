package main

import (
	"strings"
	"os"
	"bufio"
	"fmt"
	"github.com/skimura1/pokedexcli/internal/pokedexapi"
	"github.com/skimura1/pokedexcli/internal/types"
)


func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)
	config := &types.Config{
		PageNumber: 0,
	}
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
			err := command.callback(config); if err != nil {
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
	callback		func(*types.Config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"help": {
			name:		 "help",
			description: "Displays a help message",
			callback:	 commandHelp,
		},
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex",
			callback:	 commandExit,
		},
		"map": {
			name:		 "map",
			description: "Display next 20 locations",
			callback:	 pokedexapi.CommandMap,
		},
		"mapb": {
			name:		"mapb",
			description: "Diplays previous 20 locations",
			callback:	 pokedexapi.CommandMapB,
		},
	}
}
