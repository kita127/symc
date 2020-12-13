package symc

import (
	_ "fmt"
)

type Module struct {
	Statements []Statement
}
type Statement interface {
	statementNode()
}

type InvalidStatement struct {
	Contents string
}

func (v *InvalidStatement) statementNode() {}

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
	occurredInvalid := false
	for p.currentToken().tokenType != eof && !occurredInvalid {
		var s Statement
		switch p.currentToken().tokenType {
		case keyExtern:
			s = p.parseVariableDecl()
		default:
			s = p.parseVariableDef()
		}
		ss = append(ss, s)
		switch s.(type) {
		case *InvalidStatement:
			occurredInvalid = true
		}
	}
	m := &Module{ss}
	return m
}

func (p *Parser) parseVariableDef() Statement {

	// semicolon or assign or lbracket or eof の手前まで pos を進める
	n := p.peekToken()
	for n.tokenType != semicolon && n.tokenType != assign && n.tokenType != lbracket && n.tokenType != eof {
		p.pos++
		n = p.peekToken()
	}

	if n.tokenType == eof {
		s := "err parse variable def"
		return &InvalidStatement{Contents: s}
	}
	// Name
	id := p.tokens[p.pos].literal
	p.pos++
	// semicolon or assign or lbracket
	switch p.currentToken().tokenType {
	case semicolon:
		fallthrough
	case assign:
		fallthrough
	case lbracket:
		// semicolon まで進める
		for p.currentToken().tokenType != semicolon {
			p.pos++
		}
		p.pos++
		// next
	default:
		s := "err parse variable def"
		return &InvalidStatement{Contents: s}
	}
	return &VariableDef{Name: id}
}

func (p *Parser) parseVariableDecl() Statement {
	// セミコロンの手前まで pos を進める
	for p.peekToken().tokenType != semicolon {
		p.pos++
	}
	// Name
	id := p.currentToken().literal
	p.pos++
	// semicolon
	p.pos++
	// next
	return &VariableDecl{Name: id}
}

func (p *Parser) peekToken() *Token {
	// 現在位置が EOF
	if p.currentToken().tokenType == eof {
		return p.currentToken()
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) currentToken() *Token {
	return p.tokens[p.pos]
}
