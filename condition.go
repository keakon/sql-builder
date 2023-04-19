package sb

import "bytes"

const (
	opIn    = "IN"
	opNotIn = "NOT IN"
	opEq    = "="
	opNe    = "!="
)

type Cond interface {
	WriteSQL(buf *bytes.Buffer, aliasMode AliasMode)
	isCond()
}

type Condition struct {
	op string
	lv Expression
	rv Expression
}

func (c Condition) isCond() {}

func (c Condition) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	if c.lv == nil { // 正常情况不会遇到，除非手动构建
		return
	}

	if aliasMode == UseAlias { // 表达式里不输出 AS ...
		aliasMode = OnlyAlias
	}

	c.lv.WriteSQL(buf, aliasMode)
	if c.rv == nil { // "= nil" -> "IS NULL", "!= nil" -> "IS NOT NULL"
		if c.op == opEq {
			buf.WriteString(" IS NULL")
			return
		} else if c.op == opNe {
			buf.WriteString(" IS NOT NULL")
			return
		}
	}

	buf.WriteByte(' ')
	buf.WriteString(c.op)
	buf.WriteByte(' ')
	if c.rv == Placeholder {
		if c.op == opIn || c.op == opNotIn {
			buf.WriteString("(?)")
		} else {
			buf.WriteByte('?')
		}
	} else {
		needBracket := c.op == opIn || c.op == opNotIn // IN、NOT IN 需要添加括号
		if !needBracket {
			_, needBracket = c.rv.(*SelectQuery) // 子查询需要添加括号
		}
		if needBracket {
			buf.WriteByte('(')
		}
		if c.rv == nil {
			buf.WriteString("NULL")
		} else {
			c.rv.WriteSQL(buf, aliasMode)
		}
		if needBracket {
			buf.WriteByte(')')
		}
	}
}

type Conditions struct {
	conditions    []Cond
	isDisjunction bool // true: OR, false: AND
	isTopLevel    bool
}

func (c Conditions) isCond() {}

func (c Conditions) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	length := len(c.conditions)
	if length > 0 {
		lastIndex := length - 1
		if c.isTopLevel {
			buf.WriteString(" WHERE ")
		} else {
			buf.WriteByte('(')
		}
		for i := 0; i < length; i++ {
			c.conditions[i].WriteSQL(buf, aliasMode)
			if i != lastIndex {
				if c.isDisjunction {
					buf.WriteString(" OR ")
				} else {
					buf.WriteString(" AND ")
				}
			}
		}
		if !c.isTopLevel {
			buf.WriteByte(')')
		}
	}
}

func And(conditions ...Cond) Conditions {
	return Conditions{
		isDisjunction: false,
		conditions:    conditions,
	}
}

func Or(conditions ...Cond) Conditions {
	return Conditions{
		isDisjunction: true,
		conditions:    conditions,
	}
}