package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/eugene-whitaker/monkey/ast"
)

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
	STRING_OBJECT       = "STRING"
	BUILTIN_OBJECT      = "BUILTIN"
	ARRAY_OBJECT        = "ARRAY"
	HASH_OBJECT         = "HASH"
	QUOTE_OBJECT        = "QUOTE"
	MACRO_OBJECT        = "MACRO"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) HashKey() HashKey {
	var v uint64

	if b.Value {
		v = 1
	} else {
		v = 0
	}

	return HashKey{
		Type:  b.Type(),
		Value: v,
	}
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

func (n *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	ps := []string{}
	for _, p := range f.Parameters {
		ps = append(ps, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJECT
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJECT
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	es := []string{}
	for _, e := range a.Elements {
		es = append(es, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJECT
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	ps := []string{}
	for _, p := range h.Pairs {
		ps = append(ps, p.Key.Inspect()+":"+p.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString("}")

	return out.String()
}

type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType {
	return QUOTE_OBJECT
}

func (q *Quote) Inspect() string {
	return "quote(" + q.Node.String() + ")"
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (m *Macro) Inspect() string {
	var out bytes.Buffer

	ps := []string{}
	for _, p := range m.Parameters {
		ps = append(ps, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}
