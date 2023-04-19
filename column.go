package sb

import (
	"bytes"
	"reflect"
)

type Column struct {
	name  string
	alias string
	table *Table
}

var columnType = reflect.TypeOf(Column{})

func (c *Column) As(alias string) *Column {
	c.alias = alias
	return c
}

func (c Column) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	switch aliasMode {
	case OnlyAlias:
		if c.alias != "" {
			buf.WriteByte('`')
			buf.WriteString(c.alias)
			buf.WriteByte('`')
			return
		}
		// else 当成 UseAlias 处理，并不会走 c.alias != "" 的流程
		fallthrough
	case UseAlias:
		if c.table != nil { // 正常情况都 != nil，除非手动构建
			alias := c.table.getAlias()
			if alias == "" {
				alias = c.table.getName()
			}
			buf.WriteByte('`')
			buf.WriteString(alias)
			buf.WriteString("`.`")
			buf.WriteString(c.name)

			if c.alias != "" {
				buf.WriteString("` AS `")
				buf.WriteString(c.alias)
			}
			buf.WriteByte('`')
			return
		}

		if c.alias != "" {
			buf.WriteByte('`')
			buf.WriteString(c.alias)
			buf.WriteByte('`')
			return
		}
		// else 当成 NoAlias 处理
		fallthrough
	case NoAlias:
		if c.name == "" { // 正常情况不会遇到，除非手动构建
			return
		}
		buf.WriteByte('`')
		buf.WriteString(c.name)
		buf.WriteByte('`')
	case ColonPrefix:
		if c.name == "" { // 正常情况不会遇到，除非手动构建
			return
		}
		buf.WriteByte(':')
		buf.WriteString(c.name) // TODO: 是否要转义？
	}
}

func (c *Column) Eq(e Expression) Condition {
	return Condition{op: "=", lv: c, rv: e}
}

func (c *Column) Ne(e Expression) Condition {
	return Condition{op: "!=", lv: c, rv: e}
}

func (c *Column) Gt(e Expression) Condition {
	return Condition{op: ">", lv: c, rv: e}
}

func (c *Column) Ge(e Expression) Condition {
	return Condition{op: ">=", lv: c, rv: e}
}

func (c *Column) Lt(e Expression) Condition {
	return Condition{op: "<", lv: c, rv: e}
}

func (c *Column) Le(e Expression) Condition {
	return Condition{op: "<=", lv: c, rv: e}
}

func (c *Column) In(e Expression) Condition {
	return Condition{op: "IN", lv: c, rv: e}
}

func (c *Column) NotIn(e Expression) Condition {
	return Condition{op: "NOT IN", lv: c, rv: e}
}

func (c *Column) Asc() OrderBy {
	return OrderBy{column: c, desc: false}
}

func (c *Column) Desc() OrderBy {
	return OrderBy{column: c, desc: true}
}

type Columns []Column

func (c Columns) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	length := len(c)
	if length > 0 {
		lastIndex := length - 1
		for i := 0; i < length; i++ {
			c[i].WriteSQL(buf, aliasMode)
			if i != lastIndex {
				buf.WriteString(", ")
			}
		}
	}
}
