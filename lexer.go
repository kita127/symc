package symc

import (
	"fmt"
	"strings"
)

type Lexer struct {
	input string
	pos   int
}

type Token struct {
	tokenType int
	literal   string
}

const (
	eof = iota
	word
	integer
	float
	assign
	plus
	minus
	bang
	asterisk
	slash
	lt
	gt
	semicolon
	lparen
	rparen
	comma
	lbrace
	rbrace
	lbracket
	rbracket
	ampersand
	tilde
	caret
	vertical
	colon
	question
	period
	backslash
	doublequot
	keyReturn
	keyIf
	keyElse
	keyWhile
	keyDo
	keyGoto
	keyFor
	keyBreak
	keyContinue
	keySwitch
	keyCase
	keyDefault
	keyExtern
	keyVolatile
	keyConst
	comment
	illegal
)

func (t *Token) String() string {
	return fmt.Sprintf("tokenType:%v, literal:%s", t.tokenType, t.literal)
}

func NewLexer(src string) *Lexer {
	return &Lexer{input: src, pos: 0}
}

func (l *Lexer) lexicalize() []*Token {
	ts := []*Token{}
	for {
		t := l.nextToken()
		ts = append(ts, t)
		if t.tokenType == eof {
			break
		}
	}
	return ts
}

func (l *Lexer) nextToken() *Token {
	// スペースをとばす
	for {
		i := l.pos
		if i >= len(l.input) {
			break
		}
		c := l.input[i]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		l.pos++
	}

	// ソースの終端
	if l.pos >= len(l.input) {
		return &Token{tokenType: eof, literal: "eof"}
	}

	var tk *Token
	c := l.input[l.pos]
	switch c {
	case '=':
		tk = &Token{tokenType: assign, literal: "="}
		l.pos++
	case '+':
		tk = &Token{tokenType: plus, literal: "+"}
		l.pos++
	case '-':
		tk = &Token{tokenType: minus, literal: "-"}
		l.pos++
	case '!':
		tk = &Token{tokenType: bang, literal: "!"}
		l.pos++
	case '*':
		tk = &Token{tokenType: asterisk, literal: "*"}
		l.pos++
	case '/':
		tk = &Token{tokenType: slash, literal: "/"}
		l.pos++
	case '<':
		tk = &Token{tokenType: lt, literal: "<"}
		l.pos++
	case '>':
		tk = &Token{tokenType: gt, literal: ">"}
		l.pos++
	case ';':
		tk = &Token{tokenType: semicolon, literal: ";"}
		l.pos++
	case '(':
		tk = &Token{tokenType: lparen, literal: "("}
		l.pos++
	case ')':
		tk = &Token{tokenType: rparen, literal: ")"}
		l.pos++
	case ',':
		tk = &Token{tokenType: comma, literal: ","}
		l.pos++
	case '{':
		tk = &Token{tokenType: lbrace, literal: "{"}
		l.pos++
	case '}':
		tk = &Token{tokenType: rbrace, literal: "}"}
		l.pos++
	case '[':
		tk = &Token{tokenType: lbracket, literal: "["}
		l.pos++
	case ']':
		tk = &Token{tokenType: rbracket, literal: "]"}
		l.pos++
	case '&':
		tk = &Token{tokenType: ampersand, literal: "&"}
		l.pos++
	case '~':
		tk = &Token{tokenType: tilde, literal: "~"}
		l.pos++
	case '^':
		tk = &Token{tokenType: caret, literal: "^"}
		l.pos++
	case '|':
		tk = &Token{tokenType: vertical, literal: "|"}
		l.pos++
	case ':':
		tk = &Token{tokenType: colon, literal: ":"}
		l.pos++
	case '?':
		tk = &Token{tokenType: question, literal: "?"}
		l.pos++
	case '.':
		tk = &Token{tokenType: period, literal: "."}
		l.pos++
	case '\\':
		tk = &Token{tokenType: backslash, literal: "\\"}
		l.pos++
	case '"':
		tk = &Token{tokenType: doublequot, literal: "\""}
		l.pos++
	case '#':
		tk = l.readHashComment()
		l.pos++
	default:
		if isLetter(c) {
			tk = l.readWord()
		} else if isDec(c) {
			tk = l.readNumber()
		}
	}
	return tk
}

func (l *Lexer) readWord() *Token {
	// ワードの終わりの次まで pos を進める
	var next int
	for next = l.pos; next < len(l.input); next++ {
		c := l.input[next]
		if !isLetter(c) && !isDec(c) {
			break
		}
	}
	w := l.input[l.pos:next]
	tk := l.determineKeyword(w)
	l.pos = next
	return tk
}

func (l *Lexer) readNumber() *Token {
	var next int
	isFloat := false

	next = l.pos
	c := l.input[next]
	if c == '0' {
		next++
		c = l.input[next]
		switch c {
		case 'x':
			// 16進数
			next++
		case 'b':
			// 2進数
			next++
		case '.':
			// 小数
			next++
			isFloat = true
		default:
			if isDec(c) {
				// 8進数
				next++
			} else {
				// エラー
				return l.newIllegal()
			}
		}
	}

	// ワードの終わりの次まで pos を進める
	for ; next < len(l.input); next++ {
		c = l.input[next]
		if c == '.' {
			isFloat = true
		} else if c == 'u' || c == 'U' || c == 'l' || c == 'L' {
			continue
		} else if !isHex(c) {
			break
		}
	}
	w := l.input[l.pos:next]
	l.pos = next

	var tk *Token
	if isFloat {
		tk = &Token{tokenType: float, literal: w}
	} else {
		tk = &Token{tokenType: integer, literal: w}
	}
	return tk
}

func (l *Lexer) readHashComment() *Token {
	// # の次の文字に移動
	l.pos++
	var next int
	for i := l.pos; i <= len(l.input); i++ {
		next = i
		if next >= len(l.input) {
			break
		}
		c := l.input[next]
		if c == '\n' || c == '\r' {
			break
		}
	}
	tk := &Token{tokenType: comment, literal: l.input[l.pos:next]}
	l.pos = next
	return tk
}

func (l *Lexer) newIllegal() *Token {
	tk := &Token{tokenType: illegal, literal: l.input[l.pos:]}
	l.pos = len(l.input)
	return tk
}

func (l *Lexer) determineKeyword(w string) *Token {
	if strings.Compare("return", w) == 0 {
		return &Token{tokenType: keyReturn, literal: w}
	} else if strings.Compare("if", w) == 0 {
		return &Token{tokenType: keyIf, literal: w}
	} else if strings.Compare("else", w) == 0 {
		return &Token{tokenType: keyElse, literal: w}
	} else if strings.Compare("while", w) == 0 {
		return &Token{tokenType: keyWhile, literal: w}
	} else if strings.Compare("do", w) == 0 {
		return &Token{tokenType: keyDo, literal: w}
	} else if strings.Compare("goto", w) == 0 {
		return &Token{tokenType: keyGoto, literal: w}
	} else if strings.Compare("for", w) == 0 {
		return &Token{tokenType: keyFor, literal: w}
	} else if strings.Compare("break", w) == 0 {
		return &Token{tokenType: keyBreak, literal: w}
	} else if strings.Compare("continue", w) == 0 {
		return &Token{tokenType: keyContinue, literal: w}
	} else if strings.Compare("switch", w) == 0 {
		return &Token{tokenType: keySwitch, literal: w}
	} else if strings.Compare("case", w) == 0 {
		return &Token{tokenType: keyCase, literal: w}
	} else if strings.Compare("default", w) == 0 {
		return &Token{tokenType: keyDefault, literal: w}
	} else if strings.Compare("extern", w) == 0 {
		return &Token{tokenType: keyExtern, literal: w}
	} else if strings.Compare("volatile", w) == 0 {
		return &Token{tokenType: keyVolatile, literal: w}
	} else if strings.Compare("const", w) == 0 {
		return &Token{tokenType: keyConst, literal: w}
	} else {
		return &Token{tokenType: word, literal: w}
	}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func isDec(c byte) bool {
	return '0' <= c && c <= '9'
}
