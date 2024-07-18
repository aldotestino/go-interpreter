package lexer

import (
	"fmt"
	"go-interpreter/shared"
	"regexp"
	"strings"
)

type Lexer struct {
	fn              string
	text            []string
	currentPosition int
	currentChar     string
}

func NewLexer(fn, input string) *Lexer {
	lex := &Lexer{
		text:            strings.Split(input, ""),
		currentPosition: -1,
		currentChar:     "",
	}
	lex.advance()

	return lex
}

func (lex *Lexer) advance() {
	lex.currentPosition++

	if lex.currentPosition < len(lex.text) {
		lex.currentChar = lex.text[lex.currentPosition]
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
		} else if lex.currentChar == "^" {
			tokens = append(tokens, NewToken(PowerTT, lex.currentChar))
			lex.advance()
		} else {
			cc := lex.currentChar
			lex.advance()
			return nil, shared.IllegalCharError(fmt.Sprintf("'%s'", cc))
		}

	}

	tokens = append(tokens, NewToken(EOFTT, ""))
	return tokens, nil
}
