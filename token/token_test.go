package token_test

import (
	"testing"

	"github.com/chanced/caps/token"
)

func TestIsNumber(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
		rules    map[rune]func(index int, r rune, val []rune) bool
	}{
		{"123", true, nil},
		{"123.456", true, nil},
		{"123.456e7", true, nil},
		{"123.456e-7", true, nil},
		{"#123", true, nil},
		{"#123.456", true, nil},
		{"#123.456e7", true, nil},
		{"#123.456e-7", true, nil},
		{"123#", false, nil},
		{"0.1", true, nil},
		{"0.1e2", true, nil},
		{"0.1e-2", true, nil},
		{"0.1e+2", true, nil},
		{".0", true, nil},
		{".0e2", true, nil},
		{"0.1e2e3", false, nil},
		{".", false, nil},
		{".0.", false, nil},
		{"0..0", false, nil},
		{"123.", false, nil},
		{"123.456.", false, nil},
		{"123.456e7.", false, nil},
		{"v1", true, nil},
		{"v1.0", true, nil},
		{"v1.0.0", true, nil},
		{"v1.0.0-alpha", false, nil},
		{"$1", false, nil},
		{"v1.0.0-alpha.1", false, nil},
		{"vv1", false, nil},
		{"$1", true, map[rune]func(index int, r rune, val []rune) bool{
			'$': func(index int, _ rune, _ []rune) bool {
				return index == 0
			},
		}},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if token.NewFromString(test.value).IsNumber(test.rules) != test.expected {
				if test.expected {
					t.Errorf("expected \"%s\" to be a number", test.value)
				} else {
					t.Errorf("expected \"%s\" to not be a number", test.value)
				}
			}
		})
	}
}
