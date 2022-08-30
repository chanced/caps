# caps is a case conversion library

[![Go Reference](https://pkg.go.dev/badge/github.com/chanced/caps.svg)](https://pkg.go.dev/github.com/chanced/caps)
[![GoReportCard example](https://goreportcard.com/badge/github.com/chanced/caps)](https://goreportcard.com/report/github.com/chanced/caps)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/chanced/caps)

caps is a unicode aware case conversion library
for Go. It was built with the following priorites in mind: configurability,
consistency, correctness, ergonomic, and reasonable performance; in that order.

Out of the box, the following case conversion are supported:

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
	fmt.Println(caps.ToLowerCamel("SomeID", caps.Opts{ReplaceStyle: caps.ReplaceStyleLower}))
	// Output:
	// someId
}

```

[playground link](https://go.dev/play/p/DDO-ipHfgoH)

## Word boundaries

Word boundaries are determined by the `caps.Converter`. The provided implementation, `caps.ConverterImpl`,
delegates the boundary detection to `caps.Tokenizer`. The provided implementation, `caps.TokenizerImpl`,
uses the following runes as delimiters: `" _.!?:;$-(){}[]#@&+~"`.

`caps.StdConverter` also allows users to register `caps.Replacement`s for acronym replacement. The default list is:

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
options.

You can pass a new instance of `caps.StdConverter` with a new set of
`caps.Replacement` (likely preferred).

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
        fmt.Println(caps.ToCamel("ex id", caps.Opts{ Converter: converter }))
        // Output:
        // "EXId" // note: ID was not in the replacement list above
       fmt.Println(caps.ToCamel("ex id"))
        // Output:
        // ExID
    }
```

You can update `caps.DefaultConverter`. Just be aware that the
module was not built with thread-safety in mind so you should set it once.
Otherwise, you'll need guard your usage of the library accordingly.

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

Finally, if you are so inclined, you can create your own `caps.Converter`. This
could be as simple as implementing the single `Convert` method, calling
`caps.DefaultConverter.Convert`, and then modifying the result.

### Support for special case unicode (e.g. Turkish, Azeri)

caps supports Turkish and Azeri through the `token.Caser` interface. It is
satisfied by `unicode.TurkishCase` and `unicode.AzeriCase`. `token.TurkishCaser`
and `token.AzeriCaser` are available as pointers to those variables (although
you can use the unicode variables directly).

For example, to use Turkish, you would need to instantiate a few variables:

```go
package main
import (
    "github.com/chanced/caps"
    "github.com/chanced/token"

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
    // **BEFORE YOU USE IT** (caps.ConverterImpl is not guarded for thread-safety)
    //
    // caps.DefaultConverter = turkish
    //
    // otherwise, you can pass in the converter to the config for each call:

    fmt.Println(caps.ToCamel("An example", caps.Opts{Converter: turkish}))
}

```

## License

MIT
