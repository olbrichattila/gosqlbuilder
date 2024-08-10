// Package sqlbuilder helps you compose SQL filed in a programmatic way
package sqlbuilder

import (
	builder "github.com/olbrichattila/gosqlbuilder/pkg"
)

// New creates New SQL builder instance
func New() builder.Builder {
	return builder.New()
}
