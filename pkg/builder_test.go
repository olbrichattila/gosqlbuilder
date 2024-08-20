// Package builder is an SQL string builder
package builder

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (t *TestSuite) TestComplexSelect() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Fields("field1", "field2").
		Join("table2", "table1.id", "table2.table1_id", func(w Where) {
			w.Where("join1", "=", "join2").
				Where("join3", "=", "join4")
		}).
		LeftJoin("table1l", "table1l.id", "table2l.table1_id", func(w Where) {
			w.Where("join1l", "=", "join2l").
				Where("join3l", "=", "join4l")
		}).
		RightJoin("table1r", "table1r.id", "table2r.table1_id", func(w Where) {
			w.Where("join1r", "=", "join2r").
				Where("join3r", "=", "join4r")
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

	t.Nil(err)
	t.Equal(sql, "SELECT `field1`,`field2` FROM `table1` JOIN `table2` ON `table1.id`=`table2.table1_id`  AND `join1`=? AND `join3`=? LEFT JOIN `table1l` ON `table1l.id`=`table2l.table1_id`  AND `join1l`=? AND `join3l`=? RIGHT JOIN `table1r` ON `table1r.id`=`table2r.table1_id`  AND `join1r`=? AND `join3r`=? WHERE `f1`=? OR `f2`=? AND `f3`=? AND (`sf4`=? AND `sf5`=? AND `btw` BETWEEN ? AND ?  AND (`ssf6`=? OR `ssf7`=? OR `ssf8`<=? OR (`ssssf9`=?))) AND `SL1`>? OR `orb` BETWEEN ? AND ?  GROUP BY `f1`,`f2`,`f3` ORDER BY `f5`,`f99`,`f44` LIMIT 10 OFFSET 100")

	whereParams := builder.GetParams()
	t.Len(whereParams, 20)
}

func (t *TestSuite) TestSimpleInsert() {
	builder := New()
	sql, err := builder.Insert("table").
		Fields("f1", "f2").
		Values(1, 5).
		AsSQL()

	t.Nil(err)

	t.Equal(sql, "INSERT INTO `table` (`f1`,`f2`) VALUES (?,?)")
	pars := builder.GetParams()
	t.Len(pars, 2)
}

func (t *TestSuite) TestSimpleDelete() {
	builder := New()
	sql, err := builder.Delete("table").
		Where("id", "=", 5).
		AsSQL()

	t.Equal(sql, "DELETE FROM `table` WHERE `id`=?")

	t.Nil(err)

	pars := builder.GetParams()
	t.Len(pars, 1)
}

func (t *TestSuite) TestSimpleUpdate() {
	builder := New()
	sql, err := builder.Update("table").
		Fields("f1", "f2").
		Values(1, 2).
		Where("id", "=", 5).
		AsSQL()

	t.Nil(err)

	t.Equal(sql, "UPDATE `table` SET `f1`=?,`f2`=? WHERE `id`=?")
	pars := builder.GetParams()
	t.Len(pars, 3)
}

func (t *TestSuite) TestIsNull() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		IsNull("field2").
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND `field2` IS NULL")

	whereParams := builder.GetParams()
	t.Len(whereParams, 1)
}

func (t *TestSuite) TestIsNotNull() {
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
}

func (t *TestSuite) TestIsNullInWhereGroup() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		WhereGroup(func(w Where) {
			w.IsNotNull("field2")
			w.IsNotNull("field3")
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND (`field2` IS NOT NULL AND `field3` IS NOT NULL)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 1)
}

func (t *TestSuite) TestOrIsNull() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		OrIsNull("field2").
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? OR `field2` IS NULL")

	whereParams := builder.GetParams()
	t.Len(whereParams, 1)
}

func (t *TestSuite) TestOrIsNotNull() {
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
}

func (t *TestSuite) TestOrIsNullInWhereGroup() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		WhereGroup(func(w Where) {
			w.OrIsNotNull("field2")
			w.OrIsNotNull("field3")
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND (`field2` IS NOT NULL OR `field3` IS NOT NULL)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 1)
}

func (t *TestSuite) TestIn() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		In("field2", 1, 2, 3).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND `field2` IN (?,?,?)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 4)
}

func (t *TestSuite) TestInWhereGroup() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		WhereGroup(func(w Where) {
			w.In("field2", 5, 8)
			w.In("field2", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND (`field2` IN (?,?) AND `field2` IN (?,?))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 5)
}

func (t *TestSuite) TestNotIn() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		NotIn("field2", 1, 2, 3).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND `field2` NOT IN (?,?,?)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 4)
}

func (t *TestSuite) TestNotInWhereGroup() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		WhereGroup(func(w Where) {
			w.NotIn("field2", 5, 8)
			w.NotIn("field2", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? AND (`field2` NOT IN (?,?) AND `field2` NOT IN (?,?))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 5)
}

func (t *TestSuite) TestOrNotIn() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		OrNotIn("field2", 1, 2, 3).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? OR `field2` NOT IN (?,?,?)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 4)
}

func (t *TestSuite) TestOrNotInWhereGroup() {
	builder := New()
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		OrWhereGroup(func(w Where) {
			w.OrNotIn("field2", 5, 8)
			w.OrNotIn("field2", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM `table1` WHERE `field1`=? OR (`field2` NOT IN (?,?) OR `field2` NOT IN (?,?))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 5)
}

func (t *TestSuite) TestPostgresFlavour() {
	builder := New()
	builder.SetSQLFlavour(FlavourPgSQL)
	sql, err := builder.
		Select("table1").
		Where("field1", "=", 5).
		OrWhereGroup(func(w Where) {
			w.OrNotIn("field2", 5, 8)
			w.OrNotIn("field2", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT * FROM \"table1\" WHERE \"field1\"=$1 OR (\"field2\" NOT IN ($2,$3) OR \"field2\" NOT IN ($4,$5))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 5)
}

func (t *TestSuite) TestPostgresFlavourWithInsert() {
	builder := New()
	builder.SetSQLFlavour(FlavourPgSQL)

	sql, err := builder.
		Insert("users").
		Fields("name", "email", "password").
		Values(1, 2, 3).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "INSERT INTO \"users\" (\"name\",\"email\",\"password\") VALUES ($1,$2,$3)")

	whereParams := builder.GetParams()
	t.Len(whereParams, 3)
}

func (t *TestSuite) TestRawFields() {
	builder := New()
	sql, err := builder.
		Select("table1").
		RawFields("count(*) as cnt", "item_id").
		Where("field1", "=", 5).
		OrWhereGroup(func(w Where) {
			w.OrNotIn("field2", 5, 8)
			w.OrNotIn("field2", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT count(*) as cnt,item_id FROM `table1` WHERE `field1`=? OR (`field2` NOT IN (?,?) OR `field2` NOT IN (?,?))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 5)
}

func (t *TestSuite) TestRawWhere() {
	builder := New()
	sql, err := builder.
		Select("table1").
		RawFields("count(*) as cnt", "item_id").
		RawWhere("field1", "=", 5).
		RawOrWhere("field2", "=", 5).
		OrWhereGroup(func(w Where) {
			w.OrNotIn("field3", 5, 8)
			w.OrNotIn("field3", 3, 2)
		}).
		AsSQL()

	t.Nil(err)
	t.Equal(sql, "SELECT count(*) as cnt,item_id FROM `table1` WHERE field1=? OR field2=? OR (`field3` NOT IN (?,?) OR `field3` NOT IN (?,?))")

	whereParams := builder.GetParams()
	t.Len(whereParams, 6)
}
