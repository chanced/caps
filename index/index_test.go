package index_test

import (
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

func TestMatch(t *testing.T) {
	idx := index.New(nil)
	idx.Add("Abcd", "ABCD")

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

func TestAddMatchGet(t *testing.T) {
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
				idx, hasMatch := idx.Match(ts)
				if !hasMatch {
					t.Error("expected match for", ts)
				}
				if !idx.HasMatched() {
					t.Error("expected match for", ts)
				}
				_, ok := idx.Get(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			idx, hasMatch := idx.Match(ts)
			if !hasMatch {
				t.Error("expected match for", ts)
			}
			if !idx.HasPartialMatches() {
				t.Error("expected", ts, "to be in index")
			}
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
