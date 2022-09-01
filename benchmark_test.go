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
