package sb

import (
	"bytes"
	"testing"
)

func TestAssignmentWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		column    Column
		value     Expression
		aliasMode AliasMode
		expected  string
	}{
		{
			column:    Column{name: "col1"},
			value:     Expr("1"),
			aliasMode: NoAlias,
			expected:  "`col1`=1",
		},
		{
			column:    Column{name: "col1"},
			value:     Placeholder,
			aliasMode: NoAlias,
			expected:  "`col1`=?",
		},
		{
			column:    Column{name: "col1"},
			value:     Column{name: "col2"},
			aliasMode: NoAlias,
			expected:  "`col1`=`col2`",
		},
		{
			column:    Column{name: "col1", alias: "c1"},
			value:     Column{name: "col2"},
			aliasMode: UseAlias,
			expected:  "`c1`=`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			a := &Assignment{column: &test.column, value: test.value}
			a.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestAssignmentsWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		assignments Assignments
		aliasMode   AliasMode
		expected    string
	}{
		{
			assignments: Assignments{{&Column{name: "col1"}, Expr("1")}},
			aliasMode:   NoAlias,
			expected:    "`col1`=1",
		},
		{
			assignments: Assignments{{&Column{name: "col1", alias: "c1"}, Placeholder}, {&Column{name: "col2"}, (&Column{name: "col2"}).Plus(Column{name: "col3"})}},
			aliasMode:   NoAlias,
			expected:    "`col1`=?, `col2`=`col2`+`col3`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.assignments.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
