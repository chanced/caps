package caps

import (
	"sort"
	"strings"
	"unicode"

	"github.com/chanced/caps/token"
)

const (
	DEFAULT_DELIMITERS string = " _.!?:;$-(){}[]#@&+~"
)

var (
	DefaultTokenizer Tokenizer = NewTokenizer(DEFAULT_DELIMITERS)

	DefaultReplacements []Replacement = []Replacement{
		{"Http", "HTTP"},
		{"Https", "HTTPS"},
		{"Id", "ID"},
		{"Ip", "IP"},
		{"Html", "HTML"},
		{"Xml", "XML"},
		{"Json", "JSON"},
		{"Csv", "CSV"},
		{"Aws", "AWS"},
		{"Gcp", "GCP"},
		{"Sql", "SQL"},
	}

	// Formatter is the default FormatterInterface.
	//
	// # Default parameters:
	//
	// delimiters:
	//   - " _.!?:;$-(){}[]#@&+~"
	//
	// replacements:
	//  { UpperCamel: "Http",  Screaming: "HTTP"  },
	//  { UpperCamel: "Https", Screaming: "HTTPS" },
	//  { UpperCamel: "Html",  Screaming: "HTML"  },
	//  { UpperCamel: "Xml",   Screaming: "XML"   },
	//  { UpperCamel: "Json",  Screaming: "JSON"  },
	//  { UpperCamel: "Csv",   Screaming: "CSV"   },
	//  { UpperCamel: "Aws",   Screaming: "AWS"   },
	//  { UpperCamel: "Gcp",   Screaming: "GCP"   },
	//  { UpperCamel: "Sql",   Screaming: "SQL"   },
	//
	// tokenizer:
	//  DefaultTokenizer
	DefaultFormatter = NewFormatter(DefaultReplacements, DefaultTokenizer)
)

type (
	TokenizerInterface interface {
		Tokenize(value string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token
	}

	// FormatterInterface is an interface satisfied by types which can finalize the case formatting of a string.
	//
	// # Parameters
	//
	//	style:    Expected output caps.Style of the string.
	//	repStyle: The caps.ReplaceStyle to use if a word needs to be replaced.
	//	words:    A list parsed caps.Word.
	//	join:     The delimiter to use when joining the words. For CamelCase, this is an empty string.
	FormatterInterface interface {
		Format(style Style, repStyle ReplaceStyle, input string, join string) string
	}
)

func NewTokenizer(delimiters string) Tokenizer {
	d := runes(delimiters)
	sort.Sort(d)
	return Tokenizer{
		delimiters: d,
	}
}

type Tokenizer struct {
	delimiters runes
}

// Tokenize splits a string into a list of tokens based on the case of each
// rune, its delimiters, and the specified allowedSymbols.
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
func (t Tokenizer) Tokenize(str string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token {
	tokens := []token.Token{}
	pending := []token.Token{}
	foundLower := false
	current := token.Token{}
	prevNumber := false
	allowed := newRunes(allowedSymbols)

	for i, r := range str {
		rStr := string(r)
		_ = rStr
		switch {
		case unicode.IsUpper(r):
			if foundLower && current.Len() > 0 {
				tokens = append(tokens, current)
				current = token.Token{}
			}
			current = token.AppendRune(current, r)
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
			current = token.AppendRune(current, r)
			pending = nil
			foundLower = true
		case unicode.IsNumber(r):
			// if adding the number onto current makes it a valid number
			// then append this rune to current
			if token.AppendRune(current, r).IsNumber(numberRules) {
				current = token.AppendRune(current, r)
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
				current = token.AppendRune(current, r)
			}
			prevNumber = true
		default:
			if allowed.Contains(r) {
				if current.Len() > 0 {
					if current.IsNumber(numberRules) {
						// this gets a bit tricky because we need to check if adding the
						// rune to the current makes it a number. However, in all
						// default cases, a non-numeric rune must be followed either by
						// a number or a 'e' and a number. as such, we have to check if
						// both this and the next rune (and possibly the rune after
						// that) make it a number.

						n := token.AppendRune(current, r)
						runes := []rune(str)
						if n.IsNumber(numberRules) {
							current = n
						} else if i <= len(runes)-2 {
							next := token.AppendRune(n, runes[i+1])
							if next.IsNumber(numberRules) {
								current = n
							} else {
								if foundLower {
									tokens = append(tokens, current)
									current = token.New([]rune{r})
								} else {
									pending = append(pending, current)
									current = token.New([]rune{r})
								}
							}
						} else {
							if foundLower {
								tokens = append(tokens, current)
								current = token.New([]rune{r})
							} else {
								pending = append(pending, current)
								current = token.New([]rune{r})
							}
						}
					} else {
						current = token.AppendRune(current, r)
					}
				} else {
					current = token.AppendRune(current, r)
				}
			} else if t.delimiters.Contains(r) {
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

var _ TokenizerInterface = Tokenizer{}

// ReplaceStyle is used to indicate the case style the text should be transformed to
// when seeking replacement text in a Replacer.
//
// When a Replacer is configured, the expected input is:
//
//	caps.R{ Camel: "Json", Screaming: "JSON" }
//
// If the ReplaceStyle equals ReplaceStyleScreaming then an input of "MarshalJson" will return
// "MarshaalJSON".
type ReplaceStyle uint8

type (
	R           = Replacement
	Replacement struct {
		// Camelcase variant of the word which should be replaced.
		// e.g. "Http"
		Camel string
		// Screaming (all uppercase) representation of the word to replace.
		// e.g. "HTTP"
		Screaming string
	}
)

type Style uint8

const (
	StyleLower      Style = iota // The output should be lowercase (e.g. "an_example")
	StyleScreaming               // The output should be screaming (e.g. "AN_EXAMPLE")
	StyleCamel                   // The output should be camelcase (e.g. "AnExample")
	StyleLowerCamel              // The output should be lower camelcase (e.g. "anExample")
)

const (
	ReplaceStyleNotSpecified = iota
	ReplaceStyleCamel        // Text should be replaced with the Camel variant (e.g. "Json").
	ReplaceStyleScreaming    // Text should be replaced with the screaming variant (e.g. "JSON").
	ReplaceStyleLower        // Text should be replaced with the lowercase variant (e.g. "json").
)

// NewFormatter creates a new Formatter which is used to format the input text to the desired output.
//
// replacements are used to make replacements of tokens to the specified
// formatting (e.g. { "Json", "JSON"}).
//
// tokenizer is used to tokenize the input text.
func NewFormatter(replacements []Replacement, tokenizer TokenizerInterface) Formatter {
	r := Formatter{
		from:   make(map[string]string, len(replacements)),
		to:     make(map[string]string, len(replacements)),
		lookup: make(map[string]lookupResult, len(replacements)*2),
	}
	for _, v := range replacements {
		r.set(v.Camel, v.Screaming)
	}
	return r
}

// Formatter contains a table of words to their desired replacement. Words,
// without neighboring numbers, will be compared against the keys of this
// table to determine if the string should be replaced with the value of the
// table.
//
// This is primarily designed for acronyms but it could be used for other
// purposes.
//
// If Fn is set, it will be called to determine the replacement for each word
// rather than relying on the internal table to perform lookups.
// Defaults:
//
//	{ "Http",  "HTTP" },
//	{ "Https", "HTTPS" },
//	{ "Html",  "HTML" },
//	{ "Xml",   "XML" },
//	{ "Json",  "JSON" },
//	{ "Csv",   "CSV" },
//	{ "Aws",   "AWS" },
//	{ "Gcp",   "GCP" },
//	{ "Sql",   "SQL" },
type Formatter struct {
	delimiters string
	from       map[string]string
	to         map[string]string
	lookup     map[string]lookupResult
	tokenizer  Tokenizer
}

type lookupResult struct {
	from string
	to   string
}

func (f Formatter) Contains(key string) bool {
	_, ok := f.lookup[strings.ToLower(key)]
	return ok
}

// Lookup returns the Replacement for the given key, returning nil if it does
// not exist.
func (r Formatter) Lookup(key string) *Replacement {
	res, ok := r.lookup[key]
	if ok {
		return &Replacement{Camel: res.from, Screaming: res.to}
	}
	if res, ok = r.lookup[strings.ToLower(key)]; ok {
		return &Replacement{Camel: res.from, Screaming: res.to}
	} else {
		return nil
	}
}

// Table returns a representation of the internal table.
func (r Formatter) Table(key string) map[string]string {
	m := make(map[string]string, len(r.from))
	for k, v := range r.from {
		m[k] = v
	}
	return m
}

func (r Formatter) Replacements() []Replacement {
	res := make([]Replacement, 0, len(r.from))
	for upper, screaming := range r.from {
		res = append(res, Replacement{
			Camel:     upper,
			Screaming: screaming,
		})
	}
	return res
}

func (r *Formatter) set(key, value string) {
	r.lookup[strings.ToLower(key)] = lookupResult{
		from: key,
		to:   value,
	}
	r.lookup[strings.ToLower(value)] = lookupResult{
		from: key,
		to:   value,
	}
	r.from[key] = value
	r.to[value] = key
}

// Set adds the key/value pair to the table.
func (r *Formatter) Set(key, value string) {
	l := strings.ToLower(key)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from)
		delete(r.to, v.to)
		delete(r.lookup, l)
		return
	}
	l = strings.ToLower(value)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from)
		delete(r.to, v.to)
		delete(r.lookup, l)
	}
	r.set(key, value)
}

// Remove deletes the key from the map. Either variant is sufficient.
func (r *Formatter) Delete(key string) {
	l := strings.ToLower(key)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from)
		delete(r.to, v.to)
		delete(r.lookup, l)
	}
}

func (r *Formatter) resolve(str string, style ReplaceStyle) (string, bool) {
	if lookup, ok := r.lookup[str]; ok {
		switch style {
		case ReplaceStyleCamel:
			return lookup.from, true
		case ReplaceStyleScreaming:
			return lookup.to, true
		}
	}
	return str, false
}

func FormatToken(style Style, tok token.Token) string {
	switch style {
	case StyleCamel:

	case StyleLowerCamel:
		return "LowerCamel"
	case StyleScreaming:
		return "Screaming"
	case StyleLower:
		return "Lower"
	}
	return "NotSpecified"
}

func (r Formatter) Format(style Style, repStyle ReplaceStyle, input string, join string) string {
	// chain :=token.Token{}
	// parts := []string{}
	// var lookup string

	// for _, tok := range tokens {

	// for i := len(tokens) - 1; i >= 0; i-- {
	// 	token := tokens[i]
	// 	switch token.Len() {
	// 	case 0:
	// 		continue
	// 	case 1:
	// 		switch chain.Len() {
	// 		case 0:
	// 			chain = token
	// 		case 1: // found another possible link in a chain
	// 			chain = chain.Append(token)
	// 		default:
	// 			if lookup, ok = r.resolve(string(chain.lower), repStyle); ok {
	// 				parts = append(parts, lookup)
	// 				chain = token
	// 			} else {
	// 			}
	// 		}
	// 	default:
	// 		if chain.Len() > 0 {
	// 		}
	// 	}
	// }
	// if chain.Len() > 0 {
	// 	parts = append(parts, string(chain.value))
	// }
	// return strings.Join(parts, join)
	panic("not impl")
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
	// The Formatter to use.
	//
	// Default:
	// 	DefaultReplacer
	Formatter FormatterInterface

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
}

func loadOpts(opts []Opts) Opts {
	// result := Opts{
	// 	AllowedSymbols: "",
	// 	Formatter:      DefaultFormatter,
	// 	ReplaceStyle:   ReplaceStyleNotSpecified,
	// }
	// if len(opts) == 0 {
	// 	return result
	// }

	// if opts[0].AllowedSymbols != "" {
	// 	result.AllowedSymbols = opts[0].AllowedSymbols
	// }
	// if opts[0].Delimiters != "" {
	// 	result.Delimiters = opts[0].Delimiters
	// }
	// if opts[0].Formatter != nil {
	// 	result.Formatter = opts[0].Formatter
	// }
	// return result
	panic("not implemented")
}

// func tryReplacements[T ~string](replace bool, str T, replacements map[string]string) T {
// 	if !replace {
// 		return str
// 	}

// 	key := strings.Builder{}

// 	for _, r := range str {
// 		if unicode.IsNumber(r) {
// 			if key.Len() > 0 {
// 				key.Reset()
// 			} else {
// 				if v, ok := replacements[key.String()]; ok {
// 					return T(v)
// 				}
// 				return str

// 			}
// 		} else {
// 			key.WriteRune(r)
// 		}
// 	}

// 	if v, ok := replacements[key.String()]; ok {
// 		return T(v)
// 	}
// 	return str
// }

// // MakeReplacements checks to see if the str, without numbers, equals any of the
// // keys in the Replacements map. If there is a match then the replacement value,
// // along with any numbers in str, is used instead.
// func PerformReplace[T ~string](str T, replacements map[string]string) T {
// 	return tryReplacements(true, str, replacements)
// }

// UpperFirst converts the first rune of str to uppercase.
func UpperFirst[T ~string](str T) T {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return T(runes)
}

// LowerFirst converts the first rune of str to lowercase.
func LowerFirst[T ~string](str T) T {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return T(runes)
}

// Without numbers returns the string with numbers removed.
func WithoutNumbers[T ~string](s T) T {
	return T(strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, string(s)))
}

// ToCamel transforms the case of str into (lower) camelCase using delimiters to
// determine word breaks.
//
// If useReplacements is true then any entries in Replacements will be made
// for single words.
// If allowed is not empty, then those non-alphanumeric characters will be
// allowed in the output.
//
// If delimiters is empty, the DefaultDelimiters will be used, which are "-#@!$&=.+:;_~ (){}[]"
//
//	text.ToCamel("This is [an] {example}${id32}.") // thisIsAnExampleID32
//	text.ToCamel("This is [an] {example}${id32}.", text.Options{ AllowedSymbols: "$" }) // thisIsAnExample$ID32
func ToCamel[T ~string](str T, options ...Opts) T {
	panic("")
	// return T(opts.Replacer.(style, parts, ""))
}

// ToPascal transforms the case of str into PascalCase (also known as upper
// camelcase or, in some cases, simply camelcase) using delimiters to determine
// word breaks.
//
// If useReplacements is true then any entries in Replacements will be made for
// single words. If allowed is not empty, then those non-alphanumeric characters
// will be allowed in the output.
//
// If delimiters is empty, the DefaultDelimiters will be used, which are
// "-#@!$&=.+:;_~ (){}[]"
//
//	text.ToPascal("This is [an] {example}${id32}.", true, "") // ThisIsAnExampleID32
//	text.ToPascal("This is [an] {example}${id32}.", true, "$") // ThisIsAnExample$ID32
func ToPascal[T ~string](str T, opts ...Opts) T {
	camel := ToCamel(str, opts...)
	return UpperFirst(camel)
}

// ToSnake transforms the case of str into lowercase string seperated by
// delimiter, using delimiters to determine word breaks.
//
// If lowercase is false, the output will be all uppercase.
//
// # See Options for more information on available configuration
//
// # Example
//
//	text.ToDelimited("This is [an] {example}${id}.#32", '.', true) // this.is.an.example.id.#32
//	text.ToDelimited("This is [an] {example}${id32}.break32", '.', false, text.Opts{AllowedSymbols: "$"}) // THIS.IS.AN.EXAMPLE.ID.BREAK.32
//	text.ToDelimited("This is [an] {example}${id32}.v32", '.', true, text.Opts{AllowedSymbols: "$" }) // this.is.an.example.id.$.v32
func ToSnake[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '_', true, options...)
}

// ToDelimited transforms the case of str into a string seperated by delimiter,
// using either the specified Opts.Delimiters or DefaultDelimiters to determine
// word breaks.
//
// If lowercase is false, the output will be all uppercase.
//
// # See Options for more information on available configuration
//
// # Example
//
//	text.ToDelimited("This is [an] {example}${id}.#32", '.', true) // this.is.an.example.id.#32
//	text.ToDelimited("This is [an] {example}${id32}.break32", '.', false, text.Opts{AllowedSymbols: "$"}) // THIS.IS.AN.EXAMPLE.ID.BREAK.32
//	text.ToDelimited("This is [an] {example}${id32}.v32", '.', true, text.Opts{AllowedSymbols: "$" }) // this.is.an.example.id.$.v32
func ToDelimited[T ~string](value T, delimiter rune, lowercase bool, options ...Opts) T {
	panic("not impl")
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
	return res > -1 && res < len(r)
}

func newRunes(val []rune) runes {
	if len(val) == 0 {
		return nil
	}
	r := runes(val)
	sort.Sort(r)
	return r
}

var _ sort.Interface = (*runes)(nil)
