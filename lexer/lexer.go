package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/ever-eduardo/vida/token"
	"github.com/ever-eduardo/vida/verror"
)

type Lexer struct {
	LexicalError verror.VidaError
	src          []byte
	ModuleName   string
	pointer      int
	leadPointer  int
	srcLen       int
	line         uint
	c            rune
}

const bom = 0xFEFF
const eof = -1
const unexpected = -2

func New(src []byte, moduleName string) *Lexer {
	src = append(src, 10)
	lexer := Lexer{
		src:         src,
		c:           0,
		line:        1,
		pointer:     0,
		leadPointer: 0,
		srcLen:      len(src),
		ModuleName:  moduleName,
	}
	lexer.next()
	if lexer.c == bom {
		lexer.next()
	}
	return &lexer
}

func (l *Lexer) next() {
	if l.leadPointer < l.srcLen {
		l.pointer = l.leadPointer
		if l.c == '\n' {
			l.line++
		}
		r, w := rune(l.src[l.leadPointer]), 1
		if r >= utf8.RuneSelf {
			r, w = utf8.DecodeRune(l.src[l.leadPointer:])
			if r == utf8.RuneError && w == 1 {
				r = unexpected
				l.LexicalError = verror.New(l.ModuleName, fmt.Sprintf("The file %v is not has not correct utf-8 encoding", l.ModuleName), verror.FileErrMsg, l.line)
			} else if r == bom && l.pointer > 0 {
				r = unexpected
				l.LexicalError = verror.New(l.ModuleName, fmt.Sprintf("The file %v has a bom in an expectec place", l.ModuleName), verror.FileErrMsg, l.line)
			}
		}
		l.c = r
		l.leadPointer += w
	} else {
		l.c = eof
	}
}

func (l *Lexer) peek() byte {
	if l.leadPointer < l.srcLen {
		return l.src[l.leadPointer]
	}
	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.c == ' ' || l.c == '\t' || l.c == '\n' || l.c == '\r' {
		l.next()
	}
}

func lower(c rune) rune {
	return 32 | c
}

func isDecimal(c rune) bool {
	return '0' <= c && c <= '9'
}

func isOctal(c rune) bool {
	return '0' <= c && c <= '7'
}

func isBin(c rune) bool {
	return '0' <= c && c <= '1'
}

func isHex(ch rune) bool {
	return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f'
}

func isLetter(c rune) bool {
	return 'a' <= lower(c) && lower(c) <= 'z' || c == '_' || c >= utf8.RuneSelf && unicode.IsLetter(c)
}

func isDigit(c rune) bool {
	return isDecimal(c) || c >= utf8.RuneSelf && unicode.IsDigit(c)
}

func (l *Lexer) scanComment() token.Token {
	if l.c == '/' {
		l.next()
		for l.c != '\n' && l.c >= 0 {
			l.next()
		}
		if l.c == '\n' {
			l.next()
		}
		goto exit
	}
	l.next()
	for l.c >= 0 {
		ch := l.c
		if ch == '\n' {
			l.line++
		}
		l.next()
		if ch == '*' && l.c == '/' {
			l.next()
			goto exit
		}
	}
exit:
	if l.c < 0 {
		l.LexicalError = verror.New(l.ModuleName, "Unterminated comment ", verror.LexicalErrMsg, l.line)
		return token.UNEXPECTED
	}
	return token.COMMENT
}

func (l *Lexer) scanString() (token.Token, string) {
	init := l.pointer - 1
	for {
		ch := l.c
		if ch == '\n' || ch < 0 {
			l.c = unexpected
			l.LexicalError = verror.New(l.ModuleName, fmt.Sprintf("The file %v has an unterminated string literal", l.ModuleName), verror.FileErrMsg, l.line)
			return token.UNEXPECTED, ""
		}
		l.next()
		if ch == '"' {
			break
		}
		if ch == '\\' && l.c == '"' {
			l.next()
		}
	}
	return token.STRING, string(l.src[init:l.pointer])
}

func (l *Lexer) scanRawString() (token.Token, string) {
	init := l.pointer - 1
	hasCR := false
	for {
		ch := l.c
		if ch < 0 {
			l.c = unexpected
			l.LexicalError = verror.New(l.ModuleName, fmt.Sprintf("The file %v has an unterminated string literal", l.ModuleName), verror.FileErrMsg, l.line)
			return token.UNEXPECTED, ""
		}
		l.next()
		if ch == '`' {
			break
		}
		if ch == 'r' {
			hasCR = true
		}
	}
	lit := l.src[init:l.pointer]
	if hasCR {
		lit = stripCR(lit)
	}
	return token.STRING, string(lit)
}

func stripCR(b []byte) []byte {
	lb := len(b)
	c := make([]byte, lb)
	i := 0
	for j, ch := range b {
		if ch != 'r' || j+1 < lb {
			c[i] = ch
			i++
		}
	}
	return c[:i]
}

func (l *Lexer) scanIdentifier() string {
	pointer := l.pointer
	for leadPointer, b := range l.src[l.leadPointer:] {
		if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9' {
			continue
		}
		l.leadPointer += leadPointer
		if 0 < b && b < utf8.RuneSelf {
			l.c = rune(b)
			l.pointer = l.leadPointer
			l.leadPointer++
			goto exit
		}
		l.next()
		for isLetter(l.c) || isDigit(l.c) {
			l.next()
		}
		goto exit
	}
exit:
	return string(l.src[pointer:l.pointer])
}

func (l *Lexer) scanNumber() (token.Token, string) {
	init := l.pointer
	tok := token.INTEGER
	if l.c != '.' {
		if l.c == '0' {
			l.next()
			switch lower(l.c) {
			case 'x':
				l.next()
				for isHex(l.c) || l.c == '_' {
					l.next()
				}
			case 'b':
				l.next()
				for isBin(l.c) || l.c == '_' {
					l.next()
				}
			case 'o':
				l.next()
				for isOctal(l.c) || l.c == '_' {
					l.next()
				}
			case '.':
				goto fractional
			default:
				for isOctal(l.c) || l.c == '_' {
					l.next()
				}
			}
		} else {
			for isDecimal(l.c) || l.c == '_' {
				l.next()
			}
		}
	}
fractional:
	if l.c == '.' && rune(l.peek()) != '.' {
		tok = token.FLOAT
		l.next()
		for isDecimal(l.c) || l.c == '_' {
			l.next()
		}
	}

	if e := lower(l.c); e == 'e' || e == 'p' {
		l.next()
		tok = token.FLOAT
		if l.c == '+' || l.c == '-' {
			l.next()
		}
		for isDecimal(l.c) || l.c == '_' {
			l.next()
		}
	}
	return tok, string(l.src[init:l.pointer])
}

func (l *Lexer) Next() (line uint, tok token.Token, lit string) {
	l.skipWhitespace()
	line = l.line
	switch ch := l.c; {
	case isLetter(ch):
		lit = l.scanIdentifier()
		if len(lit) > 1 {
			tok = token.LookUp(lit)
		} else {
			tok = token.IDENTIFIER
		}
	case isDecimal(ch) || l.c == '.' && isDecimal(rune(l.peek())):
		tok, lit = l.scanNumber()
	default:
		l.next()
		switch ch {
		case eof:
			tok = token.EOF
		case '=':
			if l.c == '=' {
				l.next()
				tok = token.EQ
			} else if l.c == '>' {
				l.next()
				tok = token.ARROW
			} else {
				tok = token.ASSIGN
			}
		case '"':
			tok, lit = l.scanString()
		case '`':
			tok, lit = l.scanRawString()
		case '+':
			tok = token.ADD
		case '-':
			tok = token.SUB
		case '*':
			tok = token.MUL
		case '/':
			if l.c == '/' || l.c == '*' {
				tok = l.scanComment()
			} else {
				tok = token.DIV
			}
		case '%':
			tok = token.REM
		case ',':
			tok = token.COMMA
		case '.':
			if l.c == '.' {
				l.next()
				if l.c == '.' {
					l.next()
					tok = token.ELLIPSIS
				} else {
					tok = token.DOUBLE_DOT
				}
			} else {
				tok = token.DOT
			}
		case '!':
			if l.c == '=' {
				l.next()
				tok = token.NEQ
			} else {
				tok = token.UNEXPECTED
				l.LexicalError = verror.New(l.ModuleName, "Unexpected symbol "+string('!'), verror.LexicalErrMsg, l.line)
			}
		case '<':
			if l.c == '=' {
				l.next()
				tok = token.LE
			} else {
				tok = token.LT
			}
		case '>':
			if l.c == '=' {
				l.next()
				tok = token.GE
			} else {
				tok = token.GT
			}
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '{':
			tok = token.LCURLY
		case '}':
			tok = token.RCURLY
		case '[':
			tok = token.LBRACKET
		case ']':
			tok = token.RBRACKET
		case ':':
			tok = token.COLON
		default:
			tok = token.UNEXPECTED
			lit = string(ch)
			l.LexicalError = verror.New(l.ModuleName, "Unexpected symbol "+lit, verror.LexicalErrMsg, l.line)
		}
	}
	return
}
