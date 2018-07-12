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
	StringLiteral struct {
		Token token.Token
		Value string
	}
)

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

func (sl *StringLiteral) String() string {
	return sl.Token.Literal
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
