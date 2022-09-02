# caps is a case conversion library

[![Go Reference](https://pkg.go.dev/badge/github.com/chanced/caps.svg)](https://pkg.go.dev/github.com/chanced/caps)
[![GoReportCard](https://goreportcard.com/badge/github.com/chanced/caps)](https://goreportcard.com/report/github.com/chanced/caps)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/chanced/caps)
![Build Status](https://github.com/chanced/caps/actions/workflows/main.yaml/badge.svg)
[![Latest Version](https://img.shields.io/github/v/tag/chanced/caps.svg?sort=semver&color=informational)](https://img.shields.io/github/v/tag/chanced/caps.svg?sort=semver&color=informational)

caps is a unicode aware, case conversion library for Go. It
was built with the following priorites in mind: configurability, consistency,
correctness, ergonomic, and reasonable performance; in that order.

It has built in support for configurable initialism replacements (e.g. "Uuid" ->
"UUID"), can support special case unicode (e.g. Turkish, Azeri, ...), and
finally is modular in that each step in the case conversion process can be
configured, augmented, or replaced.

The API should be stable as of v0.7.0.

## Install

```
go get github.com/chanced/caps
```

## Supported Case Conversions

The following case conversions are available:

-   Camel Case (e.g. CamelCase)
-   Lower Camel Case (e.g. lowerCamelCase)
-   Snake Case (e.g. snake_case)
-   Screaming Snake Case (e.g. SCREAMING_SNAKE_CASE)
-   Kebab Case (e.g. kebab-case)
-   Screaming Kebab Case(e.g. SCREAMING-KEBAB-CASE)
-   Dot Notation Case (e.g. dot.notation.case)
-   Screaming Dot Notation Case (e.g. DOT.NOTATION.CASE)
-   Title Case (e.g. Title Case)
-   Other deliminations

## Example

```go
package main
import (
	"fmt"
	"github.com/chanced/caps"
)
func main() {
	fmt.Println(caps.ToCamel("http request"))
	// Output:
	// HTTPRequest
	fmt.Println(caps.ToLowerCamel("some_id"))
	// Output:
	// someID
	fmt.Println(caps.ToLowerCamel("SomeID", caps.WithReplaceStyleCamel()))
	// Output:
	// someId

	// Alternatively:
	capsJS := caps.New(caps.CapsOpts{
		AllowedSymbols: "$",
		ReplaceStyle:   caps.ReplaceStyleCamel,
	})
	fmt.Println(capsJS.ToCamel("SomeID"))
	// Output:
	// someId
}
```

[go playground link](https://go.dev/play/p/8iSSUFlWTTw)

## Word boundaries

Word boundaries are determined by the `caps.Converter`. The provided implementation, `caps.StdConverter`,
delegates the boundary detection to `caps.Tokenizer`. The provided implementation, `caps.StdTokenizer`,
uses the following rules:

-   The following characters are considered word breaks `" _.!?:;$-(){}[]#@&+~"` unless present in `AllowedSymbols`
-   Strings with all upper case characters are split by the above symbols or by
    numbers, unless the character is allowed in a number based on the following rules:
    -   'v' or 'V' followed by numbers
    -   '.' before/after a number and only once
    -   'e' or 'E' if in the fractional part of a number and only once
    -   '-', '+' if at the start and followed by either a number or '.' and a
        number or in the fractional part proceeded by 'e' or 'E'
    -   additional rules can be added through the number rules (e.g. `WithNumberRules`)
    -   NOTE: If '.', '+', '-' are not in the `AllowedSymbols` they are
        considered breaks even for numbers
-   When a string consists of both upper case and lower case letters, upper case
    letters are considered boundaries (e.g. `"ThisVar"` would be tokenized into `["This", "Var"]`)
-   When mixed with lower and upper case characters, sequences of upper case are
    broken up into tokens (e.g. `"SomeID"` would be tokenized into `["Some", "I", "D"]`).
-   Replacement rules are then evaluated based on the tokens, which may
    combine them based on the rules below.

## Replacements

`caps.StdConverter` also allows users to register `caps.Replacement`s for
initialism replacements. Each `Replacement` is indexed in a trie (see
[Index](https://github.com/chanced/caps/blob/main/index/index.go)).

-   Multi-rune tokens are searched independently unless followed by a number (e.g.
    `"ID"`, `"UTF8"`).
-   Sequences of single rune tokens (e.g.`["U", "U", "I", "D"]`) are
    evaluated as a potential `Replacement` until a non-match is
    found or the sequence is broken by a token with more than one rune.

### Default replacements

```go
{"Acl", "ACL"}
{"Api", "API"}
{"Ascii", "ASCII"}
{"Cpu", "CPU"}
{"Css", "CSS"}
{"Dns", "DNS"}
{"Eof", "EOF"}
{"Guid", "GUID"}
{"Html", "HTML"}
{"Http", "HTTP"}
{"Https", "HTTPS"}
{"Id", "ID"}
{"Ip", "IP"}
{"Json", "JSON"}
{"Lhs", "LHS"}
{"Qps", "QPS"}
{"Ram", "RAM"}
{"Rhs", "RHS"}
{"Rpc", "RPC"}
{"Sla", "SLA"}
{"Smtp", "SMTP"}
{"Sql", "SQL"}
{"Ssh", "SSH"}
{"Tcp", "TCP"}
{"Tls", "TLS"}
{"Ttl", "TTL"}
{"Udp", "UDP"}
{"Ui", "UI"}
{"Uid", "UID"}
{"Uuid", "UUID"}
{"Uri", "URI"}
{"Url", "URL"}
{"Utf8", "UTF8"}
{"Vm", "VM"}
{"Xml", "XML"}
{"Xmpp", "XMPP"}
{"Xsrf", "XSRF"}
{"Xss", "XSS"}
```

If you would like to add or remove entries from that list, you have a few
options. See below.

## Customizing the `Converter`

### Using caps.Caps

This is likely your best option. You can create and utilize an instance of
`caps.Caps` which as all of the conversions as methods:

```go
package main

import (
	"fmt"

	"github.com/chanced/caps"
)

func main() {
	c := caps.New(caps.CapsOpts{
		Replacements: []caps.Replacement{
			{"Ex", "EX"},
			// ... your replacements
		},
	})
	fmt.Println(c.ToCamel("some ex"))
	// Output:
	// SomeEX
}
```

[go playground link](https://go.dev/play/p/QBJSq39sMCZ)

### Creating isolated `caps.StdConverter` instances

You can pass a new instance of `caps.StdConverter` with a new set of
`caps.Replacement`.

```go
    package main
    import (
        "fmt"
        "github.com/chanced/caps"
        "github.com/chanced/caps/token"
    )
    func main() {
        replacements := []caps.Replacement{
            {"Ex", "EX" },
            // ... your replacements
        }
        converter := caps.NewConverter(replacements, caps.DefaultTokenizer, token.DefaultCaser)
        fmt.Println(caps.ToCamel("ex id", caps.WithConverter(converter)))
        // note: ID was not in the replacement list above
        // Output:
        // "EXId"
       fmt.Println(caps.ToCamel("ex id"))
        // Output:
        // ExID
    }
```

[go playground link](https://go.dev/play/p/jlFAj3ujhTW)

### Modifying the `caps.DefaultConverter` global

You can update `caps.DefaultConverter`. You should set it before you make any
conversions. Otherwise, you'll need guard your usage of the library accordingly
(e.g. a mutex).

```go
package main
import (
    "fmt"
    "github.com/chanced/caps"
)
func main() {
    caps.DefaultConverter.Set("Gcp", "GCP")
    fmt.Println(caps.ToCamel("some_gcp_var"))
    // Output:
    // SomeGCPVar
}
```

[go playground link](https://go.dev/play/p/QARyN7-fUQ5)

### Creating a custom `caps.Converter`

Finally, if you are so inclined, you can create your own `caps.Converter`. This
could be as simple as implementing the single `Convert` method, calling
`caps.DefaultConverter.Convert`, and then modifying the result.

```go
package main
import (
    "fmt"
    "github.com/chanced/caps"
)
type MyConverter struct{}
func (MyConverter) Convert(req caps.ConvertRequest) string {
    res := caps.DefaultConverter.Convert(req)
    if req.Style.IsLowerCamel() && req.ReplaceStyle.IsCamel() && res == "id" {
        return "_id"
    }
    return res
}
func main() {
    fmt.Println(caps.ToLowerCamel("ID", caps.WithReplaceStyleCamel(), caps.WithConverter(MyConverter{})))
    // Output:
    // _id
}
```

[go playground link](https://go.dev/play/p/dg19iBIsHvh)

## Support for special case unicode (e.g. Turkish, Azeri)

caps supports Turkish and Azeri through the `token.Caser` interface. It is
satisfied by `unicode.TurkishCase` and `unicode.AzeriCase`. `token.TurkishCaser`
and `token.AzeriCaser` are available as pointers to those variables (although
you can use the unicode variables directly).

For example, to use Turkish, you would need to instantiate a few variables:

```go
package main
import (
    "github.com/chanced/caps"
    "github.com/chanced/caps/token"

)
func main() {
    tokenizer := caps.NewTokenizer(caps.DEFAULT_DELIMITERS, token.TurkishCaser)
    // I suppose these would need to be specific to Turkish?
    // if not, you can just use caps.DefaultReplacements
    replacements := []caps.Replacement{
        { Camel: "Http", Screaming: "HTTP" }, // just an example
    }
    turkish := caps.NewConverter(replacements, tokenizer, token.TurkishCaser)

    // to use this as your default throughout your application
    // you can overwrite caps.DefaultConverter
    //
    // caps.DefaultConverter = turkish
    //
    // otherwise, you can pass in the converter to the config for each call:
    fmt.Println(caps.ToScreamingKebab("i ı", caps.WithConverter(turkish)))
    // Output:
    // İ-I
}
```

[go playground link](https://go.dev/play/p/aKfuU5eZJgp)

## Benchmarks

```
input: "Example Uuid."
```

```
goos: darwin
goarch: arm64
pkg: github.com/chanced/caps
```

Using a `caps.Caps` instance:

```
BenchmarkCapsToTitle
BenchmarkCapsToTitle-10                   	 2418796	       473.9 ns/op	     208 B/op	       9 allocs/op
BenchmarkCapsToCamel
BenchmarkCapsToCamel-10                   	 2572858	       466.4 ns/op	     208 B/op	       9 allocs/op
BenchmarkCapsToLowerCamel
BenchmarkCapsToLowerCamel-10              	 2607831	       458.1 ns/op	     208 B/op	       9 allocs/op
BenchmarkCapsToSnake
BenchmarkCapsToSnake-10                   	 2567318	       466.2 ns/op	     208 B/op	       9 allocs/op
BenchmarkCapsToScreamingSnake
BenchmarkCapsToScreamingSnake-10          	 2380708	       504.7 ns/op	     256 B/op	      10 allocs/op
BenchmarkCapsToKebab
BenchmarkCapsToKebab-10                   	 2408059	       501.0 ns/op	     256 B/op	      10 allocs/op
BenchmarkCapsToScreamingKebab
BenchmarkCapsToScreamingKebab-10          	 2383544	       500.0 ns/op	     256 B/op	      10 allocs/op
BenchmarkCapsToDotNotation
BenchmarkCapsToDotNotation-10             	 2418474	       495.8 ns/op	     256 B/op	      10 allocs/op
BenchmarkCapsToScreamingDotNotation
BenchmarkCapsToScreamingDotNotation-10    	 2375731	       504.0 ns/op	     256 B/op	      10 allocs/op
```

Using top-level functions:

```
BenchmarkToTitle
BenchmarkToTitle-10                       	 2349061	       505.8 ns/op	     256 B/op	      10 allocs/op
BenchmarkToCamel
BenchmarkToCamel-10                       	 2406518	       503.7 ns/op	     256 B/op	      10 allocs/op
BenchmarkToLowerCamel
BenchmarkToLowerCamel-10                  	 2391112	       503.9 ns/op	     256 B/op	      10 allocs/op
BenchmarkToSnake
BenchmarkToSnake-10                       	 2392641	       503.3 ns/op	     256 B/op	      10 allocs/op
BenchmarkToScreamingSnake
BenchmarkToScreamingSnake-10              	 2311863	       517.7 ns/op	     256 B/op	      10 allocs/op
BenchmarkToKebab
BenchmarkToKebab-10                       	 2390890	       503.5 ns/op	     256 B/op	      10 allocs/op
BenchmarkToScreamingKebab
BenchmarkToScreamingKebab-10              	 2313615	       519.0 ns/op	     256 B/op	      10 allocs/op
BenchmarkToDotNotation
BenchmarkToDotNotation-10                 	 2344407	       510.3 ns/op	     256 B/op	      10 allocs/op
BenchmarkToScreamingDotNotation
BenchmarkToScreamingDotNotation-10        	 2294145	       521.2 ns/op	     256 B/op	      10 allocs/op
```

## License

MIT
