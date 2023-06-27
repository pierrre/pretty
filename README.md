# Pretty

Pretty print Go values.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/pretty.svg)](https://pkg.go.dev/github.com/pierrre/pretty)

## Features

- [Pretty print value](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
- [String](https://pkg.go.dev/github.com/pierrre/pretty#String) / [Write](https://pkg.go.dev/github.com/pierrre/pretty#Write) / [Formatter](https://pkg.go.dev/github.com/pierrre/pretty#Formatter)
- [Config](https://pkg.go.dev/github.com/pierrre/pretty#Config): indentation, max depth, max length (string / map / slice), sort map keys
- Custom [value writer](https://pkg.go.dev/github.com/pierrre/pretty#ValueWriter) (`error`, `[]byte` hex dump, `fmt.Stringer`)
- No infinite recursion
- Fast and (almost) no memory allocation

## Usage

[Example](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
