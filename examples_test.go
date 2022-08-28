package caps_test

func ExampleToCamel() {
	// fmt.Println(caps.ToCamel("This is [an] {example}${id32}."))
	// fmt.Println(caps.ToCamel("This is [an] {example}${id32}.break32"))
	// fmt.Println(caps.ToCamel("This example allows for $ symbols", caps.Opts{AllowedSymbols: "$"}))

	// customReplacer := caps.NewFormatter([]caps.R{{"Http", "HTTP"}, {"Https", "HTTPS"}})
	// fmt.Println(caps.ToCamel("No Id just http And Https", caps.Opts{Formatter: customReplacer}))

	// Outputx:
	// thisIsAnExampleID32
	// thisIsAnExampleID32Break32
	// thisExampleAllowsFor$symbols
	// noIdJustHTTPAndHTTPS
}

func ExampleToDelimited() {
	// fmt.Println(caps.ToDelimited("A # B _ C", '.', true))
	// fmt.Println(caps.ToDelimited("$id", '.', false))
	// fmt.Println(caps.ToDelimited("$id", '.', true, caps.Opts{AllowedSymbols: "$"}))
	// fmt.Println(caps.ToDelimited("fromCamelcaseString", '.', true))
	// Outputx:
	// a.b.c
	// ID
	// $id
	// from.camelcase.string
}

func ExampleToSnake() {
	// fmt.Println(caps.ToSnake("A long string with spaces"))
	// fmt.Println(caps.ToSnake(strings.ToLower("A_SCREAMING_SNAKE_MUST_BE_LOWERED_FIRST")))
	// fmt.Println(caps.ToSnake("$word", caps.Opts{AllowedSymbols: "$"}))
	// OutputX:
	// a_long_string_with_spaces
	// a_screaming_snake_must_be_lowered_first
	// $word
}
