package sb

import (
	"bytes"
	"reflect"
	"strings"
)

type Table struct {
	name  string
	alias string
}

func (t Table) isTable() {}

func (t Table) getName() string { return t.name }

func (t Table) getAlias() string { return t.alias }

func (t Table) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	if aliasMode != NoAlias {
		buf.WriteByte('`')
		if t.alias == "" {
			buf.WriteString(t.name)
		} else {
			buf.WriteString(t.alias)
		}
		buf.WriteString("`.*")
	} else {
		buf.WriteByte('*')
	}
}

func (t Table) Select(expressions ...Expression) *SelectQuery {
	if len(expressions) == 0 {
		return Select(t).From(t)
	}
	return Select(expressions...).From(t)
}

func (t Table) Update(assignments ...Assignment) *UpdateQuery {
	return Update(t).Set(assignments...)
}

func (t Table) Insert() *InsertQuery {
	return Insert(t)
}

func (t Table) Delete() *DeleteQuery {
	return Delete(t)
}

type AnyTable interface {
	getName() string
	getAlias() string
	isTable()
}

var tableType = reflect.TypeOf(Table{})

func New[T AnyTable](alias string) *T {
	var t T
	ptr := reflect.ValueOf(&t)
	rv := reflect.Indirect(ptr)
	fieldCount := rv.NumField()
	if fieldCount == 0 {
		return &t
	}

	f0 := rv.Field(0)
	if f0.Type() != tableType {
		return &t
	}

	rt := reflect.TypeOf(t)
	name := rt.Field(0).Tag.Get("db")
	f0.Set(reflect.ValueOf(Table{name: name, alias: alias}))

	if fieldCount > 1 {
		table := f0.Interface().(Table)

		for i := 1; i < fieldCount; i++ {
			f := rv.Field(i)
			if f.Type() == columnType {
				fi := rt.Field(i)
				if fi.IsExported() {
					name := fi.Tag.Get("db")
					if name == "" {
						name = strings.ToLower(fi.Name)
					}
					f.Set(reflect.ValueOf(Column{name: name, table: &table}))
				}
			}
		}
	}
	return &t
}
