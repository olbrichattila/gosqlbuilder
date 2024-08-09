package builder

const (
	typeAnd       = 0
	typeOr        = 1
	typeBetween   = 2
	typeOrBetween = 3
)

func NewBlankWhere() Where {
	return &Wh{}
}

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

func NewWhereGroup(
	operator int,
) Where {
	return &Wh{
		operator: operator,
	}
}

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

type WhereGroupFunc func(Where)

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

type Wh struct {
	operator int
	field    string
	relation string // TODO this could be a domain object verifying = <, > ...
	value    interface{}
	value2   interface{}
	items    []Where
}

func (w *Wh) Where(field, relation string, value interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		relation: relation,
		value:    value,
		operator: typeAnd,
	})
	return w
}

func (w *Wh) OrWhere(field, relation string, value interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		relation: relation,
		value:    value,
		operator: typeOr,
	})
	return w
}

func (w *Wh) Between(field string, value1, value2 interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		value:    value1,
		value2:   value2,
		operator: typeBetween,
	})
	return w
}

func (w *Wh) OrBetween(field string, value1, value2 interface{}) Where {
	w.items = append(w.items, &Wh{
		field:    field,
		value:    value1,
		value2:   value2,
		operator: typeOrBetween,
	})
	return w
}

func (w *Wh) WhereGroup(fn WhereGroupFunc) Where {
	where := &Wh{operator: typeAnd}

	w.items = append(w.items, where)
	fn(where)

	return w
}

func (w *Wh) OrWhereGroup(fn WhereGroupFunc) Where {
	where := &Wh{operator: typeOr}

	w.items = append(w.items, where)
	fn(where)

	return w
}

func (w *Wh) GetItems() []Where {
	return w.items
}

func (w *Wh) GetOperator() int {
	return w.operator
}

func (w *Wh) GetField() string {
	return w.field
}

func (w *Wh) GetRelation() string {
	return w.relation
}

func (w *Wh) GetValue() interface{} {
	return w.value
}

func (w *Wh) GetValue2() interface{} {
	return w.value2
}

func (w *Wh) AppendItem(wh Where) {
	w.items = append(w.items, wh)
}
