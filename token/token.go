package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(tt TokenType, ch byte) Token {
	return Token{
		Type:    tt,
		Literal: string(ch),
	}
}
