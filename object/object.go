package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cijin/go-interpreter/ast"
)

const (
	INTEGER_OBJ      = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOLEAN_OBJ      = "BOOLEAN"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	NULL_OBJ         = "NULL"
	FUNCTION_OBJ     = "FUNCTION"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

// int
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

// string
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// bool
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b) }

// return
type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Message }

// null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Enviornment
type Enviornment struct {
	store map[string]Object
	outer *Enviornment
}

func NewEnviornment() *Enviornment {
	return &Enviornment{store: make(map[string]Object), outer: nil}
}

func NewEnclosedEnviornment(outer *Enviornment) *Enviornment {
	return &Enviornment{store: make(map[string]Object), outer: outer}
}

func (e *Enviornment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.store[name]
	}

	return val, ok
}

func (e *Enviornment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}

// Function
type Function struct {
	Args []*ast.Identifier
	Body *ast.BlockStatement
	Env  *Enviornment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var buf bytes.Buffer
	var params []string

	for _, arg := range f.Args {
		params = append(params, arg.String())
	}

	buf.WriteString("fn")
	buf.WriteString("(")
	buf.WriteString(strings.Join(params, ", "))
	buf.WriteString(") {\n")
	buf.WriteString(f.Body.String())
	buf.WriteString("}")

	return buf.String()
}
