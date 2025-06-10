package main

import (
	"time"

	"github.com/skimura1/pokedexcli/internal/pokedexapi"
)


func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
	}

	startRepl(cfg)
}
