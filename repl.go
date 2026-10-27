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

type config struct {
	commandRegistry map[string]cliCommand
	nextMap         string
	previousMap     *string
	cache           *pokecache.Cache
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

func commandExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage: \n")
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
				if err := command.callback(config); err != nil {
					fmt.Println("Error ocurred: ", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}
