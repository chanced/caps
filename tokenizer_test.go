package caps_test

import (
	"fmt"
	"testing"

	"github.com/chanced/caps"
	"github.com/chanced/caps/token"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		value          string
		expected       []string
		allowedSymbols []rune
		numberRules    map[rune]func(index int, r rune, val []rune) bool
	}{
		{"123", []string{"123"}, nil, nil},                                                                   // 0
		{"aLowerCamelcaseString", []string{"a", "Lower", "Camelcase", "String"}, nil, nil},                   // 1
		{"A_SCREAMING_SNAKE_STRING", []string{"A", "SCREAMING", "SNAKE", "STRING"}, nil, nil},                // 2
		{"A_SCREAMING_SNAKE_STRING", []string{"A_SCREAMING_SNAKE_STRING"}, []rune{'_'}, nil},                 // 3
		{"ACamelCaseString", []string{"A", "Camel", "Case", "String"}, nil, nil},                             // 4
		{"A_CamelCaseString", []string{"A", "Camel", "Case", "String"}, nil, nil},                            // 5
		{"123.456", []string{"123", "456"}, nil, nil},                                                        // 6
		{"123.456", []string{"123.456"}, []rune{'.'}, nil},                                                   // 7
		{"MarshalJSON", []string{"Marshal", "J", "S", "O", "N"}, nil, nil},                                   // 8
		{"a-kebab-string", []string{"a", "kebab", "string"}, nil, nil},                                       // 9
		{"a_snake_string", []string{"a", "snake", "string"}, nil, nil},                                       // 10
		{"a_scientific_n_-123.456e7", []string{"a", "scientific", "n", "-123.456e7"}, []rune{'-', '.'}, nil}, // 11
		{"my_software_v1.3.3", []string{"my", "software", "v1", "3", "3"}, nil, nil},                         // 12
		{"my_software_v1.3.3", []string{"my", "software", "v1.3.3"}, []rune{'.'}, nil},                       // 13
		{"#123", []string{"123"}, nil, nil},                                                                  // 14
		{"#123.456", []string{"123.456"}, []rune{'.'}, nil},                                                  // 15
		{"UTF8", []string{"UTF", "8"}, nil, nil},                                                             // 16
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d___%s", i, test.value), func(t *testing.T) {
			tokenizer := caps.NewTokenizer(caps.DEFAULT_DELIMITERS, token.DefaultCaser)
			tokens := tokenizer.Tokenize(test.value, test.allowedSymbols, test.numberRules)
			if len(tokens) != len(test.expected) {
				t.Logf("expected: %+v, got: %+v", test.expected, tokens)
				t.Errorf("expected %d tokens, got %d", len(test.expected), len(tokens))
			} else {
				for i, token := range tokens {
					if token.String() != test.expected[i] {
						t.Errorf("expected token %d to be \"%s\", got \"%s\"", i, test.expected[i], token.String())
					}
				}
			}
		})
	}
}