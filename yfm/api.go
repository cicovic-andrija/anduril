package yfm

import (
	"errors"
	"io"
)

// YAML Front Matter parser.

// Named errors.
var (
	ErrNotFound     = errors.New("YAML front matter not found")
	ErrInvalidInput = errors.New("input is not in valid format")
)

// Parse unmarshalles YAML Front Matter metadata with --- delimiters into v.
func Parse(r io.Reader, v interface{}) error {
	return newParser(r).parse(v)
}
