package pokeapi

import (
    "testing"
    "time"
)

func TestNewClient(t *testing.T) {
    client := NewClient(5 * time.Second)
    if client == nil {
        t.Errorf("expected client to not be nil")
    }
}
