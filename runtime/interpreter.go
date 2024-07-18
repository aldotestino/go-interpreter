package runtime

import (
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/shared"
	"math"
	"strconv"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (intr *Interpreter) visitNumberNode(node *parser.NumberNode) (RuntimeValue, error) {
	v, err := strconv.ParseFloat(node.Token.Value, 64)
	if err != nil {
		return nil, shared.RuntimeError("invalid number")
	}
	return NewNumberValue(v), nil
}

func (intr *Interpreter) visitUnOpNode(node *parser.UnOpNode, env *Environment) (RuntimeValue, error) {
	num, err := intr.Visit(node.Node, env)

	if err != nil {
		return nil, err
	}

	if node.Operator.Type == lexer.MinusTT {
		return NewNumberValue(num.(*NumberValue).Value * -1), nil
	}

	return num, nil
}

func (intr *Interpreter) visitBinOpNode(node *parser.BinOpNode, env *Environment) (RuntimeValue, error) {
	lhs, err := intr.Visit(node.Left, env)

	if err != nil {
		return nil, err
	}

	rhs, err := intr.Visit(node.Right, env)

	if err != nil {
		return nil, err
	}

	switch node.Operation.Type {
	case lexer.PlusTT:
		return NewNumberValue(lhs.(*NumberValue).Value + rhs.(*NumberValue).Value), nil
	case lexer.MinusTT:
		return NewNumberValue(lhs.(*NumberValue).Value - rhs.(*NumberValue).Value), nil
	case lexer.MultiplyTT:
		return NewNumberValue(lhs.(*NumberValue).Value * rhs.(*NumberValue).Value), nil
	case lexer.DivideTT:
		if rhs.(*NumberValue).Value == 0 {
			return nil, shared.RuntimeError("Division by 0")
		}
		return NewNumberValue(lhs.(*NumberValue).Value / rhs.(*NumberValue).Value), nil
	case lexer.PowerTT:
		return NewNumberValue(math.Pow(lhs.(*NumberValue).Value, rhs.(*NumberValue).Value)), nil
	default:
		return nil, shared.RuntimeError("Unsupported operation")
	}
}

func (intr *Interpreter) visitVarAccessNode(node *parser.VarAccessNode, env *Environment) (RuntimeValue, error) {
	varName := node.VarName.Value
	value, err := env.Get(varName)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (intr *Interpreter) visitVarAssignNode(node *parser.VarAssignNode, env *Environment) (RuntimeValue, error) {
	varName := node.VarName.Value
	value, err := intr.Visit(node.Value, env)

	if err != nil {
		return nil, err
	}

	env.Set(varName, value)
	return value, nil
}

func (intr *Interpreter) Visit(node parser.AstNode, env *Environment) (RuntimeValue, error) {
	switch node.GetType() {
	case parser.NumberNT:
		return intr.visitNumberNode(node.(*parser.NumberNode))
	case parser.UnOpNT:
		return intr.visitUnOpNode(node.(*parser.UnOpNode), env)
	case parser.BinOpNt:
		return intr.visitBinOpNode(node.(*parser.BinOpNode), env)
	case parser.VarAccessNT:
		return intr.visitVarAccessNode(node.(*parser.VarAccessNode), env)
	case parser.VarAssignNT:
		return intr.visitVarAssignNode(node.(*parser.VarAssignNode), env)
	default:
		return nil, shared.RuntimeError("Unsupported node")
	}
}
