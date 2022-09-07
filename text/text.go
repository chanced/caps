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

// Package text contains a single Text type with functions from caps and strings
// as methods.
//
// All methods return new values and do not mutate the existing Text.
package text

import (
	"sort"
	"strings"
	"unicode"

	"github.com/chanced/caps"
	"github.com/chanced/caps/token"
)

type Texts []Text

// Less implements sort.Interface
func (t Texts) Less(i int, j int) bool { return t[i] < t[j] }

// Swap implements sort.Interface
func (t Texts) Swap(i int, j int) { t[i], t[j] = t[j], t[i] }

func (t Texts) Len() int {
	return len(t)
}

// TotalLen returns the sum len of each Text in t.
func (t Texts) TotalLen() int {
	var l int
	for _, v := range t {
		l += len(v)
	}
	return l
}

func (t Texts) Join(sep string) Text {
	var b strings.Builder
	b.Grow(t.TotalLen() + len(sep)*(len(t)-1))
	for i, v := range t {
		if i > 0 && len(sep) > 0 {
			b.WriteString(sep)
		}
		b.WriteString(v.String())
	}
	return Text(b.String())
}

type Text string

func (t Text) String() string {
	return string(t)
}

// ToLower returns t with all Unicode letters mapped to their lower case.
func (t Text) ToLower() Text {
	return Text(strings.ToLower(t.String()))
}

// ToUpper returns t with all Unicode letters mapped to their upper case.
func (t Text) ToUpper(caser ...token.Caser) Text {
	return Text(strings.ToUpper(t.String()))
}

// ToCamel transforms the case of the Text t into Camel Case (e.g.
// AnExampleString) using either the provided Converter or the DefaultConverter
// otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes
// "AnExampleString". It also has a configurable set of replacements, such that
// "some_json" becomes "SomeJSON" so long as opts.ReplacementStyle is set to
// ReplaceStyleScreaming. A ReplaceStyle of ReplaceStyleCamel would result in
// "SomeJson".
func (t Text) ToCamel(opts ...caps.Opts) Text {
	return caps.ToCamel(t, opts...)
}

// ToLowerCamel transforms the case of the Text t into Lower Camel Case (e.g.
// anExampleString) using either the provided Converter or the DefaultConverter
// otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes
// "anExampleString". It also has a configurable set of replacements, such that
// "some_json" becomes "someJSON" so long as opts.ReplacementStyle is set to
// ReplaceStyleScreaming. A ReplaceStyle of ReplaceStyleCamel would result in
// "someJson".
func (t Text) ToLowerCamel(opts ...caps.Opts) Text {
	return caps.ToLowerCamel(t, opts...)
}

// ToSnake transforms the case of the Text t into Lower Snake Case (e.g.
// an_example_string) using either the provided Converter or the
// DefaultConverter otherwise.
func (t Text) ToSnake(opts ...caps.Opts) Text {
	return caps.ToSnake(t, opts...)
}

// ToScreamingSnake transforms the case of the Text t into Screaming Snake Case
// (e.g. AN_EXAMPLE_STRING) using either the provided Converter or the
// DefaultConverter otherwise.
func (t Text) ToScreamingSnake(opts ...caps.Opts) Text {
	return caps.ToScreamingSnake(t, opts...)
}

// ToKebab transforms the case of Text t into Lower Kebab Case (e.g. an-example-string) using
// either the provided Converter or the DefaultConverter otherwise.
func (t Text) ToKebab(opts ...caps.Opts) Text {
	return caps.ToKebab(t, opts...)
}

// ToScreamingKebab transforms the case of the Text t into Screaming Kebab Snake (e.g.
// AN-EXAMPLE-STRING) using either the provided Converter or the
// DefaultConverter otherwise.
func (t Text) ToScreamingKebab(opts ...caps.Opts) Text {
	return caps.ToScreamingKebab(t, opts...)
}

// ToDotNotation transforms the case of the Text t into Lower Dot Notation Case (e.g. an.example.string) using
// either the provided Converter or the DefaultConverter otherwise.
func (t Text) ToDotNotation(opts ...caps.Opts) Text {
	return caps.ToDotNotation(t, opts...)
}

// ReplaceAll returns a copy of the Text t with all non-overlapping instances
// of old replaced by new. If old is empty, it matches at the beginning of the
// string and after each UTF-8 sequence, yielding up to k+1 replacements for a
// k-rune string.
func (t Text) ReplaceAll(old string, new string) Text {
	return Text(strings.ReplaceAll(t.String(), old, new))
}

// ToScreamingKebab transforms the case of the Text t into Screaming Dot Notation Case
// (e.g. AN.EXAMPLE.STRING) using either the provided Converter or the
// DefaultConverter otherwise.
func (t Text) ToScreamingDotNotation(opts ...caps.Opts) Text {
	return Text(caps.ToScreamingDotNotation(t.String(), opts...))
}

// ToTitle transforms the case of t into Title Case (e.g. An Example String) using
// either the provided Converter or the DefaultConverter otherwise.
func (t Text) ToTitle(opts ...caps.Opts) Text {
	return caps.ToTitle(t, opts...)
}

// ToDelimited transforms the case of t into Text separated by delimiter,
// using either the provided Converter or the DefaultConverter otherwise.
//
// If lowercase is false, the output will be all uppercase.
func (t Text) ToDelimited(delimiter string, lowercase bool, opts ...caps.Opts) Text {
	return caps.ToDelimited(t, delimiter, lowercase, opts...)
}

// Replace returns a copy of the Text t with the first n non-overlapping
// instances of old replaced by new. If old is empty, it matches at the
// beginning of the string and after each UTF-8 sequence, yielding up to k+1
// replacements for a k-rune string. If n < 0, there is no limit on the number
// of replacements.
func (t Text) Replace(old, new string, n int) Text {
	return Text(strings.Replace(t.String(), old, new, n))
}

// Compare returns an integer comparing two Texts lexicographically. The result
// will be 0 if a == b, -1 if a < b, and +1 if a > b.
//
// Compare is included only for symmetry with package bytes. It is usually
// clearer and always faster to use the built-in string comparison operators ==,
// <, >, and so on.
func (t Text) Compare(other string) int {
	return strings.Compare(t.String(), other)
}

// Trim returns a slice of the Text t with all leading and
// trailing Unicode code points contained in cutset removed.
func (t Text) Trim(cutset string) Text {
	return Text(strings.Trim(t.String(), cutset))
}

// TrimLeft returns a slice of the Text t with all leading
// Unicode code points contained in cutset removed.
//
// To remove a prefix, use TrimPrefix instead.
func (t Text) TrimLeft(cutset string) Text {
	return Text(strings.TrimLeft(t.String(), cutset))
}

// TrimRight returns a slice of the Te t, with all trailing
// Unicode code points contained in cutset removed.
//
// To remove a suffix, use TrimSuffix instead.
func (t Text) TrimRight(cutset string) Text {
	return Text(strings.TrimRight(t.String(), cutset))
}

// TrimSpace returns a slice of the Text t, with all leading
// and trailing white space removed, as defined by Unicode.
func (t Text) TrimSpace() Text {
	return Text(strings.TrimSpace(t.String()))
}

// TrimPrefix returns t without the provided leading prefix string. If t doesn't
// start with prefix, t is returned unchanged.
func (t Text) TrimPrefix(prefix string) Text {
	return Text(strings.TrimPrefix(t.String(), prefix))
}

// TrimSuffix returns t without the provided trailing suffix string.
// If t doesn't end with suffix, t is returned unchanged.
func (t Text) TrimSuffix(suffix string) Text {
	return Text(strings.TrimSuffix(t.String(), suffix))
}

// TrimLeftFunc returns a slice of the Text t with all leading
// Unicode code points c satisfying f(c) removed.
func (t Text) TrimLeftFunc(f func(r rune) bool) Text {
	return Text(strings.TrimLeftFunc(t.String(), f))
}

// TrimRightFunc returns a slice of the Text t with all trailing
// Unicode code points c satisfying f(c) removed.
func (t Text) TrimRightFunc(f func(r rune) bool) Text {
	return Text(strings.TrimRightFunc(t.String(), f))
}

// EqualFold reports whether t and v, interpreted as UTF-8 strings,
// are equal under simple Unicode case-folding, which is a more general
// form of case-insensitivity.
func (t Text) EqualFold(v string) bool {
	return strings.EqualFold(t.String(), v)
}

// Index returns the index of the first instance of substr in t, or -1 if substr
// is not present in t.
func (t Text) Index(substr string) int {
	return strings.Index(t.String(), substr)
}

// IndexByte returns the index of the first instance of c in t, or -1 if c is
// not present in t.
func (t Text) IndexByte(c byte) int {
	return strings.IndexByte(t.String(), c)
}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in t. If r is utf8.RuneError, it returns the
// first instance of any invalid UTF-8 byte sequence.
func (t Text) IndexRune(r rune) int {
	return strings.IndexRune(t.String(), r)
}

// IndexRune returns the index of the first instance of the Unicode code point
// r, or -1 if rune is not present in t. If r is utf8.RuneError, it returns the
// first instance of any invalid UTF-8 byte sequence.
func (t Text) IndexFunc(fn func(r rune) bool) int {
	return strings.IndexFunc(t.String(), fn)
}

// Cut slices t around the first instance of sep,
// returning the text before and after sep.
// The found result reports whether sep appears in t.
// If sep does not appear in t, cut returns t, "", false.
func (t Text) Cut(sep string) (before, after Text, found bool) {
	b, a, f := strings.Cut(t.String(), sep)

	return Text(b), Text(a), f
}

// Clone returns a fresh copy of t.
// It guarantees to make a copy of t into a new allocation,
// which can be important when retaining only a small substring
// of a much larger string. Using Clone can help such programs
// use less memory. Of course, since using Clone makes a copy,
// overuse of Clone can make programs use more memory.
// Clone should typically be used only rarely, and only when
// profiling indicates that it is needed.
// For strings of length zero the string "" will be returned
// and no allocation is made.
func (t Text) Clone() Text {
	return Text(strings.Clone(string(t)))
}

// Contains reports whether substr is within t.
func (t Text) Contains(substr string) bool {
	return strings.Contains(string(t), substr)
}

// ContainsAny reports whether any Unicode code points in chars are within t.
func (t Text) ContainsAny(chars string) bool {
	return strings.ContainsAny(string(t), chars)
}

// ContainsRune reports whether the Unicode code point r is within t.
func (t Text) ContainsRune(r rune) bool {
	return strings.ContainsRune(string(t), r)
}

// Count counts the number of non-overlapping instances of substr in t.
func (t Text) Count(substr string) int {
	return strings.Count(string(t), substr)
}

// Fields splits the Text t around each instance of one or more consecutive
// white space characters, as defined by unicode.IsSpace, returning a slice of
// substrings of t or an empty slice if t contains only white space.
func (t Text) Fields() Texts {
	return collect(strings.Fields(string(t)))
}

// FieldsFunc splits the Text t at each run of Unicode code points c satisfying f(c)
// and returns an array of slices of t. If all code points in t satisfy f(c) or the
// Text is empty, an empty slice is returned.

// FieldsFunc makes no guarantees about the order in which it calls f(c)
// and assumes that f always returns the same value for a given c.
func (t Text) FieldsFunc(f func(rune) bool) Texts {
	return collect(strings.FieldsFunc(string(t), f))
}

// HasPrefix tests whether the Text t begins with prefix.
func (t Text) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(t), prefix)
}

// HasSuffix tests whether the Text t ends with suffix.
func (t Text) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(t), suffix)
}

// Append appends each elem to a copy of t and returns the result.
func (t Text) Append(elems ...Text) Text {
	for _, e := range elems {
		t += e
	}
	return t
}

// AppendRune append each rune in elem to a copy t and returns the result.
func (t Text) AppendRune(elems ...rune) Text {
	for _, e := range elems {
		t += Text(e)
	}

	return t
}

// LastIndex returns the index of the last instance of substr in t, or -1 if substr is not present in t.
func (t Text) LastIndex(substr string) int {
	return strings.LastIndex(string(t), substr)
}

// LastIndexAny returns the index of the last instance of any Unicode code point
// from chars in t, or -1 if no Unicode code point from chars is present in t.
func (t Text) LastIndexAny(chars string) int {
	return strings.LastIndexAny(string(t), chars)
}

// LastIndexByte returns the index of the last instance of c in t, or -1 if c is
// not present in t.
func (t Text) LastIndexByte(b byte) int {
	return strings.LastIndexByte(string(t), b)
}

// LastIndexByte returns the index of the last instance of r in t, or -1 if r is
// not present in t.
func (t Text) LastIndexRune(r rune) int {
	return strings.LastIndexFunc(string(t), func(rv rune) bool {
		return rv == r
	})
}

// LastIndexFunc returns the index into t of the last Unicode code point
// satisfying f(c), or -1 if none do.
func (t Text) LastIndexFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(string(t), f)
}

// Map returns a copy of the Text t with all its characters modified according
// to the mapping function. If mapping returns a negative value, the character
// is dropped from the string with no replacement.
func (t Text) Map(f func(r rune) rune) Text {
	return Text(strings.Map(f, string(t)))
}

// Repeat returns a new Text consisting of count copies of the Text t.
//
// It panics if count is negative or if the result of (len(t) * count) overflows.
func (t Text) Repeat(count int) Text {
	return Text(strings.Repeat(string(t), count))
}

// Split slices t into all substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// If t does not contain sep and sep is not empty, Split returns a
// slice of length 1 whose only element is t.
//
// If sep is empty, Split splits after each UTF-8 sequence. If both t
// and sep are empty, Split returns an empty slice.
func (t Text) Split(sep string) Texts {
	return collect(strings.Split(string(t), sep))
}

// SplitAfter slices t into all substrings after each instance of sep and
// returns a slice of those substrings.
//
// If t does not contain sep and sep is not empty, SplitAfter returns
// a slice of length 1 whose only element is t.
//
// If sep is empty, SplitAfter splits after each UTF-8 sequence. If
// both t and sep are empty, SplitAfter returns an empty slice.
//
// It is equivalent to SplitAfterN with a count of -1.
func (t Text) SplitAfter(sep string) Texts {
	return collect(strings.SplitAfter(string(t), sep))
}

// SplitAfterN slices t into substrings after each instance of sep and
// returns a slice of those substrings.
//
// The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
//
// Edge cases for t and sep (for example, empty strings) are handled
// as described in the documentation for SplitAfter.
func (t Text) SplitAfterN(sep string, n int) Texts {
	return collect(strings.SplitAfterN(string(t), sep, n))
}

// SplitN slices t into substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// The count determines the number of substrings to return:
//
//	n > 0: at most n substrings; the last substring will be the unsplit remainder.
//	n == 0: the result is nil (zero substrings)
//	n < 0: all substrings
//
// Edge cases for t and sep (for example, empty strings) are handled
// as described in the documentation for Split.
//
// To split around the first instance of a separator, see Cut.
func (t Text) SplitN(sep string, n int) Texts {
	return collect(strings.SplitN(string(t), sep, n))
}

// Title returns a copy of the Text t with all Unicode letters that begin words
// mapped to their title case.
func (t Text) Title() Text {
	return Text(strings.Title(string(t)))
}

// ToLowerSpecial returns a copy of the Text t with all Unicode letters mapped
// to their lower case using the case mapping specified by c.
func (t Text) ToLowerSpecial(c unicode.SpecialCase) Text {
	return Text(strings.ToLowerSpecial(c, string(t)))
}

// ToUpperSpecial returns a copy of the Text t with all Unicode letters mapped
// to their upper case using the case mapping specified by c.
func (t Text) ToUpperSpecial(c unicode.SpecialCase) Text {
	return Text(strings.ToUpperSpecial(c, string(t)))
}

// ToValidUTF8 returns a copy of the Text t with each run of invalid UTF-8 byte sequences
// replaced by the replacement string, which may be empty.
func (t Text) ToValidUTF8(replacement string) Text {
	return Text(strings.ToValidUTF8(string(t), replacement))
}

func collect(slice []string) Texts {
	res := make([]Text, len(slice))
	for i, v := range slice {
		res[i] = Text(v)
	}
	return res
}

var _ sort.Interface = (*Texts)(nil)
