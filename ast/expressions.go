package ast

import (
	"bytes"

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

type (
	PrefixExpression struct {
		Token    token.Token // expects ex. !, -
		Operator string
		Right    Expression
	}
)

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
