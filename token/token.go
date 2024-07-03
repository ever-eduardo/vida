package token

import "unicode"

type Token byte

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
	LPAREN
	RPAREN
	LCURLY
	RCURLY
	operator_end

	binary_op_init
	binary_op_end

	keyword_init
	TRUE
	FALSE
	NIL
	LOC
	AND
	OR
	NOT
	keyword_end
)

var tokens = [...]string{
	UNEXPECTED: "Unexpected",
	EOF:        "EOF",
	COMMENT:    "Comment",
	IDENTIFIER: "Identifier",
	ASSIGN:     "Assign",
	COMMA:      "Comma",
	LPAREN:     "(",
	RPAREN:     ")",
	LCURLY:     "{",
	RCURLY:     "}",
	AND:        "and",
	OR:         "or",
	NOT:        "not",
	TRUE:       "true",
	FALSE:      "false",
	NIL:        "nil",
	LOC:        "loc",
}

type TokenInfo struct {
	Lit   string
	Line  uint
	Token Token
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

func (token Token) IsBinaryOperator() bool {
	return (binary_op_init < token && token < binary_op_end) || token == AND || token == OR
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

const (
	LowestPrec  = 0
	PrefixPrec  = 6
	HighestPrec = 7
)

func (op Token) Precedence() int {
	switch op {
	case OR:
		return 1
	case AND:
		return 2
	}
	return LowestPrec
}
