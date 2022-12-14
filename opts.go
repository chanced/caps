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

package caps

import "github.com/chanced/caps/token"

// ReplaceStyle is used to indicate the case style that text should be
// transformed to when performing replacement in a Converter.
//
// For example, a call to ToCamel with a ReplaceStyleScreaming with an input of
// "MarshalJson" will return "MarshaalJSON" while ReplaceStyleCamel would return
// "MarshalJson"
type ReplaceStyle uint8

const (
	ReplaceStyleNotSpecified ReplaceStyle = iota
	ReplaceStyleCamel                     // Text should be replaced with the Camel variant (e.g. "Json").
	ReplaceStyleScreaming                 // Text should be replaced with the screaming variant (e.g. "JSON").
	ReplaceStyleLower                     // Text should be replaced with the lowercase variant (e.g. "json").
)

func (rs ReplaceStyle) String() string {
	switch rs {
	case ReplaceStyleCamel:
		return "ReplaceStyleCamel"
	case ReplaceStyleScreaming:
		return "ReplaceStyleScreaming"
	case ReplaceStyleLower:
		return "ReplaceStyleLower"
	}
	return "ReplaceStyleNotSpecified"
}

// IsCamel returns true if rs equals ReplaceStyleCamel
func (rs ReplaceStyle) IsCamel() bool {
	return rs == ReplaceStyleCamel
}

// IsScreaming returns true if rs equals ReplaceStyleScreaming
func (rs ReplaceStyle) IsScreaming() bool {
	return rs == ReplaceStyleScreaming
}

// IsLower returns true if rs equals ReplaceStyleLower
func (rs ReplaceStyle) IsLower() bool {
	return rs == ReplaceStyleLower
}

type Style uint8

const (
	StyleNotSpecified Style = iota
	StyleLower              // The output should be lowercase (e.g. "an_example")
	StyleScreaming          // The output should be screaming (e.g. "AN_EXAMPLE")
	StyleCamel              // The output should be camel case (e.g. "AnExample")
	StyleLowerCamel         // The output should be lower camel case (e.g. "anExample")
)

func (s Style) String() string {
	switch s {
	case StyleLower:
		return "StyleLower"
	case StyleScreaming:
		return "StyleScreaming"
	case StyleCamel:
		return "StyleCamel"
	case StyleLowerCamel:
		return "StyleLowerCamel"
	}
	return "StyleNotSpecified"
}

func (s Style) IsLower() bool {
	return s == StyleLower
}

func (s Style) IsScreaming() bool {
	return s == StyleScreaming
}

func (s Style) IsCamel() bool {
	return s == StyleCamel
}

func (s Style) IsLowerCamel() bool {
	return s == StyleLowerCamel
}

// Opts include configurable options for case conversion.
//
// See the documentation for the individual fields for more information.
type Opts struct {
	// Any characters within this string will be allowed in the output.
	//
	// This does not affect delimiters (e.g. '_', '-', '.') as they are added
	// post-tokenization.
	//
	// Default:
	//  ""
	AllowedSymbols string
	// The Converter to use.
	//
	// Default:
	// 	DefaultConverter
	Converter Converter

	// ReplaceStyle overwrites the way Replacements are cased.
	//
	// A typical call to ToLowerCamel for "ServeJSON" with a Converter that
	// contains {"Json": "JSON"} would result in "serveJSON" by using the
	// ReplaceStyleScreaming variant. If ReplacementStyle was set to
	// ReplaceStyleCamel then a call to ToLowerCamel then the result would be
	// "serveHttp".
	//
	// The default ReplaceStyle is dependent upon the target Style.
	ReplaceStyle ReplaceStyle
	// NumberRules are used by the DefaultTokenizer to augment the standard
	// rules for determining if a rune is part of a number.
	//
	// Note, if you add special characters here, they must be present in the
	// AllowedSymbols string for them to be part of the output.
	NumberRules token.NumberRules
}

// WithConverter sets the Converter to use
func WithConverter(converter Converter) Opts {
	return Opts{
		Converter: converter,
	}
}

// WithReplaceStyle sets the ReplaceStyle to use
//
// There are also methods for each ReplaceStyle (e.g. WithReplaceStyleCamel)
func WithReplaceStyle(style ReplaceStyle) Opts {
	return Opts{
		ReplaceStyle: style,
	}
}

// WithReplaceStyleCamel indicates Replacements should use the Camel variant
// (e.g. "Json").
func WithReplaceStyleCamel() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleCamel,
	}
}

// WithReplaceStyleScreaming indicates Replacements should use the screaming
// variant (e.g. "JSON").
func WithReplaceStyleScreaming() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleScreaming,
	}
}

// WithReplaceStyleLower indicates Replacements should use the lowercase variant
// (e.g. "json").
func WithReplaceStyleLower() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleLower,
	}
}

// WithNumberRules sets the NumberRules to use
func WithNumberRules(rules NumberRules) Opts {
	return Opts{
		NumberRules: rules,
	}
}

// WithAllowedSymbols sets the AllowedSymbols to use
func WithAllowedSymbols(symbols string) Opts {
	return Opts{
		AllowedSymbols: symbols,
	}
}

func loadOpts(opts []Opts) Opts {
	result := Opts{
		AllowedSymbols: "",
		Converter:      DefaultConverter,
		ReplaceStyle:   ReplaceStyleScreaming,
	}
	for _, opt := range opts {
		result.AllowedSymbols += opt.AllowedSymbols
		if opt.Converter != nil {
			result.Converter = opt.Converter
		}
		if opt.ReplaceStyle != ReplaceStyleNotSpecified {
			result.ReplaceStyle = opt.ReplaceStyle
		}
		if len(opt.NumberRules) > 0 {
			if result.NumberRules == nil {
				result.NumberRules = make(NumberRules)
			}
			for k, v := range opt.NumberRules {
				result.NumberRules[k] = v
			}
		}
	}
	return result
}

// Config include configurable options for Caps instances.
//
// See the documentation for the individual fields for more information.
type Config struct {
	// Any characters within this string will be allowed in the output.
	//
	// This does not affect delimiters (e.g. "_", "-", ".") as they are added
	// post-tokenization.
	//
	// Default:
	//  ""
	AllowedSymbols string
	// The Converter to use.
	//
	// Default:
	// 	A StdConverter with the Replacements, Caser, and Tokenizer.
	Converter Converter

	// If not set, this will be DefaultReplacements.
	Replacements []Replacement

	// ReplaceStyle overwrites the way words are replaced.
	//
	// A typical call to ToLowerCamel for "ServeJSON" with a Converter that
	// contains {"Json": "JSON"} would result in "serveJSON" by using the
	// ReplaceStyleScreaming variant. If ReplacementStyle was set to
	// ReplaceStyleCamel, on the call to ToLowerCamel then the result would
	// be "serveHttp".
	//
	// The default replacement style is dependent upon the target Style.
	ReplaceStyle ReplaceStyle
	// NumberRules are used by the DefaultTokenizer to augment the standard
	// rules for determining if a rune is part of a number.
	//
	// Note, if you add special characters here, they must be present in the
	// AllowedSymbols string for them to be part of the output.
	NumberRules token.NumberRules
	// Special unicode case rules.
	// See unicode.SpecialCase or token.Caser for more information.
	//
	// Default: token.DefaultCaser (which relies on the default unicode
	// functions)
	Caser token.Caser

	// If not set, uses StdTokenizer with the provided delimiters and token.Caser.
	Tokenizer Tokenizer
}

func loadConfig(opts []Config) Config {
	result := Config{
		AllowedSymbols: "",
		ReplaceStyle:   ReplaceStyleScreaming,
	}

	for _, opt := range opts {
		result.AllowedSymbols += opt.AllowedSymbols
		if opt.Converter != nil {
			result.Converter = opt.Converter
		}
		if opt.ReplaceStyle != ReplaceStyleNotSpecified {
			result.ReplaceStyle = opt.ReplaceStyle
		}
		if len(opt.NumberRules) > 0 {
			if result.NumberRules == nil {
				result.NumberRules = make(NumberRules)
			}
			for k, v := range opt.NumberRules {
				result.NumberRules[k] = v
			}
		}
		if opt.Replacements != nil {
			result.Replacements = append(result.Replacements, opt.Replacements...)
		}
		if opt.Caser != nil {
			result.Caser = opt.Caser
		}
		if opt.Tokenizer != nil {
			result.Tokenizer = opt.Tokenizer
		}
	}
	if result.Caser == nil {
		result.Caser = token.DefaultCaser
	}
	if result.Replacements == nil {
		result.Replacements = DefaultReplacements
	}
	if result.Tokenizer == nil {
		result.Tokenizer = NewTokenizer(DEFAULT_DELIMITERS, result.Caser)
	}
	if result.Converter == nil {
		result.Converter = NewConverter(result.Replacements, result.Tokenizer, result.Caser)
	}

	return result
}
