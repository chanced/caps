package caps

import (
	"strings"
	"unicode"

	"github.com/chanced/caps/token"
)

type Caps struct {
	caser          token.Caser
	allowedSymbols string
	converter      Converter
	replaceStyle   ReplaceStyle
	numberRules    token.NumberRules
}

// CapsOpts include configurable options for case conversion.
//
// See the documentation for the individual fields for more information.
type CapsOpts struct {
	// Any characters within this string will be allowed in the output.
	//
	// This does not affect delimiters (e.g. '_', '-', '.') as they are added
	// post-tokenization.
	//
	// Default:
	//  ""
	AllowedSymbols string
	// The Converter to use.
	//
	// Default:
	// 	DefaultConverter
	Converter Converter

	// If not set, this will be DefaultReplacements.
	Replacements []Replacement

	// ReplaceStyle overwrites the way words are replaced.
	//
	// A typical call to ToLowerCamel for "ServeJSON" with a Converter that
	// contains {"Json": "JSON"} would result in "serveJSON" by using the
	// ReplaceStyleScreaming variant. If ReplacementStyle was set to
	// ReplaceStyleCamel, on the call to ToLowerCamel then the result would
	// be "serveHttp".
	//
	// The default replacement style is dependent upon the target Style.
	ReplaceStyle ReplaceStyle
	// NumberRules are used by the DefaultTokenizer to augment the standard
	// rules for determining if a rune is part of a number.
	//
	// Note, if you add special characters here, they must be present in the
	// AllowedSymbols string for them to be part of the output.
	NumberRules token.NumberRules
	// Special unicode case rules.
	// See unicode.SpecialCase or token.Caser for more information.
	//
	// Default: token.DefaultCaser (which relies on the default unicode
	// functions)
	Caser token.Caser

	// If not set, uses DefaultTokenizer
	Tokenizer Tokenizer
}

func loadCapsOpts(opts []CapsOpts) CapsOpts {
	result := CapsOpts{
		AllowedSymbols: "",
		ReplaceStyle:   ReplaceStyleScreaming,
	}

	for _, opt := range opts {
		result.AllowedSymbols += opt.AllowedSymbols
		if opt.Converter != nil {
			result.Converter = opt.Converter
		}
		if opt.ReplaceStyle != ReplaceStyleNotSpecified {
			result.ReplaceStyle = opt.ReplaceStyle
		}
		if len(opt.NumberRules) > 0 {
			if result.NumberRules == nil {
				result.NumberRules = make(NumberRules)
			}
			for k, v := range opt.NumberRules {
				result.NumberRules[k] = v
			}
		}
		if opt.Replacements != nil {
			result.Replacements = append(result.Replacements, opt.Replacements...)
		}
		if opt.Caser != nil {
			result.Caser = opt.Caser
		}
		if opt.Tokenizer != nil {
			result.Tokenizer = opt.Tokenizer
		}
	}
	if result.Caser == nil {
		result.Caser = token.DefaultCaser
	}
	if result.Replacements == nil {
		result.Replacements = DefaultReplacements
	}
	if result.Tokenizer == nil {
		result.Tokenizer = NewTokenizer(DEFAULT_DELIMITERS, result.Caser)
	}
	if result.Converter == nil {
		result.Converter = NewConverter(result.Replacements, result.Tokenizer, result.Caser)
	}

	return result
}

// New returns a new Caps instance with the provided options.
//
// if caser is nil, token.DefaultCaser is used (which relies on the default
// unicode functions)
func New(options ...CapsOpts) Caps {
	opts := loadCapsOpts(options)

	return Caps{
		caser:          opts.Caser,
		allowedSymbols: opts.AllowedSymbols,
		converter:      opts.Converter,
		replaceStyle:   opts.ReplaceStyle,
		numberRules:    opts.NumberRules,
	}
}

func (c Caps) ReplaceStyle() ReplaceStyle {
	return c.replaceStyle
}

func (c Caps) NumberRules() token.NumberRules {
	return c.numberRules
}

// AllowedSymbols is
func (c Caps) AllowedSymbols() string {
	return c.allowedSymbols
}

// UpperFirst converts the first rune of str to unicode upper case.
//
// This method does not support special cases (such as Turkish and Azeri)
func (c Caps) UpperFirst(str string) string {
	return token.UpperFirst(c.caser, str)
}

// LowerFirst converts the first rune of str to lowercase.
func (c Caps) LowerFirst(str string) string {
	return token.LowerFirst(c.caser, str)
}

// Without numbers returns the string with all numeric runes removed.
//
// It does not currently use any logic to determine if a rune (e.g. '.')
// is part of a number. This may change in the future.
func (c Caps) WithoutNumbers(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, string(s))
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
func (c Caps) ToCamel(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           "",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
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
func (c Caps) ToLowerCamel(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleLowerCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           "",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}

// ToSnake transforms the case of str into Lower Snake Case (e.g. an_example_string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToSnake("This is [an] {example}${id32}.") // this_is_an_example_id_32
func (c Caps) ToSnake(str string) string {
	return c.ToDelimited(str, '_', true)
}

// ToScreamingSnake transforms the case of str into Screaming Snake Case (e.g.
// AN_EXAMPLE_STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingSnake("This is [an] {example}${id32}.") // THIS_IS_AN_EXAMPLE_ID_32
func (c Caps) ToScreamingSnake(str string) string {
	return ToDelimited(str, '_', false)
}

// ToKebab transforms the case of str into Lower Kebab Case (e.g. an-example-string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToKebab("This is [an] {example}${id32}.") // this-is-an-example-id-32
func (c Caps) ToKebab(str string) string {
	return ToDelimited(str, '-', true)
}

// ToScreamingKebab transforms the case of str into Screaming Kebab Snake (e.g.
// AN-EXAMPLE-STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingKebab("This is [an] {example}${id32}.") // THIS-IS-AN-EXAMPLE-ID-32
func (c Caps) ToScreamingKebab(str string) string {
	return ToDelimited(str, '-', false)
}

// ToDotNotation transforms the case of str into Lower Dot Notation Case (e.g. an.example.string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToDotNotation("This is [an] {example}${id32}.") // this.is.an.example.id.32
func (c Caps) ToDotNotation(str string) string {
	return ToDelimited(str, '.', true)
}

// ToScreamingDotNotation transforms the case of str into Screaming Kebab Case (e.g.
// AN.EXAMPLE.STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingDotNotation("This is [an] {example}${id32}.") // THIS.IS.AN.EXAMPLE.ID.32
func (c Caps) ToScreamingDotNotation(str string) string {
	return ToDelimited(str, '.', false)
}

// ToTitle transforms the case of str into Title Case (e.g. An Example String) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToTitle("This is [an] {example}${id32}.") // This Is An Example ID 32
func (c Caps) ToTitle(str string) string {
	return c.converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   c.replaceStyle,
		Input:          str,
		Join:           " ",
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
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
func (c Caps) ToDelimited(str string, delimiter rune, lowercase bool) string {
	var style Style
	var replacementStyle ReplaceStyle
	if lowercase {
		style = StyleLower
		replacementStyle = ReplaceStyleLower
	} else {
		style = StyleScreaming
		replacementStyle = ReplaceStyleScreaming
	}
	return c.converter.Convert(ConvertRequest{
		Style:          style,
		ReplaceStyle:   replacementStyle,
		Input:          str,
		Join:           string(delimiter),
		AllowedSymbols: c.allowedSymbols,
		NumberRules:    c.numberRules,
	})
}
