# Pretty

Go pretty print library.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/pretty.svg)](https://pkg.go.dev/github.com/pierrre/pretty)

## Features

- [Pretty print value](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
- [String](https://pkg.go.dev/github.com/pierrre/pretty#String) / [Write](https://pkg.go.dev/github.com/pierrre/pretty#Write) / [Formatter](https://pkg.go.dev/github.com/pierrre/pretty#Formatter)
- [Configuration](https://pkg.go.dev/github.com/pierrre/pretty#CommonValueWriter):
  - [Indentation](https://pkg.go.dev/github.com/pierrre/pretty#Config)
  - [Max depth](https://pkg.go.dev/github.com/pierrre/pretty#MaxDepthValueWriter)
  - [String](https://pkg.go.dev/github.com/pierrre/pretty#StringValueWriter)
  - [Slice](https://pkg.go.dev/github.com/pierrre/pretty#SliceValueWriter)
  - [Map](https://pkg.go.dev/github.com/pierrre/pretty#MapValueWriter)
- [Modular design](https://pkg.go.dev/github.com/pierrre/pretty#ValueWriter) (you can replace everything with your own implementation)
  - [`error`](https://pkg.go.dev/github.com/pierrre/pretty#ErrorValueWriter)
  - [`[]byte` hex dump](https://pkg.go.dev/github.com/pierrre/pretty#BytesHexDumpValueWriter)
  - [`fmt.Stringer`](https://pkg.go.dev/github.com/pierrre/pretty#StringerValueWriter)
  - [`pierrre/errors`](https://pkg.go.dev/github.com/pierrre/pretty/ext/pierrreerrors)
- [No infinite recursion](https://pkg.go.dev/github.com/pierrre/pretty#RecursionValueWriter)
- Fast and (almost) no memory allocation

## Usage

[Example](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
