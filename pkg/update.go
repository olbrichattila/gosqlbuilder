package builder

import (
	"fmt"
	"strings"
)

// Update initiates an SQL UPDATE statement
func (b *Build) Update(tableName string) Builder {
	b.reset()
	b.tableName = tableName
	b.sQLType = typeUpdate
	return b
}

func (b *Build) generateUpdateSQL() (string, error) {
	valueCount := len(b.values)
	if len(b.fields) != valueCount {
		return "", errFieldCountMismatch
	}

	if valueCount == 0 {
		return "", fmt.Errorf("at least one field need to be updated")
	}

	builder := &strings.Builder{}
	builderConcat(
		builder,
		"UPDATE ", b.fieldQuote, b.tableName, b.fieldQuote,
		" SET ",
	)

	for i, fn := range b.fields {
		if i > 0 {
			builder.WriteString(",")
		}
		builderConcat(
			builder,
			b.fieldQuote, fn, b.fieldQuote, "=", b.getBindingParameter(),
		)
	}

	whereSQL := b.generateWhere(b.where)
	if whereSQL != "" {
		builderConcat(
			builder,
			" ", tokenWhere, " ", whereSQL,
		)
	}

	return builder.String(), nil
}
