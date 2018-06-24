package parser

import (
	"fmt"

	"github.com/smith-30/go-monkey/ast"
	"github.com/smith-30/go-monkey/lexer"
	"github.com/smith-30/go-monkey/token"
)

type (
	Parser struct {
		l *lexer.Lexer

		currentToken token.Token
		peekToken    token.Token

		errors []string

		prefixParseFns map[token.TokenType]prefixParseFn
		infixParseFns  map[token.TokenType]infixParseFn
	}
)

type (
	// 前置構文解析
	prefixParseFn func() ast.Expression

	// 中置構文解析
	infixParseFn func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// read 2 token. both of currentToken and peekToken will be set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// Todo: We skipped reading until it reaches the semicolon.
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Todo: We skipped reading until it reaches the semicolon.
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
