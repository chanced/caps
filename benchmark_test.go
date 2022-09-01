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
	"testing"

	"github.com/chanced/caps"
)

var testCase string = "Example Uuid."

var c caps.Caps = caps.New()

func BenchmarkCapsToTitle(b *testing.B) {
	var s string
	expected := "Example UUID"
	for n := 0; n < b.N; n++ {
		s = c.ToTitle(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToCamel(b *testing.B) {
	var s string
	expected := "ExampleUUID"
	for n := 0; n < b.N; n++ {
		s = c.ToCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToLowerCamel(b *testing.B) {
	var s string
	expected := "exampleUUID"
	for n := 0; n < b.N; n++ {
		s = c.ToLowerCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToSnake(b *testing.B) {
	var s string
	expected := "example_uuid"
	for n := 0; n < b.N; n++ {
		s = c.ToSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingSnake(b *testing.B) {
	var s string
	expected := "EXAMPLE_UUID"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToKebab(b *testing.B) {
	var s string
	expected := "example-uuid"
	for n := 0; n < b.N; n++ {
		s = c.ToKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingKebab(b *testing.B) {
	var s string
	expected := "EXAMPLE-UUID"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToDotNotation(b *testing.B) {
	var s string
	expected := "example.uuid"
	for n := 0; n < b.N; n++ {
		s = c.ToDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingDotNotation(b *testing.B) {
	var s string
	expected := "EXAMPLE.UUID"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// 							Func benchmark
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

func BenchmarkToTitle(b *testing.B) {
	var s string
	expected := "Example UUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToTitle(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToCamel(b *testing.B) {
	var s string
	expected := "ExampleUUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToLowerCamel(b *testing.B) {
	var s string
	expected := "exampleUUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToLowerCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToSnake(b *testing.B) {
	var s string
	expected := "example_uuid"
	for n := 0; n < b.N; n++ {
		s = caps.ToSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingSnake(b *testing.B) {
	var s string
	expected := "EXAMPLE_UUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToKebab(b *testing.B) {
	var s string
	expected := "example-uuid"
	for n := 0; n < b.N; n++ {
		s = caps.ToKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingKebab(b *testing.B) {
	var s string
	expected := "EXAMPLE-UUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToDotNotation(b *testing.B) {
	var s string
	expected := "example.uuid"
	for n := 0; n < b.N; n++ {
		s = caps.ToDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingDotNotation(b *testing.B) {
	var s string
	expected := "EXAMPLE.UUID"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}
