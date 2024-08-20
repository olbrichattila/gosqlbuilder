// Package builder is an sql builder
package builder

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	typeSelect    = 1
	typeInsert    = 2
	typeDelete    = 3
	typeUpdate    = 4
	joinTypeInner = "JOIN"
	joinTypeLeft  = "LEFT JOIN"
	joinTypeRight = "RIGHT JOIN"

	FlavourSqLite      = 1
	FlavourMySQL       = 2
	FlavourPgSQL       = 3
	FlavourFirebirdSQL = 4
)

var (
	ErrInvalidSQLFlavour = errors.New("invalid SQL flavour")
)

// Builder is the base SQL builder interface
type Builder interface {
	SetSQLFlavour(int) error
	Where(field, relation string, value interface{}) Builder
	RawWhere(field, relation string, value interface{}) Builder
	OrWhere(field, relation string, value interface{}) Builder
	RawOrWhere(field, relation string, value interface{}) Builder
	Between(field string, value1, value2 interface{}) Builder
	OrBetween(field string, value1, value2 interface{}) Builder
	WhereGroup(fn WhereGroupFunc) Builder
	IsNull(string) Builder
	IsNotNull(string) Builder
	OrIsNull(string) Builder
	OrIsNotNull(string) Builder
	In(string, ...interface{}) Builder
	NotIn(string, ...interface{}) Builder
	OrIn(string, ...interface{}) Builder
	OrNotIn(string, ...interface{}) Builder
	OrWhereGroup(fn WhereGroupFunc) Builder
	AsSQL() (string, error)
	GetParams() []interface{}
	Delete(tableName string) Builder
	Insert(tableName string) Builder
	Fields(fields ...string) Builder
	RawFields(fields ...string) Builder
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

// New creates new SQL builder
func New() Builder {
	return &Build{
		fieldQuote:   "`",
		bindingStyle: "?",
		where:        NewBlankWhere(),
		joins:        make([]*Join, 0),
	}
}

// The Build struct holding the builder logic
type Build struct {
	sQLType        int
	tableName      string
	fieldQuote     string
	bindingStyle   string
	fields         []string
	fieldsAreRaw   bool
	values         []interface{}
	where          Where
	groupBy        []string
	orderBy        []string
	limit          int
	offset         int
	joins          []*Join
	parameterCount int
}

// SetSQLFlavour can set your preferred SQL engine, as they have different quotation mark and parameter binding
func (b *Build) SetSQLFlavour(sQLFlavour int) error {
	switch sQLFlavour {
	case FlavourSqLite, FlavourFirebirdSQL:
		b.fieldQuote = "\""
		b.bindingStyle = "?"
	case FlavourMySQL:
		b.fieldQuote = "`"
		b.bindingStyle = "?"
	case FlavourPgSQL:
		b.fieldQuote = "\""
		b.bindingStyle = "$"
	default:
		return ErrInvalidSQLFlavour
	}

	return nil
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
	b.parameterCount = 0
	b.fieldsAreRaw = false
	b.tableName = ""
	b.fields = make([]string, 0)
	b.values = make([]interface{}, 0)
	b.groupBy = make([]string, 0)
	b.orderBy = make([]string, 0)
	b.values = make([]interface{}, 0)
	b.where = NewBlankWhere()
	b.limit = 0
	b.offset = 0
	b.joins = make([]*Join, 0)
}

// GetParams returns the binding params for the last generated SQL
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

func (b *Build) getBindingParameter() string {
	b.parameterCount++

	if b.bindingStyle == "?" {
		return b.bindingStyle
	}
	return b.bindingStyle + strconv.Itoa(b.parameterCount)
}
