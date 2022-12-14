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

// Package token contains functions for working with a string as a token.
package token

import (
	"strings"
	"unicode"
)

// NumberRules are a set of rules for determining if rune r at index of val
// should be considered a number
//
// Example:
//
//	NumberRules{
//	    '$': func(i int, r rune, v string) bool {
//	        return i == 0
//	    },
//	}
type NumberRules map[rune]func(index int, r rune, val string) bool

// Append appends all of elems to t
func Append(caser Caser, t string, elems ...string) string {
	caser = CaserOrDefault(caser)
	b := strings.Builder{}
	b.WriteString(t)
	for i, e := range elems {
		if len(e) == 0 {
			continue
		}
		b.Grow(len(e))
		for y, r := range e {
			if y == 0 && i == 0 && len(t) > 0 && unicode.IsTitle(r) {
				b.WriteRune(caser.ToUpper(r))
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

// WriteUpperFirstLowerRest writes the first rune as upper case and the rest as
// lower case
func WriteUpperFirstLowerRest(b *strings.Builder, caser Caser, s string) {
	for i, r := range s {
		switch {
		case i == 0 && b.Len() == 0:
			b.WriteRune(caser.ToTitle(r))
		case i == 0 && b.Len() > 0:
			b.WriteRune(caser.ToUpper(r))
		default:
			b.WriteRune(caser.ToLower(r))
		}
	}
}

// WriteLowerFirstUpperRest writes the first rune as upper case and the rest are
// separated by sep and written as lower case
func WriteSplitLowerFirstUpperRest(b *strings.Builder, caser Caser, sep string, s string) {
	WriteSplitLowerFirstUpperRestRunes(b, caser, sep, []rune(s))
}

// WriteSplitLowerFirstUpperRestRunes writes the first rune as upper case and
// the rest are separated by sep and written as lower case
func WriteSplitLowerFirstUpperRestRunes(b *strings.Builder, caser Caser, sep string, s []rune) {
	for i, r := range s {
		if i == 0 && b.Len() == 0 {
			b.WriteRune(caser.ToLower(r))
		} else if b.Len() > 0 {
			if len(sep) > 0 {
				b.WriteString(sep)
			}
			if i == 0 {
				b.WriteRune(caser.ToLower(r))
			} else {
				b.WriteRune(caser.ToUpper(r))
			}
		}
	}
}

// WriteSplitLower writes all strings in elems separated by sep and written as lower case
func WriteSplitLower(b *strings.Builder, caser Caser, sep string, elems ...string) {
	for _, s := range elems {
		for _, r := range s {
			if b.Len() > 0 && len(sep) > 0 {
				b.WriteString(sep)
			}
			b.WriteRune(caser.ToLower(r))
		}
	}
}

// WriteSplitUpper writes all runes in elems separated by sep and written as lower case
func WriteSplitLowerRunes(b *strings.Builder, caser Caser, sep string, s []rune) {
	for _, r := range s {
		if b.Len() > 0 && len(sep) > 0 {
			b.WriteString(sep)
		}
		b.WriteRune(caser.ToLower(r))
	}
}

// WriteSplitUpper writes all strings in elems separated by sep and written as upper case
func WriteSplitUpper(b *strings.Builder, caser Caser, sep string, elems ...string) {
	for _, s := range elems {
		WriteSplitUpperRunes(b, caser, sep, []rune(s))
	}
}

// WriteSplitUpperRunes uses caser to upper case each rune and writes each to b,
// separated by sep
func WriteSplitUpperRunes(b *strings.Builder, caser Caser, sep string, s []rune) {
	for i, r := range s {
		if b.Len() > 0 && len(sep) > 0 {
			b.WriteString(sep)
		}
		if i == 0 && b.Len() == 0 {
			b.WriteRune(caser.ToTitle(r))
		} else {
			b.WriteRune(caser.ToUpper(r))
		}
	}
}

// Write writes e to b
func Write(b *strings.Builder, caser Caser, e string) {
	caser = CaserOrDefault(caser)
	if len(e) == 0 {
		return
	}
	for y, r := range e {
		if y == 0 && b.Len() > 0 && unicode.IsTitle(r) {
			b.WriteRune(caser.ToUpper(r))
		} else {
			b.WriteRune(r)
		}
	}
}

// WriteUpper uses caser to upper case the runes in s and writes to b
func WriteUpper(b *strings.Builder, caser Caser, s string) {
	for _, r := range s {
		if b.Len() == 0 {
			b.WriteRune(caser.ToTitle(r))
		} else {
			b.WriteRune(caser.ToUpper(r))
		}
	}
}

// WriteLower uses caser to lower case the runes in s and writes to b
func WriteLower(b *strings.Builder, caser Caser, s string) {
	for _, r := range s {
		b.WriteRune(caser.ToLower(r))
	}
}

// WriteRune writes the runes to the b.
func WriteRune(b *strings.Builder, caser Caser, r rune) {
	if b.Len() > 0 && unicode.IsTitle(r) {
		r = caser.ToUpper(r)
	} else if b.Len() == 0 && unicode.IsUpper(r) {
		r = caser.ToTitle(r)
	}
	b.WriteRune(r)
}

// AppendRune append the rune to the current token.
func AppendRune(caser Caser, t string, runes ...rune) string {
	caser = CaserOrDefault(caser)
	b := strings.Builder{}
	b.Grow(len(t) + len(runes))
	b.WriteString(t)

	for _, r := range runes {
		if b.Len() > 0 && unicode.IsTitle(r) {
			r = caser.ToUpper(r)
		} else if b.Len() == 0 && unicode.IsUpper(r) {
			r = caser.ToTitle(r)
		}
		b.WriteRune(r)
	}
	return b.String()
}

// ToLower lower cases each rune of s using caser and writes to b
func ToLower(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	b := strings.Builder{}
	b.Grow(len(s))
	for _, r := range s {
		b.WriteRune(caser.ToLower(r))
	}
	return b.String()
}

// ToLUpper upper cases cases each rune of s using caser and writes to b
func ToUpper(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	b := strings.Builder{}
	b.Grow(len(s))
	for i, r := range s {
		if i == 0 {
			b.WriteRune(caser.ToTitle(r))
		} else {
			b.WriteRune(caser.ToUpper(r))
		}
	}
	return b.String()
}

// IsEmpty reports whether s has a len of 0
func IsEmpty(s string) bool {
	return len(s) == 0
}

// FirstRune returns the first rune of s
func FirstRune(s string) (rune, bool) {
	for _, r := range s {
		return r, true
	}
	return 0, false
}

// UpperFirstLowerRest title cases the first rune and lower cases the rest of s
func UpperFirstLowerRest(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	if len(s) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Grow(len(s))
	for i, r := range s {
		if i == 0 {
			b.WriteRune(caser.ToTitle(r))
		} else {
			b.WriteRune(caser.ToLower(r))
		}
	}
	return b.String()
}

// UpperFirst title cases the first rune of s
func UpperFirst(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	if len(s) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Grow(len(s))
	for i, r := range s {
		if i == 0 {
			b.WriteRune(caser.ToTitle(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// LowerFirst lower cases the first rune of s
func LowerFirst(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	if len(s) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.Grow(len(s))
	for i, r := range s {
		if i == 0 {
			sb.WriteRune(caser.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// IsNumber reports true if the string is considered a valid number based on the
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
func IsNumber(s string, additionalRules NumberRules) bool {
	isDec := false
	var prev rune
	e := -1

	if len(s) == 0 {
		return false
	}

	for i, r := range s {
		if additionalRules != nil {
			if fn, ok := additionalRules[r]; ok {
				if !fn(i, r, s) {
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
				if len(s) == 1 {
					return false
				}
				if i == len(s)-1 {
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
				if i == len(s)-1 {
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

// Reverse reverses s
func Reverse(caser Caser, s string) string {
	caser = CaserOrDefault(caser)
	if len(s) == 0 {
		return ""
	}
	b := strings.Builder{}
	b.Grow(len(s))
	var r rune
	runes := []rune(s)

	for i := len(runes) - 1; i >= 0; i-- {
		r = runes[i]
		switch {
		case i == len(runes)-1 && unicode.IsUpper(r):
			b.WriteRune(caser.ToTitle(r))
		case i == 0 && unicode.IsTitle(r):
			b.WriteRune(caser.ToUpper(r))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// HasLower returns true if any rune in
// the token is a unicode lowercase letter.
func HasLower(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}
