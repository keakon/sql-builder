package sb

import "bytes"

type Expression interface {
	WriteSQL(buf *bytes.Buffer, aliasMode AliasMode)
}

type Expressions []Expression // 输出时用 ", " 分隔每个元素

func (e Expressions) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	length := len(e)
	if length > 0 {
		lastIndex := length - 1
		for i := 0; i < length; i++ {
			if e[i] == nil {
				buf.WriteString("NULL")
			} else {
				e[i].WriteSQL(buf, aliasMode)
			}
			if i != lastIndex {
				buf.WriteString(", ")
			}
		}
	}
}

type Expr string

func (e Expr) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	buf.WriteString(string(e))
}

const Placeholder = Expr("?")

type ConcatExpressions Expressions // 直接输出每个元素

func Concat(expressions ...Expression) ConcatExpressions {
	return ConcatExpressions(expressions)
}

func (e ConcatExpressions) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	for i := 0; i < len(e); i++ {
		e[i].WriteSQL(buf, aliasMode)
	}
}

type Function struct {
	Name        string
	Expressions Expressions
	Alias       string
}

func Func(name string, exps ...Expression) *Function {
	return &Function{
		Name:        name,
		Expressions: exps,
	}
}

func (f *Function) WriteSQL(buf *bytes.Buffer, aliasMode AliasMode) {
	buf.WriteString(f.Name)
	buf.WriteByte('(')
	f.Expressions.WriteSQL(buf, aliasMode)
	buf.WriteByte(')')
	if f.Alias != "" {
		buf.WriteString(" AS `")
		buf.WriteString(f.Alias)
		buf.WriteByte('`')
	}
}

func (f *Function) As(alias string) *Function {
	f.Alias = alias
	return f
}
