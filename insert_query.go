package sb

import (
	"bytes"
)

type InsertQuery struct {
	table       AnyTable
	columns     Columns
	values      Expressions
	assignments Assignments
	aliasMode   AliasMode // of values
	ignore      bool
}

func Insert(table AnyTable) *InsertQuery {
	return &InsertQuery{table: table}
}

func (q *InsertQuery) Ignore() *InsertQuery {
	q.ignore = true
	return q
}

func (q *InsertQuery) Columns(columns ...Column) *InsertQuery {
	q.columns = columns
	return q
}

func (q *InsertQuery) Values(values ...Expression) *InsertQuery {
	q.values = values
	return q
}

func (q *InsertQuery) NamedValues(values ...Expression) *InsertQuery {
	q.values = values
	q.aliasMode = ColonPrefix
	return q
}

func (q *InsertQuery) OnDuplicateKeyUpdate(assignments ...Assignment) *InsertQuery {
	q.assignments = assignments
	return q
}

func (q *InsertQuery) WriteSQL(buf *bytes.Buffer) {
	if q.ignore {
		buf.WriteString("INSERT IGNORE INTO `")
	} else {
		buf.WriteString("INSERT INTO `")
	}
	buf.WriteString(q.table.getName())
	buf.WriteString("` (")
	q.columns.WriteSQL(buf, NoAlias)
	buf.WriteString(") VALUES (")
	if q.values == nil {
		count := len(q.columns)
		if count > 0 {
			if q.aliasMode == ColonPrefix { // 使用 NamedValues() 绑定时，如果没有参数，就用 columns
				q.columns.WriteSQL(buf, q.aliasMode)
			} else { // 填充 '?'
				for i := 0; i < count; i++ {
					buf.WriteByte('?')
					if i != count-1 {
						buf.WriteString(", ")
					}
				}
			}
		}
	} else if len(q.values) > 0 {
		q.values.WriteSQL(buf, q.aliasMode)
	}
	buf.WriteByte(')')
	if len(q.assignments) > 0 {
		buf.WriteString(" ON DUPLICATE KEY UPDATE ")
		q.assignments.WriteSQL(buf, NoAlias)
	}
}

func (q *InsertQuery) String() string {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	q.WriteSQL(buf)

	sql := buf.String()
	pool.Put(buf)
	return sql
}
