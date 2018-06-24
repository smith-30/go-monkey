package ast

import (
	"github.com/smith-30/go-monkey/token"
)

type (
	ExpressionStatement struct {
		Token      token.Token // first statement token
		Expression Expression
	}
)

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
