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
			input:    "  hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  why are we   just to suffer   ",
			expected: []string{"why", "are", "we", "just", "to", "suffer"},
		},
		{
			input:    "charmander BULBASAUR pikaCHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q)	returned %d, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q)[%d] = %q, expected %q", c.input, i, actual[i], c.expected[i])
			}
		}
	}
}
