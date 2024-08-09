package builder

import (
	"strconv"
	"strings"
)

func (b *Build) Select(tableName string) Builder {
	b.reset()
	b.tableName = tableName
	b.sQLType = typeSelect
	return b
}

func (b *Build) GroupBy(fields ...string) Builder {
	b.groupBy = fields
	return b
}

func (b *Build) OrderBy(fields ...string) Builder {
	b.orderBy = fields
	return b
}

func (b *Build) Limit(l int) Builder {
	b.limit = l
	return b
}

func (b *Build) Offset(o int) Builder {
	b.offset = o
	return b
}

func (b *Build) generateSelectSQL() (string, error) {
	builder := &strings.Builder{}
	builderConcat(
		builder,
		"SELECT ",
		b.getSelectFields(),
		" FROM ",
		b.fieldQuote, b.tableName, b.fieldQuote,
	)

	builderConcat(
		builder,
		b.generateJoins(),
	)

	builderConcat(
		builder,
		" WHERE ", b.generateWhere(b.where),
	)

	groupBySql := b.getGroupBy()
	if groupBySql != "" {
		builderConcat(
			builder,
			" GROUP BY ", groupBySql,
		)
	}

	orderBySql := b.getOrderBy()
	if orderBySql != "" {
		builderConcat(
			builder,
			" ORDER BY ", orderBySql,
		)
	}

	if b.limit > 0 {
		builderConcat(
			builder,
			" LIMIT ", strconv.Itoa(b.limit),
		)
	}

	if b.offset > 0 {
		builderConcat(
			builder,
			" OFFSET ", strconv.Itoa(b.offset),
		)
	}

	return builder.String(), nil
}

func (b *Build) getSelectFields() string {
	if len(b.fields) == 0 {
		return "*"
	}

	return b.getFieldList(b.fields)
}

func (b *Build) getGroupBy() string {
	if len(b.groupBy) == 0 {
		return ""
	}

	return b.getFieldList(b.groupBy)
}

func (b *Build) getOrderBy() string {
	if len(b.orderBy) == 0 {
		return ""
	}

	return b.getFieldList(b.orderBy)
}

func (b *Build) getSelectParams() []interface{} {
	var pars []interface{}
	for _, join := range b.joins {
		if join.where != nil {
			pars = append(pars, b.getWhereParams(join.where)...)
		}
	}

	return append(pars, b.getWhereParams(b.where)...)
}

func (b *Build) getFieldList(fl []string) string {
	strBuilder := &strings.Builder{}
	for i, fn := range fl {
		if i > 0 {
			strBuilder.WriteString(",")
		}
		builderConcat(
			strBuilder,
			b.fieldQuote, fn, b.fieldQuote,
		)
	}

	return strBuilder.String()
}
