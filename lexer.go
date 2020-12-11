package symc

type Lexer struct {
	input    string
	position int
}

type Token struct {
	tokenType int
	literal   string
}

const (
	word = iota
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
	eof
)

func NewLexer(src string) *Lexer {
	return &Lexer{input: src, position: 0}
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
		i := l.position
		if i >= len(l.input) {
			break
		}
		if l.input[i] != ' ' && l.input[i] != '\t' {
			break
		}
		l.position++
	}

	// ソースの終端
	if l.position >= len(l.input) {
		return &Token{tokenType: eof, literal: "eof"}
	}

	var tk *Token
	c := l.input[l.position]
	switch c {
	case '=':
		tk = &Token{tokenType: assign, literal: "="}
		l.position++
	case '+':
		tk = &Token{tokenType: plus, literal: "+"}
		l.position++
	case '-':
		tk = &Token{tokenType: minus, literal: "-"}
		l.position++
	case '!':
		tk = &Token{tokenType: bang, literal: "!"}
		l.position++
	case '*':
		tk = &Token{tokenType: asterisk, literal: "*"}
		l.position++
	case '/':
		tk = &Token{tokenType: slash, literal: "/"}
		l.position++
	case '<':
		tk = &Token{tokenType: lt, literal: "<"}
		l.position++
	case '>':
		tk = &Token{tokenType: gt, literal: ">"}
		l.position++
	case ';':
		tk = &Token{tokenType: semicolon, literal: ";"}
		l.position++
	case '(':
		tk = &Token{tokenType: lparen, literal: "("}
		l.position++
	case ')':
		tk = &Token{tokenType: rparen, literal: ")"}
		l.position++
	case ',':
		tk = &Token{tokenType: comma, literal: ","}
		l.position++
	case '{':
		tk = &Token{tokenType: lbrace, literal: "{"}
		l.position++
	case '}':
		tk = &Token{tokenType: rbrace, literal: "}"}
		l.position++
	case '[':
		tk = &Token{tokenType: lbracket, literal: "["}
		l.position++
	case ']':
		tk = &Token{tokenType: rbracket, literal: "]"}
		l.position++
	case '#':
		tk = &Token{tokenType: hash, literal: "#"}
		l.position++
	case '&':
		tk = &Token{tokenType: ampersand, literal: "&"}
		l.position++
	case '~':
		tk = &Token{tokenType: tilde, literal: "~"}
		l.position++
	case '^':
		tk = &Token{tokenType: caret, literal: "^"}
		l.position++
	case '|':
		tk = &Token{tokenType: vertical, literal: "|"}
		l.position++
	case ':':
		tk = &Token{tokenType: colon, literal: ":"}
		l.position++
	case '?':
		tk = &Token{tokenType: question, literal: "?"}
		l.position++
	case '.':
		tk = &Token{tokenType: period, literal: "."}
		l.position++
	case '\\':
		tk = &Token{tokenType: backslash, literal: "\\"}
		l.position++
	default:
		tk = l.readWord()
	}
	return tk
}

func (l *Lexer) readWord() *Token {
	// ワードの終わりの次まで position を進める
	var next int
	for next = l.position; next < len(l.input); next++ {
		if l.input[next] == ' ' || l.input[next] == '\t' {
			break
		}
	}
	w := l.input[l.position:next]
	l.position = next
	return &Token{tokenType: word, literal: w}
}
