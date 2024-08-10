package builder

const (
	typeAnd       = 0
	typeOr        = 1
	typeBetween   = 2
	typeOrBetween = 3
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
	GetItems() []Where
	GetOperator() int
	GetField() string
	GetRelation() string
	GetValue() interface{}
	GetValue2() interface{}
	AppendItem(Where)
}

// Wh is the structure behind the Where builder
type Wh struct {
	operator int
	field    string
	relation string // TODO this could be a domain object verifying = <, > ...
	value    interface{}
	value2   interface{}
	items    []Where
}

// Where creates SQL WHERE block
func (w *Wh) Where(field, relation string, value interface{}) Where {
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

// OrWhere creates a new groups of WHERE preceded by OR operator, like WHERE `field` = ? and (`field2` = ?....). Provide the conditions in the closure where you get a Where builder
func (w *Wh) OrWhereGroup(fn WhereGroupFunc) Where {
	where := &Wh{operator: typeOr}

	w.items = append(w.items, where)
	fn(where)

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

// GetFields returns the database field name set for this WHERE clause
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

// AppendItem add a new WHERE builder object to the multiple and recursive WHERE blocks
func (w *Wh) AppendItem(wh Where) {
	w.items = append(w.items, wh)
}
