package lexer

import (
	"fmt"
	"go-interpreter/shared"
	"regexp"
	"strings"
)

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
	EOFTT        TokenType = "EOF"
)

type Token struct {
	Type     TokenType
	Value    string
	PosStart *shared.Position
	PosEnd   *shared.Position
}

func NewToken(tokenType TokenType, value string, posStart, posEnd *shared.Position) *Token {
	tok := &Token{
		Type:     tokenType,
		Value:    value,
		PosStart: nil,
		PosEnd:   nil,
	}

	if posStart != nil {
		tok.PosStart = posStart.Copy()
		tok.PosEnd = posStart.Copy()
		tok.PosEnd.Advance("")
	}

	if posEnd != nil {
		tok.PosEnd = posEnd.Copy()
	}

	return tok
}

type Lexer struct {
	fn              string
	text            []string
	currentPosition *shared.Position
	currentChar     string
}

func NewLexer(fn, input string) *Lexer {
	lex := &Lexer{
		fn:              fn,
		text:            strings.Split(input, ""),
		currentPosition: shared.NewPosition(-1, 0, -1, fn, input),
		currentChar:     "",
	}
	lex.advance()

	return lex
}

func (lex *Lexer) advance() {
	lex.currentPosition.Advance(lex.currentChar)
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

	posStart := lex.currentPosition.Copy()

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
		return NewToken(IntTT, numString, posStart, lex.currentPosition)
	} else {
		return NewToken(FloatTT, numString, posStart, lex.currentPosition)
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
			tokens = append(tokens, NewToken(PlusTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else if lex.currentChar == "-" {
			tokens = append(tokens, NewToken(MinusTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else if lex.currentChar == "*" {
			tokens = append(tokens, NewToken(MultiplyTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else if lex.currentChar == "/" {
			tokens = append(tokens, NewToken(DivideTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else if lex.currentChar == "(" {
			tokens = append(tokens, NewToken(OpenParenTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else if lex.currentChar == "/" {
			tokens = append(tokens, NewToken(CloseParenTT, lex.currentChar, lex.currentPosition, nil))
			lex.advance()
		} else {
			posStart := lex.currentPosition.Copy()
			cc := lex.currentChar
			lex.advance()
			return nil, shared.IllegalCharError(posStart, lex.currentPosition, fmt.Sprintf("'%s'", cc))
		}

	}

	tokens = append(tokens, NewToken(EOFTT, "", lex.currentPosition, nil))
	return tokens, nil
}
