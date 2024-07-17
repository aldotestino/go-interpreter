package parser

import "go-interpreter/lexer"

type NodeType string

const (
	NumberNT      NodeType = "number"
	UnOpNT        NodeType = "unop"
	BinOpNt       NodeType = "binop"
	ParseResultNT NodeType = "parse_result"
)

type AstNode interface {
	GetType() NodeType
}

// ParseResult

type ParseResult struct {
	Type  NodeType
	Error error
	Node  AstNode
}

func NewParseResult() *ParseResult {
	return &ParseResult{
		Type:  ParseResultNT,
		Error: nil,
		Node:  nil,
	}
}

func (pr *ParseResult) GetType() NodeType {
	return pr.Type
}

func (pr *ParseResult) register(res AstNode) AstNode {
	if res.GetType() == ParseResultNT {
		if res.(*ParseResult).Error != nil {
			pr.Error = res.(*ParseResult).Error
		}
		return res.(*ParseResult).Node
	}
	return res
}

func (pr *ParseResult) success(node AstNode) *ParseResult {
	pr.Node = node
	return pr
}

func (pr *ParseResult) failure(err error) *ParseResult {
	pr.Error = err
	return pr
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
