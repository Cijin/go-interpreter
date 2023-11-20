package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// operators
	ASSIGN   = "="
	EQ       = "=="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	NOT_EQ   = "!="
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LSQUIRLY = "{"
	RSQUIRLY = "}"

	// keywords
	LET      = "LET"
	FUNCTION = "FUNCTION"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if ttype, ok := keywords[ident]; ok {
		return ttype
	}

	return IDENT
}
