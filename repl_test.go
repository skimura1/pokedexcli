package main

import (
	"testing"
	"fmt"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input		string
		expected	[]string
	}{
		{
			input: "   hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input: "What is good?",
			expected:[]string{"what", "is", "good?"},
		},
		{
			input: "This is very cool    ",
			expected:[]string{"this", "is", "very", "cool"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		

		if len(actual) != len(c.expected){
			t.Errorf("Test Failed: The the lengths do not match!")
		}
		fmt.Println(actual)

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Test Failed: The words do not match the expected output")
			}
		}
	}
}
