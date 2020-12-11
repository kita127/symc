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
		tk = &Token{tokenType: assign, literal: "+"}
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
