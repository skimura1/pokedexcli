package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"pokedexcli/internal/pokecache"
)

var POKEMON_LOCATION_AREA_API = "https://pokeapi.co/api/v2/location-area/"

type config struct {
	commandRegistry map[string]cliCommand
	nextMap         string
	previousMap     *string
	cache           *pokecache.Cache
	location        string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type Map struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Result   []struct {
		Name string `json:"name"`
		URL  string `json:"URL"`
	} `json:"results"`
}

type Area struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel int `json:"min_level"`
				MaxLevel int `json:"max_level"`
				Chance   int `json:"chance"`
				Method   struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				ConditionValues []any `json:"condition_values"`
				PokemonDetails  any   `json:"pokemon_details"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: ")
	for commandName, commandMap := range config.commandRegistry {
		fmt.Printf("%s: %s\n", commandName, commandMap.description)
	}
	return nil
}

func commandMap(config *config) error {
	_, exists := config.cache.Get(config.nextMap)
	if !exists {
		res, err := http.Get(config.nextMap)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			return err
		}
		config.cache.Add(config.nextMap, body)
	}
	body, _ := config.cache.Get(config.nextMap)

	var target Map
	if err := json.Unmarshal(body, &target); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil
	}
	locations := target.Result
	for _, location := range locations {
		fmt.Println(location.Name)
	}
	config.nextMap = target.Next
	config.previousMap = target.Previous
	return nil
}

func commandMapb(config *config) error {
	if config.previousMap == nil {
		return errors.New("No previous map")
	}
	_, exists := config.cache.Get(*config.previousMap)
	if !exists {
		res, err := http.Get(*config.previousMap)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			return err
		}
		config.cache.Add(*config.previousMap, body)
	}
	body, _ := config.cache.Get(*config.previousMap)
	var target Map
	if err := json.Unmarshal(body, &target); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil
	}
	locations := target.Result
	for _, location := range locations {
		fmt.Println(location.Name)
	}
	config.nextMap = target.Next
	config.previousMap = target.Previous
	return nil
}

func commandExplore(config *config) error {
	location := config.location
	if len(config.location) == 0 {
		return errors.New("No location! Please provide a location to explore")
	}
	_, exists := config.cache.Get(location)
	if !exists {
		URL := POKEMON_LOCATION_AREA_API + location
		res, err := http.Get(URL)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			return err
		}
		config.cache.Add(location, body)
	}
	body, _ := config.cache.Get(location)
	var target Area
	if err := json.Unmarshal(body, &target); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil
	}
	fmt.Println("Exploring " + location)
	fmt.Println("Found Pokemon: ")
	pokemonEncounters := target.PokemonEncounters
	for _, encounters := range pokemonEncounters {
		fmt.Println(" - " + encounters.Pokemon.Name)
	}
	config.location = ""
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
				if inputArray[0] == "explore" {
					if len(inputArray) == 2 {
						config.location = inputArray[1]
					}
				}
				if err := command.callback(config); err != nil {
					fmt.Println("Error ocurred: ", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
