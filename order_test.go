package sb

import (
	"bytes"
	"testing"
)

func TestOrderBysWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	table1 := Table{name: "test"}
	table2 := Table{name: "test2", alias: "t2"}

	tests := []struct {
		orderBys  OrderBys
		aliasMode AliasMode
		expected  string
	}{
		{
			expected: "",
		},
		{
			orderBys:  OrderBys{{&Column{name: "col1"}, false}},
			aliasMode: NoAlias,
			expected:  " ORDER BY `col1`",
		},
		{
			orderBys:  OrderBys{{&Column{name: "col1"}, true}},
			aliasMode: NoAlias,
			expected:  " ORDER BY `col1` DESC",
		},
		{
			orderBys:  OrderBys{{&Column{name: "col1"}, true}, {&Column{name: "col2"}, false}},
			aliasMode: NoAlias,
			expected:  " ORDER BY `col1` DESC, `col2`",
		},
		{
			orderBys:  OrderBys{{&Column{name: "col1"}, true}, {&Column{name: "col2"}, false}},
			aliasMode: UseAlias,
			expected:  " ORDER BY `col1` DESC, `col2`",
		},
		{
			orderBys:  OrderBys{{&Column{name: "col1", table: &table1}, true}, {&Column{name: "col2", table: &table2}, false}},
			aliasMode: UseAlias,
			expected:  " ORDER BY `test`.`col1` DESC, `t2`.`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.orderBys.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
