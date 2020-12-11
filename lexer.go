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
	eof
)

func NewLexer(src string) *Lexer {
	return &Lexer{input: src, position: 0}
}

func (l *Lexer) lexicalize(src string) []*Token {
	ts := []*Token{}
	t := l.nextToken()
	ts = append(ts, t)
	return ts
}

func (l *Lexer) nextToken() *Token {
	// スペースをとばす
	for i := l.position; i < len(l.input); i++ {
		if l.input[i] != ' ' {
			l.position = i
			break
		}
	}
	var next int
	for next = l.position; next < len(l.input); next++ {
		if l.input[next] == ' ' {
			break
		}
	}
	w := l.input[l.position:next]
	l.position = next
	return &Token{tokenType: word, literal: w}
}
