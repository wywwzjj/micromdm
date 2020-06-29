// Package id implements a unique identifier.
package id

import (
	"crypto/rand"

	"github.com/oklog/ulid/v2"
)

// New returns a unique identifier.
// The format is a Universally Unique Lexicographically Sortable Identifier (ULID).
func New() string { return ulid.MustNew(ulid.Now(), rand.Reader).String() }
