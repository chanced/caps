package token

import (
	"unicode"
)

// Token contains a slice of runes representing the raw value and value in
// lowercase form.
//
// This is used for the
type Token struct {
	value []rune
	lower []rune
	upper []rune
	len   int
}

// Append appends all of o to t
func Append(t Token, elems ...Token) Token {
	for _, e := range elems {
		t = Token{
			value: append(t.value, e.value...),
			lower: append(t.lower, e.lower...),
			upper: append(t.upper, e.upper...),
			len:   t.len + e.len,
		}
	}
	return t
}

// AppendRune append the rune to the current token.
func AppendRune(t Token, r rune) Token {
	return Token{
		value: append(t.value, r),
		lower: append(t.lower, unicode.ToLower(r)),
		upper: append(t.upper, unicode.ToUpper(r)),
		len:   t.len + 1,
	}
}

func FromString[T ~string](value T) Token {
	return FromRunes([]rune(value))
}

func FromRune(value rune) Token {
	return Token{
		value: []rune{value},
		lower: []rune{unicode.ToLower(value)},
		upper: []rune{unicode.ToUpper(value)},
		len:   1,
	}
}

func FromRunes(value []rune) Token {
	upper := make([]rune, len(value))
	lower := make([]rune, len(value))
	for i, r := range value {
		upper[i] = unicode.ToUpper(r)
		lower[i] = unicode.ToLower(r)
	}
	return Token{
		value: value,
		lower: lower,
		upper: upper,
		len:   len(value),
	}
}

func (t Token) String() string {
	return t.Value()
}

// Len returns the number of runes in the Part.
func (t Token) Len() int {
	return t.len
}

func (t Token) Value() string {
	return string(t.value)
}

func (t Token) Lower() string {
	return string(t.lower)
}

func (t Token) Upper() string {
	return string(t.upper)
}

func (t Token) IsEmpty() bool {
	return t.Len() == 0
}

func (t Token) UpperFirstLowerRest() string {
	switch len(t.value) {
	case 0:
		return ""
	case 1:
		return string(t.upper)
	default:
		return string(append([]rune{t.upper[0]}, t.lower[1:]...))
	}
}

func (t Token) UpperFirst() string {
	switch len(t.value) {
	case 0:
		return ""
	case 1:
		return string(t.upper)
	default:
		return string(append([]rune{t.upper[0]}, t.value[1:]...))
	}
}

func (t Token) LowerFirst() string {
	switch len(t.value) {
	case 0:
		return ""
	case 1:
		return string(t.lower)
	default:
		return string(append([]rune{t.lower[0]}, t.value[1:]...))
	}
}

// IsNumber reports true if the Token is considered a valid number based on the
// following rules:
//
// - If the Token is composed only of numbers
//
// - If the Token is prefixed with any of the following: + - . v V # and
// followed by a number
//
// - Numbers may only be separated by a single '.' and '.' may be the first rune
// or proceeded by a number, '+', or '-'
//
// - A single 'e' or 'E' may only be used in the exponent portion of a number
//
// - 'e' or 'E' may be followed
//
// - ',' must be preceded by a number and followed by a number
//
// - if additionalRules is not nil and the rune is present in the map, the
// result of the provided func overrides the rules above
func (t Token) IsNumber(additionalRules ...map[rune]func(index int, r rune, val []rune) bool) bool {
	var rules map[rune]func(index int, r rune, val []rune) bool
	if additionalRules != nil {
		if len(additionalRules) == 1 {
			rules = additionalRules[0]
		} else {
			rules = make(map[rune]func(index int, r rune, val []rune) bool)
			for _, r := range additionalRules {
				for k, v := range r {
					if _, ok := rules[k]; ok {
						rules[k] = v
					}
				}
			}
		}
	}
	isDec := false
	var prev rune
	e := -1

	if len(t.value) == 0 {
		return false
	}

	for i, r := range t.value {
		if rules != nil {
			if fn, ok := rules[r]; ok {
				if !fn(i, r, t.value) {
					return false
				}
				prev = r
				continue
			}
		}
		if !unicode.IsNumber(r) {
			switch r {
			case 'v', 'V', '#':
				if i > 0 {
					return false
				}
			case '+', '-':
				if prev > 0 && (!isDec || e != i-1) {
					return false
				}
			case '.':
				if t.Len() == 1 {
					return false
				}
				if i == t.Len()-1 {
					return false
				}
				if prev > 0 && !unicode.IsNumber(prev) && prev != '-' && prev != '+' {
					return false
				}
				isDec = true
			case 'e', 'E':
				if !isDec {
					return false
				}
				if e != -1 {
					return false
				}
				if i == len(t.value)-1 {
					return false
				}
				e = i
			default:
				return false
			}
		}
		prev = r
	}

	return true
}

func (t Token) Reverse() Token {
	r := Token{
		value: make([]rune, len(t.value)),
		lower: make([]rune, len(t.lower)),
		upper: make([]rune, len(t.upper)),
		len:   t.len,
	}
	x := 0
	for i := t.len - 1; i >= 0; i-- {
		x = t.len - 1 - i
		r.value[x] = t.value[i]
		r.lower[x] = t.lower[i]
		r.upper[x] = t.upper[i]

	}
	return r
}

// Split returns the current token split into a slice of Tokens for each rune in
// the list.
func (t Token) Split() []Token {
	result := make([]Token, t.Len())
	for i, r := range t.value {
		result[i] = Token{
			value: []rune{r},
			lower: []rune{t.lower[i]},
			upper: []rune{t.upper[i]},
			len:   1,
		}
	}
	return result
}

func (t Token) FirstLowerRune() (rune, bool) {
	if t.len == 0 {
		return 0, false
	}
	return t.lower[0], true
}

func (t Token) FirstUpperRune() (rune, bool) {
	if t.len == 0 {
		return 0, false
	}
	return t.upper[0], true
}

func (t Token) FirstRune() (rune, bool) {
	if t.len == 0 {
		return 0, false
	}
	return t.value[0], true
}

func (t Token) ReverseSplit() []Token {
	result := make([]Token, t.Len())
	for i := t.len - 1; i >= 0; i-- {
		result[t.len-1-i] = Token{
			value: []rune{t.value[i]},
			lower: []rune{t.lower[i]},
			upper: []rune{t.upper[i]},
			len:   1,
		}
	}
	return result
}

func (t Token) LowerRunes() []rune {
	return t.lower
}

func (t Token) UpperRunes() []rune {
	return t.upper
}

func (t Token) Runes() []rune {
	return t.value
}

func (t Token) LowerReversedRunes() []rune {
	res := make([]rune, len(t.lower))
	for i := t.len - 1; i >= 0; i-- {
		res[t.len-1-i] = t.lower[i]
	}
	return res
}
