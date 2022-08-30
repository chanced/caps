package caps_test

import (
	"fmt"

	"github.com/chanced/caps"
)

func ExampleToCamel() {
	fmt.Println(caps.ToCamel("This is [an] {example}${id32}."))
	fmt.Println(caps.ToCamel("AN_EXAMPLE_STRING"))
	// Output:
	// ThisIsAnExampleID32
	// AnExampleString
}

func ExampleToLowerCamel() {
	fmt.Println(caps.ToLowerCamel("This is [an] {example}${id32}."))
	// Output:
	// thisIsAnExampleID32
}

func ExampleToSnake() {
	fmt.Println(caps.ToSnake("This is [an] {example}${id32}."))
	fmt.Println(caps.ToSnake("v3.2.2"))
	// Output:
	// this_is_an_example_id_32
	// v3_2_2
}

func ExampleToScreamingSnake() {
	fmt.Println(caps.ToScreamingSnake("This is [an] {example}${id32}."))
	// Output:
	// THIS_IS_AN_EXAMPLE_ID_32
}

func ExampleToKebab() {
	fmt.Println(caps.ToKebab("This is [an] {example}${id32}."))
	// Output:
	// this-is-an-example-id-32
}

func ExampleToScreamingKebab() {
	fmt.Println(caps.ToScreamingKebab("This is [an] {example}${id32}."))
	// Output:
	// THIS-IS-AN-EXAMPLE-ID-32
}

func ExampleToDotNotation() {
	fmt.Println(caps.ToScreamingDot("This is [an] {example}${id32}."))
	// Output:
	// this.is.an.example.id.32
}

func ExampleToScreamingDotNotation() {
	fmt.Println(caps.ToScreamingDot("This is [an] {example}${id32}."))
	// Output:
	// THIS.IS.AN.EXAMPLE.ID.32
}

func ExampleToDelimited() {
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.#32", '.', true))
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.break32", '.', false))
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.v32", '.', true, caps.Opts{AllowedSymbols: "$"}))

	// Output:
	// this.is.an.example.id.32
	// THIS.IS.AN.EXAMPLE.ID.BREAK.32
	// this.is.an.example.$.id.v32
}

func ExampleToTitle() {
	fmt.Println(caps.ToTitle("This is [an] {example}${id32}."))
	// Output:
	// This Is An Example ID 32
}
