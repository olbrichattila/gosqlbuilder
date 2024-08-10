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

	whereSQL := b.generateWhere(b.where)
	if whereSQL != "" {
		builderConcat(
			builder,
			" ", tokenWhere, " ", whereSQL,
		)
	}

	return builder.String(), nil
}
