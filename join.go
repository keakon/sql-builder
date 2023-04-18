package sb

import "bytes"

type JoinType uint8

const (
	InnerJoin JoinType = iota
	LeftJoin
	RightJoin
	OuterJoin
)

type FromTables struct {
	table AnyTable
	joins []Join
}

func (f *FromTables) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	if f.table == nil {
		return
	}

	buf.WriteString(" FROM `")
	buf.WriteString(f.table.getName())
	buf.WriteByte('`')
	if aliasMode != NoAlias {
		alias := f.table.getAlias()
		if alias != "" {
			buf.WriteString(" AS `")
			buf.WriteString(alias)
			buf.WriteByte('`')
		}
	}
	for _, join := range f.joins {
		join.WriteSQL(buf, aliasMode)
	}
}

type Join struct {
	typ   JoinType
	table AnyTable
	on    Condition // 必须是 table1.col1 = table2.col
}

func (j *Join) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	if j.table == nil { // 正常情况不会遇到，除非手动构建
		return
	}

	switch j.typ {
	case InnerJoin:
		buf.WriteString(" JOIN `")
	case LeftJoin:
		buf.WriteString(" LEFT JOIN `")
	case RightJoin:
		buf.WriteString(" RIGHT JOIN `")
	case OuterJoin:
		buf.WriteString(" OUTER JOIN `")
	default:
		return
	}

	buf.WriteString(j.table.getName())
	buf.WriteByte('`')
	if aliasMode != NoAlias {
		alias := j.table.getAlias()
		if alias != "" {
			buf.WriteString(" AS `")
			buf.WriteString(alias)
			buf.WriteByte('`')
		}
	}

	buf.WriteString(" ON ")
	j.on.WriteSQL(buf, aliasMode)
}

func (t *Table) InnerJoin(table AnyTable, on Condition) *FromTables {
	return &FromTables{table: t, joins: []Join{{typ: InnerJoin, table: table, on: on}}}
}

func (t *FromTables) InnerJoin(table AnyTable, on Condition) *FromTables {
	t.joins = append(t.joins, Join{typ: InnerJoin, table: table, on: on})
	return t
}

func (t *Table) LeftJoin(table AnyTable, on Condition) *FromTables {
	return &FromTables{table: t, joins: []Join{{typ: LeftJoin, table: table, on: on}}}
}

func (t *FromTables) LeftJoin(table AnyTable, on Condition) *FromTables {
	t.joins = append(t.joins, Join{typ: LeftJoin, table: table, on: on})
	return t
}

func (t *Table) RightJoin(table AnyTable, on Condition) *FromTables {
	return &FromTables{table: t, joins: []Join{{typ: RightJoin, table: table, on: on}}}
}

func (t *FromTables) RightJoin(table AnyTable, on Condition) *FromTables {
	t.joins = append(t.joins, Join{typ: RightJoin, table: table, on: on})
	return t
}

func (t *Table) OuterJoin(table AnyTable, on Condition) *FromTables {
	return &FromTables{table: t, joins: []Join{{typ: OuterJoin, table: table, on: on}}}
}

func (t *FromTables) OuterJoin(table AnyTable, on Condition) *FromTables {
	t.joins = append(t.joins, Join{typ: OuterJoin, table: table, on: on})
	return t
}
