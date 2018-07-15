package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 + リテラル
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"
	STRING = "STRING" // "<sequence of characters>"

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = ")"
	RPAREN = "("
	LBRACE = "{"
	RBRACE = "}"

	LBRACKET = "["
	RBRACKET = "]"

	COLON = ":"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)
