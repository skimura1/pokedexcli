package main

import (
	"fmt"
	"github.com/skimura1/pokedexcli/internal/types"
)

func commandHelp(cfg *types.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for _, value := range getCommands() {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}
