package builder

import (
	"strings"
)

// Join is the structure of a JOINS SQL clause
type Join struct {
	joinType  string
	tableName string
	leftCond  string
	rightCond string
	where     Where
}

// Join creates a table join clause, like JOIN `table1` ON `table1.id` = `table2.table1_id`
func (b *Build) Join(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	return b.getJoinBuilder(joinTypeInner, tableName, leftCond, rightCond, fn)
}

// LeftJoin creates a table left join clause, like LEFT JOIN `table1` ON `table1.id` = `table2.table1_id`
func (b *Build) LeftJoin(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder {
	return b.getJoinBuilder(joinTypeLeft, tableName, leftCond, rightCond, fn)
}

// RightJoin creates a table right join clause, like RIGHT JOIN `table1` ON `table1.id` = `table2.table1_id`
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
