package lexer

type TokenType string

const (
	IntTT        TokenType = "Int"
	FloatTT      TokenType = "Float"
	IdentifierTT TokenType = "Identifier"
	KeywordTT    TokenType = "Keyword"
	PlusTT       TokenType = "Plus"
	MinusTT      TokenType = "Minus"
	MultiplyTT   TokenType = "Multiply"
	DivideTT     TokenType = "Divide"
	PowerTT      TokenType = "Power"
	EqualsTT     TokenType = "Equals"
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

func (t *Token) Matches(tt TokenType, v string) bool {
	return t.Type == tt && t.Value == v
}

var KEYWORDS = []string{"var"}
