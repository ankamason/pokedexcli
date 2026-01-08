package main

import (
    "bufio"
    "fmt"
    "os"
    "time"

    "github.com/ankamason/pokedexcli/internal/pokeapi"
)

func main() {
    // Initialize the config with PokeAPI client (5 minute cache)
    cfg := &config{
        pokeapiClient: pokeapi.NewClient(5 * time.Minute),
        pokedex:       make(map[string]pokeapi.PokemonResponse),
    }

    // Create a scanner that reads from standard input
    scanner := bufio.NewScanner(os.Stdin)

    // Infinite loop - runs once per command
    for {
        // Print prompt without newline
        fmt.Print("Pokedex > ")

        // Block and wait for user input
        scanner.Scan()

        // Get the input as a string
        input := scanner.Text()

        // Clean the input
        words := cleanInput(input)

        // Check if there's at least one word
        if len(words) == 0 {
            continue
        }

        // Get the command name (first word)
        commandName := words[0]

        // Get arguments (remaining words)
        args := []string{}
        if len(words) > 1 {
            args = words[1:]
        }

        // Look up the command in the registry
        command, exists := getCommands()[commandName]

        if !exists {
            fmt.Println("Unknown command")
            continue
        }

        // Call the callback with config, args and handle any errors
        err := command.callback(cfg, args)
        if err != nil {
            fmt.Println(err)
        }
    }
}
