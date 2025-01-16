// Package must provides helpers that panic.
package must

// NoError panics if err is not nil.
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}

// Handle panics if h is false.
func Handle(h bool) {
	if !h {
		panic("not handled")
	}
}
