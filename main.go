package main

import (
	"time"

	"pokedexcli/internal/pokecache"
)

func main() {
	c := pokecache.NewCache(10 * time.Second)
	config := config{
		commandRegistry: map[string]cliCommand{
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"map": {
				name:        "map",
				description: "Displays location on the map",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Displays previous locations on the map",
				callback:    commandMapb,
			},
			"explore": {
				name:        "explore",
				description: "Displays name of all Pokemon in that area",
				callback:    commandExplore,
			},
		},
		cache:       c,
		nextMap:     "https://pokeapi.co/api/v2/location-area/",
		previousMap: nil,
		location:    "",
	}
	repl(&config)
}
