package parser

import "go-interpreter/lexer"

type NodeType string

const (
	NumberNT    NodeType = "Number"
	UnOpNT      NodeType = "UnOp"
	BinOpNt     NodeType = "BinOp"
	VarAccessNT NodeType = "VarAccess"
	VarAssignNT NodeType = "VarAssign"
	IfNT        NodeType = "If"
	ForNT       NodeType = "For"
	WhileNT     NodeType = "While"
	FuncDefNT   NodeType = "FunDef"
	CallNT      NodeType = "Call"
	StringNT    NodeType = "String"
)

type AstNode interface {
	GetType() NodeType
}

// NumerNode

type NumberNode struct {
	Type  NodeType
	Token *lexer.Token
}

func NewNumberNode(token *lexer.Token) *NumberNode {
	return &NumberNode{
		Type:  NumberNT,
		Token: token,
	}
}

func (n *NumberNode) GetType() NodeType {
	return n.Type
}

// UnOpNode

type UnOpNode struct {
	Type     NodeType
	Node     AstNode
	Operator *lexer.Token
}

func NewUnOpNode(n AstNode, o *lexer.Token) *UnOpNode {
	return &UnOpNode{
		Type:     UnOpNT,
		Node:     n,
		Operator: o,
	}
}

func (u *UnOpNode) GetType() NodeType {
	return u.Type
}

// BinOpNode

type BinOpNode struct {
	Type      NodeType
	Left      AstNode
	Right     AstNode
	Operation *lexer.Token
}

func NewBinOpNode(l, r AstNode, o *lexer.Token) *BinOpNode {
	return &BinOpNode{
		Type:      BinOpNt,
		Left:      l,
		Right:     r,
		Operation: o,
	}
}

func (n *BinOpNode) GetType() NodeType {
	return n.Type
}

// VarAccessNode

type VarAccessNode struct {
	Type    NodeType
	VarName *lexer.Token
}

func NewVarAccessNode(vn *lexer.Token) *VarAccessNode {
	return &VarAccessNode{
		Type:    VarAccessNT,
		VarName: vn,
	}
}

func (n *VarAccessNode) GetType() NodeType {
	return n.Type
}

// VarAssignNode

type VarAssignNode struct {
	Type    NodeType
	VarName *lexer.Token
	Value   AstNode
}

func NewVarAssignNode(vn *lexer.Token, v AstNode) *VarAssignNode {
	return &VarAssignNode{
		Type:    VarAssignNT,
		VarName: vn,
		Value:   v,
	}
}

func (n *VarAssignNode) GetType() NodeType {
	return n.Type
}

// IfNode

type IfNode struct {
	Type     NodeType
	Cases    [][]AstNode
	ElseCase AstNode
}

func NewIfNode(c [][]AstNode, e AstNode) *IfNode {
	return &IfNode{
		Type:     IfNT,
		Cases:    c,
		ElseCase: e,
	}
}

func (n *IfNode) GetType() NodeType {
	return n.Type
}

// ForNode

type ForNode struct {
	Type       NodeType
	VarName    *lexer.Token
	StartValue AstNode
	EndValue   AstNode
	StepValue  AstNode
	Body       AstNode
}

func NewForNode(vn *lexer.Token, s, e, st, b AstNode) *ForNode {
	return &ForNode{
		Type:       ForNT,
		VarName:    vn,
		StartValue: s,
		EndValue:   e,
		StepValue:  st,
		Body:       b,
	}
}

func (n *ForNode) GetType() NodeType {
	return n.Type
}

// WhileNode

type WhileNode struct {
	Type      NodeType
	Condition AstNode
	Body      AstNode
}

func NewWhileNode(c, b AstNode) *WhileNode {
	return &WhileNode{
		Type:      WhileNT,
		Condition: c,
		Body:      b,
	}
}

func (n *WhileNode) GetType() NodeType {
	return n.Type
}

// FuncDefNode

type FuncDefNode struct {
	Type    NodeType
	VarName *lexer.Token
	Args    []*lexer.Token
	Body    AstNode
}

func NewFuncDefNode(v *lexer.Token, a []*lexer.Token, b AstNode) *FuncDefNode {
	return &FuncDefNode{
		Type:    FuncDefNT,
		VarName: v,
		Args:    a,
		Body:    b,
	}
}

func (n *FuncDefNode) GetType() NodeType {
	return n.Type
}

// CallNode

type CallNode struct {
	Type NodeType
	Node AstNode
	Args []AstNode
}

func NewCallNode(n AstNode, a []AstNode) *CallNode {
	return &CallNode{
		Type: CallNT,
		Node: n,
		Args: a,
	}
}

func (n *CallNode) GetType() NodeType {
	return n.Type
}

// StringNode

type StringNode struct {
	Type  NodeType
	Token *lexer.Token
}

func NewStringNode(token *lexer.Token) *StringNode {
	return &StringNode{
		Type:  StringNT,
		Token: token,
	}
}

func (n *StringNode) GetType() NodeType {
	return n.Type
}
