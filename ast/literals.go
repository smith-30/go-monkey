package ast

import "github.com/smith-30/go-monkey/token"

type (
	IntegerLiteral struct {
		Token token.Token
		Value int64
	}
)

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type (
	Boolean struct {
		Token token.Token
		Value bool
	}
)

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}
