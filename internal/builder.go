// Package builder is an sql builder
package builder

import (
	"fmt"
)

const (
	typeSelect    = 1
	typeInsert    = 2
	typeDelete    = 3
	typeUpdate    = 4
	joinTypeInner = "JOIN"
	joinTypeLeft  = "LEFT JOIN"
	joinTypeRight = "RIGHT JOIN"
)

// New creates new SQL builder
func New() Builder {
	return &Build{
		fieldQuote: "`",
		where:      NewBlankWhere(),
		joins:      make([]*Join, 0),
	}
}

type Builder interface {
	Where(field, relation string, value interface{}) Builder
	OrWhere(field, relation string, value interface{}) Builder
	Between(field string, value1, value2 interface{}) Builder
	OrBetween(field string, value1, value2 interface{}) Builder
	WhereGroup(fn WhereGroupFunc) Builder
	OrWhereGroup(fn WhereGroupFunc) Builder
	AsSQL() (string, error)
	GetParams() []interface{}
	Delete(tableName string) Builder
	Insert(tableName string) Builder
	Fields(fields ...string) Builder
	Values(values ...interface{}) Builder
	Join(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder
	LeftJoin(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder
	RightJoin(tableName, leftCond, rightCond string, fn WhereGroupFunc) Builder
	Select(tableName string) Builder
	GroupBy(fields ...string) Builder
	OrderBy(fields ...string) Builder
	Limit(l int) Builder
	Offset(o int) Builder
	Update(tableName string) Builder
}

// The Builder struct holding the builder logic
type Build struct {
	sQLType    int
	tableName  string
	fieldQuote string
	fields     []string
	values     []interface{}
	where      Where
	groupBy    []string
	orderBy    []string
	limit      int
	offset     int
	joins      []*Join
}

// AsSQL returns the SQL representation of the build SQL command
func (b *Build) AsSQL() (string, error) {
	switch b.sQLType {
	case typeSelect:
		return b.generateSelectSQL()
	case typeInsert:
		return b.generateInsertSQL()
	case typeDelete:
		return b.generateDeleteSQL()
	case typeUpdate:
		return b.generateUpdateSQL()
	default:
		return "", fmt.Errorf("invalid SQL type")
	}
}

func (b *Build) reset() {
	b.tableName = ""
	b.fields = make([]string, 0)
	b.groupBy = make([]string, 0)
	b.orderBy = make([]string, 0)
	b.values = make([]interface{}, 0)
}

func (b *Build) GetParams() []interface{} {
	switch b.sQLType {
	case typeSelect:
		return b.getSelectParams()
	case typeInsert:
		return b.values
	case typeDelete:
		return b.getWhereParams(b.where)
	case typeUpdate:
		return append(b.values, b.getWhereParams(b.where)...)
	default:
		return nil
	}
}
