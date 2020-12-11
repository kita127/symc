package symc

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
	hash
	ampersand
	tilde
	caret
	vertical
	colon
	question
	period
	backslash
	doublequot
)

func NewLexer(src string) *Lexer {
	return &Lexer{input: src, pos: 0}
}

func (l *Lexer) lexicalize(src string) []*Token {
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
		if l.input[i] != ' ' && l.input[i] != '\t' {
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
	case '#':
		tk = &Token{tokenType: hash, literal: "#"}
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
	default:
		if isLetter(c) {
			tk = l.readWord()
		} else if isDigit(c) {
			tk = l.readNumber()
		}
	}
	return tk
}

func (l *Lexer) readWord() *Token {
	// ワードの終わりの次まで pos を進める
	var next int
	for next = l.pos; next < len(l.input); next++ {
		if l.input[next] == ' ' || l.input[next] == '\t' {
			break
		}
	}
	w := l.input[l.pos:next]
	l.pos = next
	return &Token{tokenType: word, literal: w}
}

func (l *Lexer) readNumber() *Token {
	// ワードの終わりの次まで pos を進める
	var next int
	for next = l.pos; next < len(l.input); next++ {
		if l.input[next] == ' ' || l.input[next] == '\t' {
			break
		}
	}
	w := l.input[l.pos:next]
	l.pos = next
	return &Token{tokenType: integer, literal: w}
}

func isLetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
