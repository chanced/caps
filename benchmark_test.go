package caps_test

import (
	"testing"

	"github.com/chanced/caps"
)

var testCase string = "Example Uuid Test Case."

func BenchmarkToTitle(b *testing.B) {
	var s string
	expected := "Example UUID Test Case"
	for n := 0; n < b.N; n++ {
		s = caps.ToTitle(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToCamel(b *testing.B) {
	var s string
	expected := "ExampleUUIDTestCase"
	for n := 0; n < b.N; n++ {
		s = caps.ToCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToLowerCamel(b *testing.B) {
	var s string
	expected := "exampleUUIDTestCase"
	for n := 0; n < b.N; n++ {
		s = caps.ToLowerCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToSnake(b *testing.B) {
	var s string
	expected := "example_uuid_test_case"
	for n := 0; n < b.N; n++ {
		s = caps.ToSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingSnake(b *testing.B) {
	var s string
	expected := "EXAMPLE_UUID_TEST_CASE"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToKebab(b *testing.B) {
	var s string
	expected := "example-uuid-test-case"
	for n := 0; n < b.N; n++ {
		s = caps.ToKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingKebab(b *testing.B) {
	var s string
	expected := "EXAMPLE-UUID-TEST-CASE"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToDotNotation(b *testing.B) {
	var s string
	expected := "example.uuid.test.case"
	for n := 0; n < b.N; n++ {
		s = caps.ToDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkToScreamingDotNotation(b *testing.B) {
	var s string
	expected := "EXAMPLE.UUID.TEST.CASE"
	for n := 0; n < b.N; n++ {
		s = caps.ToScreamingDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// 							Instance benchmark
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

var c caps.Caps = caps.New()

func BenchmarkCapsToTitle(b *testing.B) {
	var s string
	expected := "Example UUID Test Case"
	for n := 0; n < b.N; n++ {
		s = c.ToTitle(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToCamel(b *testing.B) {
	var s string
	expected := "ExampleUUIDTestCase"
	for n := 0; n < b.N; n++ {
		s = c.ToCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToLowerCamel(b *testing.B) {
	var s string
	expected := "exampleUUIDTestCase"
	for n := 0; n < b.N; n++ {
		s = c.ToLowerCamel(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToSnake(b *testing.B) {
	var s string
	expected := "example_uuid_test_case"
	for n := 0; n < b.N; n++ {
		s = c.ToSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingSnake(b *testing.B) {
	var s string
	expected := "EXAMPLE_UUID_TEST_CASE"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingSnake(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToKebab(b *testing.B) {
	var s string
	expected := "example-uuid-test-case"
	for n := 0; n < b.N; n++ {
		s = c.ToKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingKebab(b *testing.B) {
	var s string
	expected := "EXAMPLE-UUID-TEST-CASE"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingKebab(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToDotNotation(b *testing.B) {
	var s string
	expected := "example.uuid.test.case"
	for n := 0; n < b.N; n++ {
		s = c.ToDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}

func BenchmarkCapsToScreamingDotNotation(b *testing.B) {
	var s string
	expected := "EXAMPLE.UUID.TEST.CASE"
	for n := 0; n < b.N; n++ {
		s = c.ToScreamingDotNotation(testCase)
	}

	if expected != s {
		b.Fatalf("Expected %s, got %s", expected, s)
	}
}
