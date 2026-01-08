package pokeapi

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

// LocationAreaDetailResponse represents detailed info about a location area
type LocationAreaDetailResponse struct {
    ID                int    `json:"id"`
    Name              string `json:"name"`
    PokemonEncounters []struct {
        Pokemon struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"pokemon"`
    } `json:"pokemon_encounters"`
}

// GetLocationAreaDetail fetches details for a specific location area
func (c *Client) GetLocationAreaDetail(areaName string) (LocationAreaDetailResponse, error) {
    url := baseURL + "/location-area/" + areaName

    // Check cache first
    if cachedData, exists := c.cache.Get(url); exists {
        fmt.Println("(using cached data)")

        var locationArea LocationAreaDetailResponse
        err := json.Unmarshal(cachedData, &locationArea)
        if err != nil {
            return LocationAreaDetailResponse{}, fmt.Errorf("error unmarshaling cached JSON: %w", err)
        }

        return locationArea, nil
    }

    fmt.Println("(fetching from API...)")

    // Make GET request
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return LocationAreaDetailResponse{}, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    // Check for non-OK status
    if resp.StatusCode != http.StatusOK {
        return LocationAreaDetailResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return LocationAreaDetailResponse{}, fmt.Errorf("error reading response: %w", err)
    }

    // Add to cache
    c.cache.Add(url, body)

    // Unmarshal JSON into struct
    var locationArea LocationAreaDetailResponse
    err = json.Unmarshal(body, &locationArea)
    if err != nil {
        return LocationAreaDetailResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
    }

    return locationArea, nil
}
