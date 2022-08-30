package index_test

import (
	"testing"

	"github.com/chanced/caps/index"
	"github.com/chanced/caps/token"
)

func TestAddMatchGet(t *testing.T) {
	index := index.New()

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
		index.Add(token.FromString(test.camel), token.FromString(test.screaming))
		for i := range str {
			if i == 0 {
				break
			}
			ts := token.FromString(str[:i])
			if i == len(str)-1 {
				idx := index.MatchForward(ts)
				if !idx.HasMatch() {
					t.Error("expected match for", ts)
				}
				_, ok := index.GetForward(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			idx := index.MatchForward(ts)
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
				idx := index.MatchReverse(ts)
				if !idx.HasMatch() {
					t.Error("expected match for", ts)
				}
				_, ok := index.GetForward(ts)
				if !ok {
					t.Error("expected get result for for", ts)
				}
				break
			}
			idx := index.MatchReverse(ts)
			if !idx.HasPartialMatches() {
				t.Error("expected", ts, "to be in index")
			}
		}

	}
}

func TestDelete(t *testing.T) {
	index := index.New()

	tests := []struct {
		camel     string
		screaming string
	}{
		{"Jsonp", "JSONP"},
		{"Json", "JSON"},
		{"Js", "JS"},
	}
	for _, test := range tests {
		index.Add(token.FromString(test.camel), token.FromString(test.screaming))
	}

	for i := len(tests) - 1; i >= 0; i-- {
		test := tests[i]
		tok := token.FromString(test.camel)

		ok := index.ContainsForward(tok)
		if !ok {
			t.Error("expected", tok.Lower(), "to be in index")
		}
		ok = index.ContainsReverse(tok.Reverse())
		if !ok {
			t.Error("expected", tok.Reverse().Lower(), "to be in index")
		}
		// if !index.Delete(tok) {
		// 	t.Error("expected", tok, "to be deleted")
		// }

	}
}
