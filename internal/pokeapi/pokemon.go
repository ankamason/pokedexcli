package pokeapi

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// PokemonResponse represents detailed info about a Pokemon
type PokemonResponse struct {
    ID             int    `json:"id"`
    Name           string `json:"name"`
    BaseExperience int    `json:"base_experience"`
    Height         int    `json:"height"`
    Weight         int    `json:"weight"`
    Stats          []struct {
        BaseStat int `json:"base_stat"`
        Stat     struct {
            Name string `json:"name"`
        } `json:"stat"`
    } `json:"stats"`
    Types []struct {
        Type struct {
            Name string `json:"name"`
        } `json:"type"`
    } `json:"types"`
}

// GetPokemon fetches details for a specific Pokemon
func (c *Client) GetPokemon(pokemonName string) (PokemonResponse, error) {
    url := baseURL + "/pokemon/" + pokemonName

    // Check cache first
    if cachedData, exists := c.cache.Get(url); exists {
        var pokemon PokemonResponse
        err := json.Unmarshal(cachedData, &pokemon)
        if err != nil {
            return PokemonResponse{}, fmt.Errorf("error unmarshaling cached JSON: %w", err)
        }

        return pokemon, nil
    }

    // Make GET request
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return PokemonResponse{}, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    // Check for non-OK status
    if resp.StatusCode != http.StatusOK {
        return PokemonResponse{}, fmt.Errorf("pokemon not found: %s", pokemonName)
    }

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return PokemonResponse{}, fmt.Errorf("error reading response: %w", err)
    }

    // Add to cache
    c.cache.Add(url, body)

    // Unmarshal JSON into struct
    var pokemon PokemonResponse
    err = json.Unmarshal(body, &pokemon)
    if err != nil {
        return PokemonResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
    }

    return pokemon, nil
}
