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

package index_test

import (
	"strings"
	"testing"

	"github.com/chanced/caps/index"
	"github.com/chanced/caps/token"
)

func TestClone(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Abcd", "ABCD")

	if !idx.Contains("Abcd") {
		t.Error("expected idx to contain Abcd")
	}
	if idx.HasPartialMatches() {
		t.Errorf("expected idx to not contain partial matches has: %s", idx.PartialMatches())
	}

	match, ok := idx.Match("abc")
	if !ok {
		t.Error("expected match for abc")
	}
	if !match.HasPartialMatches() {
		t.Error("expected match to have partial matches")
	}
	clone := idx.Clone()
	if !clone.Contains("Abcd") {
		t.Error("expected clone to contain Abcd")
	}
	clone = match.Clone()
	if !clone.HasPartialMatches() {
		t.Error("expected clone to have partial matches")
	}
}

func TestValues(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Cat", "CAT")
	idx.Add("Dog", "DOG")
	idx.Add("Wolf", "WOLF")
	idx.Add("Dogfish", "DOGFISH")
	vals := idx.Values()

	if len(vals) != 4 {
		t.Error("expected 4 values, got", len(vals))
	}
	seen := map[string]bool{
		"Cat":     false,
		"Dog":     false,
		"Wolf":    false,
		"Dogfish": false,
	}
	for _, val := range vals {
		switch val.Camel {
		case "Cat", "Dog", "Wolf", "Dogfish":
			seen[val.Camel] = true
		default:
			t.Error("unexpected value", val)
		}
	}
	for k, v := range seen {
		if !v {
			t.Errorf("%s not in values", k)
		}
	}
}

func TestMatch(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Abcd", "ABCD")
	if _, ok := idx.Match(""); ok {
		t.Error("expected empty string to not match")
	}
	m, ok := idx.Match("abc")
	if !ok {
		t.Error("expected match for abc")
	}
	m, ok = m.Match("d")
	if !ok {
		t.Error("expected match for d")
	}
	if !m.HasMatched() {
		t.Error("expected match of ABCD")
	}
	m, ok = m.Match("z")
	if ok {
		t.Error("expected no match for z")
	}
	if !m.HasMatched() {
		t.Error("expected match of ABCD still")
	}
	if m.LastMatch().Lower != "abcd" {
		t.Error("expected last match of abcd")
	}
}

func TestAdd(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Abcd", "ABCD")
	if !idx.Contains("abcd") {
		t.Error("expected Abcd to be in index")
	}

	if _, ok := idx.Get(""); ok {
		t.Error("expected empty string to not be in index")
	}

	if _, ok := idx.Get("abcd"); !ok {
		t.Error("expected ABCD to be in index")
	}
	idx.Add("Abcd", "ABCD2")

	if !idx.Contains("abcd2") {
		t.Error("expected abcd2 to be in index")
	}
	if !idx.Contains("ABCD") {
		t.Error("expected ABCD to be in index")
	}

	idx.Add("Abcd2", "ABCD2")
	if v2, ok := idx.Get("abcd2"); ok {
		if v2.Camel != "Abcd2" {
			t.Error("expected Abcd2")
		}
	} else {
		t.Error("expected abcd2")
	}
	if idx.Contains("abcd") {
		t.Error("expected abcd to be deleted")
	}
}

func TestHasPartialMatches(t *testing.T) {
	idx := index.New(nil)

	tests := []struct {
		camel     string
		screaming string
	}{
		{"Json", "JSON"},
		{"Jsonp", "JSONP"},
		{"Js", "JS"},
		{"Utf8", "UTF8"},
	}
	for _, test := range tests {
		idx.Add(test.camel, test.screaming)
		for i := range test.camel {
			if i == 0 {
				break
			}
			ts := test.camel[:i]
			if i == len(test.camel)-1 {
				subidx, hasMatch := idx.Match(ts)

				if !hasMatch {
					t.Error("expected match for", ts)
				}
				if !subidx.HasMatched() {
					t.Error("expected match for", ts)
				}
				_, ok := subidx.Get(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			sub, hasMatch := idx.Match(ts)
			if !hasMatch {
				t.Error("expected match for", ts)
			}
			if !sub.HasPartialMatches() {
				t.Error("expected", ts, "to be in index")
			}

		}
		x, ok := idx.Match(test.camel)
		if !ok {
			t.Error("expected match for", test.camel)
		}
		v := x.Value()
		if v.Camel != test.camel {
			t.Error("expected", test.camel, "got", v.Camel)
		}
		if v.Screaming != test.screaming {
			t.Error("expected", test.screaming, "got", v.Screaming)
		}
		if v.Lower != strings.ToLower(test.camel) {
			t.Error("expected", strings.ToLower(test.camel), "got", v.Lower)
		}

	}
}

func TestDelete(t *testing.T) {
	idx := index.New(nil)

	tests := []struct {
		camel     string
		screaming string
	}{
		{"Jsonp", "JSONP"},
		{"Json", "JSON"},
		{"Js", "JS"},
	}
	for _, test := range tests {
		idx.Add(test.camel, test.screaming)
	}

	if ok := idx.Delete("json_doesnotexist"); ok {
		t.Error("expected delete to fail")
	}

	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]

		if ok := idx.Contains(test.camel); !ok {
			t.Error("expected", test.camel, "to be in index")
		}
	}
	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		idx.Delete(test.camel)
		if ok := idx.Contains(test.camel); ok {
			t.Error("expected", test.camel, "to have been deleted")
		}
	}
}

func TestPartialMatches(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Abcd", "ABCD")

	m, ok := idx.Match("abc")
	if !ok {
		t.Error("expected match for abc")
	}

	merged := m.PartialMatches()
	if token.ToLower(nil, merged) != "abc" {
		t.Error("expected abc, got", merged)
	}
}
