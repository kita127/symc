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
	// type
	p.pos++
	id := p.tokens[p.pos].literal
	p.pos++
	// semicolon
	p.pos++
	// next
	return &VariableDef{Name: id}
}

func (p *Parser) parseVariableDecl() Statement {
	// extern
	p.pos++
	// type
	p.pos++
	id := p.tokens[p.pos].literal
	p.pos++
	// semicolon
	p.pos++
	// next
	return &VariableDecl{Name: id}
}
