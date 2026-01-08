package main

import (
    "fmt"
    "os"
    "strings"
)

// Struct that describes a command
type cliCommand struct {
    name        string
    description string
    callback    func(*config, []string) error
}

// Returns the registry of available commands
func getCommands() map[string]cliCommand {
    return map[string]cliCommand{
        "help": {
            name:        "help",
            description: "Displays a help message",
            callback:    commandHelp,
        },
        "exit": {
            name:        "exit",
            description: "Exit the Pokedex",
            callback:    commandExit,
        },
        "map": {
            name:        "map",
            description: "Displays the next 20 location areas",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays the previous 20 location areas",
            callback:    commandMapb,
        },
        "explore": {
            name:        "explore",
            description: "Explore a location area for Pokemon",
            callback:    commandExplore,
        },
    }
}

// Callback for the help command
func commandHelp(cfg *config, args []string) error {
    fmt.Println()
    fmt.Println("Welcome to the Pokedex!")
    fmt.Println("Usage:")
    fmt.Println()

    for _, cmd := range getCommands() {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }

    fmt.Println()
    return nil
}

// Callback for the exit command
func commandExit(cfg *config, args []string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

// Callback for the map command
func commandMap(cfg *config, args []string) error {
    // Determine URL to use
    url := ""
    if cfg.nextURL != nil {
        url = *cfg.nextURL
    }

    // Fetch location areas
    locationAreas, err := cfg.pokeapiClient.GetLocationAreas(url)
    if err != nil {
        return err
    }

    // Update config with new pagination URLs
    cfg.nextURL = locationAreas.Next
    cfg.previousURL = locationAreas.Previous

    // Print location area names
    for _, area := range locationAreas.Results {
        fmt.Println(area.Name)
    }

    return nil
}

// Callback for the mapb command
func commandMapb(cfg *config, args []string) error {
    // Check if we're on the first page
    if cfg.previousURL == nil {
        fmt.Println("you're on the first page")
        return nil
    }

    // Fetch location areas using previous URL
    locationAreas, err := cfg.pokeapiClient.GetLocationAreas(*cfg.previousURL)
    if err != nil {
        return err
    }

    // Update config with new pagination URLs
    cfg.nextURL = locationAreas.Next
    cfg.previousURL = locationAreas.Previous

    // Print location area names
    for _, area := range locationAreas.Results {
        fmt.Println(area.Name)
    }

    return nil
}

// Callback for the explore command
func commandExplore(cfg *config, args []string) error {
    // Check if area name was provided
    if len(args) < 1 {
        return fmt.Errorf("please provide a location area name")
    }

    areaName := args[0]
    fmt.Printf("Exploring %s...\n", areaName)

    // Fetch location area details
    locationArea, err := cfg.pokeapiClient.GetLocationAreaDetail(areaName)
    if err != nil {
        return err
    }

    // Print found Pokemon
    fmt.Println("Found Pokemon:")
    for _, encounter := range locationArea.PokemonEncounters {
        fmt.Printf(" - %s\n", encounter.Pokemon.Name)
    }

    return nil
}

// Clean user input
func cleanInput(text string) []string {
    lowered := strings.ToLower(text)
    words := strings.Fields(lowered)
    return words
}
