package runtime

import (
	"fmt"
	"go-interpreter/parser"
	"go-interpreter/utils"
	"math"
	"reflect"
	"slices"
)

type ValueType string

const (
	NumberVT ValueType = "Number"
	FuncVT   ValueType = "Function"
	StringVT ValueType = "String"
	ListVT   ValueType = "List"
)

type RuntimeValue interface {
	GetType() ValueType
	GetValue() any
	Print() string

	Add(other RuntimeValue) (RuntimeValue, error)
	Subtract(other RuntimeValue) (RuntimeValue, error)
	Multiply(other RuntimeValue) (RuntimeValue, error)
	Divide(other RuntimeValue) (RuntimeValue, error)
	Mod(other RuntimeValue) (RuntimeValue, error)
	Power(other RuntimeValue) (RuntimeValue, error)
	Equals(other RuntimeValue) (RuntimeValue, error)
	NotEquals(other RuntimeValue) (RuntimeValue, error)
	LessThan(other RuntimeValue) (RuntimeValue, error)
	GreaterThan(other RuntimeValue) (RuntimeValue, error)
	LessThanEquals(other RuntimeValue) (RuntimeValue, error)
	GreaterThanEquals(other RuntimeValue) (RuntimeValue, error)
	And(other RuntimeValue) (RuntimeValue, error)
	Or(other RuntimeValue) (RuntimeValue, error)

	Execute(parentEnv *Environment, args []RuntimeValue) (RuntimeValue, error)
}

// NumberValue

type NumberValue struct {
	Type  ValueType
	Value float64
}

func NewNumberValue(num float64) *NumberValue {
	return &NumberValue{
		Type:  NumberVT,
		Value: num,
	}
}

func (nv *NumberValue) GetType() ValueType {
	return nv.Type
}

func (nv *NumberValue) GetValue() any {
	return nv.Value
}

func (nv *NumberValue) Print() string {
	return fmt.Sprintf("%v", nv.Value)
}

func (nv *NumberValue) Add(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '+'")
	}

	return NewNumberValue(nv.Value + other.GetValue().(float64)), nil
}

func (nv *NumberValue) Subtract(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '-'")
	}

	return NewNumberValue(nv.Value - other.GetValue().(float64)), nil
}

func (nv *NumberValue) Multiply(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '*'")
	}

	return NewNumberValue(nv.Value * other.GetValue().(float64)), nil
}

func (nv *NumberValue) Divide(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '/'")
	}

	if other.GetValue().(float64) == 0.0 {
		return nil, utils.RuntimeError("Division by 0")
	}

	return NewNumberValue(nv.Value / other.GetValue().(float64)), nil
}

func (nv *NumberValue) Mod(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '%'")
	}

	return NewNumberValue(math.Mod(nv.Value, other.GetValue().(float64))), nil
}

func (nv *NumberValue) Power(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '*'")
	}

	return NewNumberValue(math.Pow(nv.Value, other.GetValue().(float64))), nil
}

func (nv *NumberValue) Equals(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '=='")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value == other.GetValue().(float64))), nil
}

func (nv *NumberValue) NotEquals(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '!='")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value != other.GetValue().(float64))), nil
}

func (nv *NumberValue) LessThan(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '<'")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value < other.GetValue().(float64))), nil
}

func (nv *NumberValue) GreaterThan(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '>'")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value > other.GetValue().(float64))), nil
}

func (nv *NumberValue) LessThanEquals(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '<='")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value <= other.GetValue().(float64))), nil
}

func (nv *NumberValue) GreaterThanEquals(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '>='")
	}

	return NewNumberValue(utils.BoolToNumber(nv.Value >= other.GetValue().(float64))), nil
}

func (nv *NumberValue) And(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation 'and'")
	}

	return NewNumberValue(utils.AndNumbers(nv.Value, other.GetValue().(float64))), nil
}

func (nv *NumberValue) Or(other RuntimeValue) (RuntimeValue, error) {
	if nv.Type != NumberVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation 'or'")
	}

	return NewNumberValue(utils.OrNumbers(nv.Value, other.GetValue().(float64))), nil
}

func (nv *NumberValue) Execute(parentEnv *Environment, args []RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '()'")
}

// FuncValue

type FunctionValue struct {
	Type     ValueType
	Name     string
	Body     parser.AstNode
	ArgNames []string
}

func NewFunctionValue(n string, b parser.AstNode, a []string) *FunctionValue {
	return &FunctionValue{
		Type:     FuncVT,
		Name:     n,
		Body:     b,
		ArgNames: a,
	}
}

func (f *FunctionValue) GetType() ValueType {
	return f.Type
}

func (f *FunctionValue) GetValue() any {
	return f.Body
}

func (f *FunctionValue) Print() string {
	return fmt.Sprintf("<function %s>", f.Name)
}

func (f *FunctionValue) Add(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '+'")
}

func (f *FunctionValue) Subtract(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation ''")
}

func (f *FunctionValue) Multiply(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '*'")
}

func (f *FunctionValue) Divide(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '/'")
}

func (f *FunctionValue) Mod(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '%'")
}

func (f *FunctionValue) Power(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '^'")
}

func (f *FunctionValue) Equals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '='")
}

func (f *FunctionValue) NotEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '!='")
}

func (f *FunctionValue) LessThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<'")
}

func (f *FunctionValue) GreaterThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>'")
}

func (f *FunctionValue) LessThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<='")
}

func (f *FunctionValue) GreaterThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>='")
}

func (f *FunctionValue) And(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'and'")
}

func (f *FunctionValue) Or(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'or'")
}

func (f *FunctionValue) Execute(parentEnv *Environment, args []RuntimeValue) (RuntimeValue, error) {
	intr := NewInterpreter()
	env := NewEnvironment(parentEnv)

	argsDiff := len(args) - len(f.ArgNames)

	if argsDiff > 0 {
		return nil, utils.RuntimeError(fmt.Sprintf("%d too many args passed into '%s'", argsDiff, f.Name))
	} else if argsDiff < 0 {
		return nil, utils.RuntimeError(fmt.Sprintf("%d too few args passed into '%s'", argsDiff*-1, f.Name))
	}

	for i, arg := range args {
		argName := f.ArgNames[i]
		argValue := arg
		env.Set(argName, argValue)
	}

	return intr.Visit(f.Body, env)
}

// StringValue

type StringValue struct {
	Type  ValueType
	Value string
}

func NewStringValue(v string) *StringValue {
	return &StringValue{
		Type:  StringVT,
		Value: v,
	}
}

func (s *StringValue) GetType() ValueType {
	return s.Type
}

func (s *StringValue) GetValue() any {
	return s.Value
}

func (s *StringValue) Print() string {
	return s.Value
}

func (s *StringValue) Add(other RuntimeValue) (RuntimeValue, error) {
	if s.Type != StringVT || other.GetType() != StringVT {
		return nil, utils.RuntimeError("Illegal operation '+'")
	}

	return NewStringValue(s.Value + other.GetValue().(string)), nil
}

func (s *StringValue) Subtract(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '-'")
}

func (s *StringValue) Multiply(other RuntimeValue) (RuntimeValue, error) {
	if s.Type != StringVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '*'")
	}

	finalStr := ""
	times := int(other.GetValue().(float64))

	for i := 0; i < times; i++ {
		finalStr += s.Value
	}

	return NewStringValue(finalStr), nil
}

func (s *StringValue) Divide(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '/'")
}

func (s *StringValue) Mod(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '%'")
}

func (s *StringValue) Power(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '^'")
}

func (s *StringValue) Equals(other RuntimeValue) (RuntimeValue, error) {
	if s.Type != StringVT || other.GetType() != StringVT {
		return nil, utils.RuntimeError("Illegal operation '=='")
	}

	return NewNumberValue(utils.BoolToNumber(s.Value == other.GetValue().(string))), nil
}

func (s *StringValue) NotEquals(other RuntimeValue) (RuntimeValue, error) {
	if s.Type != StringVT || other.GetType() != StringVT {
		return nil, utils.RuntimeError("Illegal operation '!='")
	}

	return NewNumberValue(utils.BoolToNumber(s.Value != other.GetValue().(string))), nil
}

func (s *StringValue) LessThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<'")
}

func (s *StringValue) GreaterThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>'")
}

func (s *StringValue) LessThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<='")
}

func (s *StringValue) GreaterThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>='")
}

func (s *StringValue) And(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'and'")
}

func (s *StringValue) Or(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'or'")
}

func (s *StringValue) Execute(parentEnv *Environment, args []RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '()'")
}

// ListValue

type ListValue struct {
	Type     ValueType
	Elements []RuntimeValue
}

func NewListValue(els []RuntimeValue) *ListValue {
	return &ListValue{
		Type:     ListVT,
		Elements: els,
	}
}

func (l *ListValue) GetType() ValueType {
	return ListVT
}

func (l *ListValue) GetValue() any {
	return l.Elements
}

func (l *ListValue) Print() string {
	str := "["

	for i, el := range l.Elements {
		if el != nil {
			if i > 0 {
				str += ", "
			}
			str += fmt.Sprintf("%v", el.Print())
		}
	}

	str += "]"
	return str
}

// append a single element to a list
func (l *ListValue) Add(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '+'")
	}

	return NewListValue(append(l.Elements, other)), nil
}

// remove element at index other.Value from list
func (l *ListValue) Subtract(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '-'")
	}

	index := other.GetValue().(float64)

	if !utils.FloatIsInt(index) {
		return nil, utils.RuntimeError("Index must be an integer")
	}

	absIndex := int(math.Abs(index))
	length := len(l.Elements)

	if absIndex > length {
		return nil, utils.RuntimeError("Index out of bounds")
	}

	elementsCopy := make([]RuntimeValue, len(l.Elements))
	copy(elementsCopy, l.Elements)
	var newElements []RuntimeValue

	if index >= 0 {
		newElements = append(elementsCopy[:absIndex], elementsCopy[absIndex+1:]...)
	} else {
		newElements = append(elementsCopy[:length-absIndex], elementsCopy[length-absIndex+1:]...)
	}

	return NewListValue(newElements), nil
}

// concat list
func (l *ListValue) Multiply(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != ListVT {
		return nil, utils.RuntimeError("Illegal operation '*'")
	}

	return NewListValue(slices.Concat(l.Elements, other.(*ListValue).Elements)), nil
}

// get element at index other.Value
func (l *ListValue) Divide(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != NumberVT {
		return nil, utils.RuntimeError("Illegal operation '/'")
	}

	index := other.GetValue().(float64)

	if index != float64(int(index)) {
		return nil, utils.RuntimeError("Index must be an integer")
	}

	absIndex := int(math.Abs(index))
	length := len(l.Elements)

	if absIndex > length {
		return nil, utils.RuntimeError("Index out of bounds")
	}

	if index >= 0 {
		return l.Elements[absIndex], nil
	} else {
		return l.Elements[length-absIndex], nil
	}
}

func (l *ListValue) Mod(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '%'")
}

func (l *ListValue) Power(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '^'")
}

func (l *ListValue) Equals(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != ListVT {
		return nil, utils.RuntimeError("Illegal operation '=='")
	}

	return NewNumberValue(utils.BoolToNumber(reflect.DeepEqual(l.Elements, other.(*ListValue).Elements))), nil
}

func (l *ListValue) NotEquals(other RuntimeValue) (RuntimeValue, error) {
	if l.Type != ListVT || other.GetType() != ListVT {
		return nil, utils.RuntimeError("Illegal operation '!='")
	}

	return NewNumberValue(utils.BoolToNumber(!reflect.DeepEqual(l.Elements, other.(*ListValue).Elements))), nil
}

func (l *ListValue) LessThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<'")
}

func (l *ListValue) GreaterThan(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>'")
}

func (l *ListValue) LessThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '<='")
}

func (l *ListValue) GreaterThanEquals(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '>='")
}

func (l *ListValue) And(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'and'")
}

func (l *ListValue) Or(other RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation 'or'")
}

func (l *ListValue) Execute(parentEnv *Environment, args []RuntimeValue) (RuntimeValue, error) {
	return nil, utils.RuntimeError("Illegal operation '()'")
}
