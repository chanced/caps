package caps

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

const (
	ReplaceStyleNotSpecified ReplaceStyle = iota
	ReplaceStyleCamel                     // Text should be replaced with the Camel variant (e.g. "Json").
	ReplaceStyleScreaming                 // Text should be replaced with the screaming variant (e.g. "JSON").
	ReplaceStyleLower                     // Text should be replaced with the lowercase variant (e.g. "json").
)

type Style uint8

const (
	StyleNotSpecified Style = iota
	StyleLower              // The output should be lowercase (e.g. "an_example")
	StyleScreaming          // The output should be screaming (e.g. "AN_EXAMPLE")
	StyleCamel              // The output should be camel case (e.g. "AnExample")
	StyleLowerCamel         // The output should be lower camel case (e.g. "anExample")
)

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

// Deprecated: Use WithAllowedSymbols instead.
func UseAllowedSymbols(symbols string) Opts {
	return Opts{
		AllowedSymbols: symbols,
	}
}

// WithConverter sets the Converter to use
func WithConverter(converter Converter) Opts {
	return Opts{
		Converter: converter,
	}
}

// WithReplaceStyle sets the ReplaceStyle to use
//
// There are also methods for each ReplaceStyle (e.g. WithReplaceStyleCamel)
func WithReplaceStyle(style ReplaceStyle) Opts {
	return Opts{
		ReplaceStyle: style,
	}
}

// WithReplaceStyleCamel indicates Replacements should use the Camel variant
// (e.g. "Json").
func WithReplaceStyleCamel() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleCamel,
	}
}

// WithReplaceStyleScreaming indicates Replacements should use the screaming
// variant (e.g. "JSON").
func WithReplaceStyleScreaming() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleScreaming,
	}
}

// WithReplaceStyleLower indicates Replacements should use the lowercase variant
// (e.g. "json").
func WithReplaceStyleLower() Opts {
	return Opts{
		ReplaceStyle: ReplaceStyleLower,
	}
}

// WithNumberRules sets the NumberRules to use
func WithNumberRules(rules map[rune]func(index int, r rune, val []rune) bool) Opts {
	return Opts{
		NumberRules: rules,
	}
}

// UseConverter sets the Converter to use
//
// Deprecated: Use WithConverter instead.
func UseConverter(converter Converter) Opts {
	return Opts{
		Converter: converter,
	}
}

// Deprecated: Use WithReplaceStyle instead.
func UseReplaceStyle(style ReplaceStyle) Opts {
	return Opts{
		ReplaceStyle: style,
	}
}

// Deprecated: Use WithNumberRules instead.
func UseNumberRules(rules map[rune]func(index int, r rune, val []rune) bool) Opts {
	return Opts{
		NumberRules: rules,
	}
}

func loadOpts(opts []Opts) Opts {
	result := Opts{
		AllowedSymbols: "",
		Converter:      DefaultConverter,
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
				result.NumberRules = make(map[rune]func(index int, r rune, val []rune) bool)
			}
			for k, v := range opt.NumberRules {
				result.NumberRules[k] = v
			}
		}
	}
	return result
}
