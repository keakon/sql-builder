package sb

import (
	"bytes"
	"strconv"
)

type SelectLockMode uint8

const (
	SelectNoLock SelectLockMode = iota
	SelectForRead
	SelectForWrite
)

type Query interface {
}

type SelectQuery struct {
	buf         bytes.Buffer
	expressions Expressions
	from        *FromTables
	where       Cond
	groupBys    Columns
	orderBys    OrderBys
	limit       uint
	offset      uint
	lockMode    SelectLockMode
}

func Select(expressions ...Expression) *SelectQuery {
	if len(expressions) == 0 {
		return nil
	}
	return &SelectQuery{expressions: expressions}
}

func (q *SelectQuery) From(from *FromTables) *SelectQuery {
	q.from = from
	return q
}

func (q *SelectQuery) FromTable(from AnyTable) *SelectQuery {
	q.from = &FromTables{table: from}
	return q
}

func (q *SelectQuery) Where(cond Cond) *SelectQuery {
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

func (q *SelectQuery) GroupBy(columns ...Column) *SelectQuery {
	q.groupBys = append(q.groupBys, columns...)
	return q
}

func (q *SelectQuery) OrderBy(orderBy ...OrderBy) *SelectQuery {
	q.orderBys = append(q.orderBys, orderBy...)
	return q
}

func (q *SelectQuery) Limit(limit uint) *SelectQuery {
	q.limit = limit
	return q
}

func (q *SelectQuery) Offset(offset uint) *SelectQuery {
	q.offset = offset
	return q
}

func (q *SelectQuery) LockForRead() *SelectQuery {
	q.lockMode = SelectForRead
	return q
}

func (q *SelectQuery) LockForWrite() *SelectQuery {
	q.lockMode = SelectForWrite
	return q
}

func (q *SelectQuery) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	buf.WriteString("SELECT ")
	q.expressions.WriteSQL(buf, aliasMode)
	q.from.WriteSQL(buf, aliasMode)
	if q.where != nil {
		q.where.WriteSQL(buf, aliasMode)
	}
	if len(q.groupBys) > 0 {
		buf.WriteString(" GROUP BY ")
		q.groupBys.WriteSQL(buf, aliasMode)
	}
	q.orderBys.WriteSQL(buf, aliasMode)
	if q.limit > 0 || q.offset > 0 {
		buf.WriteString(" LIMIT ")
		if q.offset > 0 { // LIMIT offset, limit
			buf.WriteString(strconv.FormatUint(uint64(q.offset), 10))
			buf.WriteString(", ")
			buf.WriteString(strconv.FormatUint(uint64(q.limit), 10))
		} else { // LIMIT limit
			buf.WriteString(strconv.FormatUint(uint64(q.limit), 10))
		}
	}
	switch q.lockMode {
	case SelectForRead:
		buf.WriteString(" FOR SHARE")
	case SelectForWrite:
		buf.WriteString(" FOR UPDATE")
	}
}

func (q *SelectQuery) ToSQL() string {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()

	if len(q.from.joins) > 0 {
		q.WriteSQL(buf, UseAlias)
	} else {
		q.WriteSQL(buf, NoAlias)
	}

	sql := buf.String()
	pool.Put(buf)
	return sql
}
