package sb

import (
	"bytes"
	"testing"
)

func TestConditionWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	table := Table{name: "test"}
	table2 := Table{name: "test", alias: "t"}

	tests := []struct {
		op        string
		lv        Expression
		rv        Expression
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "",
		},
		{
			op:        "=",
			lv:        Expr("1"),
			rv:        Expr("1"),
			aliasMode: NoAlias,
			expected:  "1 = 1",
		},
		{
			op:        "=",
			lv:        Expr("1"),
			rv:        PH,
			aliasMode: NoAlias,
			expected:  "1 = ?",
		},
		{
			op:        "=",
			lv:        Expr("1"),
			rv:        nil,
			aliasMode: NoAlias,
			expected:  "1 IS NULL",
		},
		{
			op:        "!=",
			lv:        Expr("1"),
			rv:        nil,
			aliasMode: NoAlias,
			expected:  "1 IS NOT NULL",
		},
		{
			op:        ">",
			lv:        Expr("1"),
			rv:        nil,
			aliasMode: NoAlias,
			expected:  "1 > NULL",
		},
		{
			op:        "=",
			lv:        Column{name: "col1"},
			rv:        Expr("1"),
			aliasMode: NoAlias,
			expected:  "`col1` = 1",
		},
		{
			op:        "=",
			lv:        Column{name: "col1"},
			rv:        Column{name: "col2"},
			aliasMode: NoAlias,
			expected:  "`col1` = `col2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1"},
			rv:        Column{name: "col2"},
			aliasMode: UseAlias,
			expected:  "`col1` = `col2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", alias: "c1"},
			rv:        Column{name: "col2", alias: "c2"},
			aliasMode: UseAlias,
			expected:  "`c1` = `c2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table},
			rv:        Column{name: "col2", table: &table2},
			aliasMode: UseAlias,
			expected:  "`test`.`col1` = `t`.`col2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", alias: "c1", table: &table},
			rv:        Column{name: "col2", alias: "c2", table: &table2},
			aliasMode: UseAlias, // 会被转成 OnlyAlias
			expected:  "`c1` = `c2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table},
			rv:        Column{name: "col2", table: &table2},
			aliasMode: OnlyAlias,
			expected:  "`test`.`col1` = `t`.`col2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", alias: "c1", table: &table},
			rv:        Column{name: "col2", alias: "c2", table: &table2},
			aliasMode: OnlyAlias,
			expected:  "`c1` = `c2`",
		},
		{
			op:        "=",
			lv:        Column{name: "col1"},
			rv:        nil,
			aliasMode: NoAlias,
			expected:  "`col1` IS NULL",
		},
		{
			op:        "=",
			lv:        Column{name: "col1"},
			rv:        PH,
			aliasMode: NoAlias,
			expected:  "`col1` = ?",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table},
			rv:        Expr("1"),
			aliasMode: NoAlias,
			expected:  "`col1` = 1",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table},
			rv:        Expr("1"),
			aliasMode: UseAlias,
			expected:  "`test`.`col1` = 1",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table2},
			rv:        Expr("1"),
			aliasMode: UseAlias,
			expected:  "`t`.`col1` = 1",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", alias: "c1", table: &table2},
			rv:        Expr("1"),
			aliasMode: UseAlias,
			expected:  "`c1` = 1",
		},
		{
			op:        "=",
			lv:        Column{name: "col1", table: &table},
			rv:        Column{name: "col2", alias: "c2", table: &table2},
			aliasMode: UseAlias,
			expected:  "`test`.`col1` = `c2`",
		},
		{
			op:        "IN",
			lv:        Column{name: "col1"},
			rv:        PH,
			aliasMode: NoAlias,
			expected:  "`col1` IN (?)",
		},
		{
			op:        "NOT IN",
			lv:        Column{name: "col1"},
			rv:        Expressions{Expr("1"), Expr("2")},
			aliasMode: NoAlias,
			expected:  "`col1` NOT IN (1, 2)",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			c := Condition{op: test.op, lv: test.lv, rv: test.rv}
			c.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestConditionsWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	c := &Column{name: "col"}

	tests := []struct {
		conditions Conditions
		expected   string
	}{
		{
			conditions: c.Ge(Expr("1")).And(c.Le(Expr("2"))),
			expected:   "(`col` >= 1 AND `col` <= 2)",
		},
		{
			conditions: And(c.Ge(Expr("1")), c.Le(Expr("2"))),
			expected:   "(`col` >= 1 AND `col` <= 2)",
		},
		{
			conditions: c.Ge(Expr("1")).And(c.Le(Expr("2"))).Not(),
			expected:   "(NOT (`col` >= 1 AND `col` <= 2))",
		},
		{
			conditions: c.Ge(Expr("1")).Or(c.Le(Expr("2"))),
			expected:   "(`col` >= 1 OR `col` <= 2)",
		},
		{
			conditions: Or(c.Ge(Expr("1")), c.Le(Expr("2"))),
			expected:   "(`col` >= 1 OR `col` <= 2)",
		},
		{
			conditions: Not(Or(c.Ge(Expr("1")), c.Le(Expr("2")))),
			expected:   "(NOT (`col` >= 1 OR `col` <= 2))",
		},
		{
			conditions: c.Ge(Expr("1")).Not().And(c.Le(Expr("2")).Or(c.Eq(PH).Not())),
			expected:   "((NOT `col` >= 1) AND (`col` <= 2 OR (NOT `col` = ?)))",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.conditions.WriteSQL(buf, NoAlias)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
