package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/awbalessa/pokedex/internal/pokecache"
)

type PokeClient struct {
	client  *http.Client
	BaseURL string
	cache   *pokecache.Cache
}

type LocationAreasResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
}

type AreaPokemon struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string  `json:"name"`
	BaseExperience int     `json:"base_experience"`
	Weight         int     `json:"weight"`
	Height         int     `json:"height"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

func (c *PokeClient) get(url string) ([]byte, error) {
	cached, exists := c.cache.Get(url)
	if exists {
		log.Println("Cache hit for:", url)
		return cached, nil
	}
	log.Println("Cache miss for:", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching response: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading: %w", err)
	}

	c.cache.Add(url, body)
	return body, nil
}

func NewClient() PokeClient {
	return PokeClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		BaseURL: "https://pokeapi.co/api/v2",
		cache:   pokecache.NewCache(time.Second * 10),
	}
}

func (c *PokeClient) GetLocationAreas(pageURL *string) (*LocationAreasResponse, error) {
	var endpoint string
	if pageURL != nil {
		endpoint = *pageURL
	} else {
		endpoint = c.BaseURL + "/location-area?offset=0&limit=20"
	}
	res, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var locationAreas LocationAreasResponse
	err = json.Unmarshal(res, &locationAreas)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	return &locationAreas, nil
}

func (c *PokeClient) ExploreArea(name string) (*AreaPokemon, error) {
	endpoint := c.BaseURL + "/location-area/" + name
	res, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var areaPokemon AreaPokemon
	err = json.Unmarshal(res, &areaPokemon)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	return &areaPokemon, nil
}

func (c *PokeClient) GetPokemon(name string) (*Pokemon, error) {
	endpoint := c.BaseURL + "/pokemon/" + name
	res, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var pokemon Pokemon
	err = json.Unmarshal(res, &pokemon)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}
	return &pokemon, nil
}
