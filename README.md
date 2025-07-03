# Pretty

Go pretty print library.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/pretty.svg)](https://pkg.go.dev/github.com/pierrre/pretty)

## Features

- [Pretty print value](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
- [String](https://pkg.go.dev/github.com/pierrre/pretty#String) / [Write](https://pkg.go.dev/github.com/pierrre/pretty#Write) / [Formatter](https://pkg.go.dev/github.com/pierrre/pretty#Formatter)
- [Configuration](https://pkg.go.dev/github.com/pierrre/pretty#CommonWriter):
  - [Indentation](https://pkg.go.dev/github.com/pierrre/pretty#Printer.Indent)
  - [Max depth](https://pkg.go.dev/github.com/pierrre/pretty#MaxDepthWriter)
  - [String](https://pkg.go.dev/github.com/pierrre/pretty#StringWriter)
  - [Slice](https://pkg.go.dev/github.com/pierrre/pretty#SliceWriter)
  - [Map](https://pkg.go.dev/github.com/pierrre/pretty#MapWriter)
- [Modular design](https://pkg.go.dev/github.com/pierrre/pretty#ValueWriter) (you can replace everything with your own implementation)
  - [`error`](https://pkg.go.dev/github.com/pierrre/pretty#ErrorWriter)
  - [`[]byte` hex dump](https://pkg.go.dev/github.com/pierrre/pretty#BytesHexDumpWriter)
  - [`fmt.Stringer`](https://pkg.go.dev/github.com/pierrre/pretty#StringerWriter)
  - [`protobuf`](https://pkg.go.dev/github.com/pierrre/pretty/ext/protobuf/#example_)
- [No infinite recursion](https://pkg.go.dev/github.com/pierrre/pretty#RecursionWriter)
- Fast and (almost) no memory allocation

## Usage

[Example](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
