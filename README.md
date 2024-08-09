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

