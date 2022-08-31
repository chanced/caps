package token_test

import (
	"testing"
	"unicode"

	"github.com/chanced/caps/token"
)

func TestReversed(t *testing.T) {
	tok := token.FromString(nil, "abc")
	rev := tok.LowerReversedRunes()
	if string(rev) != "cba" {
		t.Error("expected \"cba\", got", rev)
	}
}

func TestTurkish(t *testing.T) {
	runes := "içğıöşü"
	for _, r := range runes {
		tok := token.FromRune(token.TurkishCaser, r)
		if []rune(tok.Upper())[0] != unicode.TurkishCase.ToUpper(r) {
			t.Errorf("expected %U, got %U", unicode.TurkishCase.ToTitle(r), []rune(tok.Upper())[0])
		}
	}
}

func TestAppend(t *testing.T) {
	var res token.Token
	// tok := token.FromString(nil, "abc")
	// res = token.Append(nil, tok, token.FromString(nil, "def"), token.FromRune(nil, 'g'), token.FromString(nil, "hij"))
	// if res.String() != "abcdefghij" {
	// 	t.Error("expected \"abcdefghij\", got", res)
	// }
	titleDZ := unicode.ToTitle('ǳ')
	upperDZ := unicode.ToUpper('ǳ')
	if unicode.IsTitle(upperDZ) {
		t.Error("expected upperDZ to not be title")
	}
	// fmt.Println("is title to upperDZ", unicode.IsTitle(upperDZ))
	// fmt.Println("is title to titleDZ", unicode.IsTitle(titleDZ))
	// fmt.Println("is title to titleDZ after ToUpper", unicode.IsTitle(unicode.ToUpper(titleDZ)))
	titleDZTok := token.FromRune(nil, titleDZ)
	res = token.Append(nil, titleDZTok, titleDZTok)
	if res.Runes()[0] != titleDZ {
		t.Errorf("expected %U to be title, got %U", titleDZ, res.Runes()[0])
	}
	if res.Runes()[1] != upperDZ {
		t.Errorf("expected %U to be upper, got %U", upperDZ, res.Runes()[1])
	}
}

func TestReverse(t *testing.T) {
	tok := token.FromString(nil, "abc")
	if string(tok.Reverse().Value()) != "cba" {
		t.Error("expected \"cba\", got", tok.Reverse())
	}
}

func TestReverseSplit(t *testing.T) {
	tok := token.FromString(nil, "abc")

	var str string
	for _, rt := range tok.ReverseSplit() {
		str = str + rt.String()
	}
	if str != "cba" {
		t.Error("expected\"cba\" but got", str)
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
		rules    map[rune]func(index int, r rune, val []rune) bool
	}{
		{"123", true, nil},
		{"123.456", true, nil},
		{"123.456e7", true, nil},
		{"123.456e-7", true, nil},
		{"#123", true, nil},
		{"#123.456", true, nil},
		{"#123.456e7", true, nil},
		{"#123.456e-7", true, nil},
		{"123#", false, nil},
		{"0.1", true, nil},
		{"0.1e2", true, nil},
		{"0.1e-2", true, nil},
		{"0.1e+2", true, nil},
		{".0", true, nil},
		{".0e2", true, nil},
		{"0.1e2e3", false, nil},
		{".", false, nil},
		{".0.", false, nil},
		{"0..0", false, nil},
		{"123.", false, nil},
		{"123.456.", false, nil},
		{"123.456e7.", false, nil},
		{"v1", true, nil},
		{"v1.0", true, nil},
		{"v1.0.0", true, nil},
		{"v1.0.0-alpha", false, nil},
		{"$1", false, nil},
		{"v1.0.0-alpha.1", false, nil},
		{"vv1", false, nil},
		{"$1", true, map[rune]func(index int, r rune, val []rune) bool{
			'$': func(index int, _ rune, _ []rune) bool {
				return index == 0
			},
		}},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if token.FromString(nil, test.value).IsNumber(test.rules) != test.expected {
				if test.expected {
					t.Errorf("expected \"%s\" to be a number", test.value)
				} else {
					t.Errorf("expected \"%s\" to not be a number", test.value)
				}
			}
		})
	}
}
