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

	t.Nil(err)
	t.Equal(sql, "SELECT `field1`,`field2` FROM `table1` JOIN `table2` ON `table1.id`=`table1.id`JOIN `table2` ON `table1.id`=`table1.id` AND `join1`=? AND `join3`=? LEFT JOIN `table1` ON `table1.id`=`table1.id`LEFT JOIN `table1` ON `table1.id`=`table1.id` AND `join1`=? AND `join3`=? RIGHT JOIN `table1` ON `table1.id`=`table1.id`RIGHT JOIN `table1` ON `table1.id`=`table1.id` AND `join1`=? AND `join3`=? WHERE `f1`=? OR `f2`=? AND `f3`=? AND (`sf4`=? AND `sf5`=? AND `btw` BETWEEN ? AND ?  AND (`ssf6`=? OR `ssf7`=? OR `ssf8`<=? OR (`ssssf9`=?))) AND `SL1`>? OR `orb` BETWEEN ? AND ?  GROUP BY `f1`,`f2`,`f3` ORDER BY `f5`,`f99`,`f44` LIMIT 10 OFFSET 100")

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

	t.Equal(sql, "INSERT INTO `table` `f1`,`f2` VALUES (?,?)")
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
