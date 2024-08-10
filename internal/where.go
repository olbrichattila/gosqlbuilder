package builder

import "fmt"

const (
	typeAnd                           = 0
	typeOr                            = 1
	typeBetween                       = 2
	typeOrBetween                     = 3
	typeIsNull                        = 4
	typeIsNotNull                     = 5
	typeOrIsNull                      = 6
	typeOrIsNotNull                   = 7
	typeIn                            = 8
	typeNotIn                         = 9
	typeOrIn                          = 10
	typeOrNotIn                       = 11
	tokenWhere                        = "WHERE"
	tokenOn                           = "ON"
	incorrectRelationshipPanicMessage = "provided relation %s is not valid"
)

// NewBlankWhere initiates a Where interface object with default values
func NewBlankWhere() Where {
	return &Wh{}
}

// NewWhere initiates a Where interface object propagating values for SQL WHERE statement
func NewWhere(
	operator int,
	field string,
	relation string,
	value interface{},
) Where {
	if !validateRelation(relation) {
		panic(fmt.Sprintf(incorrectRelationshipPanicMessage, relation))
	}

	return &Wh{
		operator: operator,
		field:    field,
		relation: relation,
		value:    value,
	}
}

// NewWhereGroup initiates a where group with AND or OR operator set
func NewWhereGroup(
	operator int,
) Where {
	return &Wh{
		operator: operator,
	}
}

// NewIsNull is to generate IS NULL
func NewIsNull(
	field string,
) Where {
	return &Wh{
		field:    field,
		operator: typeIsNull,
	}
}

// NewIsNotNull is to generate IS NOT NULL
func NewIsNotNull(
	field string,
) Where {
	return &Wh{
		field:    field,
		operator: typeIsNotNull,
	}
}

// NewOrIsNull is to generate OR IS NULL
func NewOrIsNull(
	field string,
) Where {
	return &Wh{
		field:    field,
		operator: typeOrIsNull,
	}
}

// NewOrIsNotNull is to generate OR IS NOT NULL
func NewOrIsNotNull(
	field string,
) Where {
	return &Wh{
		field:    field,
		operator: typeOrIsNotNull,
	}
}

// NewIn is to generate field_name IN (?,?)
func NewIn(
	field string,
	pars ...interface{},
) Where {
	return &Wh{
		field:    field,
		operator: typeIn,
		inValues: pars,
	}
}

// NewNotIn is to generate field_name NOT IN (?,?)
func NewNotIn(
	field string,
	pars ...interface{},
) Where {
	return &Wh{
		field:    field,
		operator: typeNotIn,
		inValues: pars,
	}
}

// NewOrIn is to generate OR field_name IN (?,?)
func NewOrIn(
	field string,
	pars ...interface{},
) Where {
	return &Wh{
		field:    field,
		operator: typeOrIn,
		inValues: pars,
	}
}

// NewOrNotIn is to generate OR field_name NOT IN (?,?)
func NewOrNotIn(
	field string,
	pars ...interface{},
) Where {
	return &Wh{
		field:    field,
		operator: typeOrNotIn,
		inValues: pars,
	}
}

// NewBetween creates a new Between where object with parameter required for an SQL BETVEEN ? ands ? statement
func NewBetween(
	operator int,
	field string,
	value interface{},
	value2 interface{},
) Where {
	return &Wh{
		operator: operator,
		field:    field,
		value:    value,
		value2:   value2,
	}
}

// WhereGroupFunc is the definition of recursive WHERE closure
type WhereGroupFunc func(Where)

// Where is the interface of recursive WHERE builder
type Where interface {
	Where(string, string, interface{}) Where
	OrWhere(string, string, interface{}) Where
	Between(string, interface{}, interface{}) Where
	OrBetween(string, interface{}, interface{}) Where
	WhereGroup(fn WhereGroupFunc) Where
	OrWhereGroup(fn WhereGroupFunc) Where
	IsNull(string) Where
	IsNotNull(string) Where
	OrIsNull(string) Where
	OrIsNotNull(string) Where
	In(string, ...interface{}) Where
	NotIn(string, ...interface{}) Where
	OrIn(string, ...interface{}) Where
	OrNotIn(string, ...interface{}) Where
	GetItems() []Where
	GetOperator() int
	GetField() string
	GetRelation() string
	GetValue() interface{}
	GetValue2() interface{}
	GetInValues() []interface{}
	AppendItem(Where)
}

// Wh is the structure behind the Where builder
type Wh struct {
	operator int
	field    string
	relation string // TODO this could be a domain object verifying = <, > ...
	value    interface{}
	value2   interface{}
	inValues []interface{}
	items    []Where
}

// Where creates SQL WHERE block
func (w *Wh) Where(field, relation string, value interface{}) Where {
	if !validateRelation(relation) {
		panic(fmt.Sprintf(incorrectRelationshipPanicMessage, relation))
	}

	w.items = append(w.items, &Wh{
		field:    field,
		relation: relation,
		value:    value,
		operator: typeAnd,
	})
	return w
}

// OrWhere creates SQL OrWhere block
func (w *Wh) OrWhere(field, relation string, value interface{}) Where {
	if !validateRelation(relation) {
		panic(fmt.Sprintf(incorrectRelationshipPanicMessage, relation))
	}

	w.items = append(w.items, &Wh{
		field:    field,
		relation: relation,
		value:    value,
		operator: typeOr,
	})
	return w
}

// Between creates SQL BETWEEN condition
func (w *Wh) Between(field string, value1, value2 interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		value:    value1,
		value2:   value2,
		operator: typeBetween,
	})
	return w
}

// OrBetween creates SQL BETWEEN with preceding OR operator
func (w *Wh) OrBetween(field string, value1, value2 interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		value:    value1,
		value2:   value2,
		operator: typeOrBetween,
	})
	return w
}

// WhereGroup creates a new groups of WHERE, lile WHERE `field` = ? and (`field2` = ?....). Provide the conditions in the closure where you get a Where builder
func (w *Wh) WhereGroup(fn WhereGroupFunc) Where {
	where := &Wh{operator: typeAnd}

	w.items = append(w.items, where)
	fn(where)

	return w
}

// OrWhereGroup creates a new groups of WHERE preceded by OR operator, like WHERE `field` = ? and (`field2` = ?....). Provide the conditions in the closure where you get a Where builder
func (w *Wh) OrWhereGroup(fn WhereGroupFunc) Where {
	where := &Wh{operator: typeOr}

	w.items = append(w.items, where)
	fn(where)

	return w
}

// IsNull generates where sql like AND `field` IS NULL
func (w *Wh) IsNull(fileName string) Where {
	where := &Wh{operator: typeIsNull, field: fileName}
	w.items = append(w.items, where)

	return w
}

// IsNotNull generates where sql like AND `field` IS NULL
func (w *Wh) IsNotNull(fileName string) Where {
	where := &Wh{operator: typeIsNotNull, field: fileName}
	w.items = append(w.items, where)

	return w
}

// OrIsNull generates where sql like OR `field` IS NULL
func (w *Wh) OrIsNull(fileName string) Where {
	where := &Wh{operator: typeOrIsNull, field: fileName}
	w.items = append(w.items, where)

	return w
}

// OrIsNotNull generates where sql like OR `field` IS NOT NULL
func (w *Wh) OrIsNotNull(fileName string) Where {
	where := &Wh{operator: typeOrIsNotNull, field: fileName}
	w.items = append(w.items, where)

	return w
}

// In generates where sql like AND `field` IN (?,?,?,?)
func (w *Wh) In(fileName string, pars ...interface{}) Where {
	where := &Wh{operator: typeIn, field: fileName, inValues: pars}
	w.items = append(w.items, where)

	return w
}

// NotIn generates where sql like AND `field` NOT IN (?,?,?,?)
func (w *Wh) NotIn(fileName string, pars ...interface{}) Where {
	where := &Wh{operator: typeNotIn, field: fileName, inValues: pars}
	w.items = append(w.items, where)

	return w
}

// OrIn generates where sql like AND OR `field` IN (?,?,?,?)
func (w *Wh) OrIn(fileName string, pars ...interface{}) Where {
	where := &Wh{operator: typeOrIn, field: fileName, inValues: pars}
	w.items = append(w.items, where)

	return w
}

// OrNotIn generates where sql like OR `field` NOT IN (?,?,?,?)
func (w *Wh) OrNotIn(fileName string, pars ...interface{}) Where {
	where := &Wh{operator: typeOrNotIn, field: fileName, inValues: pars}
	w.items = append(w.items, where)

	return w
}

// GetItems returns the child items of where
func (w *Wh) GetItems() []Where {
	return w.items
}

// GetOperator return the type of the operator (and, or)
func (w *Wh) GetOperator() int {
	return w.operator
}

// GetField returns the database field name set for this WHERE clause
func (w *Wh) GetField() string {
	return w.field
}

// GetRelation returns the relational operator between WHERE ? = ? like =, <, >, <=, >= ....
func (w *Wh) GetRelation() string {
	return w.relation
}

// GetValue returns the value for the binding params set by WHERE clause
func (w *Wh) GetValue() interface{} {
	return w.value
}

// GetValue2 returns the second binding param of BETWEEN clause
func (w *Wh) GetValue2() interface{} {
	return w.value2
}

// GetInValues returns the values for in clauses
func (w *Wh) GetInValues() []interface{} {
	return w.inValues
}

// AppendItem add a new WHERE builder object to the multiple and recursive WHERE blocks
func (w *Wh) AppendItem(wh Where) {
	w.items = append(w.items, wh)
}
