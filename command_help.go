package main

import (
	"fmt"
)

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")

	for _, value := range getCommands() {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}
