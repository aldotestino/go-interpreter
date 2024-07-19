package runtime

import (
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/shared"
	"go-interpreter/utils"
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

func (intr *Interpreter) visitUnOpNode(node *parser.UnOpNode, env *Environment) (RuntimeValue, error) {
	num, err := intr.Visit(node.Node, env)

	if err != nil {
		return nil, err
	}

	if node.Operator.Type == lexer.MinusTT {
		return NewNumberValue(num.(*NumberValue).Value * -1), nil
	} else if node.Operator.Matches(lexer.KeywordTT, "not") {
		if num.(*NumberValue).Value == 0 {
			return NewNumberValue(1), nil
		} else {
			return NewNumberValue(0), nil
		}
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

	if node.Operation.Type == lexer.PlusTT {
		return NewNumberValue(lhs.(*NumberValue).Value + rhs.(*NumberValue).Value), nil
	} else if node.Operation.Type == lexer.MinusTT {
		return NewNumberValue(lhs.(*NumberValue).Value - rhs.(*NumberValue).Value), nil
	} else if node.Operation.Type == lexer.MultiplyTT {
		return NewNumberValue(lhs.(*NumberValue).Value * rhs.(*NumberValue).Value), nil
	} else if node.Operation.Type == lexer.DivideTT {
		if rhs.(*NumberValue).Value == 0 {
			return nil, shared.RuntimeError("Division by 0")
		}
		return NewNumberValue(lhs.(*NumberValue).Value / rhs.(*NumberValue).Value), nil
	} else if node.Operation.Type == lexer.PowerTT {
		return NewNumberValue(math.Pow(lhs.(*NumberValue).Value, rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.DoubleEqualsTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value == rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.NotEqualsTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value != rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.LessThanTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value < rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.LessThanEqualsTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value <= rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.GreaterThanTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value > rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Type == lexer.GreaterThanEqualsTT {
		return NewNumberValue(utils.BoolToNumber(lhs.(*NumberValue).Value >= rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Matches(lexer.KeywordTT, "and") {
		return NewNumberValue(utils.AndNumbers(lhs.(*NumberValue).Value, rhs.(*NumberValue).Value)), nil
	} else if node.Operation.Matches(lexer.KeywordTT, "or") {
		return NewNumberValue(utils.OrNumbers(lhs.(*NumberValue).Value, rhs.(*NumberValue).Value)), nil
	}
	return nil, shared.RuntimeError("Unsupported operation")
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

func (intr *Interpreter) visitIfNode(node *parser.IfNode, env *Environment) (RuntimeValue, error) {
	for _, c := range node.Cases {

		conditionValue, err := intr.Visit(c[0], env)
		if err != nil {
			return nil, err
		}

		if conditionValue.GetValue() == 1.0 {
			return intr.Visit(c[1], env)
		}
	}

	if node.ElseCase != nil {
		return intr.Visit(node.ElseCase, env)
	}

	return nil, nil
}

func (intr *Interpreter) visitForNode(node *parser.ForNode, env *Environment) (RuntimeValue, error) {
	startValue, err := intr.Visit(node.StartValue, env)
	if err != nil {
		return nil, err
	}
	endValue, err := intr.Visit(node.EndValue, env)
	if err != nil {
		return nil, err
	}

	stepValue := NewNumberValue(1)

	if node.StepValue != nil {
		sv, err := intr.Visit(node.StepValue, env)
		if err != nil {
			return nil, err
		}
		stepValue = sv.(*NumberValue)
	}

	i := startValue.(*NumberValue).Value

	var condition = func() bool {
		return i < endValue.(*NumberValue).Value
	}

	if stepValue.Value < 0 {
		condition = func() bool {
			return i > endValue.(*NumberValue).Value
		}
	}

	for condition() {
		env.Set(node.VarName.Value, NewNumberValue(i))
		i += stepValue.Value

		_, err := intr.Visit(node.Body, env)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (intr *Interpreter) visitWhileNode(node *parser.WhileNode, env *Environment) (RuntimeValue, error) {
	for {
		condition, err := intr.Visit(node.Condition, env)
		if err != nil {
			return nil, err
		}

		if condition.GetValue() != 1.0 {
			break
		}

		_, err = intr.Visit(node.Body, env)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
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
	case parser.IfNT:
		return intr.visitIfNode(node.(*parser.IfNode), env)
	case parser.ForNT:
		return intr.visitForNode(node.(*parser.ForNode), env)
	case parser.WhileNT:
		return intr.visitWhileNode(node.(*parser.WhileNode), env)
	default:
		return nil, shared.RuntimeError("Unsupported node")
	}
}
