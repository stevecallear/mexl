package types

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	Type string

	Object interface {
		Type() Type
		Equal(o Object) bool
		Inspect() string
	}

	Integer struct {
		Value int64
	}

	Float struct {
		Value float64
	}

	String struct {
		Value string
	}

	Boolean struct {
		Value bool
	}

	Null struct{}

	Array []Object

	Map map[string]Object

	Func func(args ...Object) (Object, error)
)

const (
	TypeNull    Type = "NULL"
	TypeInteger Type = "INTEGER"
	TypeFloat   Type = "FLOAT"
	TypeString  Type = "STRING"
	TypeBoolean Type = "BOOLEAN"
	TypeArray   Type = "ARRAY"
	TypeMap     Type = "MAP"
	TypeFunc    Type = "FUNC"
)

var (
	_ Object = (*Integer)(nil)
	_ Object = (*Float)(nil)
	_ Object = (*String)(nil)
	_ Object = (*Boolean)(nil)
	_ Object = (Array)(nil)
	_ Object = (*Null)(nil)
	_ Object = (Map)(nil)
)

func (i *Integer) Equal(o Object) bool {
	if i == o {
		return true
	}
	switch t := o.(type) {
	case *Integer:
		return i.Value == t.Value
	case *Float:
		return i.Value == int64(t.Value)
	default:
		return false
	}
}

func (i *Integer) Type() Type {
	return TypeInteger
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (f *Float) Equal(o Object) bool {
	if f == o {
		return true
	}
	switch t := o.(type) {
	case *Float:
		return f.Value == t.Value
	case *Integer:
		return f.Value == float64(t.Value)
	default:
		return false
	}
}

func (f *Float) Type() Type {
	return TypeFloat
}

func (f *Float) Inspect() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

func (s *String) Type() Type {
	return TypeString
}

func (s *String) Equal(o Object) bool {
	if s == o {
		return true
	}
	if t, ok := o.(*String); ok {
		return s.Value == t.Value
	}
	return false
}

func (s *String) Inspect() string {
	return "\"" + s.Value + "\""
}

func (b *Boolean) Equal(o Object) bool {
	if b == o {
		return true
	}
	if t, ok := o.(*Boolean); ok {
		return b.Value == t.Value
	}
	return false
}

func (b *Boolean) Type() Type {
	return TypeBoolean
}

func (b *Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

func (n *Null) Equal(o Object) bool {
	if n == o {
		return true
	}
	_, ok := o.(*Null)
	return ok
}

func (n *Null) Type() Type {
	return TypeNull
}

func (n *Null) Inspect() string {
	return "null"
}

func (b Array) Type() Type {
	return TypeArray
}

func (a Array) Equal(o Object) bool {
	oa, ok := o.(Array)
	if !ok {
		return false
	}

	if len(a) != len(oa) {
		return false
	}

	for i, e := range a {
		if !e.Equal(oa[i]) {
			return false
		}
	}
	return true
}

func (a Array) Inspect() string {
	es := make([]string, len(a))
	for i, e := range a {
		es[i] = e.Inspect()
	}

	return "[" + strings.Join(es, ", ") + "]"
}

func (m Map) Type() Type {
	return TypeMap
}

func (m Map) Equal(o Object) bool {
	om, ok := o.(Map)
	if !ok {
		return false
	}

	if len(m) != len(om) {
		return false
	}

	for k, v := range m {
		if !v.Equal(om[k]) {
			return false
		}
	}
	return true
}

func (m Map) Inspect() string {
	es := make([]string, 0, len(m))
	for k, v := range m {
		es = append(es, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}

	return "{" + strings.Join(es, ", ") + "}"
}

func (f Func) Type() Type {
	return TypeFunc
}

func (f Func) Equal(o Object) bool {
	// func types are not comparable
	return false
}

func (f Func) Inspect() string {
	return "func"
}
