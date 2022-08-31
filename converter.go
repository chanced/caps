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
	Convert(
		style Style,
		repStyle ReplaceStyle,
		input string,
		join string,
		allowedSymbols []rune,
		numberRules map[rune]func(index int, r rune, val []rune) bool,
	) string
}

// NewConverter creates a new Converter which is used to convert the input text to the desired output.
//
// replacements are used to make replacements of tokens to the specified
// formatting (e.g. { "Json", "JSON"}).
//
// tokenizer is used to tokenize the input text.
func NewConverter(replacements []Replacement, tokenizer Tokenizer, caser token.Caser) ConverterImpl {
	ci := ConverterImpl{
		index:     index.New(caser),
		tokenizer: tokenizer,
		caser:     token.CaserOrDefault(caser),
	}
	for _, v := range replacements {
		ci.set(v.Camel, v.Screaming)
	}
	return ci
}

// ConverterImpl contains a table of words to their desired replacement. Tokens
// will be compared against the keys of this table to determine if the string
// should be replaced with the value of the table.
//
// This is primarily designed for acronyms but it could be used for other
// purposes.
//
// The default Replacements can be found in the DefaultReplacements variable.
type ConverterImpl struct {
	index     *index.Index
	tokenizer Tokenizer
	caser     token.Caser
}

func (ci ConverterImpl) Index() index.Index {
	return *ci.index
}

// Contains reports whether a key is in the Converter's replacement table.
func (ci ConverterImpl) Contains(key string) bool {
	return ci.index.Contains(token.FromString(ci.caser, key))
}

// Replacements returns a slice of Replacement in the lookup trie.
func (ci ConverterImpl) Replacements() []Replacement {
	indexedVals := ci.index.Values()
	res := make([]Replacement, len(indexedVals))
	for i, v := range indexedVals {
		res[i] = Replacement{
			Camel:     v.Camel.Value(),
			Screaming: v.Screaming.Value(),
		}
	}
	return res
}

func (ci *ConverterImpl) set(key, value string) {
	ci.index.Add(token.FromString(ci.caser, key), token.FromString(ci.caser, value))
}

// Set adds the key/value pair to the table.
func (ci *ConverterImpl) Set(key, value string) {
	kstr, keyHasLower := lowerAndCheck(key)
	vstr, valueHasLower := lowerAndCheck(value)
	ci.Delete(kstr)
	ci.Delete(vstr)

	// checking to see if we need to swap these.
	if !keyHasLower && valueHasLower {
		ci.set(value, key)
	} else {
		ci.set(key, value)
	}
}

// Remove deletes the key from the map. Either variant is sufficient.
func (ci *ConverterImpl) Delete(key string) {
	tok := token.FromString(ci.caser, key)
	ci.index.Delete(tok)
}

// Convert formats the string with the desired style.
func (ci ConverterImpl) Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string {
	tokens := ci.tokenizer.Tokenize(input, allowedSymbols, numberRules)
	var parts []string
	var ok bool
	var addedAsNumber bool
	idx := ci.Index()
	for i, tok := range tokens {
		switch tok.Len() {
		case 0:
			continue
		case 1:

			if idx, ok = idx.Match(tok); !ok {
				if idx.LastMatch().HasValue() {
					// appending the last match
					parts = append(parts, formatIndexedReplacement(style, repStyle, len(parts), idx.LastMatch()))
				}
				if idx.HasPartialMatches() {
					// checking to make sure it isn't a number
					accum := token.Append(ci.caser, tok, idx.PartialMatches()...)
					if accum.IsNumber() {
						parts = append(parts, FormatToken(style, len(parts), accum))
						addedAsNumber = true
					} else {
						for _, partok := range idx.PartialMatches() {
							parts = append(parts, FormatToken(style, len(parts), partok))
						}
						addedAsNumber = false
					}
				}
				if !addedAsNumber {
					parts = append(parts, FormatToken(style, len(parts), tok))
				}
				// resetting the index
				idx = ci.Index()
			}
		default:
			if idx.HasMatch() {
				parts = append(parts, formatIndexedReplacement(style, repStyle, len(parts), idx.LastMatch()))
			}
			if idx.HasPartialMatches() {
				for _, partok := range idx.PartialMatches() {
					parts = append(parts, FormatToken(style, len(parts), partok))
				}
			}
			if idx.HasMatch() || idx.HasPartialMatches() {
				// resetting index
				idx = ci.Index()
			}
			if rep, ok := idx.Get(tok); ok {
				parts = append(parts, formatIndexedReplacement(style, repStyle, len(parts), rep))
			} else if isNextTokenNumber(tokens, i) {
				if idx, ok = idx.Match(tok); !ok {
					parts = append(parts, FormatToken(style, len(parts), tok))
					idx = ci.Index()
				}
			} else {
				parts = append(parts, FormatToken(style, len(parts), tok))
			}
		}
	}
	result := strings.Builder{}
	result.Grow(len(input))

	shouldWriteDelimiter := false
	if idx.HasMatch() {
		parts = append(parts, formatIndexedReplacement(style, repStyle, len(parts), idx.LastMatch()))
	}

	if idx.HasPartialMatches() {
		for _, partok := range idx.PartialMatches() {
			parts = append(parts, FormatToken(style, len(parts), partok))
		}
	}
	for _, part := range parts {
		if shouldWriteDelimiter {
			result.WriteString(join)
		}
		result.WriteString(part)
		if !shouldWriteDelimiter {
			shouldWriteDelimiter = len(part) > 0 && len(join) > 0
		}
	}

	return result.String()
}

// FormatToken formats the token with the desired style.
func FormatToken(style Style, index int, tok token.Token) string {
	switch style {
	case StyleCamel:
		return string(tok.UpperFirstLowerRest())
	case StyleLowerCamel:
		if index == 0 {
			return string(tok.Lower())
		}
		return string(tok.UpperFirstLowerRest())
	case StyleScreaming:
		return tok.Upper()
	case StyleLower:
		return tok.Lower()
	}
	return tok.String()
}

func formatIndexedReplacement(style Style, replaceStyle ReplaceStyle, index int, rep index.IndexedReplacement) string {
	switch replaceStyle {
	case ReplaceStyleCamel:
		if index == 0 && style == StyleLowerCamel {
			return rep.Camel.Lower()
		}
		return rep.Camel.String()
	case ReplaceStyleScreaming:
		return rep.Screaming.String()
	case ReplaceStyleLower:
		return rep.Lower.String()
	default:
		return rep.Screaming.String()
	}
}

func isNextTokenNumber(tokens []token.Token, i int) bool {
	if i+1 < len(tokens) {
		return unicode.IsNumber(tokens[i+1].Runes()[0])
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
