package lexer

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
	Type  TokenType
	Value string
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}
