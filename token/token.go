package token

import "unicode"

type Token int

const (
	UNEXPECTED Token = iota
	EOF
	COMMENT

	literal_init
	IDENTIFIER
	literal_end

	operator_init
	ASSIGN
	COMMA
	operator_end

	keyword_init
	TRUE
	FALSE
	NIL
	LOC
	keyword_end
)

var tokens = [...]string{
	UNEXPECTED: "Unexpected",
	EOF:        "EOF",
	COMMENT:    "Comment",
	IDENTIFIER: "Identifier",
	ASSIGN:     "Assign",
	COMMA:      "Comma",
	TRUE:       "true",
	FALSE:      "false",
	NIL:        "nil",
	LOC:        "loc",
}

type TokenInfo struct {
	Lit   string
	Token Token
	Line  uint
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, keyword_end-(keyword_init+1))
	for i := keyword_init + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func (token Token) String() string {
	return tokens[token]
}

func (token Token) IsLiteral() bool {
	return literal_init < token && token < literal_end
}

func (token Token) IsKeyword() bool {
	return keyword_init < token && token < keyword_end
}

func (token Token) IsOperator() bool {
	return operator_init < token && token < operator_end
}

func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

func LookUp(name string) Token {
	if token, is_keyword := keywords[name]; is_keyword {
		return token
	}
	return IDENTIFIER
}

func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return true
}
