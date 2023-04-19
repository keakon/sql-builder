package sb

import "bytes"

type Operation struct {
	op string
	lv Expression
	rv Expression
}

func (o Operation) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	if o.lv == nil { // 正常情况不会遇到，除非手动构建
		return
	}

	o.lv.WriteSQL(buf, aliasMode)
	buf.WriteString(o.op)
	if o.rv == nil {
		buf.WriteString("NULL")
	} else {
		o.rv.WriteSQL(buf, aliasMode)
	}
}
