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
	"testing"

	"github.com/chanced/caps"
	"github.com/chanced/caps/token"
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
		numberRules    token.NumberRules
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
			params := caps.ConvertRequest{
				Style:          test.style,
				ReplaceStyle:   test.repStyle,
				Input:          string(test.input),
				Join:           test.join,
				AllowedSymbols: string(test.allowedSymbols),
				NumberRules:    test.numberRules,
			}

			output := converter.Convert(params)
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

	// this just checks to see if we StdConverter.Set will swap incase the user
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

type MyConverter struct{}

func (MyConverter) Convert(req caps.ConvertRequest) string {
	res := caps.DefaultConverter.Convert(req)
	if req.Style == caps.StyleLowerCamel && req.ReplaceStyle == caps.ReplaceStyleCamel && res == "id" {
		return "_id"
	}
	return res
}

func TestCustomConverter(t *testing.T) {
	res := caps.ToLowerCamel("ID", caps.WithReplaceStyleCamel(), caps.WithConverter(MyConverter{}))
	if res != "_id" {
		t.Errorf("expected _id, got %s", res)
	}
}
