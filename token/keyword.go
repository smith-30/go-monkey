package token

var (
	keywords = map[string]TokenType{
		"fn":  FUNCTION,
		"let": LET,
	}
)

func LookUpIdent(i string) TokenType {
	if t, ok := keywords[i]; ok {
		return t
	}
	return IDENT
}
