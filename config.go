package main

import "github.com/ankamason/pokedexcli/internal/pokeapi"

// Config holds the state for the REPL
type config struct {
    pokeapiClient *pokeapi.Client
    nextURL       *string
    previousURL   *string
}
