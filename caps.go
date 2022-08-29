package caps

import (
	"sort"
	"strings"
	"unicode"

	"github.com/chanced/caps/token"
)

const (
	// DEFAULT_DELIMITERS is the default set of delimiters in string convert.
	DEFAULT_DELIMITERS string = " _.!?:;$-(){}[]#@&+~"
)

var (
	// DefaultTokenizer is the default Tokenizer.
	DefaultTokenizer TokenizerImpl = NewTokenizer(DEFAULT_DELIMITERS)
	// DefaultReplacements is the list of Replacements passed to DefaultConverter.
	//
	//  {"Http", "HTTP"},
	//  {"Https", "HTTPS"},
	//  {"Id", "ID"},
	//  {"Ip", "IP"},
	//  {"Html", "HTML"},
	//  {"Xml", "XML"},
	//  {"Json", "JSON"},
	//  {"Csv", "CSV"},
	//  {"Aws", "AWS"},
	//  {"Gcp", "GCP"},
	//  {"Sql", "SQL"},
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

	// DefaultConverter is the default Converter instance.
	//
	// # Default parameters:
	//

	//
	// replacements:
	//  DefaultReplacements:
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
	DefaultConverter = NewConverter(DefaultReplacements, DefaultTokenizer)
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
	Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string
}

// NewTokenizer creates and returns a new TokenizerImpl which implements the
// Tokenizer interface.
//
// Tokenizers are used by ConverterImpl to tokenize the input text into
// token.Tokens that are then formatted.
func NewTokenizer(delimiters string) TokenizerImpl {
	d := runes(delimiters)
	sort.Sort(d)
	return TokenizerImpl{
		delimiters: d,
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
func (t TokenizerImpl) Tokenize(str string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) []token.Token {
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
						} else if i <= len(runes)-2 && (unicode.IsNumber(runes[i+1]) || unicode.IsLetter(runes[i+1])) || allowed.Contains(runes[i+1]) {
							next := token.AppendRune(n, runes[i+1])
							if next.IsNumber(numberRules) {
								current = n
							} else {
								if foundLower {
									tokens = append(tokens, current)
									current = token.FromRunes([]rune{r})
								} else {
									pending = append(pending, current)
									current = token.FromRunes([]rune{r})
								}
							}
						} else {
							if foundLower {
								tokens = append(tokens, current)
								current = token.FromRunes([]rune{r})
							} else {
								pending = append(pending, current)
								current = token.FromRunes([]rune{r})
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
func NewConverter(replacements []Replacement, tokenizer Tokenizer) ConverterImpl {
	r := ConverterImpl{
		from:      make(map[string]token.Token, len(replacements)),
		to:        make(map[string]token.Token, len(replacements)),
		lookup:    make(map[string]lookupResult, len(replacements)*2),
		tokenizer: tokenizer,
	}
	for _, v := range replacements {
		r.set(v.Camel, v.Screaming)
	}
	return r
}

// ConverterImpl contains a table of words to their desired replacement. Tokens
// will be compared against the keys of this table to determine if the string
// should be replaced with the value of the table.
//
// This is primarily designed for acronyms but it could be used for other
// purposes.
//
// The default Replacements:
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
type ConverterImpl struct {
	from      map[string]token.Token
	to        map[string]token.Token
	lookup    map[string]lookupResult
	tokenizer Tokenizer
}

type lookupResult struct {
	from token.Token
	to   token.Token
}

// Contains reports whether a key is in the Converter's replacement table.
func (f ConverterImpl) Contains(key string) bool {
	_, ok := f.lookup[strings.ToLower(key)]
	return ok
}

// Lookup returns the Replacement for the given key, returning nil if it does
// not exist.
func (r ConverterImpl) Lookup(key string) *Replacement {
	res, ok := r.lookup[key]
	if ok {
		return &Replacement{Camel: res.from.String(), Screaming: res.to.String()}
	}
	if res, ok = r.lookup[strings.ToLower(key)]; ok {
		return &Replacement{Camel: res.from.String(), Screaming: res.to.String()}
	} else {
		return nil
	}
}

// Table returns a representation of the internal table.
func (r ConverterImpl) Table(key string) map[string]token.Token {
	m := make(map[string]token.Token, len(r.from))
	for k, v := range r.from {
		m[k] = v
	}
	return m
}

// Replacements returns a slice of Replacement in the lookup table.
func (r ConverterImpl) Replacements() []Replacement {
	res := make([]Replacement, 0, len(r.from))
	for upper, screaming := range r.from {
		res = append(res, Replacement{
			Camel:     upper,
			Screaming: string(screaming.Value()),
		})
	}
	return res
}

func (r *ConverterImpl) set(key, value string) {
	from := token.FromString(key)
	to := token.FromString(value)
	r.lookup[string(from.Lower())] = lookupResult{
		from: from,
		to:   to,
	}
	r.lookup[string(to.Lower())] = lookupResult{
		from: from,
		to:   to,
	}
	r.from[key] = from
	r.to[value] = to
}

// Set adds the key/value pair to the table.
func (r *ConverterImpl) Set(key, value string) {
	l := strings.ToLower(key)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from.String())
		delete(r.to, v.to.String())
		delete(r.lookup, l)
		return
	}
	l = strings.ToLower(value)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from.String())
		delete(r.to, v.to.String())
		delete(r.lookup, l)
	}
	r.set(key, value)
}

// Remove deletes the key from the map. Either variant is sufficient.
func (r *ConverterImpl) Delete(key string) {
	l := strings.ToLower(key)
	if v, ok := r.lookup[l]; ok {
		delete(r.from, v.from.String())
		delete(r.to, v.to.String())
		delete(r.lookup, l)
	}
}

func (r *ConverterImpl) resolve(tok token.Token, style ReplaceStyle) (token.Token, bool) {
	l := string(tok.Lower())
	if lookup, ok := r.lookup[l]; ok {
		switch style {
		case ReplaceStyleCamel:
			return lookup.from, true
		case ReplaceStyleScreaming:
			return lookup.to, true
		case ReplaceStyleLower:
			return token.FromString(lookup.to.Lower()), true

		}
	}
	return token.Token{}, false
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

// Convert formats the string with the desired style.
func (r ConverterImpl) Convert(style Style, repStyle ReplaceStyle, input string, join string, allowedSymbols []rune, numberRules map[rune]func(index int, r rune, val []rune) bool) string {
	tokens := r.tokenizer.Tokenize(input, allowedSymbols, numberRules)
	var chain token.Token
	var parts []string
	var lookup token.Token
	brokeChain := false
	var ok bool
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		switch tok.Len() {
		case 0:
			continue
		case 1:
			chain = token.Append(tok, chain)
			if lookup, ok = r.resolve(chain, repStyle); ok {
				parts = append(parts, lookup.Value())
				chain = token.Token{}
			}
		default:
			brokeChain = true
			if chain.Len() > 0 {
				split := chain.Split()
				for z := len(split) - 1; z >= 0; z-- {
					letter := split[z]
					if i == 0 && z == 0 {
						parts = append(parts, FormatToken(style, 0, letter))
					} else {
						parts = append(parts, FormatToken(style, z, letter))
					}
				}
				chain = token.Token{}
			}
			if lookup, ok = r.resolve(tok, repStyle); ok {
				parts = append(parts, lookup.Value())
			} else {
				parts = append(parts, FormatToken(style, i, tok))
			}
		}
	}
	result := strings.Builder{}
	result.Grow(len(input))

	var part string
	shouldWriteDelimiter := false
	if brokeChain && chain.Len() > 0 {
		for i, letter := range chain.Split() {
			if shouldWriteDelimiter {
				result.WriteString(join)
			}
			result.WriteString(FormatToken(style, i, letter))
			if !shouldWriteDelimiter {
				shouldWriteDelimiter = len(part) > 0 && len(join) > 0
			}
		}
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

	if !brokeChain && chain.Len() > 0 {
		for i, letter := range chain.Split() {
			if shouldWriteDelimiter {
				result.WriteString(join)
			}
			// if unicode.IsLetter(letter.Value()[0]) || unicode.IsDigit(letter.Value()[0]) || allow {
			result.WriteString(FormatToken(style, len(parts)+i, letter))
			if !shouldWriteDelimiter {
				shouldWriteDelimiter = len(part) > 0 && len(join) > 0
			}
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

// UpperFirst converts the first rune of str to uppercase.
func UpperFirst[T ~string](str T) T {
	t := token.FromString(str)
	return T(t.UpperFirst())
}

// LowerFirst converts the first rune of str to lowercase.
func LowerFirst[T ~string](str T) T {
	t := token.FromString(str)
	return T(t.LowerFirst())
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

// ToDelimited transforms the case of str into a string seperated by delimiter,
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

var _ sort.Interface = (*runes)(nil)
