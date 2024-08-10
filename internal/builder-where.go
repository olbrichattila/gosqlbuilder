package builder

import (
	"strings"
)

// Where creates SQL WHERE block
func (b *Build) Where(field, relation string, value interface{}) Builder {
	b.where.AppendItem(
		NewWhere(typeAnd, field, relation, value),
	)

	return b
}

// OrWhere creates SQL OrWhere block
func (b *Build) OrWhere(field, relation string, value interface{}) Builder {
	b.where.AppendItem(NewWhere(typeOr, field, relation, value))

	return b
}

// Between creates SQL BETWEEN condition
func (b *Build) Between(field string, value1, value2 interface{}) Builder {
	b.where.AppendItem(
		NewBetween(typeBetween, field, value1, value2),
	)

	return b
}

// OrBetween creates SQL BETWEEN with preceding OR operator
func (b *Build) OrBetween(field string, value1, value2 interface{}) Builder {
	b.where.AppendItem(
		NewBetween(typeOrBetween, field, value1, value2),
	)

	return b
}

// WhereGroup creates a new groups of WHERE, lile WHERE `field` = ? and (`field2` = ?....). Provide the conditions in the closure where you get a Where builder
func (b *Build) WhereGroup(fn WhereGroupFunc) Builder {
	where := NewWhereGroup(typeAnd)
	b.where.AppendItem(where)
	fn(where)

	return b
}

// OrWhereGroup creates a new groups of WHERE preceded by OR operator, like WHERE `field` = ? and (`field2` = ?....). Provide the conditions in the closure where you get a Where builder
func (b *Build) OrWhereGroup(fn WhereGroupFunc) Builder {
	where := NewWhereGroup(typeOr)
	b.where.AppendItem(where)

	fn(where)

	return b
}

func (b *Build) generateWhere(w Where) string {
	strBuilder := &strings.Builder{}
	isFirst := true
	for _, item := range w.GetItems() {
		operator := b.getWhereOperator(item.GetOperator())

		if item.GetItems() != nil {
			builderConcat(
				strBuilder,
				" ", operator, " (", b.generateWhere(item), ")",
			)

		} else {
			if isFirst {
				isFirst = false
			} else {
				builderConcat(
					strBuilder,
					" ", operator, " ",
				)
			}

			builderConcat(
				strBuilder,
				b.fieldQuote, item.GetField(), b.fieldQuote,
			)

			if item.GetOperator() == typeBetween || item.GetOperator() == typeOrBetween {
				strBuilder.WriteString(" BETWEEN ? AND ? ")
			} else {
				builderConcat(
					strBuilder,
					item.GetRelation(), "?",
				)
			}

		}
	}
	return strBuilder.String()
}

func (b *Build) getWhereOperator(t int) string {
	switch t {
	case typeAnd, typeBetween:
		return "AND"
	case typeOr, typeOrBetween:
		return "OR"
	default:
		return "AND"
	}
}

func (b *Build) getWhereParams(w Where) []interface{} {
	var pars []interface{}

	for _, item := range w.GetItems() {
		if item.GetItems() != nil {
			pars = append(pars, b.getWhereParams(item)...)
		} else {
			pars = append(pars, item.GetValue())
			if item.GetOperator() == typeBetween || item.GetOperator() == typeOrBetween {
				pars = append(pars, item.GetValue2())
			}
		}
	}

	return pars
}
