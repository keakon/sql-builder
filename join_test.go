package sb

import (
	"bytes"
	"testing"
)

func TestFromTablesWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	table1 := Table{name: "test"}
	table2 := Table{name: "test2", alias: "t2"}

	tests := []struct {
		table     AnyTable
		joins     []Join
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "",
		},
		{
			table:     table1,
			aliasMode: NoAlias,
			expected:  " FROM `test`",
		},
		{
			table:     table2,
			aliasMode: NoAlias,
			expected:  " FROM `test2`",
		},
		{
			table:     table2,
			aliasMode: UseAlias,
			expected:  " FROM `test2` AS `t2`",
		},
		{
			table:     table1,
			joins:     []Join{{table: table2, on: (&Column{name: "col1"}).Eq(Column{name: "col2"})}},
			aliasMode: NoAlias,
			expected:  " FROM `test` JOIN `test2` ON `col1` = `col2`",
		},
		{
			table:     table1,
			joins:     []Join{{table: table2, on: (&Column{name: "col1"}).Eq(Column{name: "col2"})}},
			aliasMode: UseAlias,
			expected:  " FROM `test` JOIN `test2` AS `t2` ON `col1` = `col2`",
		},
		{
			table:     table1,
			joins:     []Join{{table: table2, on: (&Column{name: "col1", table: &table1}).Eq(Column{name: "col2"})}},
			aliasMode: UseAlias,
			expected:  " FROM `test` JOIN `test2` AS `t2` ON `test`.`col1` = `col2`",
		},
		{
			table:     table1,
			joins:     []Join{{table: table2, on: (&Column{name: "col1", table: &table1}).Eq(Column{name: "col2", table: &table2})}},
			aliasMode: UseAlias,
			expected:  " FROM `test` JOIN `test2` AS `t2` ON `test`.`col1` = `t2`.`col2`",
		},
		{
			table: table1,
			joins: []Join{
				{table: table2, on: (&Column{name: "col1", table: &table1}).Eq(Column{name: "col2", table: &table2})},
				{typ: LeftJoin, table: table1, on: (&Column{name: "col1", table: &table2}).Eq(Column{name: "col2", table: &table1})},
			},
			aliasMode: UseAlias,
			expected:  " FROM `test` JOIN `test2` AS `t2` ON `test`.`col1` = `t2`.`col2` LEFT JOIN `test` ON `t2`.`col1` = `test`.`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			from := FromTables{test.table, test.joins}
			from.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestJoinWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	table1 := Table{name: "test"}
	table2 := Table{name: "test2", alias: "t2"}

	tests := []struct {
		typ       JoinType
		table     AnyTable
		on        Condition
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "",
		},
		{
			table:     table1,
			on:        (&Column{name: "col1"}).Eq(Column{name: "col2"}),
			aliasMode: NoAlias,
			expected:  " JOIN `test` ON `col1` = `col2`",
		},
		{
			typ:       LeftJoin,
			table:     table1,
			on:        (&Column{name: "col1"}).Eq(Column{name: "col2"}),
			aliasMode: NoAlias,
			expected:  " LEFT JOIN `test` ON `col1` = `col2`",
		},
		{
			typ:       RightJoin,
			table:     table1,
			on:        (&Column{name: "col1"}).Eq(Column{name: "col2"}),
			aliasMode: NoAlias,
			expected:  " RIGHT JOIN `test` ON `col1` = `col2`",
		},
		{
			typ:       OuterJoin,
			table:     table1,
			on:        (&Column{name: "col1"}).Eq(Column{name: "col2"}),
			aliasMode: NoAlias,
			expected:  " OUTER JOIN `test` ON `col1` = `col2`",
		},
		{
			table:     table1,
			on:        (&Column{name: "col1"}).Eq(Column{name: "col2"}),
			aliasMode: UseAlias,
			expected:  " JOIN `test` ON `col1` = `col2`",
		},
		{
			table:     table2,
			on:        (&Column{name: "col1", table: &table2}).Eq(Column{name: "col2"}),
			aliasMode: UseAlias,
			expected:  " JOIN `test2` AS `t2` ON `t2`.`col1` = `col2`",
		},
		{
			table:     table1,
			on:        (&Column{name: "col1", alias: "c1", table: &table2}).Eq(Column{name: "col2"}),
			aliasMode: UseAlias,
			expected:  " JOIN `test` ON `c1` = `col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			join := Join{typ: test.typ, table: test.table, on: test.on}
			join.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
