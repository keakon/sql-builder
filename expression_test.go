package sb

import (
	"bytes"
	"testing"
)

func TestExpressionsWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		expressions Expressions
		aliasMode   AliasMode
		expected    string
	}{
		{
			expected: "",
		},
		{
			expressions: []Expression{Expr("1")},
			aliasMode:   NoAlias,
			expected:    "1",
		},
		{
			expressions: []Expression{Expr("1"), Expr("2")},
			aliasMode:   NoAlias,
			expected:    "1, 2",
		},
		{
			expressions: []Expression{Column{name: "col1"}, Column{name: "col2"}},
			aliasMode:   NoAlias,
			expected:    "`col1`, `col2`",
		},
		{
			expressions: []Expression{Expr("1"), Column{name: "col1", alias: "c1"}, Column{name: "col2", table: &Table{name: "test"}}},
			aliasMode:   UseAlias,
			expected:    "1, `c1`, `test`.`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.expressions.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}

func TestConcatExpressionsWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		expressions ConcatExpressions
		aliasMode   AliasMode
		expected    string
	}{
		{
			expected: "",
		},
		{
			expressions: ConcatExpressions{Expr("1")},
			aliasMode:   NoAlias,
			expected:    "1",
		},
		{
			expressions: ConcatExpressions{Expr("1"), Expr("2")},
			aliasMode:   NoAlias,
			expected:    "12",
		},
		{
			expressions: ConcatExpressions{Column{name: "col1"}, Column{name: "col2"}},
			aliasMode:   NoAlias,
			expected:    "`col1``col2`",
		},
		{
			expressions: NewConcatExpressions(Expr("1"), Column{name: "col1", alias: "c1"}, Column{name: "col2", table: &Table{name: "test"}}),
			aliasMode:   UseAlias,
			expected:    "1`c1``test`.`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			test.expressions.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
