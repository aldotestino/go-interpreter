package runtime

import (
	"go-interpreter/lexer"
	"go-interpreter/parser"
	"go-interpreter/utils"
	"strconv"
)

type Interpreter struct{}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (intr *Interpreter) visitNumberNode(node *parser.NumberNode) (RuntimeValue, error) {
	v, err := strconv.ParseFloat(node.Token.Value, 64)
	if err != nil {
		return nil, utils.RuntimeError("invalid number")
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
		return lhs.Add(rhs)
	} else if node.Operation.Type == lexer.MinusTT {
		return lhs.Subtract(rhs)
	} else if node.Operation.Type == lexer.MultiplyTT {
		return lhs.Multiply(rhs)
	} else if node.Operation.Type == lexer.DivideTT {
		return lhs.Divide(rhs)
	} else if node.Operation.Type == lexer.PowerTT {
		return lhs.Power(rhs)
	} else if node.Operation.Type == lexer.DoubleEqualsTT {
		return lhs.Equals(rhs)
	} else if node.Operation.Type == lexer.NotEqualsTT {
		return lhs.NotEquals(rhs)
	} else if node.Operation.Type == lexer.LessThanTT {
		return lhs.LessThan(rhs)
	} else if node.Operation.Type == lexer.LessThanEqualsTT {
		return lhs.LessThanEquals(rhs)
	} else if node.Operation.Type == lexer.GreaterThanTT {
		return lhs.GreaterThan(rhs)
	} else if node.Operation.Type == lexer.GreaterThanEqualsTT {
		return lhs.GreaterThanEquals(rhs)
	} else if node.Operation.Matches(lexer.KeywordTT, "and") {
		return lhs.And(rhs)
	} else if node.Operation.Matches(lexer.KeywordTT, "or") {
		return lhs.Or(rhs)
	}
	return nil, utils.RuntimeError("Unsupported operation")
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

func (intr *Interpreter) visitFuncDefNode(node *parser.FuncDefNode, env *Environment) (RuntimeValue, error) {
	funcName := "<anonymous>"

	if node.VarName != nil {
		funcName = node.VarName.Value
	}

	argNames := make([]string, 0)
	for _, arg := range node.Args {
		argNames = append(argNames, arg.Value)
	}

	funcValue := NewFunctionValue(funcName, node.Body, argNames)

	if node.VarName != nil {
		env.Set(funcName, funcValue)
	}

	return funcValue, nil
}

func (intr *Interpreter) visitCallNode(node *parser.CallNode, env *Environment) (RuntimeValue, error) {
	args := make([]RuntimeValue, 0)

	funcToCall, err := intr.Visit(node.Node, env)

	if err != nil {
		return nil, err
	}

	for _, arg := range node.Args {
		evalArg, err := intr.Visit(arg, env)
		if err != nil {
			return nil, err
		}
		args = append(args, evalArg)
	}

	return funcToCall.Execute(env, args)
}

func (intr *Interpreter) visitStringNode(node *parser.StringNode, env *Environment) (RuntimeValue, error) {
	return NewStringValue(node.Token.Value), nil
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
	case parser.FuncDefNT:
		return intr.visitFuncDefNode(node.(*parser.FuncDefNode), env)
	case parser.CallNT:
		return intr.visitCallNode(node.(*parser.CallNode), env)
	case parser.StringNT:
		return intr.visitStringNode(node.(*parser.StringNode), env)
	default:
		return nil, utils.RuntimeError("Unsupported node")
	}
}
