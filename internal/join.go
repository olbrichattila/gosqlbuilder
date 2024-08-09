package builder

import "strings"

type Joiner interface {
}

type Join struct {
	joinType  string
	tableName string
	leftCond  string
	rightCond string
	where     Where
}

func (b *Build) Join(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	return b.getJoinBuilder(joinTypeInner, tableName, leftCond, rightCond, fn)
}

func (b *Build) LeftJoin(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	return b.getJoinBuilder(joinTypeLeft, tableName, leftCond, rightCond, fn)
}

func (b *Build) RightJoin(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	return b.getJoinBuilder(joinTypeRight, tableName, leftCond, rightCond, fn)
}

func (b *Build) getJoinBuilder(joinType string, tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	where := NewBlankWhere()
	join := &Join{
		joinType:  joinType,
		tableName: tableName,
		leftCond:  leftCond,
		rightCond: rightCond,
		where:     where,
	}
	fn(where)
	b.joins = append(b.joins, join)
	return b
}

func (b *Build) generateJoins() string {
	builder := &strings.Builder{}
	for _, join := range b.joins {
		builder.WriteString(b.generateJoin(join))
	}
	return builder.String()

}

func (b *Build) generateJoin(j *Join) string {
	builder := &strings.Builder{}
	builderConcat(
		builder,
		" ", j.joinType, " ",
		b.fieldQuote, j.tableName, b.fieldQuote,
		" ON ",
		b.fieldQuote, j.leftCond, b.fieldQuote,
		"=",
		b.fieldQuote, j.leftCond, b.fieldQuote,
		j.joinType,
		" ",
		b.fieldQuote, j.tableName, b.fieldQuote,
		" ON ",
		b.fieldQuote, j.leftCond, b.fieldQuote,
		"=",
		b.fieldQuote, j.leftCond, b.fieldQuote,
	)
	if j.where != nil {
		builderConcat(builder, " AND ", b.generateWhere(j.where))

	}

	return builder.String()
}
