package sb

import (
	"bytes"
	"strconv"
)

type UpdateQuery struct {
	table       AnyTable
	assignments Assignments
	where       Cond
	orderBys    OrderBys
	limit       uint64
}

func Update(table AnyTable) *UpdateQuery {
	return &UpdateQuery{table: table}
}

func (q *UpdateQuery) Set(assignments ...Assignment) *UpdateQuery {
	q.assignments = assignments
	return q
}

func (q *UpdateQuery) Where(cond Cond) *UpdateQuery {
	switch cond := cond.(type) {
	case Condition:
		q.where = Conditions{
			conditions: []Cond{cond},
			isTopLevel: true,
		}
	case Conditions:
		cond.isTopLevel = true
		q.where = cond
	}
	return q
}

func (q *UpdateQuery) OrderBy(orderBy ...OrderBy) *UpdateQuery {
	q.orderBys = append(q.orderBys, orderBy...)
	return q
}

func (q *UpdateQuery) Limit(limit uint64) *UpdateQuery {
	q.limit = limit
	return q
}

func (q *UpdateQuery) WriteSQL(buf *bytes.Buffer) {
	buf.WriteString("UPDATE `")
	buf.WriteString(q.table.getName())
	buf.WriteString("` SET ")
	q.assignments.WriteSQL(buf, NoAlias)
	if q.where != nil {
		q.where.WriteSQL(buf, NoAlias)
	}
	q.orderBys.WriteSQL(buf, NoAlias)
	if q.limit > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.FormatUint(q.limit, 10))
	}
}

func (q *UpdateQuery) String() string {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	q.WriteSQL(buf)

	sql := buf.String()
	pool.Put(buf)
	return sql
}
