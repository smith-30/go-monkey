package ast

import "github.com/smith-30/go-monkey/token"

type (
	Node interface {
		TokenLiteral() string
	}

	Statement interface {
		Node
		statementNode()
	}

	Expression interface {
		Node
		expressionNode()
	}
)

type (
	LetStatement struct {
		Token token.Token // expects token.LET
		Name  *Identifier
		Value Expression
	}
)

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type (
	Identifier struct {
		Token token.Token // expects token.IDENT
		Value string
	}
)

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
