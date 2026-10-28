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
			input:    "	hello world		",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "diglet",
			expected: []string{"diglet"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice
		// if they do not match, use t.Errorf and continue to the next case
		if len(actual) != len(c.expected) {
			t.Errorf("The length of the expected and actual slice do not match")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check if each word in the slice
			if word != expectedWord {
				// if they do not match, use t.Errorf to print and error message
				// and fail the test
				t.Errorf("The expected word and actual word do not match")
			}

		}
	}
}
