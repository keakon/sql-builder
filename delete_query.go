package sb

import (
	"bytes"
	"strconv"
)

type DeleteQuery struct {
	table    AnyTable
	where    Cond
	orderBys OrderBys
	limit    uint64
}

func Delete(table AnyTable) *DeleteQuery {
	return &DeleteQuery{table: table}
}

func (q *DeleteQuery) Where(cond Cond) *DeleteQuery {
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

func (q *DeleteQuery) OrderBy(orderBy ...OrderBy) *DeleteQuery {
	q.orderBys = append(q.orderBys, orderBy...)
	return q
}

func (q *DeleteQuery) Limit(limit uint64) *DeleteQuery {
	q.limit = limit
	return q
}

func (q *DeleteQuery) WriteSQL(buf *bytes.Buffer) {
	buf.WriteString("DELETE `")
	buf.WriteString(q.table.getName())
	buf.WriteByte('`')
	if q.where != nil {
		q.where.WriteSQL(buf, NoAlias)
	}
	q.orderBys.WriteSQL(buf, NoAlias)
	if q.limit > 0 {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.FormatUint(q.limit, 10))
	}
}

func (q *DeleteQuery) String() string {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	q.WriteSQL(buf)

	sql := buf.String()
	pool.Put(buf)
	return sql
}
