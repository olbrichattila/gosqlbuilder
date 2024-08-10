package builder

import (
	"strings"
)

// Delete initiates a DELETE FROM SQL
func (b *Build) Delete(tableName string) Builder {
	b.reset()
	b.tableName = tableName
	b.sQLType = typeDelete
	return b
}

func (b *Build) generateDeleteSQL() (string, error) {
	builder := &strings.Builder{}
	builderConcat(
		builder,
		"DELETE FROM ", b.fieldQuote, b.tableName, b.fieldQuote,
	)

	builderConcat(
		builder,
		" ", tokenWhere, " ", b.generateWhere(b.where),
	)

	return builder.String(), nil
}
