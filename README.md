# caps

[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/chanced/caps)
[![GoReportCard example](https://goreportcard.com/badge/github.com/chanced/caps)](https://goreportcard.com/report/github.com/chanced/caps)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/chanced/caps)

caps is a case conversion library for Go. It was built with the following
priorites in mind: configurability, consistency, correctness, ergonomic, and
reasonable performance; in that order.

Out of the box, the following case conversion are supported:

-   Camel Case (e.g. CamelCase)
-   Lower Camel Case (e.g. lowerCamelCase)
-   Snake Case (e.g. snake_case)
-   Screaming Snake Case (e.g. SCREAMING_SNAKE_CASE)
-   Kebab Case (e.g. kebab-case)
-   Screaming Kebab Case(e.g. SCREAMING-KEBAB-CASE)
-   Dot Notation Case (e.g. dot.notation.case)
-   Title Case (e.g. Title Case)
-   Other deliminations

Word boundaries are determined by the `caps.Formatter`. The provided implementation, `caps.FormatterImpl`,
delegates the boundary detection to `caps.Tokenizer`. The provided implementation, `caps.TokenizerImpl`,
uses the following tokens as delimiters: `" _.!?:;$-(){}[]#@&+~"`.

`caps.StdFormatter` also allows users to register `caps.Replacement`s for acronym replacement. The default list is:

```go
{"Http", "HTTP"}
{"Https", "HTTPS"}
{"Id", "ID"}
{"Ip", "IP"}
{"Html", "HTML"}
{"Xml", "XML"}
{"Json", "JSON"}
{"Csv", "CSV"}
{"Aws", "AWS"}
{"Gcp", "GCP"}
{"Sql", "SQL"}
```

If you would like to add or remove entries from that list, you have a few
options.

You can pass a new instance of `caps.StdFormatter` with a new set of
`caps.Replacement` (likely preferred).

You can create your own `caps.Formatter`. This could be as simple as
implementing the single `Format` method, calling `caps.DefaultFormatter.Format`,
and then modifying the result.

Finally, if you are so inclined, you can update `caps.DefaultFormatter`. Just be aware that the
module was not built with thread-safety in mind so you should set it once.
Otherwise, you'll need guard your usage of the library accordingly.

## License

MIT
