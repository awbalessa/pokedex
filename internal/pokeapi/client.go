package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type PokeClient struct {
	client  *http.Client
	BaseURL string
}

type LocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (c *PokeClient) get(url string) ([]byte, error) {
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

	return body, nil
}

func NewClient() PokeClient {
	return PokeClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		BaseURL: "https://pokeapi.co/api/v2",
	}
}

func (c *PokeClient) GetLocationAreas(pageURL *string) (*LocationAreasResponse, error) {
	var endpoint string
	if pageURL != nil {
		endpoint = *pageURL
	} else {
		endpoint = c.BaseURL + "/location-area"
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
