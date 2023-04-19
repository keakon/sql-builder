package sb

import (
	"bytes"
	"testing"
)

func TestOperationWriteSQL(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))

	tests := []struct {
		op        string
		lv        Expression
		rv        Expression
		aliasMode AliasMode
		expected  string
	}{
		{
			op:        "+",
			lv:        Column{name: "col1"},
			rv:        Expr("1"),
			aliasMode: NoAlias,
			expected:  "`col1`+1",
		},
		{
			op:        "-",
			lv:        Column{name: "col1"},
			rv:        Placeholder,
			aliasMode: NoAlias,
			expected:  "`col1`-?",
		},
		{
			op:        "*",
			lv:        Column{name: "col1"},
			rv:        Column{name: "col2"},
			aliasMode: NoAlias,
			expected:  "`col1`*`col2`",
		},
		{
			op:        "/",
			lv:        Column{name: "col1", alias: "c1"},
			rv:        Column{name: "col2"},
			aliasMode: UseAlias,
			expected:  "`c1`/`col2`",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			buf.Reset()
			o := Operation{op: test.op, lv: test.lv, rv: test.rv}
			o.WriteSQL(buf, test.aliasMode)
			if got := buf.String(); got != test.expected {
				t.Errorf("got %s, want %s", got, test.expected)
			}
		})
	}
}
