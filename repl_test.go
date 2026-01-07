package main

import (
    "testing"
)

func TestCleanInput(t *testing.T) {
    cases := []struct {
        input    string
        expected []string
    }{
        {
            input:    "  hello  world  ",
            expected: []string{"hello", "world"},
        },
        {
            input:    "Charmander Bulbasaur PIKACHU",
            expected: []string{"charmander", "bulbasaur", "pikachu"},
        },
        {
            input:    "hello",
            expected: []string{"hello"},
        },
        {
            input:    "   UPPERCASE   lowercase   MiXeD   ",
            expected: []string{"uppercase", "lowercase", "mixed"},
        },
        {
            input:    "",
            expected: []string{},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)

        // Check the length of the actual slice against the expected slice
        if len(actual) != len(c.expected) {
            t.Errorf("Test Failed: %s\n  Expected length: %d, Got: %d",
                c.input, len(c.expected), len(actual))
            continue
        }

        // Check each word in the slice
        for i := range actual {
            word := actual[i]
            expectedWord := c.expected[i]
            if word != expectedWord {
                t.Errorf("Test Failed: %s\n  Expected: %s, Got: %s",
                    c.input, expectedWord, word)
            }
        }
    }
}
