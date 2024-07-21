package lexer

type TokenType string

const (
	IntTT               TokenType = "Int"
	FloatTT             TokenType = "Float"
	IdentifierTT        TokenType = "Identifier"
	KeywordTT           TokenType = "Keyword"
	PlusTT              TokenType = "Plus"
	MinusTT             TokenType = "Minus"
	MultiplyTT          TokenType = "Multiply"
	DivideTT            TokenType = "Divide"
	ModTT               TokenType = "Mod"
	PowerTT             TokenType = "Power"
	EqualsTT            TokenType = "Equals"
	OpenParenTT         TokenType = "OpenParen"
	CloseParenTT        TokenType = "CloseParen"
	OpenBracketTT       TokenType = "OpenBracket"
	CloseBracketTT      TokenType = "CloseBracket"
	DoubleEqualsTT      TokenType = "DoubleEquals"
	NotEqualsTT         TokenType = "NotEquals"
	LessThanTT          TokenType = "LessThan"
	GreaterThanTT       TokenType = "GreaterThan"
	LessThanEqualsTT    TokenType = "LessThanEquals"
	GreaterThanEqualsTT TokenType = "GreaterThanEquals"
	CommaTT             TokenType = "Comma"
	ArrowTT             TokenType = "Arrow"
	StringTT            TokenType = "String"
	EOFTT               TokenType = "EOF"
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

var KEYWORDS = []string{"var", "and", "or", "not", "if", "then", "elif", "else", "for", "to", "step", "while", "fun"}
