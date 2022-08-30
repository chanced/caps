package caps

import (
	"sort"
	"strings"
	"unicode"

	"github.com/chanced/caps/index"
	"github.com/chanced/caps/token"
)

const (
	// DEFAULT_DELIMITERS is the default set of delimiters in string convert.
	DEFAULT_DELIMITERS string = " _.!?:;$-(){}[]#@&+~"
)

var (
	// DefaultTokenizer is the default Tokenizer.
	DefaultTokenizer TokenizerImpl = NewTokenizer(DEFAULT_DELIMITERS, token.DefaultCaser)

	// DefaultReplacements is the list of Replacements passed to DefaultConverter.
	// 	{"Acl", "ACL"},
	// 	{"Api", "API"},
	// 	{"Ascii", "ASCII"},
	// 	{"Cpu", "CPU"},
	// 	{"Css", "CSS"},
	// 	{"Dns", "DNS"},
	// 	{"Eof", "EOF"},
	// 	{"Guid", "GUID"},
	// 	{"Html", "HTML"},
	// 	{"Http", "HTTP"},
	// 	{"Https", "HTTPS"},
	// 	{"Id", "ID"},
	// 	{"Ip", "IP"},
	// 	{"Json", "JSON"},
	// 	{"Lhs", "LHS"},
	// 	{"Qps", "QPS"},
	// 	{"Ram", "RAM"},
	// 	{"Rhs", "RHS"},
	// 	{"Rpc", "RPC"},
	// 	{"Sla", "SLA"},
	// 	{"Smtp", "SMTP"},
	// 	{"Sql", "SQL"},
	// 	{"Ssh", "SSH"},
	// 	{"Tcp", "TCP"},
	// 	{"Tls", "TLS"},
	// 	{"Ttl", "TTL"},
	// 	{"Udp", "UDP"},
	// 	{"Ui", "UI"},
	// 	{"Uid", "UID"},
	// 	{"Uuid", "UUID"},
	// 	{"Uri", "URI"},
	// 	{"Url", "URL"},
	// 	{"Utf8", "UTF8"},
	// 	{"Vm", "VM"},
	// 	{"Xml", "XML"},
	// 	{"Xmpp", "XMPP"},
	// 	{"Xsrf", "XSRF"},
	// 	{"Xss", "XSS"},
	DefaultReplacements []Replacement = []Replacement{
		{"Acl", "ACL"},
		{"Api", "API"},
		{"Ascii", "ASCII"},
		{"Cpu", "CPU"},
		{"Css", "CSS"},
		{"Dns", "DNS"},
		{"Eof", "EOF"},
		{"Guid", "GUID"},
		{"Html", "HTML"},
		{"Http", "HTTP"},
		{"Https", "HTTPS"},
		{"Id", "ID"},
		{"Ip", "IP"},
		{"Json", "JSON"},
		{"Lhs", "LHS"},
		{"Qps", "QPS"},
		{"Ram", "RAM"},
		{"Rhs", "RHS"},
		{"Rpc", "RPC"},
		{"Sla", "SLA"},
		{"Smtp", "SMTP"},
		{"Sql", "SQL"},
		{"Ssh", "SSH"},
		{"Tcp", "TCP"},
		{"Tls", "TLS"},
		{"Ttl", "TTL"},
		{"Udp", "UDP"},
		{"Ui", "UI"},
		{"Uid", "UID"},
		{"Uuid", "UUID"},
		{"Uri", "URI"},
		{"Url", "URL"},
		{"Utf8", "UTF8"},
		{"Vm", "VM"},
		{"Xml", "XML"},
		{"Xmpp", "XMPP"},
		{"Xsrf", "XSRF"},
		{"Xss", "XSS"},
	}

	// DefaultConverter is the default Converter instance.
	//
	// # Default parameters:
	//

	//
	// replacements:
	//  DefaultReplacements
	//
	// tokenizer:
	//  DefaultTokenizer
	DefaultConverter = NewConverter(DefaultReplacements, DefaultTokenizer, token.DefaultCaser)
)

// Tokenizer is an interface satisfied by tyeps which can
type Tokenizer interface {
	Tokenize(value string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token
}

// Converter is an interface satisfied by types which can convert the case of a
// string.
//
// ConverterImpl is provided as a default implementation. If you have edge cases which require custom formatting,
// you can implement your own Converter by wrapping ConverterImpl:
//
//	type MyConverter struct {}
//	func(MyConverter) Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string{
//		formatted := caps.DefaultConverter.Convert(style, repStyle, input, join, allowedSymbols, numberRules)
//		if formatted == "something_unusual" {
//		  	return "replaced"
//	 	}
//	 	return formatted
//	}
//
// # Parameters
//
//	style:          Expected output caps.Style of the string.
//	repStyle:       The caps.ReplaceStyle to use if a word needs to be replaced.
//	join:           The delimiter to use when joining the words. For Camel, this is an empty string.
//	allowedSymbols: The set of allowed symbols. If set, these should take precedence over any delimiters
//	numberRules:    Any custom rules dictating how to handle special characters in numbers.
type Converter interface {
	Convert(
		style Style,
		repStyle ReplaceStyle,
		input string,
		join string,
		allowedSymbols []rune,
		numberRules map[rune]func(index int, r rune, val []rune) bool,
	) string
}

// NewTokenizer creates and returns a new TokenizerImpl which implements the
// Tokenizer interface.
//
// Tokenizers are used by ConverterImpl to tokenize the input text into
// token.Tokens that are then formatted.
func NewTokenizer(delimiters string, caser token.Caser) TokenizerImpl {
	d := runes(delimiters)
	sort.Sort(d)
	return TokenizerImpl{
		delimiters: d,
		caser:      token.CaserOrDefault(caser),
	}
}

// TokenizerImpl is the provided implementation of the Tokenizer interface.
//
// TokenizerImpl tokenizes the input text into token.Tokens based on a set of
// delimiters (runes).
//
// If you need custom logic, consider wrapping the logic by implementing
// Tokenizer and then calling a TokenizerImpl's Tokenize method.
//
// # Example:
type TokenizerImpl struct {
	delimiters runes
	caser      token.Caser
}

// Tokenize splits a string into a list of token.Tokens based on the case of each
// rune, it's delimiters, and the specified allowedSymbols.
//
// For example:
//
//	t.Tokenize("ASnakecaseVariable", nil) -> ["A", "Snakecase", "Variable"]
//
// Tokenizer attempts to detect formatting, such that a screaming snakecase str
// will be inferred accordingly.
//
// For example:
//
//	t.Tokenize("A_SCREAMING_SNAKECASE_VARIABLE", nil) -> ["A", "SCREAMING", "SNAKECASE", "VARIABLE"]
//
// If allowedSymbols is not nil, then those symbols will be treated as
// non-delimiters even if they are in the delimiters list.
//
// For example:
//
//	t := caps.token.Newizer("_")
//	t.Tokenize("A_SCREAMING_SNAKECASE_VARIABLE", []rune{'_'}) -> ["A_SCREAMING_SNAKECASE_VARIABLE"]
func (ti TokenizerImpl) Tokenize(str string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token {
	tokens := []token.Token{}
	pending := []token.Token{}
	foundLower := false
	current := token.Token{}
	prevNumber := false
	allowed := newRunes(allowedSymbols)

	for i, r := range str {
		switch {
		case unicode.IsUpper(r):
			if foundLower && current.Len() > 0 {
				tokens = append(tokens, current)
				current = token.Token{}
			}
			current = token.AppendRune(ti.caser, current, r)
			prevNumber = false
		case unicode.IsLower(r):
			if !foundLower && current.Len() > 0 {
				// we have to break up the pending first
				for _, tok := range pending {
					if tok.IsNumber(numberRules) {
						tokens = append(tokens, tok)
					} else {
						tokens = append(tokens, tok.Split()...)
					}
				}
				// need to break up the current token if it isn't a number
				if prevNumber {
					tokens = append(tokens, current)
					current = token.Token{}
				} else {
					split := current.Split()
					// current becomes the last upper token before discovering the lowercase token
					current = split[len(split)-1]

					// all other tokens are added to the token list
					tokens = append(tokens, split[:len(split)-1]...)
				}
			}
			current = token.AppendRune(ti.caser, current, r)
			pending = nil
			foundLower = true
		case unicode.IsNumber(r):
			// if adding the number onto current makes it a valid number
			// then append this rune to current
			if token.AppendRune(ti.caser, current, r).IsNumber(numberRules) {
				current = token.AppendRune(ti.caser, current, r)
			} else {
				// otherwise it is not a number and so we add the current token
				// to the token or pending list depending on whether or not we
				// have found a lowercase rune
				if current.Len() > 0 && foundLower {
					tokens = append(tokens, current)
					current = token.Token{}
				} else if current.Len() > 0 { // otherwise, we have to push the current token into a pending state
					pending = append(pending, current)
				}
				current = token.AppendRune(ti.caser, current, r)
			}
			prevNumber = true
		default:
			if allowed.Contains(r) {
				if current.Len() > 0 {
					if current.IsNumber(numberRules) {
						// this gets a bit tricky because we need to check if adding the
						// rune to the current makes it a number. However, in all
						// default cases, a non-numeric rune must be followed either by
						// a number or an 'e' and a number. as such, we have to check if
						// both this and the next rune (and possibly the rune after
						// that) make it a number.
						n := token.AppendRune(ti.caser, current, r)
						runes := []rune(str)
						if n.IsNumber(numberRules) {
							current = n
						} else if i <= len(runes)-2 && (unicode.IsNumber(runes[i+1]) || unicode.IsLetter(runes[i+1])) || allowed.Contains(runes[i+1]) {
							next := token.AppendRune(ti.caser, n, runes[i+1])
							if next.IsNumber(numberRules) {
								current = n
							} else {
								if foundLower {
									tokens = append(tokens, current)
									current = token.FromRunes(ti.caser, []rune{r})
								} else {
									pending = append(pending, current)
									current = token.FromRunes(ti.caser, []rune{r})
								}
							}
						} else {
							if foundLower {
								tokens = append(tokens, current)
								current = token.FromRunes(ti.caser, []rune{r})
							} else {
								pending = append(pending, current)
								current = token.FromRunes(ti.caser, []rune{r})
							}
						}
					} else {
						current = token.AppendRune(ti.caser, current, r)
					}
				} else {
					current = token.AppendRune(ti.caser, current, r)
				}
			} else if ti.delimiters.Contains(r) {
				if current.Len() > 0 {
					if foundLower {
						tokens = append(tokens, current)
						current = token.Token{}
					} else {
						pending = append(pending, current)
						current = token.Token{}
					}
				}
			}
		}
	}
	if current.Len() > 0 {
		if current.IsNumber(numberRules) {
			if foundLower {
				tokens = append(tokens, current)
			} else {
				pending = append(pending, current)
			}
		} else if !foundLower {
			pending = append(pending, current)
		} else {
			tokens = append(tokens, current)
		}
	}
	if foundLower {
		for _, tok := range pending {
			if tok.IsNumber(numberRules) {
				tokens = append(tokens, tok)
			} else {
				tokens = append(tokens, tok.Split()...)
			}
		}
		return tokens
	} else {
		return pending
	}
}

var _ Tokenizer = TokenizerImpl{}

// ReplaceStyle is used to indicate the case style the text should be transformed to
// when seeking replacement text in a Converter.
//
// When a Replacer is configured, the expected input is:
//
//	caps.Replacement{ Camel: "Json", Screaming: "JSON" }
//
// If the ReplaceStyle equals ReplaceStyleScreaming then an input of "MarshalJson" will return
// "MarshaalJSON" with the above caps.Replacement.
type ReplaceStyle uint8

type (
	Replacement struct {
		// Camel case variant of the word which should be replaced.
		// e.g. "Http"
		Camel string
		// Screaming (all upper case) representation of the word to replace.
		// e.g. "HTTP"
		Screaming string
	}
)

type Style uint8

const (
	StyleLower      Style = iota // The output should be lowercase (e.g. "an_example")
	StyleScreaming               // The output should be screaming (e.g. "AN_EXAMPLE")
	StyleCamel                   // The output should be camel case (e.g. "AnExample")
	StyleLowerCamel              // The output should be lower camel case (e.g. "anExample")
)

const (
	ReplaceStyleNotSpecified = iota
	ReplaceStyleCamel        // Text should be replaced with the Camel variant (e.g. "Json").
	ReplaceStyleScreaming    // Text should be replaced with the screaming variant (e.g. "JSON").
	ReplaceStyleLower        // Text should be replaced with the lowercase variant (e.g. "json").
)

// NewConverter creates a new Converter which is used to convert the input text to the desired output.
//
// replacements are used to make replacements of tokens to the specified
// formatting (e.g. { "Json", "JSON"}).
//
// tokenizer is used to tokenize the input text.
func NewConverter(replacements []Replacement, tokenizer Tokenizer, caser token.Caser) ConverterImpl {
	ci := ConverterImpl{
		index:     index.New(caser),
		tokenizer: tokenizer,
		caser:     token.CaserOrDefault(caser),
	}
	for _, v := range replacements {
		ci.set(v.Camel, v.Screaming)
	}
	return ci
}

// ConverterImpl contains a table of words to their desired replacement. Tokens
// will be compared against the keys of this table to determine if the string
// should be replaced with the value of the table.
//
// This is primarily designed for acronyms but it could be used for other
// purposes.
//
// The default Replacements can be found in the DefaultReplacements variable.
type ConverterImpl struct {
	index     *index.Index
	tokenizer Tokenizer
	caser     token.Caser
}

func (ci ConverterImpl) Index() index.Index {
	return *ci.index
}

// Contains reports whether a key is in the Converter's replacement table.
func (ci ConverterImpl) Contains(key string) bool {
	return ci.index.ContainsForward(token.FromString(ci.caser, key))
}

// Replacements returns a slice of Replacement in the lookup trie.
func (ci ConverterImpl) Replacements() []Replacement {
	indexedVals := ci.index.Values()
	res := make([]Replacement, len(indexedVals))
	for i, v := range indexedVals {
		res[i] = Replacement{
			Camel:     v.Camel.Value(),
			Screaming: v.Screaming.Value(),
		}
	}
	return res
}

func (ci *ConverterImpl) set(key, value string) {
	ci.index.Add(token.FromString(ci.caser, key), token.FromString(ci.caser, value))
}

func lowerAndCheck(input string) (string, bool) {
	bldr := strings.Builder{}
	bldr.Grow(len(input))
	foundLower := false
	for _, r := range input {
		if !foundLower && unicode.IsLower(r) {
			foundLower = true
		}
		bldr.WriteRune(r)
	}
	return bldr.String(), foundLower
}

// Set adds the key/value pair to the table.
func (ci *ConverterImpl) Set(key, value string) {
	kstr, keyHasLower := lowerAndCheck(key)
	vstr, valueHasLower := lowerAndCheck(value)
	ci.Delete(kstr)
	ci.Delete(vstr)

	// checking to see if we need to swap these.
	if !keyHasLower && valueHasLower {
		ci.set(value, key)
	} else {
		ci.set(key, value)
	}
}

// Remove deletes the key from the map. Either variant is sufficient.
func (ci *ConverterImpl) Delete(key string) {
	tok := token.FromString(ci.caser, key)
	ci.index.Delete(tok)
}

// FormatToken formats the token with the desired style.
func FormatToken(style Style, index int, tok token.Token) string {
	switch style {
	case StyleCamel:
		return string(tok.UpperFirstLowerRest())
	case StyleLowerCamel:
		if index == 0 {
			return string(tok.Lower())
		}
		return string(tok.UpperFirstLowerRest())
	case StyleScreaming:
		return tok.Upper()
	case StyleLower:
		return tok.Lower()
	}
	return tok.String()
}

func FormatIndexedReplacement(style Style, replaceStyle ReplaceStyle, index int, rep index.IndexedReplacement) string {
	switch replaceStyle {
	case ReplaceStyleCamel:
		if index == 0 && style == StyleLowerCamel {
			return rep.Camel.Lower()
		}
		return rep.Camel.String()
	case ReplaceStyleScreaming:
		return rep.Screaming.String()
	case ReplaceStyleLower:
		return rep.Lower.String()
	default:
		return rep.Screaming.String()
	}
}

// Convert formats the string with the desired style.
func (ci ConverterImpl) Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string {
	tokens := ci.tokenizer.Tokenize(input, allowedSymbols, numberRules)
	var parts []string
	var ok bool
	idx := ci.Index()
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		switch tok.Len() {
		case 0:
			continue
		case 1:
			if idx, ok = idx.MatchReverse(tok); !ok {
				if idx.LastMatch().HasValue() {
					parts = append(parts, FormatIndexedReplacement(style, repStyle, i+1, idx.LastMatch()))
				}
				if idx.HasPartialMatches() {
					for _, partok := range idx.PartialMatches() {
						parts = append(parts, FormatToken(style, i+1, partok))
					}
				}
				parts = append(parts, FormatToken(style, i, tok))

				// resetting the index
				idx = ci.Index()
			}
		default:
			if idx.HasMatch() {
				parts = append(parts, FormatIndexedReplacement(style, repStyle, i+1, idx.LastMatch()))
			}
			if idx.HasPartialMatches() {
				for _, partok := range idx.PartialMatches() {
					parts = append(parts, FormatToken(style, i+1, partok))
				}
			}
			if idx.HasMatch() || idx.HasPartialMatches() {
				// resetting index
				idx = ci.Index()
			}
			if rep, ok := idx.GetForward(tok); ok {
				parts = append(parts, FormatIndexedReplacement(style, repStyle, i, rep))
			} else {
				parts = append(parts, FormatToken(style, i, tok))
			}
		}
	}
	result := strings.Builder{}
	result.Grow(len(input))

	var part string
	shouldWriteDelimiter := false
	if idx.HasPartialMatches() {
		i := 0
		if idx.HasMatch() {
			i = 1
		}
		for y, partok := range idx.PartialMatches() {
			parts = append(parts, FormatToken(style, i+y, partok))
		}
	}
	if idx.HasMatch() {
		parts = append(parts, FormatIndexedReplacement(style, repStyle, 0, idx.LastMatch()))
	}
	for i := len(parts) - 1; i >= 0; i-- {
		part = parts[i]
		if shouldWriteDelimiter {
			result.WriteString(join)
		}
		result.WriteString(part)
		if !shouldWriteDelimiter {
			shouldWriteDelimiter = len(part) > 0 && len(join) > 0
		}
	}

	return result.String()
}

// Opts include configurable options for case conversion.
//
// See the documentation for the individual fields for more information.
type Opts struct {
	// Any characters within this string will be allowed in the output.
	//
	// Default:
	// 	""
	AllowedSymbols string
	// The Converter to use.
	//
	// Default:
	// 	DefaultConverter
	Converter Converter

	// Styles overwrites the way words are replaced.
	//
	// A typical call to ToLowerCamel for "ServeJSON" with a Replacer that
	// contains {"Json": "JSON"} would result in "serveJSON" by using the
	// StyleScreaming variant. If ReplacementStyle was set to
	// ReplaceStyleUpperCamel, on the call to ToLowerCamel then the result would
	// be "serveHttp".
	//
	// The default replacement style is dependent upon the target casing.
	ReplaceStyle ReplaceStyle
	// NumberRules are used by the DefaultTokenizer to augment the standard
	// rules for determining if a rune is part of a number.
	//
	// Note, if you add special characters here, they must be present in the
	// AllowedSymbols string for them to be part of the output.
	NumberRules map[rune]func(index int, r rune, val []rune) bool
}

func loadOpts(opts []Opts) Opts {
	result := Opts{
		AllowedSymbols: "",
		Converter:      DefaultConverter,
		ReplaceStyle:   ReplaceStyleScreaming,
	}
	if len(opts) == 0 {
		return result
	}

	if opts[0].AllowedSymbols != "" {
		result.AllowedSymbols = opts[0].AllowedSymbols
	}
	if opts[0].Converter != nil {
		result.Converter = opts[0].Converter
	}
	if opts[0].ReplaceStyle != ReplaceStyleNotSpecified {
		result.ReplaceStyle = opts[0].ReplaceStyle
	}
	if opts[0].NumberRules != nil {
		result.NumberRules = opts[0].NumberRules
	}

	return result
}

// UpperFirst converts the first rune of str to unicode upper case.
//
// This method does not support special cases (such as Turkish and Azeri)
func UpperFirst[T ~string](str T) T {
	runes := []rune(str)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToTitle(runes[0])
	return T(runes)
}

// LowerFirst converts the first rune of str to lowercase.
func LowerFirst[T ~string](str T) T {
	runes := []rune(str)
	switch len(runes) {
	case 0:
		return ""
	case 1:
		return T(unicode.ToLower(runes[0]))
	default:
		runes[0] = unicode.ToLower(runes[0])
		return T(runes)
	}
}

// Without numbers returns the string with all numeric runes removed.
//
// It does not currently use any logic to determine if a rune (e.g. '.')
// is part of a number. This may change in the future.
func WithoutNumbers[T ~string](s T) T {
	return T(strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, string(s)))
}

// ToCamel transforms the case of str into Camel Case (e.g. AnExampleString) using
// either the provided Converter or the DefaultConverter otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes "AnExampleString".
// It also has a configurable set of replacements, such that "some_json" becomes "SomeJSON"
// so long as opts.ReplacementStyle is set to ReplaceStyleScreaming. A ReplaceStyle of
// ReplaceStyleCamel would result in "SomeJson".
//
//	caps.ToCamel("This is [an] {example}${id32}.") // ThisIsAnExampleID32
//	caps.ToCamel("AN_EXAMPLE_STRING", ) // AnExampleString
func ToCamel[T ~string](str T, options ...Opts) T {
	opts := loadOpts(options)
	return T(opts.Converter.Convert(StyleCamel, opts.ReplaceStyle, string(str), "", []rune(opts.AllowedSymbols), opts.NumberRules))
}

// ToLowerCamel transforms the case of str into Lower Camel Case (e.g. anExampleString) using
// either the provided Converter or the DefaultConverter otherwise.
//
// The default Converter detects case so that "AN_EXAMPLE_STRING" becomes "anExampleString".
// It also has a configurable set of replacements, such that "some_json" becomes "someJSON"
// so long as opts.ReplacementStyle is set to ReplaceStyleScreaming. A ReplaceStyle of
// ReplaceStyleCamel would result in "someJson".
//
//	caps.ToLowerCamel("This is [an] {example}${id32}.") // thisIsAnExampleID32
func ToLowerCamel[T ~string](str T, options ...Opts) T {
	opts := loadOpts(options)
	return T(opts.Converter.Convert(StyleLowerCamel, opts.ReplaceStyle, string(str), "", []rune(opts.AllowedSymbols), opts.NumberRules))
}

// ToSnake transforms the case of str into Lower Snake Case (e.g. an_example_string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToSnake("This is [an] {example}${id32}.") // this_is_an_example_id_32
func ToSnake[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '_', true, options...)
}

// ToScreamingSnake transforms the case of str into Screaming Snake Case (e.g.
// AN_EXAMPLE_STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingSnake("This is [an] {example}${id32}.") // THIS_IS_AN_EXAMPLE_ID_32
func ToScreamingSnake[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '_', false, options...)
}

// ToKebab transforms the case of str into Lower Kebab Case (e.g. an-example-string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToKebab("This is [an] {example}${id32}.") // this-is-an-example-id-32
func ToKebab[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '-', true, options...)
}

// ToScreamingKebab transforms the case of str into Screaming Kebab Snake (e.g.
// AN-EXAMPLE-STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingKebab("This is [an] {example}${id32}.") // THIS-IS-AN-EXAMPLE-ID-32
func ToScreamingKebab[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '-', false, options...)
}

// ToDot transforms the case of str into Lower Dot Notation Case (e.g. an.example.string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToDot("This is [an] {example}${id32}.") // this.is.an.example.id.32
func ToDot[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '.', true, options...)
}

// ToScreamingKebab transforms the case of str into Screaming Kebab Case (e.g.
// AN-EXAMPLE-STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingDot("This is [an] {example}${id32}.") // THIS.IS.AN.EXAMPLE.ID.32
func ToScreamingDot[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '.', false, options...)
}

// ToTitle transforms the case of str into Title Case (e.g. An Example String) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToTitle("This is [an] {example}${id32}.") // This Is An Example ID 32
func ToTitle[T ~string](str T, options ...Opts) T {
	opts := loadOpts(options)
	return T(opts.Converter.Convert(StyleCamel, opts.ReplaceStyle, string(str), " ", []rune(opts.AllowedSymbols), opts.NumberRules))
}

// ToDelimited transforms the case of str into a string separated by delimiter,
// using either the provided Converter or the DefaultConverter otherwise.
//
// If lowercase is false, the output will be all uppercase.
//
// # See Options for more information on available configuration
//
// # Example
//
//	caps.ToDelimited("This is [an] {example}${id}.#32", '.', true) // this.is.an.example.id.32
//	caps.ToDelimited("This is [an] {example}${id}.break32", '.', false) // THIS.IS.AN.EXAMPLE.ID.BREAK.32
//	caps.ToDelimited("This is [an] {example}${id}.v32", '.', true, caps.Opts{AllowedSymbols: "$"}) // this.is.an.example.$.id.v32
func ToDelimited[T ~string](str T, delimiter rune, lowercase bool, options ...Opts) T {
	opts := loadOpts(options)
	var style Style
	var replacementStyle ReplaceStyle
	if lowercase {
		style = StyleLower
		replacementStyle = ReplaceStyleLower
	} else {
		style = StyleScreaming
		replacementStyle = ReplaceStyleScreaming
	}
	return T(opts.Converter.Convert(style, replacementStyle, string(str), string(delimiter), []rune(opts.AllowedSymbols), opts.NumberRules))
}

// ToLower returns s with all Unicode letters mapped to their lower case.
func ToLower[T ~string](str T) T {
	return T(strings.ToLower((string(str))))
}

// ToUpper returns s with all Unicode letters mapped to their upper case.
func ToUpper[T ~string](str T) T {
	return T(strings.ToUpper((string(str))))
}

type runes []rune

// Len implements sort.Interface
func (r runes) Len() int {
	return len(r)
}

// Less implements sort.Interface
func (r runes) Less(i int, j int) bool {
	return r[i] < r[j]
}

// Swap implements sort.Interface
func (r runes) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r runes) Contains(c rune) bool {
	if len(r) == 0 {
		return false
	}
	res := sort.Search(len(r), func(i int) bool {
		return r[i] >= c
	})
	return res > -1 && res < len(r) && r[res] == c
}

func newRunes(val []rune) runes {
	if len(val) == 0 {
		return nil
	}
	r := runes(val)
	sort.Sort(r)
	return r
}

type resolvedReplacement struct {
	resolved       token.Token
	partialMatches []token.Token
}

var _ sort.Interface = (*runes)(nil)
