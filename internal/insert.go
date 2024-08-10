package builder

import (
	"fmt"
	"strings"
)

// Insert initiates a new INSERT INTO SQL
func (b *Build) Insert(tableName string) Builder {
	b.reset()
	b.tableName = tableName
	b.sQLType = typeInsert
	return b
}

// Fields add fields for insert into or other SQL types
func (b *Build) Fields(fields ...string) Builder {
	b.fields = fields
	return b
}

// Values are adding the binding values for Fields
func (b *Build) Values(values ...interface{}) Builder {
	b.values = values
	return b
}

func (b *Build) generateInsertSQL() (string, error) {
	valueCount := len(b.values)
	if len(b.fields) != valueCount {
		return "", fmt.Errorf("select filed and value count does not match")
	}

	if valueCount == 0 {
		return "", fmt.Errorf("at least one field need to be inserted")
	}

	builder := &strings.Builder{}

	builderConcat(
		builder,
		"INSERT INTO ", b.fieldQuote, b.tableName, b.fieldQuote,
		" (", b.getSelectFields(), ")",
		" VALUES (?",
		strings.Repeat(",?", valueCount-1),
		")",
	)

	return builder.String(), nil
}
