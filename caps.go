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

// Package caps is a unicode aware case conversion library.
package caps

import (
	"strings"
	"unicode"

	"github.com/chanced/caps/token"
)

// Caps provides configured case conversion methods.
type Caps struct {
	caser          token.Caser
	allowedSymbols string
	converter      Converter
	replaceStyle   ReplaceStyle
	numberRules    token.NumberRules
}

// New returns a new Caps instance with the provided options.
//
// if caser is nil, token.DefaultCaser is used (which relies on the default
// unicode functions)
func New(options ...Config) Caps {
	opts := loadConfig(options)

	return Caps{
		caser:          opts.Caser,
		allowedSymbols: opts.AllowedSymbols,
		converter:      opts.Converter,
		replaceStyle:   opts.ReplaceStyle,
		numberRules:    opts.NumberRules,
	}
}

// ReplaceStyle returns the configured ReplaceStyle of c
func (c Caps) ReplaceStyle() ReplaceStyle {
	return c.replaceStyle
}

// NumberRules returns the configured NumberRules of c
func (c Caps) NumberRules() token.NumberRules {
	return c.numberRules
}

// AllowedSymbols returns the configured AllowedSymbols of c
func (c Caps) AllowedSymbols() string {
	return c.allowedSymbols
}

// Converter returns the provided Converter of c
func (c Caps) Converter() Converter {
	return c.converter
}

// UpperFirst converts the first rune of str to unicode upper case.
//
// This method does not support special cases (such as Turkish and Azeri)
func (c Caps) UpperFirst(str string) string {
	return token.UpperFirst(c.caser, str)
}

// LowerFirst converts the first rune of str to lowercase.
func (c Caps) LowerFirst(str string) string {
	return token.LowerFirst(c.caser, str)
}

// Without numbers returns the string with all numeric runes removed.
//
// It does not currently use any logic to determine if a rune (e.g. ".")
// is part of a number. This may change in the future.
func (c Caps) WithoutNumbers(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, string(s))
}

// ToCamel transforms the case of str into Camel Case (e.g. AnExampleString) using
// either the provided Converter or the DefaultConverter otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes "AnExampleString".
// It also has a configurable set of replacements, such that "some_json" becomes "SomeJSON"
// so long as opts.ReplacementStyle is set to ReplaceStyleScreaming. A ReplaceStyle of
// ReplaceStyleCamel would result in "SomeJson".
//
//	caps.ToCamel("This is [an] {example}${id32}.") // ThisIsAnExampleID32
//	caps.ToCamel("AN_EXAMPLE_STRING", ) // AnExampleString
func (c Caps) ToCamel(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           "",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}

// ToLowerCamel transforms the case of str into Lower Camel Case (e.g. anExampleString) using
// either the provided Converter or the DefaultConverter otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes "anExampleString".
// It also has a configurable set of replacements, such that "some_json" becomes "someJSON"
// so long as opts.ReplacementStyle is set to ReplaceStyleScreaming. A ReplaceStyle of
// ReplaceStyleCamel would result in "someJson".
//
//	caps.ToLowerCamel("This is [an] {example}${id32}.") // thisIsAnExampleID32
func (c Caps) ToLowerCamel(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleLowerCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           "",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}

// ToSnake transforms the case of str into Lower Snake Case (e.g. an_example_string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToSnake("This is [an] {example}${id32}.") // this_is_an_example_id_32
func (c Caps) ToSnake(str string) string {
	return c.ToDelimited(str, "_", true)
}

// ToScreamingSnake transforms the case of str into Screaming Snake Case (e.g.
// AN_EXAMPLE_STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingSnake("This is [an] {example}${id32}.") // THIS_IS_AN_EXAMPLE_ID_32
func (c Caps) ToScreamingSnake(str string) string {
	return c.ToDelimited(str, "_", false)
}

// ToKebab transforms the case of str into Lower Kebab Case (e.g. an-example-string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToKebab("This is [an] {example}${id32}.") // this-is-an-example-id-32
func (c Caps) ToKebab(str string) string {
	return c.ToDelimited(str, "-", true)
}

// ToScreamingKebab transforms the case of str into Screaming Kebab Snake (e.g.
// AN-EXAMPLE-STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingKebab("This is [an] {example}${id32}.") // THIS-IS-AN-EXAMPLE-ID-32
func (c Caps) ToScreamingKebab(str string) string {
	return c.ToDelimited(str, "-", false)
}

// ToDotNotation transforms the case of str into Lower Dot Notation Case (e.g. an.example.string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToDotNotation("This is [an] {example}${id32}.") // this.is.an.example.id.32
func (c Caps) ToDotNotation(str string) string {
	return c.ToDelimited(str, ".", true)
}

// ToScreamingDotNotation transforms the case of str into Screaming Kebab Case (e.g.
// AN.EXAMPLE.STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingDotNotation("This is [an] {example}${id32}.") // THIS.IS.AN.EXAMPLE.ID.32
func (c Caps) ToScreamingDotNotation(str string) string {
	return c.ToDelimited(str, ".", false)
}

// ToTitle transforms the case of str into Title Case (e.g. An Example String) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToTitle("This is [an] {example}${id32}.") // This Is An Example ID 32
func (c Caps) ToTitle(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           " ",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}

// ToDelimited transforms the case of str into a string separated by delimiter,
// using either the provided Converter or the DefaultConverter otherwise.
//
// If lowercase is false, the output will be all uppercase.
//
// # See Options for more information on available configuration
//
// # Example
//
//	caps.ToDelimited("This is [an] {example}${id}.#32", ".", true) // this.is.an.example.id.32
//	caps.ToDelimited("This is [an] {example}${id}.break32", ".", false) // THIS.IS.AN.EXAMPLE.ID.BREAK.32
//	caps.ToDelimited("This is [an] {example}${id}.v32", ".", true, caps.Opts{AllowedSymbols: "$"}) // this.is.an.example.$.id.v32
func (c Caps) ToDelimited(str string, delimiter string, lowercase bool) string {
	var style Style
	var replacementStyle ReplaceStyle
	if lowercase {
		style = StyleLower
		replacementStyle = ReplaceStyleLower
	} else {
		style = StyleScreaming
		replacementStyle = ReplaceStyleScreaming
	}
	return c.converter.Convert(ConvertRequest{
		Style:          style,
		ReplaceStyle:   replacementStyle,
		Input:          str,
		Join:           delimiter,
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}
