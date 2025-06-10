package main

import (
	"fmt"
	"os"
	"github.com/skimura1/pokedexcli/internal/types"
)

func commandExit(cfg *types.Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil 
}
