package parser

import (
	"go-interpreter/lexer"
	"go-interpreter/shared"
)

// expr: term ((PLUS|MINUS) term)*
// term : factor ((MUL|DIV) factor)*
// factor: INT|FLOAT
//       : (PLUS|MINUS) factor
//       : OpenParen expr CloseParen

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

func (pars *Parser) factor() *ParseResult {
	res := NewParseResult()
	token := pars.currentToken

	if token.Type == lexer.PlusTT || token.Type == lexer.MinusTT {
		// res.register(pars.advance()) // for now it does nothing
		pars.advance()
		factor := res.register(pars.factor())

		if res.Error != nil {
			return res
		}

		return res.success(NewUnOpNode(factor, token))
	} else if token.Type == lexer.IntTT || token.Type == lexer.FloatTT {
		//res.register(pars.advance()) // for now it does nothing
		pars.advance()
		return res.success(NewNumberNode(token))
	} else if token.Type == lexer.OpenParenTT {
		// res.register(pars.advanace())
		pars.advance()

		expr := res.register(pars.expr())

		if res.Error != nil {
			return res
		}

		if pars.currentToken.Type == lexer.CloseParenTT {
			// res.register(pars.advance())
			pars.advance()

			return res.success(expr)
		} else {
			return res.failure(shared.InvalidSyntaxError(pars.currentToken.PosStart, pars.currentToken.PosEnd, "Expected ')'"))
		}
	}

	return res.failure(shared.InvalidSyntaxError(token.PosStart, token.PosEnd, "Expected int or float"))
}

func (pars *Parser) term() *ParseResult {
	res := NewParseResult()
	left := res.register(pars.factor())

	if res.Error != nil {
		return res
	}

	for pars.currentToken.Type == lexer.MultiplyTT || pars.currentToken.Type == lexer.DivideTT {
		opToken := pars.currentToken
		//res.register(pars.advance()) // for now it does nothing
		pars.advance()
		right := res.register(pars.factor())

		if res.Error != nil {
			return res
		}

		left = NewBinOpNode(left, right, opToken)
	}

	return res.success(left)
}

func (pars *Parser) expr() *ParseResult {
	res := NewParseResult()
	left := res.register(pars.term())

	if res.Error != nil {
		return res
	}

	for pars.currentToken.Type == lexer.PlusTT || pars.currentToken.Type == lexer.MinusTT {
		opToken := pars.currentToken
		//res.register(pars.advance()) // for now it does nothing
		pars.advance()
		right := res.register(pars.term())

		if res.Error != nil {
			return res
		}

		left = NewBinOpNode(left, right, opToken)
	}

	return res.success(left)
}

func (pars *Parser) Parse() *ParseResult {
	res := pars.expr()

	if res.Error == nil && pars.currentToken.Type != lexer.EOFTT {
		return res.failure(shared.InvalidSyntaxError(pars.currentToken.PosStart, pars.currentToken.PosEnd, "Expected '+', '-', '*', '/'"))
	}

	return res
}
