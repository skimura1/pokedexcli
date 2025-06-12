package pokeapi

import (
	"net/http"
	"time"

	"github.com/skimura1/pokedexcli/internal/pokedexcache"
)

type Client struct {
	httpClient http.Client
	cache pokedexcache.Cache
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokedexcache.NewCache(30 * time.Second),
	}
}
