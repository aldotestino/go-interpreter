package parser

import "go-interpreter/lexer"

type NodeType string

const (
	NumberNT NodeType = "number"
	UnOpNT   NodeType = "unop"
	BinOpNt  NodeType = "binop"
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
