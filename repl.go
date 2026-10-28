package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"pokedexcli/internal/pokecache"
)

var (
	POKEMON_LOCATION_AREA_API = "https://pokeapi.co/api/v2/location-area/"
	POKEMON_POKEMON_API       = "https://pokeapi.co/api/v2/pokemon/"
)

type config struct {
	commandRegistry map[string]cliCommand
	nextMap         string
	previousMap     *string
	cache           *pokecache.Cache
	location        string
	pokemon         string
	pokedex         map[string]Pokemon
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

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
			Ability  any  `json:"ability"`
		} `json:"abilities"`
	} `json:"past_abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order any `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		Other struct {
			Home struct {
				FrontShiny       string `json:"front_shiny"`
				FrontFemale      string `json:"front_female"`
				FrontDefault     string `json:"front_default"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			Showdown struct {
				BackShiny        string `json:"back_shiny"`
				BackFemale       string `json:"back_female"`
				FrontShiny       string `json:"front_shiny"`
				BackDefault      string `json:"back_default"`
				FrontFemale      string `json:"front_female"`
				FrontDefault     string `json:"front_default"`
				BackShinyFemale  string `json:"back_shiny_female"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
			DreamWorld struct {
				FrontFemale  any    `json:"front_female"`
				FrontDefault string `json:"front_default"`
			} `json:"dream_world"`
			OfficialArtwork struct {
				Versions struct {
					GenerationI struct {
						RedAndBlue struct {
							FrontDefault string `json:"front_default"`
						} `json:"red-and-blue"`
						RedAndGreen struct {
							FrontDefault string `json:"front_default"`
						} `json:"red-and-green"`
					} `json:"generation-i"`
					GenerationIi struct {
						GoldAndSilver struct {
							FrontDefault any `json:"front_default"`
						} `json:"gold-and-silver"`
					} `json:"generation-ii"`
				} `json:"versions"`
				FrontShiny   string `json:"front_shiny"`
				FrontDefault string `json:"front_default"`
			} `json:"official-artwork"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				Yellow struct {
					BackGbc              string `json:"back_gbc"`
					BackGray             string `json:"back_gray"`
					FrontGbc             string `json:"front_gbc"`
					FrontGray            string `json:"front_gray"`
					BackDefault          string `json:"back_default"`
					FrontDefault         string `json:"front_default"`
					BackTransparent      string `json:"back_transparent"`
					FrontTransparent     string `json:"front_transparent"`
					BackTransparentGray  string `json:"back_transparent_gray"`
					FrontTransparentGray string `json:"front_transparent_gray"`
				} `json:"yellow"`
				RedBlue struct {
					BackGray             string `json:"back_gray"`
					FrontGray            string `json:"front_gray"`
					BackDefault          string `json:"back_default"`
					FrontDefault         string `json:"front_default"`
					BackTransparent      string `json:"back_transparent"`
					FrontTransparent     string `json:"front_transparent"`
					BackTransparentGray  string `json:"back_transparent_gray"`
					FrontTransparentGray string `json:"front_transparent_gray"`
				} `json:"red-blue"`
				RedGreenJapan struct {
					BackGray     string `json:"back_gray"`
					FrontGray    string `json:"front_gray"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"red-green-japan"`
			} `json:"generation-i"`
			GenerationV struct {
				Icons struct {
					Animated struct {
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				BlackWhite struct {
					Animated struct {
						BackShiny        string `json:"back_shiny"`
						BackFemale       string `json:"back_female"`
						FrontShiny       string `json:"front_shiny"`
						BackDefault      string `json:"back_default"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationIi struct {
				Gold struct {
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"gold"`
				Silver struct {
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"silver"`
				Crystal struct {
					Animated struct {
						FrontShiny   string `json:"front_shiny"`
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					BackShiny             string `json:"back_shiny"`
					FrontShiny            string `json:"front_shiny"`
					BackDefault           string `json:"back_default"`
					FrontDefault          string `json:"front_default"`
					BackTransparent       string `json:"back_transparent"`
					FrontTransparent      string `json:"front_transparent"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
				} `json:"crystal"`
			} `json:"generation-ii"`
			GenerationIv struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				Platinum struct {
					Animated struct {
						FrontShiny       string `json:"front_shiny"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
				DiamondPearl struct {
					Animated struct {
						FrontShiny       string `json:"front_shiny"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					Animated struct {
						FrontShiny       string `json:"front_shiny"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
			} `json:"generation-iv"`
			GenerationIx struct {
				Champions struct {
					Icons struct {
						FrontShiny   string `json:"front_shiny"`
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"champions"`
				LegendsZa struct {
					Icons struct {
						FrontShiny       string `json:"front_shiny"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"legends-za"`
				ScarletViolet struct {
					Icons struct {
						FrontFemale  string `json:"front_female"`
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"scarlet-violet"`
			} `json:"generation-ix"`
			GenerationVi struct {
				XY struct {
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
				Icons struct {
					FrontFemale  any    `json:"front_female"`
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				OmegaRubyAlphaSapphire struct {
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omega-ruby-alpha-sapphire"`
			} `json:"generation-vi"`
			GenerationIii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				Emerald struct {
					Animated struct {
						BackShiny    string `json:"back_shiny"`
						FrontShiny   string `json:"front_shiny"`
						BackDefault  string `json:"back_default"`
						FrontDefault string `json:"front_default"`
					} `json:"animated"`
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"emerald"`
				RubySapphire struct {
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"ruby-sapphire"`
				FireredLeafgreen struct {
					BackShiny    string `json:"back_shiny"`
					FrontShiny   string `json:"front_shiny"`
					BackDefault  string `json:"back_default"`
					FrontDefault string `json:"front_default"`
				} `json:"firered-leafgreen"`
			} `json:"generation-iii"`
			GenerationVii struct {
				Icons struct {
					FrontFemale  any    `json:"front_female"`
					FrontDefault string `json:"front_default"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
				LetsGoPikachuLetsGoEevee struct {
					Icons struct {
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"lets-go-pikachu-lets-go-eevee"`
			} `json:"generation-vii"`
			GenerationViii struct {
				SwordShield struct {
					Icons struct {
						FrontFemale  string `json:"front_female"`
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"sword-shield"`
				LegendsArceus struct {
					Icons struct {
						FrontShiny       string `json:"front_shiny"`
						FrontFemale      string `json:"front_female"`
						FrontDefault     string `json:"front_default"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"legends-arceus"`
				BrilliantDiamondShiningPearl struct {
					Icons struct {
						FrontDefault string `json:"front_default"`
					} `json:"icons"`
					BackShiny        string `json:"back_shiny"`
					BackFemale       string `json:"back_female"`
					FrontShiny       string `json:"front_shiny"`
					BackDefault      string `json:"back_default"`
					FrontFemale      string `json:"front_female"`
					FrontDefault     string `json:"front_default"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"brilliant-diamond-shining-pearl"`
			} `json:"generation-viii"`
		} `json:"versions"`
		BackShiny        string `json:"back_shiny"`
		BackFemale       string `json:"back_female"`
		FrontShiny       string `json:"front_shiny"`
		BackDefault      string `json:"back_default"`
		FrontFemale      string `json:"front_female"`
		FrontDefault     string `json:"front_default"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontShinyFemale string `json:"front_shiny_female"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	PastStats []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		} `json:"stats"`
	} `json:"past_stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []any `json:"past_types"`
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

func commandCatch(config *config) error {
	pokemon := config.pokemon

	_, exists := config.cache.Get(pokemon)
	if !exists {
		URL := POKEMON_POKEMON_API + pokemon
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
		config.cache.Add(pokemon, body)
	}
	body, _ := config.cache.Get(pokemon)
	var target Pokemon
	if err := json.Unmarshal(body, &target); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return nil
	}
	name := target.Name
	chance := 100 - (target.BaseExperience / 100)
	n := rand.Intn(100)
	if _, ok := config.pokedex[name]; ok {
		fmt.Println(name + " already caught")
	} else {
		fmt.Println("Throwing a Pokeball at " + name + "...")
		if n < chance {
			fmt.Println(name + " was caught")
			config.pokedex[name] = target
		} else {
			fmt.Println(name + " escaped!")
		}
	}
	config.pokemon = ""
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
				if inputArray[0] == "catch" {
					if len(inputArray) == 2 {
						config.pokemon = inputArray[1]
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
