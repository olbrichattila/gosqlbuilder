package builder

import (
	"strings"
)

const (
	operatorAnd = "AND"
	operatorOr  = "OR"
)

// Where creates SQL WHERE block
func (b *Build) Where(field, relation string, value interface{}) Builder {
	b.where.AppendItem(
		NewWhere(false, typeAnd, field, relation, value),
	)

	return b
}

// Where creates SQL WHERE block
func (b *Build) RawWhere(field, relation string, value interface{}) Builder {
	b.where.AppendItem(
		NewWhere(true, typeAnd, field, relation, value),
	)

	return b
}

// OrWhere creates SQL OrWhere block
func (b *Build) OrWhere(field, relation string, value interface{}) Builder {
	b.where.AppendItem(NewWhere(false, typeOr, field, relation, value))

	return b
}

// OrWhere creates SQL OrWhere block
func (b *Build) RawOrWhere(field, relation string, value interface{}) Builder {
	b.where.AppendItem(NewWhere(true, typeOr, field, relation, value))

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

// IsNull crete IS NULL SQL clause
func (b *Build) IsNull(fieldName string) Builder {
	where := NewIsNull(fieldName)
	b.where.AppendItem(where)

	return b
}

// IsNotNull crete IS NOT NULL SQL clause
func (b *Build) IsNotNull(fieldName string) Builder {
	where := NewIsNotNull(fieldName)
	b.where.AppendItem(where)

	return b
}

// OrIsNull crete OR IS NULL SQL clause
func (b *Build) OrIsNull(fieldName string) Builder {
	where := NewOrIsNull(fieldName)
	b.where.AppendItem(where)

	return b
}

// OrIsNotNull crete OR IS NOT NULL SQL clause
func (b *Build) OrIsNotNull(fieldName string) Builder {
	where := NewOrIsNotNull(fieldName)
	b.where.AppendItem(where)

	return b
}

// In creates SQL IN (?,?)
func (b *Build) In(fieldName string, pars ...interface{}) Builder {
	where := NewIn(fieldName, pars...)
	b.where.AppendItem(where)

	return b
}

// NotIn creates SQL NOT IN (?,?)
func (b *Build) NotIn(fieldName string, pars ...interface{}) Builder {
	where := NewNotIn(fieldName, pars...)
	b.where.AppendItem(where)

	return b
}

// OrIn creates SQL OR IN (?,?)
func (b *Build) OrIn(fieldName string, pars ...interface{}) Builder {
	where := NewOrIn(fieldName, pars...)
	b.where.AppendItem(where)

	return b
}

// OrNotIn creates SQL OR NOT IN (?,?)
func (b *Build) OrNotIn(fieldName string, pars ...interface{}) Builder {
	where := NewOrNotIn(fieldName, pars...)
	b.where.AppendItem(where)

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

			if item.GetIsRaw() {
				builderConcat(strBuilder, item.GetField())
			} else {
				builderConcat(
					strBuilder,
					b.fieldQuote, item.GetField(), b.fieldQuote,
				)
			}

			switch item.GetOperator() {
			case typeBetween, typeOrBetween:
				strBuilder.WriteString(" BETWEEN ")
				strBuilder.WriteString(b.getBindingParameter())
				strBuilder.WriteString(" AND ")
				strBuilder.WriteString(b.getBindingParameter())
				strBuilder.WriteString(" ")
			case typeIsNull, typeOrIsNull:
				strBuilder.WriteString(" IS NULL")
			case typeIsNotNull, typeOrIsNotNull:
				strBuilder.WriteString(" IS NOT NULL")
			case typeIn, typeOrIn:
				builderConcat(
					strBuilder,
					" IN (", b.getBindingParameter(), strings.Repeat(","+b.getBindingParameter(), len(item.GetInValues())-1), ") ",
				)
			case typeNotIn, typeOrNotIn:
				builderConcat(
					strBuilder,
					" NOT IN (", b.getBindingParameter(), strings.Repeat(","+b.getBindingParameter(), len(item.GetInValues())-1), ") ",
				)
			default:
				builderConcat(
					strBuilder,
					item.GetRelation(), b.getBindingParameter(),
				)
			}

		}
	}
	return strBuilder.String()
}

func (b *Build) getWhereOperator(t int) string {
	switch t {
	case typeAnd, typeBetween:
		return operatorAnd
	case typeOr, typeOrBetween, typeOrIsNotNull, typeOrIsNull, typeOrIn, typeOrNotIn:
		return operatorOr
	default:
		return operatorAnd
	}
}

func (b *Build) getWhereParams(w Where) []interface{} {
	var pars []interface{}

	for _, item := range w.GetItems() {
		if item.GetItems() != nil {
			pars = append(pars, b.getWhereParams(item)...)
		} else {
			switch item.GetOperator() {
			case typeIsNull, typeIsNotNull, typeOrIsNull, typeOrIsNotNull:
				// do nothing, no parameter
			case typeIn, typeNotIn, typeOrIn, typeOrNotIn:
				pars = append(pars, item.GetInValues()...)
			case typeBetween, typeOrBetween:
				pars = append(pars, item.GetValue(), item.GetValue2())
			default:
				pars = append(pars, item.GetValue())
			}
		}
	}

	return pars
}
