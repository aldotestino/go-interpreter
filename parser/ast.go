package parser

import "go-interpreter/lexer"

type NodeType string

const (
	NumberNT    NodeType = "Number"
	UnOpNT      NodeType = "UnOp"
	BinOpNt     NodeType = "BinOp"
	VarAccessNT NodeType = "VarAccess"
	VarAssignNT NodeType = "VarAssign"
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
