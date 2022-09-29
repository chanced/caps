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

package token

import "unicode"

var DefaultCaser = Unicode{}

// Unicode is a Caser which uses the default unicode casing functions.
type Unicode struct{}

var (
	// Turkish is a Caser which uses unicode.TurkishCase rules.
	TurkishCaser = &unicode.TurkishCase
	// AzeriCaser is a Caser which uses unicode.AzeriCase rules.
	AzeriCaser = &unicode.AzeriCase
)

// ToTitle maps the rune to title case using unicode.ToTitle.
func (Unicode) ToTitle(r rune) rune {
	return unicode.ToTitle(r)
}

// ToUpper maps the rune to lower case using unicode.ToLower
func (Unicode) ToLower(r rune) rune {
	return unicode.ToLower(r)
}

// ToUpper maps the rune to upper case using unicode.ToUpper
func (Unicode) ToUpper(r rune) rune {
	return unicode.ToUpper(r)
}

// Caser is satisfied by types which can map runes to their lowercase and
// uppercase equivalents.
type Caser interface {
	// ToLower maps the rune to lower case
	ToLower(r rune) rune
	// ToUpper maps the rune to upper case
	ToUpper(r rune) rune
	// ToTitle maps the rune to title case
	ToTitle(r rune) rune
}

// CaserOrDefault returns the default caser if caser is nil.
func CaserOrDefault(caser Caser) Caser {
	if caser == nil {
		return DefaultCaser
	}
	return caser
}
