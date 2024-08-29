package pretty

import (
	"github.com/pierrre/pretty/internal/indent"
)

// DefaultConfig is the default [Config].
var DefaultConfig = NewConfig()

// Config is a configuration used to pretty print values.
//
// It should be created with [NewConfig].
type Config struct {
	// Indent is the string used to indent.
	// Default: "\t".
	Indent string
}

// NewConfig creates a new [Config] initialized with default values.
func NewConfig() *Config {
	return &Config{
		Indent: indent.Default,
	}
}
