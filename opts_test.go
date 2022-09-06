package caps_test

import (
	"testing"

	"github.com/chanced/caps"
)

func TestReplaceStyleString(t *testing.T) {
	t.Run("ReplaceStyleCamel", func(t *testing.T) {
		if caps.ReplaceStyleCamel.String() != "ReplaceStyleCamel" {
			t.Error("expected ReplaceStyleCamel.String() to return \"ReplaceStyleCamel\"")
		}
	})
	t.Run("ReplaceStyleSnake", func(t *testing.T) {
		if caps.ReplaceStyleLower.String() != "ReplaceStyleLower" {
			t.Error("expected ReplaceStyleLower.String() to return \"ReplaceStyleLower\"")
		}
	})
	t.Run("ReplaceStyleScreaming", func(t *testing.T) {
		if caps.ReplaceStyleScreaming.String() != "ReplaceStyleScreaming" {
			t.Error("expected ReplaceStyleScreaming.String() to return \"ReplaceStyleScreaming\"")
		}
	})
	t.Run("ReplaceStyleNotSpecified", func(t *testing.T) {
		if caps.ReplaceStyleNotSpecified.String() != "ReplaceStyleNotSpecified" {
			t.Error("expected ReplaceStyleNotSpecified.String() to return \"ReplaceStyleNotSpecified\"")
		}
	})
}

func TestReplaceStyleIsCamel(t *testing.T) {
	t.Run("ReplaceStyleNotSpecified", func(t *testing.T) {
		if caps.ReplaceStyleNotSpecified.IsCamel() {
			t.Error("expected ReplaceStyleNotSpecified.IsCamel() to return false")
		}
	})
	t.Run("ReplaceStyleScreaming", func(t *testing.T) {
		if caps.ReplaceStyleScreaming.IsCamel() {
			t.Error("expected ReplaceStyleScreaming.IsCamel() to return false")
		}
	})
	t.Run("ReplaceStyleLower", func(t *testing.T) {
		if caps.ReplaceStyleLower.IsCamel() {
			t.Error("expected ReplaceStyleLower.IsCamel() to return false")
		}
	})
	t.Run("ReplaceStyleCamel", func(t *testing.T) {
		if !caps.ReplaceStyleCamel.IsCamel() {
			t.Error("expected ReplaceStyleCamel.IsCamel() to return true")
		}
	})
}

func TestReplaceStyleIsScreaming(t *testing.T) {
	t.Run("ReplaceStyleNotSpecified", func(t *testing.T) {
		if caps.ReplaceStyleNotSpecified.IsScreaming() {
			t.Error("expected ReplaceStyleNotSpecified.IsScreaming() to return false")
		}
	})
	t.Run("ReplaceStyleScreaming", func(t *testing.T) {
		if !caps.ReplaceStyleScreaming.IsScreaming() {
			t.Error("expected ReplaceStyleScreaming.IsScreaming() to return true")
		}
	})
	t.Run("ReplaceStyleLower", func(t *testing.T) {
		if caps.ReplaceStyleLower.IsScreaming() {
			t.Error("expected ReplaceStyleLower.IsScreaming() to return false")
		}
	})
	t.Run("ReplaceStyleCamel", func(t *testing.T) {
		if caps.ReplaceStyleCamel.IsScreaming() {
			t.Error("expected ReplaceStyleCamel.IsScreaming() to return false")
		}
	})
}

func TestReplaceStyleIsLower(t *testing.T) {
	t.Run("ReplaceStyleNotSpecified", func(t *testing.T) {
		if caps.ReplaceStyleNotSpecified.IsLower() {
			t.Error("expected ReplaceStyleNotSpecified.IsLower() to return false")
		}
	})
	t.Run("ReplaceStyleScreaming", func(t *testing.T) {
		if caps.ReplaceStyleScreaming.IsLower() {
			t.Error("expected ReplaceStyleScreaming.IsLower() to return false")
		}
	})
	t.Run("ReplaceStyleLower", func(t *testing.T) {
		if !caps.ReplaceStyleLower.IsLower() {
			t.Error("expected ReplaceStyleLower.IsLower() to return true")
		}
	})
	t.Run("ReplaceStyleCamel", func(t *testing.T) {
		if caps.ReplaceStyleCamel.IsLower() {
			t.Error("expected ReplaceStyleCamel.IsLower() to return false")
		}
	})
}

func TestStyleString(t *testing.T) {
	t.Run("StyleCamel", func(t *testing.T) {
		if caps.StyleCamel.String() != "StyleCamel" {
			t.Error("expected StyleCamel.String() to return \"StyleCamel\"")
		}
	})
	t.Run("StyleLower", func(t *testing.T) {
		if caps.StyleLower.String() != "StyleLower" {
			t.Error("expected StyleLower.String() to return \"StyleLower\"")
		}
	})
	t.Run("StyleScreaming", func(t *testing.T) {
		if caps.StyleScreaming.String() != "StyleScreaming" {
			t.Error("expected StyleScreaming.String() to return \"StyleScreaming\"")
		}
	})
	t.Run("StyleLowerCamel", func(t *testing.T) {
		if caps.StyleLowerCamel.String() != "StyleLowerCamel" {
			t.Error("expected StyleLowerCamel.String() to return \"StyleLowerCamel\"")
		}
	})
	t.Run("StyleNotSpecified", func(t *testing.T) {
		if caps.StyleNotSpecified.String() != "StyleNotSpecified" {
			t.Error("expected StyleNotSpecified.String() to return \"StyleNotSpecified\"")
		}
	})
}

func TestStyleIsLower(t *testing.T) {
	t.Run("StyleNotSpecified", func(t *testing.T) {
		if caps.StyleNotSpecified.IsLower() {
			t.Error("expected StyleNotSpecified.IsLower() to return false")
		}
	})
	t.Run("StyleScreaming", func(t *testing.T) {
		if caps.StyleScreaming.IsLower() {
			t.Error("expected StyleScreaming.IsLower() to return false")
		}
	})
	t.Run("StyleLower", func(t *testing.T) {
		if !caps.StyleLower.IsLower() {
			t.Error("expected StyleLower.IsLower() to return true")
		}
	})
	t.Run("StyleCamel", func(t *testing.T) {
		if caps.StyleCamel.IsLower() {
			t.Error("expected StyleCamel.IsLower() to return false")
		}
	})
	t.Run("StyleLowerCamel", func(t *testing.T) {
		if caps.StyleLowerCamel.IsLower() {
			t.Error("expected StyleLowerCamel.IsLower() to return false")
		}
	})
}

func TestStyleIsScreaming(t *testing.T) {
	t.Run("StyleNotSpecified", func(t *testing.T) {
		if caps.StyleNotSpecified.IsScreaming() {
			t.Error("expected StyleNotSpecified.IsScreaming() to return false")
		}
	})
	t.Run("StyleScreaming", func(t *testing.T) {
		if !caps.StyleScreaming.IsScreaming() {
			t.Error("expected StyleScreaming.IsScreaming() to return true")
		}
	})
	t.Run("StyleLower", func(t *testing.T) {
		if caps.StyleLower.IsScreaming() {
			t.Error("expected StyleLower.IsScreaming() to return false")
		}
	})
	t.Run("StyleCamel", func(t *testing.T) {
		if caps.StyleCamel.IsScreaming() {
			t.Error("expected StyleCamel.IsScreaming() to return false")
		}
	})
	t.Run("StyleLowerCamel", func(t *testing.T) {
		if caps.StyleLowerCamel.IsScreaming() {
			t.Error("expected StyleLowerCamel.IsScreaming() to return false")
		}
	})
}

func TestStyleIsCamel(t *testing.T) {
	t.Run("StyleNotSpecified", func(t *testing.T) {
		if caps.StyleNotSpecified.IsCamel() {
			t.Error("expected StyleNotSpecified.IsCamel() to return false")
		}
	})
	t.Run("StyleScreaming", func(t *testing.T) {
		if caps.StyleScreaming.IsCamel() {
			t.Error("expected StyleScreaming.IsCamel() to return false")
		}
	})
	t.Run("StyleLower", func(t *testing.T) {
		if caps.StyleLower.IsCamel() {
			t.Error("expected StyleLower.IsCamel() to return false")
		}
	})
	t.Run("StyleCamel", func(t *testing.T) {
		if !caps.StyleCamel.IsCamel() {
			t.Error("expected StyleCamel.IsCamel() to return true")
		}
	})
	t.Run("StyleLowerCamel", func(t *testing.T) {
		if caps.StyleLowerCamel.IsCamel() {
			t.Error("expected StyleLowerCamel.IsCamel() to return false")
		}
	})
}

func TestStyleIsLowerCamel(t *testing.T) {
	t.Run("StyleNotSpecified", func(t *testing.T) {
		if caps.StyleNotSpecified.IsLowerCamel() {
			t.Error("expected StyleNotSpecified.IsLowerCamel() to return false")
		}
	})
	t.Run("StyleScreaming", func(t *testing.T) {
		if caps.StyleScreaming.IsLowerCamel() {
			t.Error("expected StyleScreaming.IsLowerCamel() to return false")
		}
	})
	t.Run("StyleLower", func(t *testing.T) {
		if caps.StyleLower.IsLowerCamel() {
			t.Error("expected StyleLower.IsLowerCamel() to return false")
		}
	})
	t.Run("StyleCamel", func(t *testing.T) {
		if caps.StyleCamel.IsLowerCamel() {
			t.Error("expected StyleCamel.IsLowerCamel() to return false")
		}
	})
	t.Run("StyleLowerCamel", func(t *testing.T) {
		if !caps.StyleLowerCamel.IsLowerCamel() {
			t.Error("expected StyleLowerCamel.IsLowerCamel() to return true")
		}
	})
}

func TestWithConverter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if caps.WithConverter(nil).Converter != nil {
			t.Error("expected WithConverter(nil) to return a caps.Opts with a nil Converter")
		}
	})
	t.Run("not nil", func(t *testing.T) {
		c := caps.NewConverter(nil, nil, nil)
		if caps.WithConverter(c).Converter != c {
			t.Error("expected WithConverter(c) to return an Opts with Converter == c")
		}
	})
}

func TestWithReplaceStyle(t *testing.T) {
	if caps.WithReplaceStyle(caps.ReplaceStyleNotSpecified).ReplaceStyle != caps.ReplaceStyleNotSpecified {
		t.Error("expected ReplaceStyleNotSpecified")
	}
	if caps.WithReplaceStyle(caps.ReplaceStyleCamel).ReplaceStyle != caps.ReplaceStyleCamel {
		t.Error("expected ReplaceStyleCamel")
	}
}

func TestWithReplaceStyleLower(t *testing.T) {
	if caps.WithReplaceStyleLower().ReplaceStyle != caps.ReplaceStyleLower {
		t.Error("expected ReplaceStyleLower")
	}
}

func TestWithReplaceStyleScreaming(t *testing.T) {
	if caps.WithReplaceStyleScreaming().ReplaceStyle != caps.ReplaceStyleScreaming {
		t.Error("expected ReplaceStyleScreaming")
	}
}

func TestWithReplaceStyleCamel(t *testing.T) {
	if caps.WithReplaceStyleCamel().ReplaceStyle != caps.ReplaceStyleCamel {
		t.Error("expected ReplaceStyleCamel")
	}
}
