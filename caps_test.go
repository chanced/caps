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

type testcase struct {
	input    string
	expected string
	opts     Opts
}

type testcases []testcase

var usd caps.NumberRules = caps.NumberRules{
	'$': func(index int, r rune, val string) bool {
		return index == 0
	},
}

var titleTestCases = testcases{
	{"", "", nil},
	{"a", "A", nil},
	{"aA", "A A", nil},
	{"aAa", "A Aa", nil},
	{"id_a", "ID A", nil},
	{"id", "ID", nil},
	{"ID 3", "ID 3", nil},
	{"test_from_snake", "Test From Snake", nil},
	{"TestFromCamel", "Test From Camel", nil},
	{"test-with-number-123", "Test With Number 123", nil},
	{"test-with-usd-$123.34", "Test With Usd $123.34", Opts{caps.WithNumberRules(usd), caps.WithAllowedSymbols("$.")}},
	{"test with number -123", "Test With Number -123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "Test With Number -123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "Test With Number -123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms uuid id guid http https", "Test Initialisms UUID ID GUID HTTP HTTPS", nil},
	{"test camel initialisms uuid id guid http https", "Test Camel Initialisms Uuid Id Guid Http Https", Opts{caps.WithReplaceStyleCamel()}},
}

func TestToTitle(t *testing.T) {
	for _, test := range titleTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToTitle(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}

	for _, test := range titleTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToTitle(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
		func(test testcase) {
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToTitle(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var camelTestCases = testcases{
	{"", "", nil},
	{"a", "A", nil},
	{"aA", "AA", nil},
	{"aAa", "AAa", nil},
	{"test_from_snake", "TestFromSnake", nil},
	{"TestFromCamel", "TestFromCamel", nil},
	{"testCamelFromLowerCamel", "TestCamelFromLowerCamel", nil},
	{"test-with-number-123", "TestWithNumber123", nil},
	{"test with number -123", "TestWithNumber-123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "TestWithNumber-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "TestWithNumber-123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ip html eof ascii cpu", "TestInitialismsIPHTMLEOFASCIICPU", nil},
	{"test camel initialisms ip html eof ascii cpu", "TestCamelInitialismsIpHtmlEofAsciiCpu", Opts{caps.WithReplaceStyleCamel()}},
}

func TestToCamel(t *testing.T) {
	for _, test := range camelTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToCamel(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToCamel(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var lowerCamelTestCases = testcases{
	{"", "", nil},
	{"a", "a", nil},
	{"aA", "aA", nil},
	{"aAa", "aAa", nil},
	{"AaA", "aaA", nil},
	{"test_from_snake", "testFromSnake", nil},
	{"TestFromCamel", "testFromCamel", nil},
	{"testFromLowerCamel", "testFromLowerCamel", nil},
	{"test-with-number-123", "testWithNumber123", nil},
	{"test with number -123", "testWithNumber-123", Opts{caps.WithAllowedSymbols("-")}}, // 7
	{"test with number -123.456", "testWithNumber-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "testWithNumber-123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ip html eof ascii cpu", "testInitialismsIPHTMLEOFASCIICPU", nil},
	{"test camel initialisms ip html eof ascii cpu", "testCamelInitialismsIpHtmlEofAsciiCpu", Opts{caps.WithReplaceStyleCamel()}},
}

func TestToLowerCamel(t *testing.T) {
	for _, test := range lowerCamelTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToLowerCamel(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToLowerCamel(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var kebabTestCases = testcases{
	{"", "", nil},
	{"a", "a", nil},
	{"aA", "a-a", nil},
	{"aAa", "a-aa", nil},
	{"AaA", "aa-a", nil},
	{"test_from_snake", "test-from-snake", nil},
	{"TestFromCamel", "test-from-camel", nil},
	{"testFromLowerCamel", "test-from-lower-camel", nil},
	{"test-with-number-123", "test-with-number-123", nil},
	{"test with number -123", "test-with-number--123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "test-with-number--123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "test-with-number--123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "test-initialisms-ram-tcp-ttl-ascii", nil},
}

func TestToKebab(t *testing.T) {
	for _, test := range kebabTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToKebab(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToKebab(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var screamingKebabTestCases = testcases{
	{"", "", nil},
	{"a", "A", nil},
	{"aA", "A-A", nil},
	{"aAa", "A-AA", nil},
	{"AaA", "AA-A", nil},
	{"test_from_snake", "TEST-FROM-SNAKE", nil},
	{"TestFromCamel", "TEST-FROM-CAMEL", nil},
	{"testFromLowerCamel", "TEST-FROM-LOWER-CAMEL", nil},
	{"test-with-number-123", "TEST-WITH-NUMBER-123", nil},
	{"test with number -123", "TEST-WITH-NUMBER--123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "TEST-WITH-NUMBER--123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "TEST-WITH-NUMBER--123.456E-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "TEST-INITIALISMS-RAM-TCP-TTL-ASCII", nil},
}

func TestToScreamingKebab(t *testing.T) {
	for _, test := range screamingKebabTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToScreamingKebab(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToScreamingKebab(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var snakeTestCases = testcases{
	{"", "", nil},
	{"a", "a", nil},
	{"aA", "a_a", nil},
	{"aAa", "a_aa", nil},
	{"AaA", "aa_a", nil},
	{"test_from_snake", "test_from_snake", nil},
	{"TestFromCamel", "test_from_camel", nil},
	{"testFromLowerCamel", "test_from_lower_camel", nil},
	{"test_with_number_123", "test_with_number_123", nil},
	{"test with number -123", "test_with_number_-123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "test_with_number_-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "test_with_number_-123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "test_initialisms_ram_tcp_ttl_ascii", nil},
}

func TestToSnake(t *testing.T) {
	for _, test := range snakeTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToSnake(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToSnake(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var screamingSnakeTestCases = testcases{
	{"", "", nil},
	{"a", "A", nil},
	{"aA", "A_A", nil},
	{"aAa", "A_AA", nil},
	{"AaA", "AA_A", nil},
	{"test_from_snake", "TEST_FROM_SNAKE", nil},
	{"TestFromCamel", "TEST_FROM_CAMEL", nil},
	{"testFromLowerCamel", "TEST_FROM_LOWER_CAMEL", nil},
	{"test_with_number_123", "TEST_WITH_NUMBER_123", nil},
	{"test with number -123", "TEST_WITH_NUMBER_-123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "TEST_WITH_NUMBER_-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "TEST_WITH_NUMBER_-123.456E-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "TEST_INITIALISMS_RAM_TCP_TTL_ASCII", nil},
}

func TestToScreamingSnake(t *testing.T) {
	for _, test := range screamingSnakeTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToScreamingSnake(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToScreamingSnake(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var dotNotationTestCases = testcases{
	{"", "", nil},
	{"a", "a", nil},
	{"aA", "a.a", nil},
	{"aAa", "a.aa", nil},
	{"AaA", "aa.a", nil},
	{"test.from.dotnotation", "test.from.dotnotation", nil},
	{"TestFromCamel", "test.from.camel", nil},
	{"testFromLowerCamel", "test.from.lower.camel", nil},
	{"test.with.number.123", "test.with.number.123", nil},
	{"test with number -123", "test.with.number.-123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "test.with.number.-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "test.with.number.-123.456e-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "test.initialisms.ram.tcp.ttl.ascii", nil},
}

func TestToDotNotation(t *testing.T) {
	for _, test := range dotNotationTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToDotNotation(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToDotNotation(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

var screamingDotNotationTestCases = testcases{
	{"", "", nil},
	{"a", "A", nil},
	{"aA", "A.A", nil},
	{"aAa", "A.AA", nil},
	{"AaA", "AA.A", nil},
	{"test.from.dotnotation", "TEST.FROM.DOTNOTATION", nil},
	{"TestFromCamel", "TEST.FROM.CAMEL", nil},
	{"testFromLowerCamel", "TEST.FROM.LOWER.CAMEL", nil},
	{"test.with.number.123", "TEST.WITH.NUMBER.123", nil},
	{"test with number -123", "TEST.WITH.NUMBER.-123", Opts{caps.WithAllowedSymbols("-")}},
	{"test with number -123.456", "TEST.WITH.NUMBER.-123.456", Opts{caps.WithAllowedSymbols("-.")}},
	{"test with number -123.456e-2", "TEST.WITH.NUMBER.-123.456E-2", Opts{caps.WithAllowedSymbols("-.")}},
	{"test initialisms ram tcp ttl ascii", "TEST.INITIALISMS.RAM.TCP.TTL.ASCII", nil},
}

func TestToDelimited(t *testing.T) {
	for _, test := range dotNotationTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToDelimited(test.input, ".", true, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToDelimited(test.input, ".", true)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

func TestToScreamingDotNotation(t *testing.T) {
	for _, test := range screamingDotNotationTestCases {
		func(test testcase) {
			t.Run(test.input, func(t *testing.T) {
				t.Parallel()
				output := caps.ToScreamingDotNotation(test.input, test.opts...)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				t.Parallel()
				c := caps.New(test.opts.toConfig())
				output := c.ToScreamingDotNotation(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

func TestCapsAccessors(t *testing.T) {
	t.Run("ReplaceStyle", func(t *testing.T) {
		t.Parallel()
		c := caps.New()
		rs := c.ReplaceStyle()
		if rs != caps.ReplaceStyleScreaming {
			t.Errorf("expected %s, got %s", caps.ReplaceStyleScreaming, rs)
		}

		c = caps.New(caps.Config{
			ReplaceStyle: caps.ReplaceStyleLower,
		})

		rs = c.ReplaceStyle()
		if rs != caps.ReplaceStyleLower {
			t.Errorf("expected %s, got %s", caps.ReplaceStyleLower, rs)
		}
	})
	t.Run("NumberRules", func(t *testing.T) {
		t.Parallel()
		c := caps.New()
		nr := c.NumberRules()
		if nr != nil {
			t.Errorf("expected nil, got %v", nr)
		}

		c = caps.New(caps.Config{
			NumberRules: usd,
		})
		if len(c.NumberRules()) != 1 || c.NumberRules()['$'] == nil {
			t.Errorf("expected %v, got %v", usd, c.NumberRules())
		}
	})
	t.Run("AllowedSymbols", func(t *testing.T) {
		t.Parallel()
		c := caps.New()
		as := c.AllowedSymbols()
		if as != "" {
			t.Errorf("expected \"\", got %v", as)
		}
		c = caps.New(caps.Config{
			AllowedSymbols: "$",
		})
		if c.AllowedSymbols() != "$" {
			t.Errorf("expected \"$\", got %v", c.AllowedSymbols())
		}
	})
}

func TestToUpper(t *testing.T) {
	res := caps.ToUpper("test")
	if res != "TEST" {
		t.Errorf("expected \"TEST\", got \"%s\"", res)
	}
}

func TestToLower(t *testing.T) {
	res := caps.ToLower("TEST")
	if res != "test" {
		t.Errorf("expected \"test\", got \"%s\"", res)
	}
}

func TestWithoutNumbers(t *testing.T) {
	type Test struct {
		input    string
		expected string
	}
	tests := []Test{
		{"123.456", "."},
		{"a12bc3d", "abcd"},
		{"a", "a"},
		{"", ""},
	}
	c := caps.New()
	for _, test := range tests {
		func(test Test) {
			t.Run(test.input, func(t *testing.T) {
				output := caps.WithoutNumbers(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				output := c.WithoutNumbers(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

func TestLowerFirst(t *testing.T) {
	type Test struct {
		input    string
		expected string
	}
	tests := []Test{
		{"anExampleString", "anExampleString"},
		{"ANEXAMPLESTRING", "aNEXAMPLESTRING"},
		{"AnExampleString, ", "anExampleString, "},
		{"a", "a"},
		{"", ""},
		{"123", "123"},
	}

	for _, test := range tests {
		func(test Test) {
			t.Run(test.input, func(t *testing.T) {
				output := caps.LowerFirst(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
			t.Run("Caps::"+test.input, func(t *testing.T) {
				output := caps.New().LowerFirst(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

func TestUpperFirst(t *testing.T) {
	type Test struct {
		input    string
		expected string
	}

	tests := []Test{
		{"anExampleString", "AnExampleString"},
		{"ANEXAMPLESTRING", "ANEXAMPLESTRING"},
		{"a", "A"},
		{"", ""},
		{"123", "123"},
	}
	for _, test := range tests {
		func(test Test) {
			t.Run(test.input, func(t *testing.T) {
				output := caps.UpperFirst(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})

			t.Run("Caps::"+test.input, func(t *testing.T) {
				output := caps.New().UpperFirst(test.input)
				if output != test.expected {
					t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
				}
			})
		}(test)
	}
}

func TestFormatToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		style    caps.Style
		index    int
		caser    token.Caser
	}{
		{"test", "Test", caps.StyleCamel, 0, token.DefaultCaser},
		{"test", "test", caps.StyleLowerCamel, 0, token.DefaultCaser},
		{"test", "TEST", caps.StyleScreaming, 0, token.DefaultCaser},
		{"uuid", "Uuid", caps.StyleCamel, 0, token.DefaultCaser},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.FormatToken(test.caser, test.style, test.index, test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

type Opts []caps.Opts

func loadOpts(opts Opts) caps.Opts {
	result := caps.Opts{}
	for _, opt := range opts {
		if opt.AllowedSymbols != "" {
			result.AllowedSymbols = result.AllowedSymbols + opt.AllowedSymbols
		}
		if opt.NumberRules != nil {
			result.NumberRules = opt.NumberRules
		}
		if opt.ReplaceStyle != 0 {
			result.ReplaceStyle = opt.ReplaceStyle
		}
		if opt.Converter != nil {
			result.Converter = opt.Converter
		}
	}
	return result
}

func (o Opts) toConfig() caps.Config {
	capopts := caps.Config{}
	opts := loadOpts(o)
	capopts.AllowedSymbols = opts.AllowedSymbols
	capopts.NumberRules = opts.NumberRules
	capopts.ReplaceStyle = opts.ReplaceStyle
	capopts.Converter = opts.Converter
	return capopts
}
