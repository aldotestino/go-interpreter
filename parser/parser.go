package parser

import (
	"go-interpreter/lexer"
	"go-interpreter/utils"
	"slices"
)

// expr      : KEYWORD:var IDENTIFIER EQ expr
//           : comp ((KEYWORD:and|KEYWORD:or) comp)*

// comp      : KEYWORD:not comp
//           : mathExpr ((EE|NE|LT|LTE|GT|GTE) mathExpr)*

// math-expr : term ((PLUS|MINUS) term)*

// term      : factor ((MUL|DIV|MOD) factor)*

// factor    : (PLUS|MINUS) factor
//           : powerExpr

// power-expr: call (POW factor)*

// call      : atom (OpenParen (expr (COMMA expr)*)? RightParen)?

// atom      : INT|FLOAT|STRING|IDENTIFIER
//	         : OpenParen expr CloseParen
//           : list
//           : if-expr
//           : for-expr
//           : while-expr
//           : func-def

// list-expr : OpenBracket (expr, (COMMA expr)*)? CloseBracket

// if-expr   : KEYOWRD:if expr KEYWORD:then expr
//           : (KEYWORD:elif expr KEYWORD:then expr)*
//           : (KEYWORD: else expr)?

// for-expr  : KEYWORD:for IDENTIFIER EQ expr KEYWORD:to expr
//           : (KEYWORD:step expr)? KEYWORD:then expr

// while-expr: KEYWORD:while expr KEYWORD:then expr

// func-def  : KEYWORD:fun IDENTIFIER?
//           : OpenParen (IDENTIFIER (COMMA IDENTIFIER)*)? CloseParen
//           : ARROW expr

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

func (pars *Parser) forExpr() (AstNode, error) {
	if !pars.currentToken.Matches(lexer.KeywordTT, "for") {
		return nil, utils.InvalidSyntaxError("Expected 'for'")
	}

	pars.advance()

	if pars.currentToken.Type != lexer.IdentifierTT {
		return nil, utils.InvalidSyntaxError("Expected identifier")
	}

	varName := pars.currentToken
	pars.advance()

	if pars.currentToken.Type != lexer.EqualsTT {
		return nil, utils.InvalidSyntaxError("Expected '='")
	}

	pars.advance()

	startValue, err := pars.expr()
	if err != nil {
		return nil, err
	}

	if !pars.currentToken.Matches(lexer.KeywordTT, "to") {
		return nil, utils.InvalidSyntaxError("Expected 'to'")
	}

	pars.advance()

	endValue, err := pars.expr()
	if err != nil {
		return nil, err
	}

	var stepValue AstNode = nil

	if pars.currentToken.Matches(lexer.KeywordTT, "step") {
		pars.advance()
		stepValue, err = pars.expr()

		if err != nil {
			return nil, err
		}
	}

	if !pars.currentToken.Matches(lexer.KeywordTT, "then") {
		return nil, utils.InvalidSyntaxError("Expected 'then'")
	}

	pars.advance()

	body, err := pars.expr()

	if err != nil {
		return nil, err
	}

	return NewForNode(varName, startValue, endValue, stepValue, body), nil
}

func (pars *Parser) whileExpr() (AstNode, error) {

	if !pars.currentToken.Matches(lexer.KeywordTT, "while") {
		return nil, utils.InvalidSyntaxError("Expected 'while'")
	}

	pars.advance()

	condition, err := pars.expr()
	if err != nil {
		return nil, err
	}

	if !pars.currentToken.Matches(lexer.KeywordTT, "then") {
		return nil, utils.InvalidSyntaxError("Expected 'then'")
	}

	pars.advance()

	body, err := pars.expr()

	if err != nil {
		return nil, err
	}

	return NewWhileNode(condition, body), nil
}

func (pars *Parser) funcDef() (AstNode, error) {
	if !pars.currentToken.Matches(lexer.KeywordTT, "fun") {
		return nil, utils.InvalidSyntaxError("Expected 'fun'")
	}

	pars.advance()

	var varName *lexer.Token = nil
	if pars.currentToken.Type == lexer.IdentifierTT {
		varName = pars.currentToken

		pars.advance()

		if pars.currentToken.Type != lexer.OpenParenTT {
			return nil, utils.InvalidSyntaxError("Expected '('")
		}
	} else {
		if pars.currentToken.Type != lexer.OpenParenTT {
			return nil, utils.InvalidSyntaxError("Expected '('")
		}
	}

	pars.advance()

	args := make([]*lexer.Token, 0)

	if pars.currentToken.Type == lexer.IdentifierTT {
		args = append(args, pars.currentToken)
		pars.advance()

		for pars.currentToken.Type == lexer.CommaTT {
			pars.advance()

			if pars.currentToken.Type != lexer.IdentifierTT {
				return nil, utils.InvalidSyntaxError("Expected identifier")
			}

			args = append(args, pars.currentToken)
			pars.advance()
		}

		if pars.currentToken.Type != lexer.CloseParenTT {
			return nil, utils.InvalidSyntaxError("Expected ')'")
		}
	} else {
		if pars.currentToken.Type != lexer.CloseParenTT {
			return nil, utils.InvalidSyntaxError("Expected identifier or ')'")
		}
	}

	pars.advance()

	if pars.currentToken.Type != lexer.ArrowTT {
		return nil, utils.InvalidSyntaxError("Expected '->'")
	}

	pars.advance()

	node, err := pars.expr()

	if err != nil {
		return nil, err
	}

	return NewFuncDefNode(varName, args, node), nil
}

func (pars *Parser) call() (AstNode, error) {
	atom, err := pars.atom()
	if err != nil {
		return nil, err
	}

	if pars.currentToken.Type == lexer.OpenParenTT {
		pars.advance()
		args := make([]AstNode, 0)

		if pars.currentToken.Type == lexer.CloseParenTT {
			pars.advance()
		} else {
			newArg, err := pars.expr()
			if err != nil {
				return nil, err
			}
			args = append(args, newArg)

			for pars.currentToken.Type == lexer.CommaTT {
				pars.advance()

				newArg, err = pars.expr()
				if err != nil {
					return nil, err
				}
				args = append(args, newArg)
			}

			if pars.currentToken.Type != lexer.CloseParenTT {
				return nil, utils.InvalidSyntaxError("Expected ')'")
			}

			pars.advance()
		}

		return NewCallNode(atom, args), nil
	}

	return atom, nil
}

func (pars *Parser) ifExpr() (AstNode, error) {
	var cases = make([][]AstNode, 0)
	var elseCase AstNode = nil

	if !pars.currentToken.Matches(lexer.KeywordTT, "if") {
		return nil, utils.InvalidSyntaxError("Expected 'if'")
	}

	pars.advance()

	condition, err := pars.expr()
	if err != nil {
		return nil, err
	}

	if !pars.currentToken.Matches(lexer.KeywordTT, "then") {
		return nil, utils.InvalidSyntaxError("Expected 'then'")
	}

	pars.advance()

	expr, err := pars.expr()
	if err != nil {
		return nil, err
	}
	cases = append(cases, []AstNode{condition, expr})

	for pars.currentToken.Matches(lexer.KeywordTT, "elif") {
		pars.advance()

		condition, err = pars.expr()
		if err != nil {
			return nil, err
		}

		if !pars.currentToken.Matches(lexer.KeywordTT, "then") {
			return nil, utils.InvalidSyntaxError("Expected 'then'")
		}

		pars.advance()

		expr, err = pars.expr()
		if err != nil {
			return nil, err
		}

		cases = append(cases, []AstNode{condition, expr})
	}

	if pars.currentToken.Matches(lexer.KeywordTT, "else") {
		pars.advance()

		elseCase, err = pars.expr()
		if err != nil {
			return nil, err
		}
	}

	return NewIfNode(cases, elseCase), nil
}

func (pars *Parser) listExpr() (AstNode, error) {
	elements := make([]AstNode, 0)

	if pars.currentToken.Type != lexer.OpenBracketTT {
		return nil, utils.InvalidSyntaxError("Expected ']'")
	}

	pars.advance()

	if pars.currentToken.Type == lexer.CloseBracketTT {
		pars.advance()
	} else {
		newEl, err := pars.expr()
		if err != nil {
			return nil, err
		}
		elements = append(elements, newEl)

		for pars.currentToken.Type == lexer.CommaTT {
			pars.advance()

			newEl, err = pars.expr()
			if err != nil {
				return nil, err
			}
			elements = append(elements, newEl)
		}

		if pars.currentToken.Type != lexer.CloseBracketTT {
			return nil, utils.InvalidSyntaxError("Expected ']'")
		}

		pars.advance()
	}

	return NewListNode(elements), nil
}

func (pars *Parser) atom() (AstNode, error) {
	token := pars.currentToken

	if token.Type == lexer.IntTT || token.Type == lexer.FloatTT {
		pars.advance()
		return NewNumberNode(token), nil
	} else if token.Type == lexer.StringTT {
		pars.advance()
		return NewStringNode(token), nil
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
			return nil, utils.InvalidSyntaxError("Expected ')'")
		}
	} else if token.Type == lexer.OpenBracketTT {
		return pars.listExpr()
	} else if token.Matches(lexer.KeywordTT, "if") {
		return pars.ifExpr()
	} else if token.Matches(lexer.KeywordTT, "for") {
		return pars.forExpr()
	} else if token.Matches(lexer.KeywordTT, "while") {
		return pars.whileExpr()
	} else if token.Matches(lexer.KeywordTT, "fun") {
		return pars.funcDef()
	}

	var errMsg string

	if pars.currentPosition > 1 { // we advanced, so we don't expect the 'var' keyword
		errMsg = "Expected int, float, identifier, '+', '-', '(', '[', 'if', 'for', 'while' or 'fun'"
	} else {
		errMsg = "Expected int, float, identifier, 'var', '+', '-', '(', '[', '!', 'if', 'for', 'while' or 'fun'"
	}

	return nil, utils.InvalidSyntaxError(errMsg)
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

	return pars.powerExpr()
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

func (pars *Parser) mathExpr() (AstNode, error) {
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

	return pars.binOp(pars.mathExpr, pars.mathExpr, func(t *lexer.Token) bool {
		return slices.Contains([]lexer.TokenType{lexer.DoubleEqualsTT, lexer.NotEqualsTT, lexer.LessThanTT, lexer.LessThanEqualsTT, lexer.GreaterThanTT, lexer.GreaterThanEqualsTT}, t.Type)
	})
}

func (pars *Parser) powerExpr() (AstNode, error) {
	return pars.binOp(pars.call, pars.factor, func(t *lexer.Token) bool {
		return t.Type == lexer.PowerTT
	})
}

func (pars *Parser) term() (AstNode, error) {
	return pars.binOp(pars.factor, pars.factor, func(t *lexer.Token) bool {
		return slices.Contains([]lexer.TokenType{lexer.MultiplyTT, lexer.DivideTT, lexer.ModTT}, t.Type)
	})
}

func (pars *Parser) expr() (AstNode, error) {
	if pars.currentToken.Matches(lexer.KeywordTT, "var") {
		pars.advance()

		if pars.currentToken.Type != lexer.IdentifierTT {
			return nil, utils.InvalidSyntaxError("Expected identifier")
		}

		varName := pars.currentToken
		pars.advance()

		if pars.currentToken.Type != lexer.EqualsTT {
			return nil, utils.InvalidSyntaxError("Expected '='")
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
		return nil, utils.InvalidSyntaxError("Expected '+', '-', '*', '/', '^'")
	}

	return res, nil
}
