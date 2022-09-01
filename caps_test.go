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

func TestWithoutNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123.456", "."},
		{"a12bc3d", "abcd"},
		{"a", "a"},
		{"", ""},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.WithoutNumbers(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"anExampleString", "anExampleString"},
		{"ANEXAMPLESTRING", "aNEXAMPLESTRING"},
		{"AnExampleString, ", "anExampleString, "},
		{"a", "a"},
		{"", ""},
		{"123", "123"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.LowerFirst(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"anExampleString", "AnExampleString"},
		{"ANEXAMPLESTRING", "ANEXAMPLESTRING"},
		{"a", "A"},
		{"", ""},
		{"123", "123"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			output := caps.UpperFirst(test.input)
			if output != test.expected {
				t.Errorf("expected \"%s\", got \"%s\"", test.expected, output)
			}
		})
	}
}
