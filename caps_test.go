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

func TestConverterConvert(t *testing.T) {
	converter := caps.NewConverter(caps.DefaultReplacements, caps.DefaultTokenizer, nil)

	tests := []struct {
		input          string
		expected       string
		join           string
		style          caps.Style
		repStyle       caps.ReplaceStyle
		allowedSymbols []rune
		converter      caps.Converter
		numberRules    map[rune]func(index int, r rune, val []rune) bool
	}{
		{"An example string", "AnExampleString", "", caps.StyleCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"An example string", "anExampleString", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"aCamelCaseExample", "ACamelCaseExample", "", caps.StyleCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"serveHttp", "ServeHTTP", "", caps.StyleCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"A_SCREAMING_SNAKECASE_STRING", "aScreamingSnakecaseString", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"#12.34", "12.34", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, []rune{'.'}, nil, nil},
		{"MarshalJSON", "marshalJSON", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"MarshalJSON", "marshalJson", "", caps.StyleLowerCamel, caps.ReplaceStyleCamel, nil, nil, nil},
		{"marshal_json", "marshalJSON", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"MarshalJSON", "marshal_json", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
		{"aABC", "aABC", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"MarshalJS", "marshal_j_s", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
		{"An[example]_split#with(other).symbols", "anExampleSplitWithOtherSymbols", "", caps.StyleLowerCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"SomeUUID", "some_uuid", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
		{"SomeUID", "some_uid", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
		{"FULLURI", "fulluri", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
		{"FullURI", "full_uri", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := converter.Convert(test.style, test.repStyle, test.input, test.join, test.allowedSymbols, test.numberRules)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

func TestConverterTableOps(t *testing.T) {
	c := caps.NewConverter(caps.DefaultReplacements, caps.DefaultTokenizer, nil)

	hasCamel := false
	for _, v := range caps.DefaultReplacements {
		if v.Camel == "Http" {
			hasCamel = true
			break
		}
	}
	if !hasCamel {
		t.Error("expected \"Http\" to be in the default replacements")
	}
	if !c.Contains("http") {
		t.Errorf("expected caps.ConverterImpl to contain \"http\"")
	}

	c.Delete("http")
	if c.Contains("http") || c.Contains("Http") || c.Contains("HTTP") {
		t.Errorf("expected caps.ConverterImpl to have removed \"http\"")
	}

	c.Set("Tcp", "TCP")
	if !(c.Contains("Tcp") && c.Contains("TCP") && c.Contains("tcp")) {
		t.Errorf("expected caps.ConverterImpl to contain \"http\"")
	}

	// tcp := c.Lookup("tcp")
	// if tcp.Camel != "Tcp" {
	// 	t.Errorf("expected \"Tcp\", got \"%s\"", tcp.Camel)
	// }
	// if tcp.Screaming != "TCP" {
	// 	t.Errorf("expected \"TCP\", got \"%s\"", tcp.Screaming)
	// }

	// this just checks to see if we StdReplacer.Set will swap incase the user
	// flips the order of the strings

	c.Set("WSS", "Wss")
	// wss := c.Lookup("wss")
	// if wss.Camel != "Wss" {
	// 	t.Errorf("expected \"Wss\", got \"%s\"", wss.Camel)
	// }
	// if wss.Screaming != "WSS" {
	// 	t.Errorf("expected \"WSS\", got \"%s\"", wss.Screaming)
	// }
}

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
