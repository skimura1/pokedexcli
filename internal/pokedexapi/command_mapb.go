package pokedexapi

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"github.com/skimura1/pokedexcli/internal/types"
)

func CommandMapB(cfg *types.Config) error {
	if cfg.PageNumber == 0 {
		fmt.Println("you're on the first page")
		return nil
	} else {
		cfg.PageNumber -= 1
	}
	offset := cfg.PageNumber * 20
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?limit=20&offset=%d", offset)
	res, err := http.Get(url)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Println("Response failed with status code: %d and \n body: %s\n", res.StatusCode, body)
	}
	locationArea := types.LocationArea{}
	err = json.Unmarshal(body, &locationArea)
	for _, location := range locationArea.Results {
		fmt.Println(location.Name)
	}
	return err
}
