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
	"sort"
	"strings"
	"unicode"

	"github.com/chanced/caps/token"
)

const (
	// DEFAULT_DELIMITERS is the default set of delimiters in string convert.
	DEFAULT_DELIMITERS string = " _.!?:;$-(){}[]#@&+~"
)

// DefaultTokenizer is the default Tokenizer.
var DefaultTokenizer StdTokenizer = NewTokenizer(DEFAULT_DELIMITERS, token.DefaultCaser)

type NumberRules = token.NumberRules

// Tokenizer is an interface satisfied by tyeps which can
type Tokenizer interface {
	Tokenize(value string, allowedSymbols string, numberRules NumberRules) []string
}

// NewTokenizer creates and returns a new TokenizerImpl which implements the
// Tokenizer interface.
//
// Tokenizers are used by ConverterImpl to tokenize the input text into
// token.Tokens that are then formatted.
func NewTokenizer(delimiters string, caser token.Caser) StdTokenizer {
	d := runes(delimiters)
	sort.Sort(d)
	return StdTokenizer{
		delimiters: d,
		caser:      token.CaserOrDefault(caser),
	}
}

// StdTokenizer is the provided implementation of the Tokenizer interface.
//
// StdTokenizer tokenizes the input text into token.Tokens based on a set of
// delimiters (runes).
//
// If you need custom logic, consider wrapping the logic by implementing
// Tokenizer and then calling a StdTokenizer's Tokenize method.
//
// # Example:
type StdTokenizer struct {
	delimiters runes
	caser      token.Caser
}

// Tokenize splits a string into a list of token.Tokens based on the case of each
// rune, it's delimiters, and the specified allowedSymbols.
//
// For example:
//
//	t.Tokenize("ASnakecaseVariable", nil) -> ["A", "Snakecase", "Variable"]
//
// Tokenizer attempts to detect formatting, such that a screaming snakecase str
// will be inferred accordingly.
//
// For example:
//
//	t.Tokenize("A_SCREAMING_SNAKECASE_VARIABLE", nil) -> ["A", "SCREAMING", "SNAKECASE", "VARIABLE"]
//
// If allowedSymbols is not nil, then those symbols will be treated as
// non-delimiters even if they are in the delimiters list.
//
// For example:
//
//	t := caps.token.Newizer("_")
//	t.Tokenize("A_SCREAMING_SNAKECASE_VARIABLE", []rune{'_'}) -> ["A_SCREAMING_SNAKECASE_VARIABLE"]
func (ti StdTokenizer) Tokenize(str string, allowedSymbols string, numberRules NumberRules) []string {
	var tokens []string
	var pending []string

	switch {
	case len(str) == 0:
		return nil
	case len(str) < 6:
		tokens = make([]string, 0, 4)
	default:
		tokens = make([]string, 0, 8)
	}

	foundLower := false
	current := strings.Builder{}
	prevNumber := false
	allowed := newRunes(allowedSymbols)

	for i, r := range str {
		switch {
		case unicode.IsUpper(r):
			if foundLower && current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			current.WriteRune(r)
			prevNumber = false
		case unicode.IsLower(r):
			if !foundLower && current.Len() > 0 {
				// we have to break up the pending first
				for _, tok := range pending {
					if token.IsNumber(tok, numberRules) {
						tokens = append(tokens, tok)
					} else {
						tokens = append(tokens, strings.Split(tok, "")...)
					}
				}
				pending = nil
				// need to break up the current token if it isn't a number
				if prevNumber {
					tokens = append(tokens, current.String())
					current.Reset()
				} else {
					split := strings.Split(current.String(), "")
					// current becomes the last upper token before discovering the lowercase token
					current.Reset()
					current.WriteString(split[len(split)-1])
					// all other tokens are added to the token list
					tokens = append(tokens, split[:len(split)-1]...)
				}
			}
			tokens = append(tokens, pending...)
			current.WriteRune(r)
			pending = nil
			foundLower = true
		case unicode.IsNumber(r):
			// if adding the number onto current makes it a valid number
			// then append this rune to current
			if token.IsNumber(token.AppendRune(ti.caser, current.String(), r), numberRules) {
				current.WriteRune(r)
			} else {
				// otherwise it is not a number and so we add the current token
				// to the token or pending list depending on whether or not we
				// have found a lowercase rune
				if current.Len() > 0 && foundLower {
					tokens = append(tokens, current.String())
					current.Reset()
				} else if current.Len() > 0 { // otherwise, we have to push the current token into a pending state

					pending = append(pending, current.String())
					current.Reset()
				}
				current.WriteRune(r)
			}
			prevNumber = true
		default:
			if allowed.Contains(r) {
				if current.Len() > 0 {
					if token.IsNumber(current.String(), numberRules) {
						// this gets a bit tricky because we need to check if adding the
						// rune to current makes it a number. However, in all
						// default cases, a non-numeric rune must be followed either by
						// a number or an 'e' and a number. as such, we have to check if
						// both this and the next rune (and possibly the rune after
						// that) make it a number.
						n := token.AppendRune(ti.caser, current.String(), r)

						runes := []rune(str)
						if token.IsNumber(n, numberRules) {
							current.Reset()
							current.WriteString(n)
						} else if i < len(runes)-1 && canCheckNext(runes[i+1], allowed) {
							if token.IsNumber(token.AppendRune(ti.caser, n, runes[i+1]), numberRules) {
								current.Reset()
								current.WriteString(n)
							} else {
								if foundLower {
									tokens = append(tokens, current.String())
								} else {
									pending = append(pending, current.String())
								}
								current.Reset()
								current.WriteRune(r)
							}
						} else {
							if foundLower {
								tokens = append(tokens, current.String())
							} else {
								pending = append(pending, current.String())
							}
							current.Reset()
							current.WriteRune(r)
						}
					} else {
						current.WriteRune(r)
					}
				} else {
					current.Reset()
					current.WriteRune(r)
				}
			} else if ti.delimiters.Contains(r) || unicode.IsSpace(r) {
				if current.Len() > 0 {
					if foundLower {
						tokens = append(tokens, pending...)
						tokens = append(tokens, current.String())
					} else {
						pending = append(pending, current.String())
					}
					current.Reset()
				}
			}
		}
	}
	if current.Len() > 0 {
		if token.IsNumber(current.String(), numberRules) {
			if foundLower {
				tokens = append(tokens, current.String())
			} else {
				pending = append(pending, current.String())
			}
		} else if !foundLower {
			pending = append(pending, current.String())
		} else {
			tokens = append(tokens, current.String())
		}
	}
	if foundLower {
		for _, tok := range pending {
			if token.IsNumber(tok, numberRules) {
				tokens = append(tokens, tok)
			} else {
				tokens = append(tokens, strings.Split(tok, "")...)
			}
		}
		return tokens
	} else {
		return pending
	}
}

// Deprecated: Use StdTokenizer.
type TokenizerImpl = StdTokenizer

type runes []rune

// Len implements sort.Interface
func (r runes) Len() int {
	return len(r)
}

// Less implements sort.Interface
func (r runes) Less(i int, j int) bool {
	return r[i] < r[j]
}

// Swap implements sort.Interface
func (r runes) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r runes) Contains(c rune) bool {
	if len(r) == 0 {
		return false
	}
	res := sort.Search(len(r), func(i int) bool {
		return r[i] >= c
	})
	return res > -1 && res < len(r) && r[res] == c
}

func newRunes(val string) runes {
	if len(val) == 0 {
		return nil
	}
	r := runes([]rune(val))
	sort.Sort(r)
	return r
}

func canCheckNext(r rune, allowed runes) bool {
	return unicode.IsNumber(r) || unicode.IsLetter(r) || allowed.Contains(r)
}

var _ sort.Interface = (*runes)(nil)

var _ Tokenizer = StdTokenizer{}
