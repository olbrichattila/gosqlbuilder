# Golang SQL Builder

This is an experimental SQL builder component, work in progress, not fully tested!


### Usage:

## Select
This example is not a valid SQL, but shows all kind of combinations you can use in your select
- Select("table_name")
- Fields("field1", "field2") (if not set, the builder returns '*'
- Join("table2", "table1.id", "table2.table1_id", func(w Where) {
			w.Where("join1", "=", "join2").
				Where("join3", "=", "join4")
		})
- LeftJoin("table1", "table1.id", "table2.table1_id", func(w Where) {
			w.Where("join1", "=", "join2").
				Where("join3", "=", "join4")
		})
- RightJoin("table1", "table1.id", "table2.table1_id", func(w Where) {
			w.Where("join1", "=", "join2").
				Where("join3", "=", "join4")
		}).        

- Where("f1", "=", 1)
- OrWhere("f2", "=", 2)
- Where("f3", "=", 3)
- WhereGroup(func(w Where) {}) (where w is another where builder instance, capable of doing nested where statements /in brackets/)
- Between("orb", "orb1", "orb2")
- OrBetween("orb", "orb1", "orb2")
- GroupBy("f1", "f2", "f3")
- OrderBy("f5", "f99", "f44")
- Limit(10)
- Offset(100)
```
builder := sqlbuilder.New()
sql, err := builder.
    Select("table1").
    Fields("field1", "field2").
    Join("table2", "table1.id", "table2.table1_id", func(w Where) {
        w.Where("join1", "=", "join2").
            Where("join3", "=", "join4")
    }).
    LeftJoin("table1", "table1.id", "table2.table1_id", func(w Where) {
        w.Where("join1", "=", "join2").
            Where("join3", "=", "join4")
    }).
    RightJoin("table1", "table1.id", "table2.table1_id", func(w Where) {
        w.Where("join1", "=", "join2").
            Where("join3", "=", "join4")
    }).
    Where("f1", "=", 1).
    OrWhere("f2", "=", 2).
    Where("f3", "=", 3).
    WhereGroup(func(w Where) {
        w.Where("sf4", "=", 4).
            Where("sf5", "=", 5).
            Between("btw", "bt1", "bt2").
            WhereGroup(func(w Where) {
                w.Where("ssf6", "=", 6).
                    OrWhere("ssf7", "=", 7).
                    OrWhere("ssf8", "<=", 8).
                    OrWhereGroup(func(w Where) {
                        w.Where("ssssf9", "=", 9)
                    })
            })
    }).
    Where("SL1", ">", 10).
    OrBetween("orb", "orb1", "orb2").
    GroupBy("f1", "f2", "f3").
    OrderBy("f5", "f99", "f44").
    Limit(10).
    Offset(100).
    AsSQL()


bindParams := builder.GetParams()
```

## Insert
```
builder := sqlbuilder.New()
sql, err := builder.Insert("table").
    Fields("f1", "f2").
    Values(1, 5).
    AsSQL()

bindParams := builder.GetParams()        
```

## Delete
```
builder := sqlbuilder.New()
sql, err := builder.Delete("table").
    Where("id", "=", 5).
    AsSQL()

bindParams := builder.GetParams()
```

## Update
```
builder := sqlbuilder.New()
sql, err := builder.Update("table").
    Fields("f1", "f2").
    Values(1, 2).
    Where("id", "=", 5).
    AsSQL()

bindPars := builder.GetParams()
```

## Is null
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    IsNull("field2").
    AsSQL()


whereParams := builder.GetParams()
```
## Is not null
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    IsNotNull("field2").
    AsSQL()

t.Nil(err)
t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND `field2` IS NOT NULL")

whereParams := builder.GetParams()
t.Len(whereParams, 1)
```
## Is null in where group
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    WhereGroup(func(w Where) {
        w.IsNotNull("field2")
        w.IsNotNull("field3")
    }).
    AsSQL()


whereParams := builder.GetParams()
```
## Is not null
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    OrIsNull("field2").
    AsSQL()

whereParams := builder.GetParams()
```
## Is is not null
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    OrIsNotNull("field2").
    AsSQL()

t.Nil(err)
t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? OR `field2` IS NOT NULL")

whereParams := builder.GetParams()
t.Len(whereParams, 1)
```
## Or is not null in where group
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    WhereGroup(func(w Where) {
        w.OrIsNotNull("field2")
        w.OrIsNotNull("field3")
    }).
    AsSQL()

whereParams := builder.GetParams()
```	
## In
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    In("field2", 1, 2, 3).
    AsSQL()


whereParams := builder.GetParams()
```
## In in where group
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    WhereGroup(func(w Where) {
        w.In("field2", 5, 8)
        w.In("field2", 3, 2)
    }).
    AsSQL()

whereParams := builder.GetParams()
```
## Not in
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    NotIn("field2", 1, 2, 3).
    AsSQL()

whereParams := builder.GetParams()
```
## Not in in where group
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    WhereGroup(func(w Where) {
        w.NotIn("field2", 5, 8)
        w.NotIn("field2", 3, 2)
    }).
    AsSQL()

whereParams := builder.GetParams()
```
## Or not in
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    OrNotIn("field2", 1, 2, 3).
    AsSQL()

whereParams := builder.GetParams()
````
## Or not in in where group
```
builder := New()
sql, err := builder.
    Select("table1").
    Where("field1", "=", 5).
    OrWhereGroup(func(w Where) {
        w.OrNotIn("field2", 5, 8)
        w.OrNotIn("field2", 3, 2)
    }).
    AsSQL()

whereParams := builder.GetParams()
```

## Raw select, where and orWhere (fields ar not quoted, so functions can be used, like count(*))
Example:
```
Select("table1").
		RawFields("count(*) as cnt", "item_id").
		RawWhere("field1", "=", 5).
		RawOrWhere("field2", "=", 5).
		OrWhereGroup(func(w Where) {
			w.OrNotIn("field3", 5, 8)
			w.OrNotIn("field3", 3, 2)
		}).
		AsSQL()
```

> Where can be used in any combination as in the select SQL shown, for update and delete SQLs as well.
