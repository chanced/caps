/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Chance Dinkins <chanceusc@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, Subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or Substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
		allowedSymbols string
		numberRules    token.NumberRules
	}{
		{"123", []string{"123"}, "", nil},                                                        // 0
		{"aLowerCamelcaseString", []string{"a", "Lower", "Camelcase", "String"}, "", nil},        // 1
		{"A_SCREAMING_SNAKE_STRING", []string{"A", "SCREAMING", "SNAKE", "STRING"}, "", nil},     // 2
		{"A_SCREAMING_SNAKE_STRING", []string{"A_SCREAMING_SNAKE_STRING"}, "_", nil},             // 3
		{"ACamelCaseString", []string{"A", "Camel", "Case", "String"}, "", nil},                  // 4
		{"A_CamelCaseString", []string{"A", "Camel", "Case", "String"}, "", nil},                 // 5
		{"123.456", []string{"123", "456"}, "", nil},                                             // 6
		{"123.456", []string{"123.456"}, ".", nil},                                               // 7
		{"MarshalJSON", []string{"Marshal", "J", "S", "O", "N"}, "", nil},                        // 8
		{"a-kebab-string", []string{"a", "kebab", "string"}, "", nil},                            // 9
		{"a_snake_string", []string{"a", "snake", "string"}, "", nil},                            // 10
		{"a_scientific_n_-123.456e7", []string{"a", "scientific", "n", "-123.456e7"}, "-.", nil}, // 11
		{"my_software_v1.3.3", []string{"my", "software", "v1", "3", "3"}, "", nil},              // 12
		{"my_software_v1.3.3", []string{"my", "software", "v1.3.3"}, ".", nil},                   // 13
		{"#123", []string{"123"}, "", nil},                                                       // 14
		{"#123.456", []string{"123.456"}, ".", nil},                                              // 15
		{"UTF8", []string{"UTF", "8"}, "", nil},                                                  // 16
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
					if token != test.expected[i] {
						t.Errorf("expected token %d to be \"%s\", got \"%s\"", i, test.expected[i], token)
					}
				}
			}
		})
	}
}
