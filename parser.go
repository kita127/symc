package symc

import (
	"fmt"
)

type Module struct {
	Statements []Statement
}

func (m *Module) String() string {
	s := "Module : Statements={ "
	sep := ""
	for _, v := range m.Statements {
		s += sep
		sep = ", "
		s += v.String()
	}
	s += " }"
	return s
}

type Statement interface {
	statementNode()
	fmt.Stringer
}

type InvalidStatement struct {
	Contents string
}

func (v *InvalidStatement) statementNode() {}
func (v *InvalidStatement) String() string {
	return fmt.Sprintf("InvalidStatement : Contents=%s", v.Contents)
}

type VariableDef struct {
	Name string
}

func (v *VariableDef) statementNode() {}
func (v *VariableDef) String() string {
	return fmt.Sprintf("VariableDef : Name=%s", v.Name)
}

type VariableDecl struct {
	Name string
}

func (v *VariableDecl) statementNode() {}
func (v *VariableDecl) String() string {
	return fmt.Sprintf("VariableDecl : Name=%s", v.Name)
}

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
	for p.curToken().tokenType != eof && !occurredInvalid {
		var s Statement = &InvalidStatement{Contents: "unknown syntax"}
		switch p.curToken().tokenType {
		case keyExtern:
			s = p.parseVariableDecl(s)
		default:
			s = p.parseVariableDef(s)
		}
		ss = append(ss, s)
		if _, invalid := s.(*InvalidStatement); invalid {
			occurredInvalid = true
		}
	}
	m := &Module{ss}
	return m
}

func (p *Parser) parseVariableDef(s Statement) Statement {
	if _, invalid := s.(*InvalidStatement); !invalid {
		// 既に解析済みの場合はリターン
		return s
	}

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
	switch p.curToken().tokenType {
	case semicolon:
		fallthrough
	case assign:
		fallthrough
	case lbracket:
		// semicolon まで進める
		for p.curToken().tokenType != semicolon {
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

func (p *Parser) parseVariableDecl(s Statement) Statement {
	if _, invalid := s.(*InvalidStatement); !invalid {
		// 既に解析済みの場合はリターン
		return s
	}
	// セミコロンの手前まで pos を進める
	for p.peekToken().tokenType != semicolon {
		p.pos++
	}
	// Name
	id := p.curToken().literal
	p.pos++
	// semicolon
	p.pos++
	// next
	return &VariableDecl{Name: id}
}

func (p *Parser) peekToken() *Token {
	// 現在位置が EOF
	if p.curToken().tokenType == eof {
		return p.curToken()
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) curToken() *Token {
	return p.tokens[p.pos]
}
