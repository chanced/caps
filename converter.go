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

import (
	"strings"
	"unicode"

	"github.com/chanced/caps/index"
	"github.com/chanced/caps/token"
)

// DefaultConverter is the default Converter instance.
var DefaultConverter = NewConverter(DefaultReplacements, DefaultTokenizer, token.DefaultCaser)

// Converter is an interface satisfied by types which can convert the case of a
// string.
//
// ConverterImpl is provided as a default implementation. If you have edge cases which require custom formatting,
// you can implement your own Converter by wrapping ConverterImpl:
//
//	type MyConverter struct {}
//	func(MyConverter) Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string{
//		formatted := caps.DefaultConverter.Convert(style, repStyle, input, join, allowedSymbols, numberRules)
//		if formatted == "something_unusual" {
//		  	return "replaced"
//	 	}
//	 	return formatted
//	}
//
// # Parameters
//
//	style:          Expected output caps.Style of the string.
//	repStyle:       The caps.ReplaceStyle to use if a word needs to be replaced.
//	join:           The delimiter to use when joining the words. For Camel, this is an empty string.
//	allowedSymbols: The set of allowed symbols. If set, these should take precedence over any delimiters
//	numberRules:    Any custom rules dictating how to handle special characters in numbers.
type Converter interface {
	Convert(req ConvertRequest) string
}

type ConvertRequest struct {
	Style          Style
	ReplaceStyle   ReplaceStyle
	Input          string
	Join           string
	AllowedSymbols string
	NumberRules    map[rune]func(index int, r rune, val string) bool
}

// NewConverter creates a new Converter which is used to convert the input text to the desired output.
//
// replacements are used to make replacements of tokens to the specified
// formatting (e.g. { "Json", "JSON"}).
//
// tokenizer is used to tokenize the input text.
func NewConverter(replacements []Replacement, tokenizer Tokenizer, caser token.Caser) StdConverter {
	sc := StdConverter{
		index:     index.New(caser),
		tokenizer: tokenizer,
		caser:     token.CaserOrDefault(caser),
	}
	for _, v := range replacements {
		sc.set(v.Camel, v.Screaming)
	}
	return sc
}

// StdConverter contains a table of words to their desired replacement. Tokens
// will be compared against the keys of this table to determine if the string
// should be replaced with the value of the table.
//
// This is primarily designed for acronyms but it could be used for other
// purposes.
//
// The default Replacements can be found in the DefaultReplacements variable.
type StdConverter struct {
	index     *index.Index
	tokenizer Tokenizer
	caser     token.Caser
}

func (sc StdConverter) Index() index.Index {
	return *sc.index
}

// Contains reports whether a key is in the Converter's replacement table.
func (sc StdConverter) Contains(key string) bool {
	return sc.index.Contains(key)
}

// Replacements returns a slice of Replacement in the lookup trie.
func (sc StdConverter) Replacements() []Replacement {
	indexedVals := sc.index.Values()
	res := make([]Replacement, len(indexedVals))
	for i, v := range indexedVals {
		res[i] = Replacement{
			Camel:     v.Camel,
			Screaming: v.Screaming,
		}
	}
	b := strings.Builder{}
	b.WriteByte('.')
	return res
}

func (sc *StdConverter) set(key, value string) {
	sc.index.Add(key, value)
}

// Set adds the key/value pair to the table.
func (sc *StdConverter) Set(key, value string) {
	kstr, keyHasLower := lowerAndCheck(key)
	vstr, valueHasLower := lowerAndCheck(value)
	sc.Delete(kstr)
	sc.Delete(vstr)

	// checking to see if we need to swap these.
	if !keyHasLower && valueHasLower {
		sc.set(value, key)
	} else {
		sc.set(key, value)
	}
}

// Remove deletes the key from the map. Either variant is sufficient.
func (sc *StdConverter) Delete(key string) {
	sc.index.Delete(key)
}

func (StdConverter) writeIndexReplacement(b *strings.Builder, style Style, repStyle ReplaceStyle, join string, rep index.IndexedReplacement) {
	if len(join) > 0 && b.Len() > 0 {
		b.WriteString(join)
	}
	switch repStyle {
	case ReplaceStyleCamel:
		if b.Len() == 0 && style == StyleLowerCamel {
			b.WriteString(rep.Lower)
		} else {
			b.WriteString(rep.Camel)
		}
	case ReplaceStyleScreaming:
		b.WriteString(rep.Screaming)
	case ReplaceStyleLower:
		b.WriteString(rep.Lower)
	default:
		b.WriteString(rep.Screaming)
	}
}

func (sc StdConverter) writeToken(b *strings.Builder, style Style, join string, tok string) {
	if len(join) > 0 && b.Len() > 0 {
		b.WriteString(join)
	}
	switch style {
	case StyleCamel:
		token.WriteUpperFirstLowerRest(b, sc.caser, tok)
	case StyleLowerCamel:
		if b.Len() == 0 {
			token.WriteLower(b, sc.caser, tok)
		} else {
			token.WriteUpperFirstLowerRest(b, sc.caser, tok)
		}
	case StyleScreaming:
		token.WriteUpper(b, sc.caser, tok)
	case StyleLower:
		token.WriteLower(b, sc.caser, tok)
	default:
		b.WriteString(tok)
	}
}

func (sc StdConverter) writeReplaceSplit(b *strings.Builder, style Style, join string, s string) {
	switch style {
	case StyleCamel:
		token.WriteSplitUpper(b, sc.caser, join, s)
	case StyleLowerCamel:
		if b.Len() == 0 {
			token.WriteSplitLowerFirstUpperRest(b, sc.caser, join, s)
		} else {
			token.WriteSplitUpper(b, sc.caser, join, s)
		}
	case StyleScreaming:
		token.WriteSplitUpper(b, sc.caser, join, s)
	case StyleLower:
		token.WriteSplitLower(b, sc.caser, join, s)
	default:
		token.WriteSplitUpper(b, sc.caser, join, s)
	}
}

// Convert formats the string with the desired style.
func (sc StdConverter) Convert(req ConvertRequest) string {
	tokens := sc.tokenizer.Tokenize(req.Input, req.AllowedSymbols, req.NumberRules)
	b := strings.Builder{}
	if len(tokens) == 0 {
		return ""
	}
	if len(req.Join) > 0 {
		b.Grow(len(req.Input) + len(req.Join)*(len(tokens)-1))
	} else {
		b.Grow(len(req.Input))
	}
	var ok bool
	var addedAsNumber bool
	idx := sc.Index()
	for i, tok := range tokens {
		switch len(tok) {
		case 0:
			continue
		case 1:
			if idx, ok = idx.Match(tok); !ok {
				if idx.LastMatch().HasValue() {
					// appending the last match
					//  formatIndexedReplacement(req.Style, req.ReplaceStyle, b.Len(), idx.LastMatch()), req.Join
					sc.writeIndexReplacement(&b, req.Style, req.ReplaceStyle, req.Join, idx.LastMatch())
				}
				if idx.HasPartialMatches() {
					// checking to make sure it isn't a number
					if token.IsNumber(token.Append(sc.caser, tok, idx.PartialMatches()), req.NumberRules) {
						b.WriteString(FormatToken(sc.caser, req.Style, b.Len(), token.Append(sc.caser, tok, idx.PartialMatches())))
						addedAsNumber = true
					} else {
						sc.writeReplaceSplit(&b, req.Style, req.Join, idx.PartialMatches())
						addedAsNumber = false
					}
				}
				if !addedAsNumber {
					sc.writeToken(&b, req.Style, req.Join, tok)
				}
				// resetting the index
				idx = sc.Index()
			}
		default:
			if idx.HasMatched() {
				sc.writeIndexReplacement(&b, req.Style, req.ReplaceStyle, req.Join, idx.LastMatch())
			}
			if idx.HasPartialMatches() {
				sc.writeReplaceSplit(&b, req.Style, req.Join, idx.PartialMatches())
			}
			if idx.HasMatched() || idx.HasPartialMatches() {
				// resetting index
				idx = sc.Index()
			}
			if rep, ok := idx.Get(tok); ok {
				sc.writeIndexReplacement(&b, req.Style, req.ReplaceStyle, req.Join, rep)
			} else if isNextTokenNumber(tokens, i) {
				if idx, ok = idx.Match(tok); !ok {
					sc.writeToken(&b, req.Style, req.Join, tok)
					idx = sc.Index()
				}
			} else {
				sc.writeToken(&b, req.Style, req.Join, tok)
			}
		}
	}
	if idx.HasMatched() {
		sc.writeIndexReplacement(&b, req.Style, req.ReplaceStyle, req.Join, idx.LastMatch())
		// parts = append(parts, formatIndexedReplacement(req.Style, req.ReplaceStyle, len(parts), idx.LastMatch()))
	}

	if idx.HasPartialMatches() {
		sc.writeReplaceSplit(&b, req.Style, req.Join, idx.PartialMatches())
	}
	// for _, part := range parts {
	// 	if shouldWriteDelimiter {
	// 		result.WriteString(req.Join)
	// 	}
	// 	result.WriteString(part)
	// 	if !shouldWriteDelimiter {
	// 		shouldWriteDelimiter = len(part) > 0 && len(req.Join) > 0
	// 	}
	// }

	return b.String()
}

// FormatToken formats the str with the desired style.
func FormatToken(caser token.Caser, style Style, index int, tok string) string {
	switch style {
	case StyleCamel:
		return token.UpperFirstLowerRest(caser, tok)
	case StyleLowerCamel:
		if index == 0 {
			return token.ToLower(caser, tok)
		}
		return token.UpperFirstLowerRest(caser, tok)
	case StyleScreaming:
		return token.ToUpper(caser, tok)
	case StyleLower:
		return token.ToLower(caser, tok)
	}
	return tok
}

func isNextTokenNumber(tokens []string, i int) bool {
	if i+1 < len(tokens) {
		if r, ok := token.FirstRune(tokens[i+1]); ok {
			return unicode.IsNumber(r)
		}
	}
	return false
}

func lowerAndCheck(input string) (string, bool) {
	bldr := strings.Builder{}
	bldr.Grow(len(input))
	foundLower := false
	for _, r := range input {
		if !foundLower && unicode.IsLower(r) {
			foundLower = true
		}
		bldr.WriteRune(r)
	}
	return bldr.String(), foundLower
}

// Deprecated: Use StdConverter
type ConverterImpl = StdConverter
