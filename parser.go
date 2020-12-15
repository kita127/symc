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

type PrototypeDecl struct {
	Name string
}

func (v *PrototypeDecl) statementNode() {}
func (v *PrototypeDecl) String() string {
	return fmt.Sprintf("PrototypeDecl : Name=%s", v.Name)
}

type FunctionDef struct {
	Name string
}

func (v *FunctionDef) statementNode() {}
func (v *FunctionDef) String() string {
	return fmt.Sprintf("FunctionDef : Name=%s", v.Name)
}

// 構文解析器
type Parser struct {
	lexer   *Lexer
	tokens  []*Token
	pos     int
	prevPos int
}

func NewParser(l *Lexer) *Parser {
	tks := l.lexicalize()
	return &Parser{lexer: l, tokens: tks, pos: 0, prevPos: 0}
}

func (p *Parser) Parse() *Module {
	ast := p.parseModule()
	return ast
}

func (p *Parser) parseModule() *Module {
	ss := []Statement{}
	for p.curToken().tokenType != eof {
		var s Statement = &InvalidStatement{Contents: "unknown syntax"}
		switch p.curToken().tokenType {
		case keyExtern:
			s = p.parsePrototypeDecl(s)
			s = p.parseVariableDecl(s)
		default:
			s = p.parseFunctionDef(s)
			s = p.parsePrototypeDecl(s)
			s = p.parseVariableDef(s)
		}
		ss = append(ss, s)
		if _, invalid := s.(*InvalidStatement); invalid {
			break
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
		p.pos = p.prevPos
		return &InvalidStatement{Contents: "err parse variable def"}
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
		p.pos = p.prevPos
		return &InvalidStatement{Contents: "err parse variable def"}
	}
	p.prevPos = p.pos
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

	p.prevPos = p.pos
	return &VariableDecl{Name: id}
}

func (p *Parser) parsePrototypeDecl(s Statement) Statement {
	if _, invalid := s.(*InvalidStatement); !invalid {
		// 既に解析済みの場合はリターン
		return s
	}
	errMsg := "err parse prototype decl"

	// lparen or eof の手前まで pos を進める
	p.progUntilPrev(lparen)
	if p.peekToken().tokenType == eof {
		p.pos = p.prevPos
		return &InvalidStatement{Contents: errMsg}
	}

	// Name
	id := p.curToken().literal

	// rparen or eof まで pos を進める
	p.progUntil(rparen)
	if p.curToken().tokenType == eof {
		p.pos = p.prevPos
		return &InvalidStatement{Contents: errMsg}
	}
	p.pos++
	if p.curToken().tokenType != semicolon {
		p.pos = p.prevPos
		return &InvalidStatement{Contents: errMsg}
	}
	p.pos++
	// next
	p.prevPos = p.pos
	return &PrototypeDecl{Name: id}

}

func (p *Parser) parseFunctionDef(s Statement) Statement {
	if _, invalid := s.(*InvalidStatement); !invalid {
		// 既に解析済みの場合はリターン
		return s
	}
	errMsg := "err parse function def"

	// lparen or eof の手前まで pos を進める
	p.progUntilPrev(lparen)
	if p.peekToken().tokenType == eof {
		p.posReset()
		return &InvalidStatement{Contents: errMsg}
	}

	// Name
	id := p.curToken().literal

	// rparen or eof まで pos を進める
	p.progUntil(rparen)
	if p.curToken().tokenType == eof {
		p.posReset()
		return &InvalidStatement{Contents: errMsg}
	}

	// lbrace or eof まで pos を進める
	p.progUntil(lbrace)
	if p.curToken().tokenType == eof {
		p.posReset()
		return &InvalidStatement{Contents: errMsg}
	}

	// rbrace or eof まで pos を進める
	p.progUntil(rbrace)
	if p.curToken().tokenType == eof {
		p.posReset()
		return &InvalidStatement{Contents: errMsg}
	}
	p.pos++
	// next

	p.prevPos = p.pos
	return &FunctionDef{Name: id}
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

func (p *Parser) progUntil(tkType int) {
	t := p.curToken()
	for t.tokenType != tkType && t.tokenType != eof {
		p.pos++
		t = p.curToken()
	}
}

func (p *Parser) progUntilPrev(tkType int) {
	n := p.peekToken()
	for n.tokenType != tkType && n.tokenType != eof {
		p.pos++
		n = p.peekToken()
	}
}

func (p *Parser) posReset() {
	p.pos = p.prevPos
}
