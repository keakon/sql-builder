package sb

import "bytes"

type Assignment struct {
	column *Column
	value  Expression
}

func (a Assignment) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	a.column.WriteSQL(buf, aliasMode)
	if a.value == nil {
		buf.WriteString("=NULL")
	} else {
		buf.WriteByte('=')
		a.value.WriteSQL(buf, aliasMode)
	}
}

type Assignments []Assignment

func (a Assignments) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	length := len(a)
	if length > 0 {
		lastIndex := length - 1
		for i := 0; i < length; i++ {
			a[i].WriteSQL(buf, aliasMode)
			if i != lastIndex {
				buf.WriteString(", ")
			}
		}
	}
}
