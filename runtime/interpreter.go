package runtime

import (
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/shared"
	"math"
	"strconv"
)

type Interpreter struct{}

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

func (intr *Interpreter) visitUnOpNode(node *parser.UnOpNode) (RuntimeValue, error) {
	num, err := intr.Visit(node.Node)

	if err != nil {
		return nil, err
	}

	if node.Operator.Type == lexer.MinusTT {
		return NewNumberValue(num.(*NumberValue).Value * -1), nil
	}

	return num, nil
}

func (intr *Interpreter) visitBinOpNode(node *parser.BinOpNode) (RuntimeValue, error) {
	lhs, err := intr.Visit(node.Left)

	if err != nil {
		return nil, err
	}

	rhs, err := intr.Visit(node.Right)

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

func (intr *Interpreter) Visit(node parser.AstNode) (RuntimeValue, error) {
	switch node.GetType() {
	case parser.NumberNT:
		return intr.visitNumberNode(node.(*parser.NumberNode))
	case parser.UnOpNT:
		return intr.visitUnOpNode(node.(*parser.UnOpNode))
	case parser.BinOpNt:
		return intr.visitBinOpNode(node.(*parser.BinOpNode))
	default:
		return nil, shared.RuntimeError("Unsupported node")
	}
}
