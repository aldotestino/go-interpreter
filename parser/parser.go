package parser

import (
	"go-interpreter/lexer"
	"go-interpreter/shared"
	"slices"
)

// expr  :  term ((PLUS|MINUS) term)*
// term  :  factor ((MUL|DIV) factor)*
// factor:  (PLUS|MINUS) factor
//       :  power
// power :  atom (POW factor)*
// atom  :  INT|FLOAT
//	     :  OpenParen expr CloseParen

type Parser struct {
	tokens          []*lexer.Token
	currentPosition int
	currentToken    *lexer.Token
}

func NewParser(tokens []*lexer.Token) *Parser {
	pars := &Parser{
		tokens:          tokens,
		currentPosition: -1,
		currentToken:    nil,
	}
	pars.advance()

	return pars
}

func (pars *Parser) advance() *lexer.Token {
	pars.currentPosition++

	if pars.currentPosition < len(pars.tokens) {
		pars.currentToken = pars.tokens[pars.currentPosition]
	}

	return pars.currentToken
}

func (pars *Parser) atom() (AstNode, error) {
	token := pars.currentToken

	if token.Type == lexer.IntTT || token.Type == lexer.FloatTT {
		pars.advance()
		return NewNumberNode(token), nil
	} else if token.Type == lexer.OpenParenTT {
		pars.advance()
		expr, err := pars.expr()

		if err != nil {
			return nil, err
		}

		if pars.currentToken.Type == lexer.CloseParenTT {
			pars.advance()
			return expr, nil
		} else {
			return nil, shared.InvalidSyntaxError("Expected ')'")
		}
	}

	return nil, shared.InvalidSyntaxError("Expected int, float, '+', '-' or ')'")
}

func (pars *Parser) factor() (AstNode, error) {
	token := pars.currentToken

	if token.Type == lexer.PlusTT || token.Type == lexer.MinusTT {
		pars.advance()
		factor, err := pars.factor()

		if err != nil {
			return nil, err
		}

		return NewUnOpNode(factor, token), nil
	}

	return pars.power()
}

func (pars *Parser) binOp(lf, rf func() (AstNode, error), ops []lexer.TokenType) (AstNode, error) {
	left, err := lf()

	if err != nil {
		return nil, err
	}

	for slices.Contains(ops, pars.currentToken.Type) {
		opToken := pars.currentToken
		pars.advance()
		right, err := rf()

		if err != nil {
			return nil, err
		}

		left = NewBinOpNode(left, right, opToken)
	}

	return left, nil
}

func (pars *Parser) power() (AstNode, error) {
	return pars.binOp(pars.atom, pars.factor, []lexer.TokenType{lexer.PowerTT})
}

func (pars *Parser) term() (AstNode, error) {
	return pars.binOp(pars.factor, pars.factor, []lexer.TokenType{lexer.MultiplyTT, lexer.DivideTT})
}

func (pars *Parser) expr() (AstNode, error) {
	return pars.binOp(pars.term, pars.term, []lexer.TokenType{lexer.PlusTT, lexer.MinusTT})
}

func (pars *Parser) Parse() (AstNode, error) {
	res, err := pars.expr()

	if err != nil {
		return nil, err
	}

	if pars.currentToken.Type != lexer.EOFTT {
		return nil, shared.InvalidSyntaxError("Expected '+', '-', '*', '/', '^")
	}

	return res, nil
}
