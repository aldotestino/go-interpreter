package parser

import (
	"go-interpreter/lexer"
	"go-interpreter/shared"
	"slices"
)

// expr  :  KEYWORD:var IDENTIFIER EQ expr
//       :  comp ((KEYWORD:and|KEYWORD:or) comp)*
// comp  :  KEYWORD:not comp
//       :  math ((EE|NE|LT|LTE|GT|GTE) math)*
// math  :  term ((PLUS|MINUS) term)*
// term  :  factor ((MUL|DIV) factor)*
// factor:  (PLUS|MINUS) factor
//       :  power
// power :  atom (POW factor)*
// atom  :  INT|FLOAT|IDENTIFIER
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
	} else if token.Type == lexer.IdentifierTT {
		pars.advance()
		return NewVarAccessNode(token), nil
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

	var errMsg string

	if pars.currentPosition > 1 { // we advanced, so we don't expect the 'var' keyword
		errMsg = "Expected int, float, identifier, '+', '-' or ')'"
	} else {
		errMsg = "Expected int, float, identifier, 'var', '+', '-', ')' or '!'"
	}

	return nil, shared.InvalidSyntaxError(errMsg)
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

func (pars *Parser) binOp(lf, rf func() (AstNode, error), cond func(t *lexer.Token) bool) (AstNode, error) {
	left, err := lf()

	if err != nil {
		return nil, err
	}

	for cond(pars.currentToken) {
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

func (pars *Parser) math() (AstNode, error) {
	return pars.binOp(pars.term, pars.term, func(t *lexer.Token) bool {
		return slices.Contains([]lexer.TokenType{lexer.PlusTT, lexer.MinusTT}, t.Type)
	})
}

func (pars *Parser) comp() (AstNode, error) {
	if pars.currentToken.Matches(lexer.KeywordTT, "not") {
		opToken := pars.currentToken
		pars.advance()

		node, err := pars.expr()

		if err != nil {
			return nil, err
		}

		return NewUnOpNode(node, opToken), nil
	}

	return pars.binOp(pars.math, pars.math, func(t *lexer.Token) bool {
		return slices.Contains([]lexer.TokenType{lexer.DoubleEqualsTT, lexer.NotEqualsTT, lexer.LessThanTT, lexer.LessThanEqualsTT, lexer.GreaterThanTT, lexer.GreaterThanEqualsTT}, t.Type)
	})
}

func (pars *Parser) power() (AstNode, error) {
	return pars.binOp(pars.atom, pars.factor, func(t *lexer.Token) bool {
		return t.Type == lexer.PowerTT
	})
}

func (pars *Parser) term() (AstNode, error) {
	return pars.binOp(pars.factor, pars.factor, func(t *lexer.Token) bool {
		return slices.Contains([]lexer.TokenType{lexer.PlusTT, lexer.MinusTT}, t.Type)
	})
}

func (pars *Parser) expr() (AstNode, error) {
	if pars.currentToken.Matches(lexer.KeywordTT, "var") {
		pars.advance()

		if pars.currentToken.Type != lexer.IdentifierTT {
			return nil, shared.InvalidSyntaxError("Expected identifier")
		}

		varName := pars.currentToken
		pars.advance()

		if pars.currentToken.Type != lexer.EqualsTT {
			return nil, shared.InvalidSyntaxError("Expected '='")
		}

		pars.advance()
		expr, err := pars.expr()

		if err != nil {
			return nil, err
		}

		return NewVarAssignNode(varName, expr), nil
	}

	return pars.binOp(pars.comp, pars.comp, func(t *lexer.Token) bool {
		return t.Matches(lexer.KeywordTT, "and") || t.Matches(lexer.KeywordTT, "or")
	})
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
