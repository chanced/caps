package index_test

import (
	"testing"

	"github.com/chanced/caps/index"
	"github.com/chanced/caps/token"
)

func TestAddMatchGet(t *testing.T) {
	idx := index.New()

	tests := []struct {
		camel     string
		screaming string
	}{
		{"Json", "JSON"},
		{"Jsonp", "JSONP"},
		{"Js", "JS"},
	}
	for _, test := range tests {
		tok := token.FromString(test.camel)
		str := tok.Lower()
		idx.Add(token.FromString(test.camel), token.FromString(test.screaming))
		for i := range str {
			if i == 0 {
				break
			}
			ts := token.FromString(str[:i])
			if i == len(str)-1 {
				idx, hasMatch := idx.MatchForward(ts)
				if !hasMatch {
					t.Error("expected match for", ts)
				}
				if !idx.HasMatch() {
					t.Error("expected match for", ts)
				}
				_, ok := idx.GetForward(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			idx, hasMatch := idx.MatchForward(ts)
			if !hasMatch {
				t.Error("expected match for", ts)
			}
			if !idx.HasPartialMatches() {
				t.Error("expected", ts, "to be in index")
			}
		}

		// testing reverse

		tok = tok.Reverse()
		str = tok.Lower()
		for i := range str {
			if i == 0 {
				break
			}
			ts := token.FromString(str[:i])
			if i == len(str)-1 {
				idx, hasMatch := idx.MatchReverse(ts)
				if !hasMatch {
					t.Error("expected match for", ts)
				}
				if !idx.HasMatch() {
					t.Error("expected match for", ts)
				}
				_, ok := idx.GetForward(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			idx, hasMatch := idx.MatchReverse(ts)
			if !hasMatch {
				t.Error("expected", ts, "to be in index")
			}
			if !idx.HasPartialMatches() {
				t.Error("expected", ts, "to be in index")
			}
		}

	}
}

func TestDelete(t *testing.T) {
	idx := index.New()

	tests := []struct {
		camel     string
		screaming string
	}{
		{"Jsonp", "JSONP"},
		{"Json", "JSON"},
		{"Js", "JS"},
	}
	for _, test := range tests {
		idx.Add(token.FromString(test.camel), token.FromString(test.screaming))
	}

	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		tok := token.FromString(test.camel)

		if ok := idx.ContainsForward(tok); !ok {
			t.Error("expected", tok.Lower(), "to be in index")
		}
		tok = tok.Reverse()
		if ok := idx.ContainsReverse(tok); !ok {
			t.Error("expected", tok.Lower(), "to be in index")
		}
		tok = tok.Reverse()
		idx.Delete(tok)
	}
	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		tok := token.FromString(test.camel)

		if ok := idx.ContainsForward(tok); ok {
			t.Error("expected", tok.Lower(), "to have been deleted")
		}
		tok = tok.Reverse()
		if ok := idx.ContainsReverse(tok); ok {
			t.Error("expected", tok.Lower(), "to have been deleted")
		}
	}
}

func TestPartialMatches(t *testing.T) {
	idx := index.New()
	idx.Add(token.FromString("Abcd"), token.FromString("ABCD"))

	m, ok := idx.MatchForward(token.FromString("abc"))
	if !ok {
		t.Error("expected match for abc")
	}

	merged := token.Merge(m.PartialMatches()...)
	if merged.Lower() != "abc" {
		t.Error("expected abc, got", merged.Lower())
	}

	rm, ok := idx.MatchReverse(token.FromString("dcb"))
	if !ok {
		t.Error("expected match for abc")
	}

	merged = token.Merge(rm.PartialMatches()...)
	if merged.Lower() != "bcd" {
		t.Error("expected bcd, got", merged.Lower())
	}
}
