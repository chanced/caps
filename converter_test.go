package caps_test

import (
	"testing"

	"github.com/chanced/caps"
)

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
		{"some_uri8_text", "SomeURI8Text", "", caps.StyleCamel, caps.ReplaceStyleScreaming, nil, nil, nil},
		{"UTF8", "utf8", "_", caps.StyleLower, caps.ReplaceStyleLower, nil, nil, nil},
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
