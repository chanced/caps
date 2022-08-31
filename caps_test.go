package caps_test

import (
	"testing"

	"github.com/chanced/caps"
)

func TestWithoutNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123.456", "."},
		{"a12bc3d", "abcd"},
		{"a", "a"},
		{"", ""},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.WithoutNumbers(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"anExampleString", "anExampleString"},
		{"ANEXAMPLESTRING", "aNEXAMPLESTRING"},
		{"AnExampleString, ", "anExampleString, "},
		{"a", "a"},
		{"", ""},
		{"123", "123"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.LowerFirst(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"anExampleString", "AnExampleString"},
		{"ANEXAMPLESTRING", "ANEXAMPLESTRING"},
		{"a", "A"},
		{"", ""},
		{"123", "123"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.UpperFirst(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}
