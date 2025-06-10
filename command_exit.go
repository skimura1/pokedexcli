package main

import (
	"fmt"
	"os"
)

func commandExit(cfg *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil 
}
