package token_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/chanced/caps/token"
)

func TestToLower(t *testing.T) {
	strs := []string{"ABC", "abc", "", "a", "!!", "Abc"}

	for _, str := range strs {
		t.Run(str, func(t *testing.T) {
			if token.ToLower(token.DefaultCaser, str) != strings.ToLower(str) {
				t.Errorf("expected \"%s\" to be %s", str, strings.ToLower(str))
			}
			if token.ToLower(token.TurkishCaser, str) != strings.ToLowerSpecial(unicode.TurkishCase, str) {
				t.Errorf("expected \"%s\" to be %s", str, strings.ToLowerSpecial(unicode.TurkishCase, str))
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	strs := []string{"ABC", "abc", "", "a", "!!", "Abc"}

	for _, str := range strs {
		t.Run(str, func(t *testing.T) {
			if token.ToUpper(token.DefaultCaser, str) != strings.ToUpper(str) {
				t.Errorf("expected \"%s\" to be %s", str, strings.ToUpper(str))
			}
			if token.ToUpper(token.TurkishCaser, str) != strings.ToUpperSpecial(unicode.TurkishCase, str) {
				t.Errorf("expected \"%s\" to be %s", str, strings.ToUpperSpecial(unicode.TurkishCase, str))
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	if !token.IsEmpty("") {
		t.Error("expected true, got false")
	}
	if token.IsEmpty("a") {
		t.Error("expected false, got true")
	}
}

func TestFirstRune(t *testing.T) {
	f, ok := token.FirstRune("abc")
	if !ok {
		t.Error("expected true, got false")
	}
	if f != 'a' {
		t.Error("expected 'a', got", f)
	}

	if _, ok := token.FirstRune(""); ok {
		t.Error("expected false, got true")
	}
}

func TestUpperFirstLowerRest(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"abc", "Abc"},
		{"Abc", "Abc"},
		{"aBc", "Abc"},
		{"aBC", "Abc"},
		{"aBCD", "Abcd"},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if token.UpperFirstLowerRest(token.DefaultCaser, test.in) != test.out {
				t.Errorf("expected %s, got %s", test.out, token.UpperFirstLowerRest(token.DefaultCaser, test.in))
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "a"},
		{"abc", "abc"},
		{"Abc", "abc"},
		{"aBc", "aBc"},
		{"aBC", "aBC"},
		{"aBCD", "aBCD"},
		{"ABC", "aBC"},
		{"aBCD", "aBCD"},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if token.LowerFirst(token.DefaultCaser, test.in) != test.out {
				t.Errorf("expected %s, got %s", test.out, token.LowerFirst(token.DefaultCaser, test.in))
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "A"},
		{"A", "A"},
		{"abc", "Abc"},
		{"Abc", "Abc"},
		{"aBc", "ABc"},
		{"aBC", "ABC"},
		{"aBCD", "ABCD"},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if token.UpperFirst(token.DefaultCaser, test.in) != test.out {
				t.Errorf("expected %s, got %s", test.out, token.UpperFirst(token.DefaultCaser, test.in))
			}
		})
	}
}

func TestHasLower(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"", false},
		{"a", true},
		{"A", false},
		{"abc", true},
		{"Abc", true},
		{"aBc", true},
		{"aBC", true},
		{"ABC", false},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if token.HasLower(test.in) != test.out {
				t.Errorf("expected %t, got %t", test.out, token.HasLower(test.in))
			}
		})
	}
}

func TestAppendRune(t *testing.T) {
	var res string
	res = token.AppendRune(nil, "", 'a', 'b', 'c')
	if res != "abc" {
		t.Error("expected \"abc\", got", res)
	}
	tDZ := unicode.ToTitle('ǳ')
	uDZ := unicode.ToUpper('ǳ')

	res = token.AppendRune(nil, "", tDZ, uDZ)
	if res != "ǲǱ" {
		t.Error("expected \"ǲǱ\", got", res)
	}
	res = token.AppendRune(nil, "", uDZ, tDZ)
	if res != "ǲǱ" {
		t.Error("expected \"ǲǱ\", got", res)
	}

	res = token.AppendRune(token.DefaultCaser, "123", '.')
	if res != "123." {
		t.Error("expected \"123.\", got", res)
	}
}

func TestTurkish(t *testing.T) {
	runes := "içğıöşü"
	for _, r := range runes {
		res := []rune(token.ToUpper(token.TurkishCaser, string(r)))[0]
		if res != unicode.TurkishCase.ToUpper(r) {
			t.Errorf("expected %U, got %U", unicode.TurkishCase.ToTitle(r), res)
		}
	}
}

func TestAppend(t *testing.T) {
	var res string
	titleDZ := unicode.ToTitle('ǳ')
	upperDZ := unicode.ToUpper('ǳ')
	if unicode.IsTitle(upperDZ) {
		t.Error("expected upperDZ to not be title")
	}
	// fmt.Println("is title to upperDZ", unicode.IsTitle(upperDZ))
	// fmt.Println("is title to titleDZ", unicode.IsTitle(titleDZ))
	// fmt.Println("is title to titleDZ after ToUpper", unicode.IsTitle(unicode.ToUpper(titleDZ)))
	titleDZStr := string(titleDZ)
	res = token.Append(nil, titleDZStr, titleDZStr)
	if []rune(res)[0] != titleDZ {
		t.Errorf("expected %U to be title, got %U", titleDZ, []rune(res)[0])
	}
	if []rune(res)[1] != upperDZ {
		t.Errorf("expected %U to be upper to be upper, got %U", upperDZ, []rune(res)[1])
	}

	res = token.Append(nil, "a", "", "b", "c", "")
	if res != "abc" {
		t.Error("expected \"abc\", got", res)
	}
}

func TestReverse(t *testing.T) {
	tDZ := unicode.ToTitle('ǳ')
	uDZ := unicode.ToUpper('ǳ')
	titleCheck := string([]rune{tDZ, 'ŷ', uDZ})
	tests := []struct {
		in  string
		out string
	}{
		{"", ""},
		{"a", "a"},
		{"A", "A"},
		{"abc", "cba"},
		{"Abc", "cbA"},
		{titleCheck, titleCheck},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if token.Reverse(token.DefaultCaser, test.in) != test.out {
				t.Errorf("expected %s, got %s", test.out, token.Reverse(token.DefaultCaser, test.in))
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
		rules    map[rune]func(index int, r rune, val string) bool
	}{
		{"", false, nil},
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
		{"0.1e", false, nil},
		{"1e", false, nil},
		{"0.1e+2", true, nil},
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
		{"$1", true, map[rune]func(index int, r rune, val string) bool{
			'$': func(index int, _ rune, _ string) bool {
				return index == 0
			},
		}},
		{"$1.045", true, map[rune]func(index int, r rune, val string) bool{
			'$': func(index int, _ rune, _ string) bool {
				return index == 0
			},
			'e': func(index int, _ rune, _ string) bool {
				return false
			},
		}},
		{"$1.0e45", false, map[rune]func(index int, r rune, val string) bool{
			'$': func(index int, _ rune, _ string) bool {
				return index == 0
			},
			'e': func(index int, _ rune, _ string) bool {
				return false
			},
		}},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			if token.IsNumber(test.value, test.rules) != test.expected {
				if test.expected {
					t.Errorf("expected \"%s\" to be a number", test.value)
				} else {
					t.Errorf("expected \"%s\" to not be a number", test.value)
				}
			}
			// if test.rules != nil {
			// 	var rules []map[rune]func(index int, r rune, val string) bool
			// 	for k, v := range test.rules {
			// 		rules = append(rules, map[rune]func(index int, r rune, val string) bool{k: v})
			// 	}
			// 	if token.IsNumber(test.value, rules) != test.expected {
			// 		if test.expected {
			// 			t.Errorf("expected \"%s\" to be a number", test.value)
			// 		} else {
			// 			t.Errorf("expected \"%s\" to not be a number", test.value)
			// 		}
			// 	}
			// }
		})
	}
}
