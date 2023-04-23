package sb

import "bytes"

const (
	opIn    = "IN"
	opNotIn = "NOT IN"
	opEq    = "="
	opNe    = "!="
)

type boolOp uint8

const (
	and boolOp = iota
	or
	not
)

type Cond interface {
	WriteSQL(buf *bytes.Buffer, aliasMode AliasMode)
	And(cond Cond) Conditions
	Or(cond Cond) Conditions
}

type Condition struct {
	op string
	lv Expression
	rv Expression
}

func (c Condition) And(cond Cond) Conditions {
	conds, ok := cond.(Conditions)
	if ok && conds.op == and { // 类型相同时合并
		conds.conditions = Prepend[Cond](conds.conditions, c)
		return conds
	}
	return Conditions{
		conditions: []Cond{c, cond},
		op:         and,
	}
}

func (c Condition) Or(cond Cond) Conditions {
	conds, ok := cond.(Conditions)
	if ok && conds.op == or { // 类型相同时合并
		conds.conditions = Prepend[Cond](conds.conditions, c)
		return conds
	}
	return Conditions{
		conditions: []Cond{c, cond},
		op:         or,
	}
}

func (c Condition) Not() Conditions {
	return Conditions{
		conditions: []Cond{c},
		op:         not,
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
	conditions []Cond
	op         boolOp
	isTopLevel bool
}

func (c Conditions) And(cond Cond) Conditions {
	if c.op != and {
		return Conditions{
			conditions: []Cond{c, cond},
			op:         and,
		}
	}

	conds, ok := cond.(Conditions)
	if ok && c.op == conds.op { // 类型相同时合并元素
		c.conditions = append(c.conditions, conds.conditions...)
	} else {
		c.conditions = append(c.conditions, cond)
	}
	return c
}

func (c Conditions) Or(cond Cond) Conditions {
	if c.op != or {
		return Conditions{
			conditions: []Cond{c, cond},
			op:         or,
		}
	}

	conds, ok := cond.(Conditions)
	if ok && c.op == conds.op { // 类型相同时合并元素
		c.conditions = append(c.conditions, conds.conditions...)
	} else {
		c.conditions = append(c.conditions, cond)
	}
	return c
}

func (c Conditions) Not() Conditions {
	return Conditions{
		conditions: []Cond{c},
		op:         not,
	}
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
		if c.op == not {
			buf.WriteString("NOT ")
		}
		for i := 0; i < length; i++ {
			c.conditions[i].WriteSQL(buf, aliasMode)
			if i != lastIndex {
				if c.op == or {
					buf.WriteString(" OR ")
				} else if c.op == and {
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
		op:         and,
		conditions: conditions,
	}
}

func Or(conditions ...Cond) Conditions {
	return Conditions{
		op:         or,
		conditions: conditions,
	}
}

func Not(cond Cond) Conditions {
	return Conditions{
		conditions: []Cond{cond},
		op:         not,
	}
}
