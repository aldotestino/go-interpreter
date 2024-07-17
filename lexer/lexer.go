package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func BaseError(posStart *Position, posEnd *Position, errorName string, details string) error {
	errorMsg := fmt.Sprintf("%s: %s\nFile %s, line %d", errorName, details, posStart.Fn, posStart.Ln+1)
	return errors.New(errorMsg)
}

func IllegalCharError(posStart *Position, posEnd *Position, details string) error {
	return BaseError(posStart, posEnd, "Illegal Character", details)
}

type TokenType string

const (
	IntTT        TokenType = "Int"
	FloatTT      TokenType = "Float"
	PlusTT       TokenType = "Plus"
	MinusTT      TokenType = "Minus"
	MultiplyTT   TokenType = "Multiply"
	DivideTT     TokenType = "Divide"
	OpenParenTT  TokenType = "OpenParen"
	CloseParenTT TokenType = "CloseParen"
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}

type Position struct {
	Idx  int    // index of character
	Ln   int    // line
	Col  int    // column
	Fn   string // filename
	Ftxt string // text
}

func NewPosition(i, l, c int, fn, ftxt string) *Position {
	return &Position{
		Idx:  i,
		Ln:   l,
		Col:  c,
		Fn:   fn,
		Ftxt: ftxt,
	}
}

func (p *Position) advance(currentChar string) *Position {
	p.Idx++
	p.Col++

	if currentChar == "\n" {
		p.Ln++
		p.Col = 0
	}

	return p
}

func (p *Position) copy() *Position {
	return NewPosition(p.Idx, p.Ln, p.Col, p.Fn, p.Ftxt)
}

type Lexer struct {
	fn              string
	text            []string
	currentPosition *Position
	currentChar     string
}

func NewLexer(fn, input string) *Lexer {
	lex := &Lexer{
		fn:              fn,
		text:            strings.Split(input, ""),
		currentPosition: NewPosition(-1, 0, -1, fn, input),
		currentChar:     "",
	}
	lex.advance()

	return lex
}

func (lex *Lexer) advance() {
	lex.currentPosition.advance(lex.currentChar)
	if lex.currentPosition.Idx < len(lex.text) {
		lex.currentChar = lex.text[lex.currentPosition.Idx]
	} else {
		lex.currentChar = ""
	}
}

func (lex *Lexer) isSkippable(char string) bool {
	return char == " " || char == "\t" || char == "\n" || char == "\r"
}

func (lex *Lexer) isDigit(char string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(char)
}

func (lex *Lexer) makeNumber() *Token {
	numString := ""
	dotCount := 0

	for lex.currentChar != "" && (lex.isDigit(lex.currentChar) || lex.currentChar == ".") {
		if lex.currentChar == "." {
			if dotCount == 1 {
				break
			}
			dotCount++
			numString += "."
		} else {
			numString += lex.currentChar
		}
		lex.advance()
	}

	if dotCount == 0 {
		return NewToken(IntTT, numString)
	} else {
		return NewToken(FloatTT, numString)
	}
}

func (lex *Lexer) Tokenize() ([]*Token, error) {
	tokens := make([]*Token, 0)

	for lex.currentChar != "" {
		if lex.isSkippable(lex.currentChar) {
			lex.advance()
		} else if lex.isDigit(lex.currentChar) {
			tokens = append(tokens, lex.makeNumber())
		} else if lex.currentChar == "+" {
			tokens = append(tokens, NewToken(PlusTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "-" {
			tokens = append(tokens, NewToken(MinusTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "*" {
			tokens = append(tokens, NewToken(MultiplyTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "/" {
			tokens = append(tokens, NewToken(DivideTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "(" {
			tokens = append(tokens, NewToken(OpenParenTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "/" {
			tokens = append(tokens, NewToken(CloseParenTT, lex.currentChar))
			lex.advance()
		} else {
			posStart := lex.currentPosition.copy()
			cc := lex.currentChar
			lex.advance()
			return nil, IllegalCharError(posStart, lex.currentPosition, fmt.Sprintf("'%s'", cc))
		}

	}

	return tokens, nil
}
