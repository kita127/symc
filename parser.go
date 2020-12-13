package symc

import _ "fmt"

type Module struct {
	Statements []Statement
}
type Statement interface {
	statementNode()
}

type VariableDef struct {
	Name string
}

func (v *VariableDef) statementNode() {}

type VariableDecl struct {
	Name string
}

func (v *VariableDecl) statementNode() {}

type Parser struct {
	lexer  *Lexer
	tokens []*Token
	pos    int
}

func NewParser(l *Lexer) *Parser {
	tks := l.lexicalize()
	return &Parser{lexer: l, tokens: tks, pos: 0}
}

func (p *Parser) Parse() *Module {
	ast := p.parseModule()
	return ast
}

func (p *Parser) parseModule() *Module {
	ss := []Statement{}
	t := p.tokens[p.pos]
	for t.tokenType != eof {
		var s Statement
		switch t.tokenType {
		case keyExtern:
			s = p.parseVariableDecl()
		default:
			s = p.parseVariableDef()
		}
		ss = append(ss, s)
		t = p.tokens[p.pos]
	}
	m := &Module{ss}
	return m
}

func (p *Parser) parseVariableDef() Statement {

	// semicolon or assign の手前まで pos を進める
	t := p.peekToken()
	for t.tokenType != semicolon && t.tokenType != assign {
		p.pos++
		t = p.peekToken()
	}
	// Name
	id := p.tokens[p.pos].literal
	p.pos++
	// semicolon or assign
	t = p.tokens[p.pos]
	switch t.tokenType {
	case semicolon:
		p.pos++
		// next
	case assign:
		p.pos++
		// value
		p.pos++
		// semicolon
		p.pos++
		// next
	}
	return &VariableDef{Name: id}
}

func (p *Parser) parseVariableDecl() Statement {
	// セミコロンの手前まで pos を進める
	for p.peekToken().tokenType != semicolon {
		p.pos++
	}
	// Name
	id := p.tokens[p.pos].literal
	p.pos++
	// semicolon
	p.pos++
	// next
	return &VariableDecl{Name: id}
}

func (p *Parser) peekToken() *Token {
	// 現在位置が EOF
	t := p.tokens[p.pos]
	if t.tokenType == eof {
		return t
	}
	// 次の位置が EOF
	n := p.tokens[p.pos+1]
	if n.tokenType == eof {
		return n
	}
	return n
}
