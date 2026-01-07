package pokeapi

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/ankamason/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// Client for PokeAPI
type Client struct {
    httpClient *http.Client
    cache      *pokecache.Cache
}

// NewClient creates a new PokeAPI client with caching
func NewClient(cacheInterval time.Duration) *Client {
    return &Client{
        httpClient: &http.Client{},
        cache:      pokecache.NewCache(cacheInterval),
    }
}

// LocationAreasResponse represents the response from location-area endpoint
type LocationAreasResponse struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

// GetLocationAreas fetches location areas from the given URL
func (c *Client) GetLocationAreas(url string) (LocationAreasResponse, error) {
    // Use default URL if none provided
    if url == "" {
        url = baseURL + "/location-area"
    }

    // Check cache first
    if cachedData, exists := c.cache.Get(url); exists {
        fmt.Println("(using cached data)")

        var locationAreas LocationAreasResponse
        err := json.Unmarshal(cachedData, &locationAreas)
        if err != nil {
            return LocationAreasResponse{}, fmt.Errorf("error unmarshaling cached JSON: %w", err)
        }

        return locationAreas, nil
    }

    fmt.Println("(fetching from API...)")

    // Make GET request
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return LocationAreasResponse{}, fmt.Errorf("error making request: %w", err)
    }
    defer resp.Body.Close()

    // Check for non-OK status
    if resp.StatusCode != http.StatusOK {
        return LocationAreasResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return LocationAreasResponse{}, fmt.Errorf("error reading response: %w", err)
    }

    // Add to cache
    c.cache.Add(url, body)

    // Unmarshal JSON into struct
    var locationAreas LocationAreasResponse
    err = json.Unmarshal(body, &locationAreas)
    if err != nil {
        return LocationAreasResponse{}, fmt.Errorf("error unmarshaling JSON: %w", err)
    }

    return locationAreas, nil
}
