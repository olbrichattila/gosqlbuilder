package builder

import (
	"strings"
)

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
		" WHERE ", b.generateWhere(b.where),
	)

	return builder.String(), nil
}
