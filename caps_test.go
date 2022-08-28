package caps_test

import (
	"testing"

	"github.com/chanced/caps"
)

func TestTokenizer(t *testing.T) {
	// tokenizer := caps.NewTokenizer(caps.DEFAULT_DELIMITERS)
	// tokens := tokenizer.Tokenize("a_scientific_number_-123.456e7", []rune{'-', '.'}, nil)
	// expected := []string{"a", "scientific", "number", "-123.456e7"}
	// if len(tokens) != len(expected) {
	// 	t.Logf("expected: %+v, got: %+v", expected, tokens)
	// 	t.Errorf("expected %d tokens, got %d", len(expected), len(tokens))
	// } else {
	// 	for i, token := range tokens {
	// 		if token.String() != expected[i] {
	// 			t.Errorf("expected token %d to be \"%s\", got \"%s\"", i, expected[i], token.String())
	// 		}
	// 	}
	// }

	tests := []struct {
		value          string
		expected       []string
		allowedSymbols []rune
		numberRules    map[rune]func(index int, r rune, val []rune) bool
	}{
		{"123", []string{"123"}, nil, nil},
		{"aLowerCamelcaseString", []string{"a", "Lower", "Camelcase", "String"}, nil, nil},
		{"A_SCREAMING_SNAKE_STRING", []string{"A", "SCREAMING", "SNAKE", "STRING"}, nil, nil},
		{"A_SCREAMING_SNAKE_STRING", []string{"A_SCREAMING_SNAKE_STRING"}, []rune{'_'}, nil},
		{"ACamelCaseString", []string{"A", "Camel", "Case", "String"}, nil, nil},
		{"A_CamelCaseString", []string{"A", "Camel", "Case", "String"}, nil, nil},
		{"123.456", []string{"123", "456"}, nil, nil},
		{"123.456", []string{"123.456"}, []rune{'.'}, nil},
		{"MarshalJSON", []string{"Marshal", "J", "S", "O", "N"}, nil, nil},
		{"a-kebab-string", []string{"a", "kebab", "string"}, nil, nil},
		{"a_snake_string", []string{"a", "snake", "string"}, nil, nil},
		{"a_scientific_number_-123.456e7", []string{"a", "scientific", "number", "-123.456e7"}, []rune{'-', '.'}, nil},
		{"my_software_v1.3.3", []string{"my", "software", "v1", "3", "3"}, nil, nil},
		{"my_software_v1.3.3", []string{"my", "software", "v1.3.3"}, []rune{'.'}, nil},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			tokenizer := caps.NewTokenizer(caps.DEFAULT_DELIMITERS)
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

// func TestReplacer(t *testing.T) {
// 	assert := require.New(t)
// }

// func TestToCamel(t *testing.T) {
// 	assert := require.New(t)

// 	runs := []struct {
// 		input    string
// 		expected string
// 		opts     *caps.Opts
// 	}{
// 		{
// 			"This is [an] {example}${id32}.",
// 			"thisIsAnExampleID32",
// 			nil,
// 		},
// 		{
// 			"This is [an] {example}${id32}.break32",
// 			"thisIsAnExampleID32Break32",
// 			nil,
// 		},
// 		{
// 			"This is an_example_with !Custom Replacements: http And Https",
// 			"thisIsAnExampleWithCustomReplacementsHTTPAndHTTPS",
// 			&caps.Opts{
// 				Formatter: caps.NewFormatter([]caps.R{
// 					{"Http", "HTTP"},
// 				}),
// 			},
// 		},
// 		{"$word", "$word", &caps.Opts{AllowedSymbols: "$"}},
// 	}

// 	for _, run := range runs {
// 		var opts []caps.Opts
// 		if run.opts != nil {
// 			opts = []caps.Opts{*run.opts}
// 		}
// 		actual := caps.ToCamel(run.input, opts...)
// 		assert.Equal(run.expected, actual)
// 	}
// }

// func TestToDelimited(t *testing.T) {
// 	assert := require.New(t)

// 	runs := []struct {
// 		input     string
// 		expected  string
// 		delimiter rune
// 		opts      *caps.Opts
// 	}{
// 		{
// 			"This is [an] {example}${id32}.",
// 			"this.is.an.example.id.32",
// 			'.',
// 			nil,
// 		},
// 		{
// 			"This is [an] {example}${id32}.break v32",
// 			"this.is.an.example.id.32.break.v32",
// 			'.',
// 			nil,
// 		},
// 		{
// 			"This is an_example_with !Custom Replacements: http And Https",
// 			"this.is.an.example.with.custom.replacements.http.and.https",
// 			'.',
// 			&caps.Opts{
// 				Formatter: caps.NewFormatter([]caps.R{
// 					{"Http", "HTTP"},
// 					{"Https", "HTTPS"},
// 				}),
// 			},
// 		},
// 		{"$word", "$word", '.', &caps.Opts{AllowedSymbols: "$"}},
// 		{"using a different delimiter", "using_a_different_delimiter", '_', nil},
// 	}

// 	for _, run := range runs {
// 		var opts []caps.Opts
// 		if run.opts != nil {
// 			opts = []caps.Opts{*run.opts}
// 		}
// 		lower := caps.ToDelimited(run.input, run.delimiter, true, opts...)
// 		assert.Equal(run.expected, lower)
// 		upper := caps.ToDelimited(run.input, run.delimiter, false, opts...)

// 		assert.Equal(strings.ToUpper(run.expected), upper)

// 	}
// }
