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
		return "", fmt.Errorf("select filed and value count does not match")
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
			b.fieldQuote, fn, b.fieldQuote, "=?",
		)
	}

	builderConcat(
		builder,
		" WHERE ", b.generateWhere(b.where),
	)

	return builder.String(), nil
}
