/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Chance Dinkins <chanceusc@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, Subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or Substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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
	fmt.Println(caps.ToLowerCamel("entity id"))
	fmt.Println(caps.ToLowerCamel("entity id", caps.WithReplaceStyleCamel()))
	// Output:
	// thisIsAnExampleID32
	// entityID
	// entityId
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
	fmt.Println(caps.ToDotNotation("This is [an] {example}${id32}."))
	// Output:
	// this.is.an.example.id.32
}

func ExampleToScreamingDotNotation() {
	fmt.Println(caps.ToScreamingDotNotation("This is [an] {example}${id32}."))
	// Output:
	// THIS.IS.AN.EXAMPLE.ID.32
}

func ExampleToDelimited() {
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.#32", ".", true))
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.break32", ".", false))
	fmt.Println(caps.ToDelimited("This is [an] {example}${id}.v32", ".", true, caps.Opts{AllowedSymbols: "$"}))

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
