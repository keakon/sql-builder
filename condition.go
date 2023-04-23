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

func (c Condition) And(cond Cond) Conditions {
	conds, ok := cond.(Conditions)
	if ok && !conds.isDisjunction { // 类型相同时合并
		conds.conditions = Prepend[Cond](conds.conditions, c)
		return conds
	}
	return Conditions{
		conditions: []Cond{c, cond},
	}
}

func (c Condition) Or(cond Cond) Conditions {
	conds, ok := cond.(Conditions)
	if ok && conds.isDisjunction { // 类型相同时合并
		conds.conditions = Prepend[Cond](conds.conditions, c)
		return conds
	}
	return Conditions{
		conditions:    []Cond{c, cond},
		isDisjunction: true,
	}
}

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
	if c.rv == PH {
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

func (c Conditions) And(cond Cond) Conditions {
	if c.isDisjunction {
		return Conditions{
			conditions: []Cond{c, cond},
		}
	}

	conds, ok := cond.(Conditions)
	if ok { // 类型相同时合并元素
		c.conditions = append(c.conditions, conds.conditions...)
	} else {
		c.conditions = append(c.conditions, cond)
	}
	return c
}

func (c Conditions) Or(cond Cond) Conditions {
	if !c.isDisjunction {
		return Conditions{
			conditions:    []Cond{c, cond},
			isDisjunction: true,
		}
	}

	conds, ok := cond.(Conditions)
	if ok { // 类型相同时合并元素
		c.conditions = append(c.conditions, conds.conditions...)
	} else {
		c.conditions = append(c.conditions, cond)
	}
	return c
}

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

// TODO: 合并相同类型
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
