package sb

import (
	"bytes"
	"testing"
)

type TestTable struct {
	Table    `db:"test"`
	ID       Column `json:"id" db:"id"`
	UserName Column
	age      Column
}

func TestTableWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		name      string
		alias     string
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "*",
		},
		{
			name:      "test",
			alias:     "t",
			aliasMode: NoAlias,
			expected:  "*",
		},
		{
			name:      "test",
			alias:     "t",
			aliasMode: UseAlias,
			expected:  "`t`.*",
		},
		{
			name:      "test",
			alias:     "",
			aliasMode: OnlyAlias,
			expected:  "`test`.*",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			table := Table{name: test.name, alias: test.alias}
			table.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestNewTable(t *testing.T) {
	aliases := []string{"", "test", "t"}
	for _, alias := range aliases {
		t.Run(alias, func(t *testing.T) {
			table := New[TestTable](alias)
			if table.name != "test" {
				t.Errorf("got %s, want %s", table.name, "test")
			}
			if table.alias != alias {
				t.Errorf("got %s, want %s", table.alias, alias)
			}
			if table.ID.name != "id" {
				t.Errorf("got %s, want %s", table.ID.name, "id")
			}
			if *table.ID.table != table.Table {
				t.Errorf("got %v, want %v", *table.ID.table, table.Table)
			}
			if table.UserName.name != "username" {
				t.Errorf("got %s, want %s", table.ID.name, "username")
			}
			if *table.UserName.table != table.Table {
				t.Errorf("got %v, want %v", *table.UserName.table, table.Table)
			}
			if table.age.name != "" {
				t.Errorf("got %s, want %s", table.age.name, "")
			}
			if table.age.table != nil {
				t.Errorf("got %v, want nil", table.age.table)
			}
		})
	}
}
