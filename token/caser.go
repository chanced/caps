package token

import "unicode"

var DefaultCaser = Unicode{}

type Unicode struct{}

var (
	TurkishCaser = &unicode.TurkishCase
	AzeriCaser   = &unicode.AzeriCase
)

// ToTitle maps the rune to title case using unicode.ToTitle.
func (Unicode) ToTitle(r rune) rune {
	return unicode.ToTitle(r)
}

// ToUpper maps the rune to lower case using unicode.ToLower
func (Unicode) ToLower(r rune) rune {
	return unicode.ToLower(r)
}

// ToUpper maps the rune to upper case using unicode.ToUpper
func (Unicode) ToUpper(r rune) rune {
	return unicode.ToUpper(r)
}

// Caser is satisfied by types which can map runes to their lowercase and
// uppercase equivalents.
type Caser interface {
	// ToLower maps the rune to lower case
	ToLower(r rune) rune
	// ToUpper maps the rune to upper case
	ToUpper(r rune) rune
	// ToTitle maps the rune to title case
	ToTitle(r rune) rune
}

var unicodeCaser Caser = Unicode{}

func CaserOrDefault(caser Caser) Caser {
	if caser == nil {
		return unicodeCaser
	}
	return caser
}
