package caps

import (
	"sort"
	"unicode"

	"github.com/chanced/caps/token"
)

const (
	// DEFAULT_DELIMITERS is the default set of delimiters in string convert.
	DEFAULT_DELIMITERS string = " _.!?:;$-(){}[]#@&+~"
)

// DefaultTokenizer is the default Tokenizer.
var DefaultTokenizer StdTokenizer = NewTokenizer(DEFAULT_DELIMITERS, token.DefaultCaser)

// Tokenizer is an interface satisfied by tyeps which can
type Tokenizer interface {
	Tokenize(value string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token
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
func (ti StdTokenizer) Tokenize(str string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token {
	tokens := []token.Token{}
	pending := []token.Token{}
	foundLower := false
	current := token.Token{}
	prevNumber := false
	allowed := newRunes(allowedSymbols)

	for i, r := range str {
		switch {
		case unicode.IsUpper(r):
			if foundLower && current.Len() > 0 {
				tokens = append(tokens, current)
				current = token.Token{}
			}
			current = token.AppendRune(ti.caser, current, r)
			prevNumber = false
		case unicode.IsLower(r):
			if !foundLower && current.Len() > 0 {
				// we have to break up the pending first
				for _, tok := range pending {
					if tok.IsNumber(numberRules) {
						tokens = append(tokens, tok)
					} else {
						tokens = append(tokens, tok.Split()...)
					}
				}
				// need to break up the current token if it isn't a number
				if prevNumber {
					tokens = append(tokens, current)
					current = token.Token{}
				} else {
					split := current.Split()
					// current becomes the last upper token before discovering the lowercase token
					current = split[len(split)-1]

					// all other tokens are added to the token list
					tokens = append(tokens, split[:len(split)-1]...)
				}
			}
			current = token.AppendRune(ti.caser, current, r)
			pending = nil
			foundLower = true
		case unicode.IsNumber(r):

			// if adding the number onto current makes it a valid number
			// then append this rune to current
			if token.AppendRune(ti.caser, current, r).IsNumber(numberRules) {
				current = token.AppendRune(ti.caser, current, r)
			} else {
				// otherwise it is not a number and so we add the current token
				// to the token or pending list depending on whether or not we
				// have found a lowercase rune
				if current.Len() > 0 && foundLower {
					tokens = append(tokens, current)
					current = token.Token{}
				} else if current.Len() > 0 { // otherwise, we have to push the current token into a pending state
					pending = append(pending, current)
					current = token.Token{}
				}
				current = token.AppendRune(ti.caser, current, r)
			}
			prevNumber = true
		default:
			if allowed.Contains(r) {
				if current.Len() > 0 {
					if current.IsNumber(numberRules) {
						// this gets a bit tricky because we need to check if adding the
						// rune to the current makes it a number. However, in all
						// default cases, a non-numeric rune must be followed either by
						// a number or an 'e' and a number. as such, we have to check if
						// both this and the next rune (and possibly the rune after
						// that) make it a number.
						n := token.AppendRune(ti.caser, current, r)
						runes := []rune(str)
						if n.IsNumber(numberRules) {
							current = n
						} else if i <= len(runes)-2 && (unicode.IsNumber(runes[i+1]) || unicode.IsLetter(runes[i+1])) || allowed.Contains(runes[i+1]) {
							next := token.AppendRune(ti.caser, n, runes[i+1])
							if next.IsNumber(numberRules) {
								current = n
							} else {
								if foundLower {
									tokens = append(tokens, current)
									current = token.FromRunes(ti.caser, []rune{r})
								} else {
									pending = append(pending, current)
									current = token.FromRunes(ti.caser, []rune{r})
								}
							}
						} else {
							if foundLower {
								tokens = append(tokens, current)
								current = token.FromRunes(ti.caser, []rune{r})
							} else {
								pending = append(pending, current)
								current = token.FromRunes(ti.caser, []rune{r})
							}
						}
					} else {
						current = token.AppendRune(ti.caser, current, r)
					}
				} else {
					current = token.AppendRune(ti.caser, current, r)
				}
			} else if ti.delimiters.Contains(r) || unicode.IsSpace(r) {
				if current.Len() > 0 {
					if foundLower {
						tokens = append(tokens, current)
						current = token.Token{}
					} else {
						pending = append(pending, current)
						current = token.Token{}
					}
				}
			}
		}
	}
	if current.Len() > 0 {
		if current.IsNumber(numberRules) {
			if foundLower {
				tokens = append(tokens, current)
			} else {
				pending = append(pending, current)
			}
		} else if !foundLower {
			pending = append(pending, current)
		} else {
			tokens = append(tokens, current)
		}
	}
	if foundLower {
		for _, tok := range pending {
			if tok.IsNumber(numberRules) {
				tokens = append(tokens, tok)
			} else {
				tokens = append(tokens, tok.Split()...)
			}
		}
		return tokens
	} else {
		return pending
	}
}

// Deprecated: Use StdTokenizer.
type TokenizerImpl = StdTokenizer

var _ Tokenizer = StdTokenizer{}
