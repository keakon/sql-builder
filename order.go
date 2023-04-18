package sb

import "bytes"

type OrderBy struct {
	column *Column
	desc   bool
}

type OrderBys []OrderBy

func (o OrderBys) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	length := len(o)
	if length > 0 {
		buf.WriteString(" ORDER BY ")
		lastIndex := length - 1
		for i := 0; i < len(o); i++ {
			o[i].column.WriteSQL(buf, aliasMode)
			if o[i].desc {
				buf.WriteString(" DESC")
			}
			if i != lastIndex {
				buf.WriteString(", ")
			}
		}
	}
}
