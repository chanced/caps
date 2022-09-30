package text

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/chanced/caps"
	"github.com/chanced/caps/token"
)

func TestTexts_Contains(t *testing.T) {
	type args struct {
		val Text
	}
	tests := []struct {
		name string
		tr   Texts
		args args
		want bool
	}{
		{"contains", Texts{"a", "b"}, args{"a"}, true},
		{"not contains", Texts{"a", "b"}, args{"c"}, false},
		{"not contains case sensitive", Texts{"a", "b"}, args{"A"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Contains(tt.args.val); got != tt.want {
				t.Errorf("Texts.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexts_ContainsFold(t *testing.T) {
	type args struct {
		val Text
	}
	tests := []struct {
		name string
		tr   Texts
		args args
		want bool
	}{
		{"contains", Texts{"a", "b"}, args{"a"}, true},
		{"not contains", Texts{"a", "b"}, args{"c"}, false},
		{"contains case insensitive", Texts{"a", "b"}, args{"A"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ContainsFold(tt.args.val); got != tt.want {
				t.Errorf("Texts.ContainsFold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexts_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		tr   Texts
		args args
		want bool
	}{
		{"lt", Texts{"a", "b"}, args{0, 1}, true},
		{"gt", Texts{"b", "a"}, args{0, 1}, false},
		{"eq", Texts{"b", "b"}, args{0, 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Texts.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexts_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		tr   Texts
		args args
		want Texts
	}{
		{"two", Texts{"a", "b"}, args{0, 1}, Texts{"b", "a"}},
		{"multi", Texts{"b", "b", "c", "d"}, args{0, 1}, Texts{"b", "a", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.Swap(tt.args.i, tt.args.j)
		})
	}
}

func TestTexts_Len(t *testing.T) {
	tests := []struct {
		name string
		tr   Texts
		want int
	}{
		{"empty", Texts{}, 0},
		{"one", Texts{"a"}, 1},
		{"two", Texts{"a", "b"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Len(); got != tt.want {
				t.Errorf("Texts.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexts_TotalLen(t *testing.T) {
	tests := []struct {
		name string
		tr   Texts
		want int
	}{
		{"empty", Texts{}, 0},
		{"one", Texts{"a"}, 1},
		{"two", Texts{"a", "b"}, 2},
		{"alpha", Texts{"abc", "def", "ghi", "jkl", "mno", "pqr", "stu", "vwx", "yz"}, 26},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TotalLen(); got != tt.want {
				t.Errorf("Texts.TotalLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTexts_Join(t *testing.T) {
	type args struct {
		sep Text
	}
	tests := []struct {
		name string
		tr   Texts
		args args
		want Text
	}{
		{"empty", Texts{}, args{""}, ""},
		{"one", Texts{"a"}, args{""}, "a"},
		{"two", Texts{"a", "b"}, args{""}, "ab"},
		{"sep", Texts{"a", "b"}, args{" "}, "a b"},
		{"multi", Texts{"a", "b", "c", "d"}, args{" "}, "a b c d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Join(tt.args.sep); got != tt.want {
				t.Errorf("Texts.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToLower(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		want Text
	}{
		{"empty", "", ""},
		{"one_lower", "a", "a"},
		{"one_upper", "A", "a"},
		{"two_lower", "ab", "ab"},
		{"two_upper", "AB", "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToLower(); got != tt.want {
				t.Errorf("Text.ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToUpper(t *testing.T) {
	type args struct {
		caser []token.Caser
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "ab", args{}, "AB"},
		{"two_upper", "AB", args{}, "AB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToUpper(tt.args.caser...); got != tt.want {
				t.Errorf("Text.ToUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToCamel(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "ab", args{}, "Ab"},
		{"two_upper", "AB", args{}, "Ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToCamel(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToLowerCamel(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "a"},
		{"one_upper", "A", args{}, "a"},
		{"two_lower", "a_b", args{}, "aB"},
		{"two_upper", "A_B", args{}, "aB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToLowerCamel(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToLowerCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToSnake(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "a"},
		{"one_upper", "A", args{}, "a"},
		{"two_lower", "aB", args{}, "a_b"},
		{"two_upper", "AaBb", args{}, "aa_bb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToSnake(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToScreamingSnake(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "aB", args{}, "A_B"},
		{"two_upper", "AaBb", args{}, "AA_BB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToScreamingSnake(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToScreamingSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToKebab(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "a"},
		{"one_upper", "A", args{}, "a"},
		{"two_lower", "aB", args{}, "a-b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToKebab(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToKebab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToScreamingKebab(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "aB", args{}, "A-B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToScreamingKebab(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToScreamingKebab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToDotNotation(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "a"},
		{"one_upper", "A", args{}, "a"},
		{"two_lower", "aB", args{}, "a.b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToDotNotation(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToDotNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToScreamingDotNotation(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "aB", args{}, "A.B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToScreamingDotNotation(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToScreamingDotNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToTitle(t *testing.T) {
	type args struct {
		opts []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{}, ""},
		{"one_lower", "a", args{}, "A"},
		{"one_upper", "A", args{}, "A"},
		{"two_lower", "aB", args{}, "A B"},
		{"words", "a b", args{}, "A B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToTitle(tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToDelimited(t *testing.T) {
	type args struct {
		delimiter Text
		lowercase bool
		opts      []caps.Opts
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{",", true, nil}, ""},
		{"one_lower", "a", args{",", true, nil}, "a"},
		{"two_lower", "aB", args{",", true, nil}, "a,b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToDelimited(tt.args.delimiter, tt.args.lowercase, tt.args.opts...); got != tt.want {
				t.Errorf("Text.ToDelimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ReplaceAll(t *testing.T) {
	type args struct {
		old Text
		new Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{"", ""}, ""},
		{"one", "a", args{"a", "b"}, "b"},
		{"sentence", "the quick brown fox", args{"fox", "hound"}, "the quick brown hound"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ReplaceAll(tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("Text.ReplaceAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Replace(t *testing.T) {
	type args struct {
		old Text
		new Text
		n   int
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{"", "", 1}, ""},
		{"one", "a", args{"a", "b", 1}, "b"},
		{"sentence", "the quick fox fox", args{"fox", "brown", 1}, "the quick brown fox"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Replace(tt.args.old, tt.args.new, tt.args.n); got != tt.want {
				t.Errorf("Text.Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Compare(t *testing.T) {
	type args struct {
		other Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{""}, 0},
		{"one", "a", args{"a"}, 0},
		{"two", "a", args{"b"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Compare(tt.args.other); got != tt.want {
				t.Errorf("Text.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Trim(t *testing.T) {
	type args struct {
		cutset Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{" "}, ""},
		{"slash", "/a/", args{"/"}, "a"},
		{"space", "  a   ", args{" "}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Trim(tt.args.cutset); got != tt.want {
				t.Errorf("Text.Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimLeft(t *testing.T) {
	type args struct {
		cutset Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{""}, ""},
		{"space", "  a  ", args{" "}, "a  "},
		{"slash", "/a/", args{"/"}, "a/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimLeft(tt.args.cutset); got != tt.want {
				t.Errorf("Text.TrimLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimRight(t *testing.T) {
	type args struct {
		cutset Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{""}, ""},
		{"space", "  a  ", args{" "}, "  a"},
		{"slash", "/a/", args{"/"}, "/a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimRight(tt.args.cutset); got != tt.want {
				t.Errorf("Text.TrimRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimSpace(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		want Text
	}{
		{"empty", "", ""},
		{"space", "  a  ", "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimSpace(); got != tt.want {
				t.Errorf("Text.TrimSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimPrefix(t *testing.T) {
	type args struct {
		prefix Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{""}, ""},
		{"space", "  a  ", args{" "}, " a  "},
		{"slash", "/a/", args{"/"}, "a/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimPrefix(tt.args.prefix); got != tt.want {
				t.Errorf("Text.TrimPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimSuffix(t *testing.T) {
	type args struct {
		suffix Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{""}, ""},
		{"space", "  a  ", args{" "}, "  a "},
		{"slash", "/a/", args{"/"}, "/a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimSuffix(tt.args.suffix); got != tt.want {
				t.Errorf("Text.TrimSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimLeftFunc(t *testing.T) {
	type args struct {
		f func(r rune) bool
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{func(r rune) bool { return r == ' ' }}, ""},
		{"space", "  a  ", args{func(r rune) bool { return r == ' ' }}, "a  "},
		{"slash", "/a/", args{func(r rune) bool { return r == '/' }}, "a/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimLeftFunc(tt.args.f); got != tt.want {
				t.Errorf("Text.TrimLeftFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_TrimRightFunc(t *testing.T) {
	type args struct {
		f func(r rune) bool
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{func(r rune) bool { return r == ' ' }}, ""},
		{"space", "  a  ", args{func(r rune) bool { return r == ' ' }}, "  a"},
		{"slash", "/a/", args{func(r rune) bool { return r == '/' }}, "/a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.TrimRightFunc(tt.args.f); got != tt.want {
				t.Errorf("Text.TrimRightFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_EqualFold(t *testing.T) {
	type args struct {
		v Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{""}, true},
		{"lower_to_upper", "a", args{"A"}, true},
		{"emoji", "👍", args{"👍"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.EqualFold(tt.args.v); got != tt.want {
				t.Errorf("Text.EqualFold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Index(t *testing.T) {
	type args struct {
		substr Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{""}, 0},
		{"space", "  a  ", args{"a"}, 2},
		{"slash", "/a/", args{"/"}, 0},
		{"not_found", "/a/", args{"b"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Index(tt.args.substr); got != tt.want {
				t.Errorf("Text.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_IndexByte(t *testing.T) {
	type args struct {
		c byte
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{' '}, -1},
		{"exists", "abcdefg", args{'a'}, 0},
		{"duplicate", "abcbefg", args{'b'}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IndexByte(tt.args.c); got != tt.want {
				t.Errorf("Text.IndexByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_IndexRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{' '}, -1},
		{"exists", "abcdefg", args{'a'}, 0},
		{"duplicate", "abcbefg", args{'b'}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IndexRune(tt.args.r); got != tt.want {
				t.Errorf("Text.IndexRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_IndexFunc(t *testing.T) {
	type args struct {
		fn func(r rune) bool
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{func(r rune) bool { return r == ' ' }}, -1},
		{"exists", "abcdefg", args{func(r rune) bool { return r == 'a' }}, 0},
		{"duplicate", "abcbefg", args{func(r rune) bool { return r == 'b' }}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.IndexFunc(tt.args.fn); got != tt.want {
				t.Errorf("Text.IndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Cut(t *testing.T) {
	type args struct {
		sep Text
	}
	tests := []struct {
		name       string
		tr         Text
		args       args
		wantBefore Text
		wantAfter  Text
		wantFound  bool
	}{
		{"empty", "", args{""}, "", "", true},
		{"empty_sep", "a", args{""}, "", "a", true},
		{"slash", "a/b", args{"/"}, "a", "b", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBefore, gotAfter, gotFound := tt.tr.Cut(tt.args.sep)
			if gotBefore != tt.wantBefore {
				t.Errorf("Text.Cut() gotBefore = %v, want %v", gotBefore, tt.wantBefore)
			}
			if gotAfter != tt.wantAfter {
				t.Errorf("Text.Cut() gotAfter = %v, want %v", gotAfter, tt.wantAfter)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Text.Cut() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestText_Clone(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		want Text
	}{
		{"empty", "", ""},
		{"a", "a", "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Clone(); got != tt.want {
				t.Errorf("Text.Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Contains(t *testing.T) {
	type args struct {
		substr Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{""}, true},
		{"empty_substr", "a", args{""}, true},
		{"exists", "abcdefg", args{"a"}, true},
		{"not_exists", "abcdefg", args{"z"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Contains(tt.args.substr); got != tt.want {
				t.Errorf("Text.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ContainsAny(t *testing.T) {
	type args struct {
		chars Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{""}, false},
		{"empty_chars", "a", args{""}, false},
		{"exists", "abcdefg", args{"a"}, true},
		{"not_exists", "abcdefg", args{"z"}, false},
		{"multi-check", "abcdefg", args{"az"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ContainsAny(tt.args.chars); got != tt.want {
				t.Errorf("Text.ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ContainsRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{'a'}, false},
		{"exists", "abcdefg", args{'a'}, true},
		{"not_exists", "abcdefg", args{'z'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ContainsRune(tt.args.r); got != tt.want {
				t.Errorf("Text.ContainsRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Count(t *testing.T) {
	type args struct {
		substr Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{""}, 1},
		{"empty_substr", "a", args{"a"}, 1},
		{"exists", "abcdefg", args{"a"}, 1},
		{"not_exists", "abcdefg", args{"z"}, 0},
		{"multi", "abcabc", args{"abc"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Count(tt.args.substr); got != tt.want {
				t.Errorf("Text.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Fields(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		want Texts
	}{
		{"empty", "", Texts{}},
		{"a", "a", Texts{"a"}},
		{"a b", "a b", Texts{"a", "b"}},
		{"a b ", "a b ", Texts{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Fields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_FieldsFunc(t *testing.T) {
	type args struct {
		f func(rune) bool
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Texts
	}{
		{"empty", "", args{func(r rune) bool { return r == ' ' }}, Texts{}},
		{"a", "a", args{func(r rune) bool { return r == ' ' }}, Texts{"a"}},
		{"a b", "a b", args{func(r rune) bool { return r == ' ' }}, Texts{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.FieldsFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.FieldsFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_HasPrefix(t *testing.T) {
	type args struct {
		prefix Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{""}, true},
		{"empty_prefix", "a", args{""}, true},
		{"exists", "abcdefg", args{"a"}, true},
		{"not_exists", "abcdefg", args{"z"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.HasPrefix(tt.args.prefix); got != tt.want {
				t.Errorf("Text.HasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_HasSuffix(t *testing.T) {
	type args struct {
		suffix Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want bool
	}{
		{"empty", "", args{""}, true},
		{"empty_suffix", "a", args{""}, true},
		{"exists", "abcdefg", args{"g"}, true},
		{"not_exists", "abcdefg", args{"z"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.HasSuffix(tt.args.suffix); got != tt.want {
				t.Errorf("Text.HasSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Append(t *testing.T) {
	type args struct {
		elems []Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{[]Text{}}, ""},
		{"empty_elems", "a", args{[]Text{}}, "a"},
		{"append_single", "a", args{[]Text{"b"}}, "ab"},
		{"append_multi", "a", args{[]Text{"b", "c"}}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Append(tt.args.elems...); got != tt.want {
				t.Errorf("Text.Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_AppendRune(t *testing.T) {
	type args struct {
		elems []rune
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{[]rune{}}, ""},
		{"empty_elems", "a", args{[]rune{}}, "a"},
		{"append_single", "a", args{[]rune{'b'}}, "ab"},
		{"append_multi", "a", args{[]rune{'b', 'c'}}, "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.AppendRune(tt.args.elems...); got != tt.want {
				t.Errorf("Text.AppendRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_LastIndex(t *testing.T) {
	type args struct {
		substr Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{""}, 0},
		{"empty_substr", "a", args{""}, 1},
		{"exists", "abcdefg", args{"a"}, 0},
		{"not_exists", "abcdefg", args{"z"}, -1},
		{"multi", "abcabc", args{"a"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.LastIndex(tt.args.substr); got != tt.want {
				t.Errorf("Text.LastIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_LastIndexAny(t *testing.T) {
	type args struct {
		chars Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{""}, -1},
		{"empty_chars", "a", args{""}, -1},
		{"exists", "abcdefg", args{"a"}, 0},
		{"not_exists", "abcdefg", args{"z"}, -1},
		{"multi", "abcabc", args{"az"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.LastIndexAny(tt.args.chars); got != tt.want {
				t.Errorf("Text.LastIndexAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_LastIndexByte(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{0}, -1},
		{"exists", "abcdefg", args{'a'}, 0},
		{"not_exists", "abcdefg", args{'z'}, -1},
		{"multi", "abcabc", args{'a'}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.LastIndexByte(tt.args.b); got != tt.want {
				t.Errorf("Text.LastIndexByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_LastIndexRune(t *testing.T) {
	type args struct {
		r rune
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{0}, -1},
		{"exists", "abcdefg", args{'a'}, 0},
		{"not_exists", "abcdefg", args{'z'}, -1},
		{"multi", "abcabc", args{'a'}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.LastIndexRune(tt.args.r); got != tt.want {
				t.Errorf("Text.LastIndexRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_LastIndexFunc(t *testing.T) {
	type args struct {
		f func(rune) bool
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want int
	}{
		{"empty", "", args{func(r rune) bool { return false }}, -1},
		{"exists", "abcdefg", args{func(r rune) bool { return r == 'a' }}, 0},
		{"not_exists", "abcdefg", args{func(r rune) bool { return r == 'z' }}, -1},
		{"multi", "abcabc", args{func(r rune) bool { return r == 'a' }}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.LastIndexFunc(tt.args.f); got != tt.want {
				t.Errorf("Text.LastIndexFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Map(t *testing.T) {
	type args struct {
		f func(r rune) rune
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{func(r rune) rune { return r }}, ""},
		{"exists", "abcdefg", args{func(r rune) rune { return r + 1 }}, "bcdefgh"},
		{"not_exists", "abcdefg", args{func(r rune) rune { return r + 1 }}, "bcdefgh"},
		{"multi", "abcabc", args{func(r rune) rune { return r + 1 }}, "bcdbcd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Map(tt.args.f); got != tt.want {
				t.Errorf("Text.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Repeat(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{0}, ""},
		{"one", "abc", args{1}, "abc"},
		{"two", "abc", args{2}, "abcabc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Repeat(tt.args.count); got != tt.want {
				t.Errorf("Text.Repeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Split(t *testing.T) {
	type args struct {
		sep Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Texts
	}{
		{"empty", "", args{""}, Texts{}},
		{"one", "abc", args{""}, Texts{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Split(tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_SplitAfter(t *testing.T) {
	type args struct {
		sep Text
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Texts
	}{
		{"empty", "", args{""}, Texts{}},
		{"one", "abc", args{""}, Texts{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.SplitAfter(tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.SplitAfter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_SplitAfterN(t *testing.T) {
	type args struct {
		sep Text
		n   int
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Texts
	}{
		{"empty", "", args{"", 0}, Texts{}},
		{"zero", "abc", args{"", 0}, Texts{}},
		{"one_empty", "abc", args{"", 1}, Texts{"abc"}},
		{"one", "a/b/c", args{"/", 2}, Texts{"a/", "b/c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.SplitAfterN(tt.args.sep, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.SplitAfterN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_SplitN(t *testing.T) {
	type args struct {
		sep Text
		n   int
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Texts
	}{
		{"empty", "", args{"", 0}, Texts{}},
		{"zero", "abc", args{"", 0}, Texts{}},
		{"one_empty", "abc", args{"", 1}, Texts{"abc"}},
		{"one", "a/b/c", args{"/", 2}, Texts{"a", "b/c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.SplitN(tt.args.sep, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text.SplitN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_Title(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		want Text
	}{
		{"empty", "", ""},
		{"one", "abc", "Abc"},
		{"two", "abc abc", "Abc Abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Title(); got != tt.want {
				t.Errorf("Text.Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToLowerSpecial(t *testing.T) {
	type args struct {
		c unicode.SpecialCase
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{unicode.TurkishCase}, ""},
		{"turkish", "İÇĞIÖŞÜ", args{unicode.TurkishCase}, "içğıöşü"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToLowerSpecial(tt.args.c); got != tt.want {
				t.Errorf("Text.ToLowerSpecial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToUpperSpecial(t *testing.T) {
	type args struct {
		c unicode.SpecialCase
	}
	tests := []struct {
		name string
		tr   Text
		args args
		want Text
	}{
		{"empty", "", args{unicode.TurkishCase}, ""},
		{"turkish", "içğıöşü", args{unicode.TurkishCase}, "İÇĞIÖŞÜ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToUpperSpecial(tt.args.c); got != tt.want {
				t.Errorf("Text.ToUpperSpecial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_ToValidUTF8(t *testing.T) {
	tests := []struct {
		name string
		tr   Text
		repl Text
		want Text
	}{
		{"empty", "", "", ""},
		{"empty_with_repl", "", "\uFDDD", ""},
		{"chinese", "a☺\xffb☺\xC0\xAFc☺\xff", "日本語", "a☺日本語b☺日本語c☺日本語"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ToValidUTF8(tt.repl); got != tt.want {
				t.Errorf("Text.ToValidUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}
