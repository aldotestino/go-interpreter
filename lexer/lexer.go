package lexer

import (
	"fmt"
	"go-interpreter/utils"
	"regexp"
	"slices"
	"strings"
)

type Lexer struct {
	fn              string
	text            []string
	currentPosition int
	currentChar     string
}

func NewLexer(input string) *Lexer {
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

func (lex *Lexer) isAlpha(char string) bool {
	return regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString(char)
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

func (lex *Lexer) makeIdentifier() *Token {
	idString := ""

	for lex.currentChar != "" && (lex.isAlpha(lex.currentChar) || lex.isDigit(lex.currentChar)) {
		idString += lex.currentChar
		lex.advance()
	}

	if slices.Contains(KEYWORDS, idString) {
		return NewToken(KeywordTT, idString)
	}

	return NewToken(IdentifierTT, idString)
}

func (lex *Lexer) makeNotEquals() (*Token, error) {
	lex.advance()

	if lex.currentChar == "=" {
		lex.advance()
		return NewToken(NotEqualsTT, "!="), nil
	}

	return nil, utils.ExpectedCharError("'=' (after '!')")
}

func (lex *Lexer) makeEquals() *Token {
	lex.advance()

	tt := EqualsTT
	val := "="

	if lex.currentChar == "=" {
		lex.advance()
		tt = DoubleEqualsTT
		val = "=="
	}

	return NewToken(tt, val)
}

func (lex *Lexer) makeLessThan() *Token {
	lex.advance()

	tt := LessThanTT
	val := "<"

	if lex.currentChar == "=" {
		lex.advance()
		tt = LessThanEqualsTT
		val = "<="
	}

	return NewToken(tt, val)
}

func (lex *Lexer) makeGreaterThan() *Token {
	lex.advance()

	tt := GreaterThanTT
	val := ">"

	if lex.currentChar == "=" {
		lex.advance()
		tt = GreaterThanEqualsTT
		val = ">="
	}

	return NewToken(tt, val)
}

func (lex *Lexer) makeMinusOrArrow() *Token {
	lex.advance()
	tt := MinusTT
	val := "-"

	if lex.currentChar == ">" {
		lex.advance()
		tt = ArrowTT
		val = "->"
	}

	return NewToken(tt, val)

}

func (lex *Lexer) Tokenize() ([]*Token, error) {
	tokens := make([]*Token, 0)

	for lex.currentChar != "" {
		if lex.isSkippable(lex.currentChar) {
			lex.advance()
		} else if lex.isDigit(lex.currentChar) {
			tokens = append(tokens, lex.makeNumber())
		} else if lex.isAlpha(lex.currentChar) {
			tokens = append(tokens, lex.makeIdentifier())
		} else if lex.currentChar == "+" {
			tokens = append(tokens, NewToken(PlusTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "-" {
			tokens = append(tokens, lex.makeMinusOrArrow())
		} else if lex.currentChar == "*" {
			tokens = append(tokens, NewToken(MultiplyTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "/" {
			tokens = append(tokens, NewToken(DivideTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "^" {
			tokens = append(tokens, NewToken(PowerTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "(" {
			tokens = append(tokens, NewToken(OpenParenTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == ")" {
			tokens = append(tokens, NewToken(CloseParenTT, lex.currentChar))
			lex.advance()
		} else if lex.currentChar == "!" {
			neToken, err := lex.makeNotEquals()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, neToken)
		} else if lex.currentChar == "=" { // creates '=' or '=='
			tokens = append(tokens, lex.makeEquals())
		} else if lex.currentChar == "<" { // creates '<' or '<='
			tokens = append(tokens, lex.makeLessThan())
		} else if lex.currentChar == ">" { // creates '>' or '>='
			tokens = append(tokens, lex.makeGreaterThan())
		} else if lex.currentChar == "," {
			tokens = append(tokens, NewToken(CommaTT, lex.currentChar))
			lex.advance()
		} else {
			cc := lex.currentChar
			lex.advance()
			return nil, utils.IllegalCharError(fmt.Sprintf("'%s'", cc))
		}

	}

	tokens = append(tokens, NewToken(EOFTT, ""))
	return tokens, nil
}
