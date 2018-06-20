package token

var (
	keywords = map[string]TokenType{
		"fn":  FUNCTION,
		"let": LET,
		"true":  TRUE,
		"false": FALSE,
		"if":  IF,
		"else": ELSE,
		"return":  RETURN,
	}
)

func LookUpIdent(i string) TokenType {
	if t, ok := keywords[i]; ok {
		return t
	}
	return IDENT
}
