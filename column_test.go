package sb

import (
	"bytes"
	"testing"
)

func TestColumnAs(t *testing.T) {
	c := Column{name: "test"}
	if c.alias != "" {
		t.Errorf("got %s, want %s", c.alias, "")
	}
	c.As("t")
	if c.alias != "t" {
		t.Errorf("got %s, want %s", c.alias, "t")
	}
}

func TestColumnWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	table := &Table{name: "test"}
	table2 := &Table{name: "test", alias: "t"}

	tests := []struct {
		name      string
		alias     string
		table     *Table
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "",
		},
		{
			name:      "col",
			alias:     "c",
			aliasMode: NoAlias,
			expected:  "`col`",
		},
		{
			name:      "col",
			alias:     "c",
			aliasMode: UseAlias,
			expected:  "`c`",
		},
		{
			name:      "col",
			alias:     "c",
			table:     table,
			aliasMode: UseAlias,
			expected:  "`test`.`col` AS `c`",
		},
		{
			name:      "col",
			alias:     "c",
			table:     table2,
			aliasMode: UseAlias,
			expected:  "`t`.`col` AS `c`",
		},
		{
			name:      "col",
			alias:     "c",
			table:     table,
			aliasMode: OnlyAlias,
			expected:  "`c`",
		},
		{
			name:      "col",
			alias:     "",
			aliasMode: UseAlias,
			expected:  "`col`",
		},
		{
			name:      "col",
			alias:     "",
			table:     &Table{name: "test"},
			aliasMode: NoAlias,
			expected:  "`col`",
		},
		{
			name:      "col",
			alias:     "",
			table:     &Table{name: "test"},
			aliasMode: UseAlias,
			expected:  "`test`.`col`",
		},
		{
			name:      "col",
			alias:     "",
			table:     &Table{name: "test", alias: "t"},
			aliasMode: UseAlias,
			expected:  "`t`.`col`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			c := Column{name: test.name, alias: test.alias, table: test.table}
			c.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestColumnsWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		columns   Columns
		aliasMode AliasMode
		expected  string
	}{
		{
			columns:   Columns{Column{name: "col1"}},
			aliasMode: NoAlias,
			expected:  "`col1`",
		},
		{
			columns:   Columns{Column{name: "col1"}, Column{name: "col2"}},
			aliasMode: NoAlias,
			expected:  "`col1`, `col2`",
		},
		{
			columns:   Columns{Column{name: "col1", alias: "c1"}, Column{name: "col2", table: &Table{name: "test"}}},
			aliasMode: UseAlias,
			expected:  "`c1`, `test`.`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.columns.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
