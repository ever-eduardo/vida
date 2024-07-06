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
	case isDecimal(ch):
		tok, lit = l.scanNumber()
	default:
		l.next()
		switch ch {
		case eof:
			tok = token.EOF
		case '=':
			tok = token.ASSIGN
		case '+':
			tok = token.ADD
		case '-':
			tok = token.SUB
		case '*':
			tok = token.MUL
		case '/':
			tok = token.DIV
		case '%':
			tok = token.REM
		// case ',':
		// 	tok = token.COMMA
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '{':
			tok = token.LCURLY
		case '}':
			tok = token.RCURLY
		default:
			tok = token.UNEXPECTED
			lit = string(ch)
			l.LexicalError = verror.New(l.ModuleName, "Unexpected symbol "+lit, verror.LexicalErrMsg, l.line)
		}
	}
	return
}
