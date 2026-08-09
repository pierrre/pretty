# Pretty

Go pretty print library.

[![Go Reference](https://pkg.go.dev/badge/github.com/pierrre/pretty.svg)](https://pkg.go.dev/github.com/pierrre/pretty)

## Features

- [Pretty print value](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
- [String](https://pkg.go.dev/github.com/pierrre/pretty#String) / [Write](https://pkg.go.dev/github.com/pierrre/pretty#Write) / [Formatter](https://pkg.go.dev/github.com/pierrre/pretty#Formatter)
- [Configuration](https://pkg.go.dev/github.com/pierrre/pretty#CommonWriter):
  - [Indentation](https://pkg.go.dev/github.com/pierrre/pretty#Printer.Indent)
  - [Max depth](https://pkg.go.dev/github.com/pierrre/pretty#MaxDepthWriter)
  - [Unwrap interfaces](https://pkg.go.dev/github.com/pierrre/pretty#UnwrapInterfaceWriter)
  - [Recursion protection](https://pkg.go.dev/github.com/pierrre/pretty#RecursionWriter)
  - [Type filtering](https://pkg.go.dev/github.com/pierrre/pretty#FilterWriter)
  - [String](https://pkg.go.dev/github.com/pierrre/pretty#StringWriter)
  - [Slice](https://pkg.go.dev/github.com/pierrre/pretty#SliceWriter)
  - [Map](https://pkg.go.dev/github.com/pierrre/pretty#MapWriter)
- [Modular design](https://pkg.go.dev/github.com/pierrre/pretty#ValueWriter) (you can replace everything with your own implementation):
  - [`time`](https://pkg.go.dev/github.com/pierrre/pretty#TimeWriter)
  - [`error`](https://pkg.go.dev/github.com/pierrre/pretty#ErrorWriter)
  - [`[]byte` hex dump](https://pkg.go.dev/github.com/pierrre/pretty#BytesHexDumpWriter)
  - [`math/big`](https://pkg.go.dev/github.com/pierrre/pretty#MathBigWriter)
  - [`reflect`](https://pkg.go.dev/github.com/pierrre/pretty#ReflectWriter)
  - [`weak.Pointer`](https://pkg.go.dev/github.com/pierrre/pretty#WeakPointerWriter)
  - [`iter.Seq` / `iter.Seq2`](https://pkg.go.dev/github.com/pierrre/pretty#IterWriter)
  - [`Range` method (e.g. `sync.Map`)](https://pkg.go.dev/github.com/pierrre/pretty#RangeWriter)
  - [`fmt.Stringer`](https://pkg.go.dev/github.com/pierrre/pretty#StringerWriter)
  - [`fmt.GoStringer`](https://pkg.go.dev/github.com/pierrre/pretty#GoStringerWriter)
- [Extensions](https://pkg.go.dev/github.com/pierrre/pretty/ext/):
  - [`protobuf`](https://pkg.go.dev/github.com/pierrre/pretty/ext/protobuf/#example-package)
- Fast and (almost) no memory allocation

## Usage

[Example](https://pkg.go.dev/github.com/pierrre/pretty#example-package)
