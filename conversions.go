package caps

import (
	"strings"
	"unicode"
)

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
	return T(opts.Converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   opts.ReplaceStyle,
		Input:          string(str),
		Join:           "",
		AllowedSymbols: opts.AllowedSymbols,
		NumberRules:    opts.NumberRules,
	}))
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

	return T(opts.Converter.Convert(ConvertRequest{
		Style:          StyleLowerCamel,
		ReplaceStyle:   opts.ReplaceStyle,
		Input:          string(str),
		Join:           "",
		AllowedSymbols: opts.AllowedSymbols,
		NumberRules:    opts.NumberRules,
	}))
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

// ToDotNotation transforms the case of str into Lower Dot Notation Case (e.g. an.example.string) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToDotNotation("This is [an] {example}${id32}.") // this.is.an.example.id.32
func ToDotNotation[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '.', true, options...)
}

// ToScreamingDotNotation transforms the case of str into Screaming Kebab Case (e.g.
// AN.EXAMPLE.STRING) using either the provided Converter or the
// DefaultConverter otherwise.
//
//	caps.ToScreamingDotNotation("This is [an] {example}${id32}.") // THIS.IS.AN.EXAMPLE.ID.32
func ToScreamingDotNotation[T ~string](str T, options ...Opts) T {
	return ToDelimited(str, '.', false, options...)
}

// ToTitle transforms the case of str into Title Case (e.g. An Example String) using
// either the provided Converter or the DefaultConverter otherwise.
//
//	caps.ToTitle("This is [an] {example}${id32}.") // This Is An Example ID 32
func ToTitle[T ~string](str T, options ...Opts) T {
	opts := loadOpts(options)
	return T(opts.Converter.Convert(ConvertRequest{
		Style:          StyleCamel,
		ReplaceStyle:   opts.ReplaceStyle,
		Input:          string(str),
		Join:           " ",
		AllowedSymbols: opts.AllowedSymbols,
		NumberRules:    opts.NumberRules,
	}))
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
	return T(opts.Converter.Convert(ConvertRequest{
		Style:          style,
		ReplaceStyle:   replacementStyle,
		Input:          string(str),
		Join:           string(delimiter),
		AllowedSymbols: opts.AllowedSymbols,
		NumberRules:    opts.NumberRules,
	}))
}

// ToLower returns s with all Unicode letters mapped to their lower case.
func ToLower[T ~string](str T) T {
	return T(strings.ToLower((string(str))))
}

// ToUpper returns s with all Unicode letters mapped to their upper case.
func ToUpper[T ~string](str T) T {
	return T(strings.ToUpper((string(str))))
}
