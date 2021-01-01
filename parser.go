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
	Tk       *Token
}

func (v *InvalidStatement) statementNode() {}
func (v *InvalidStatement) String() string {
	return fmt.Sprintf("InvalidStatement : Contents=%s, Tk=%s", v.Contents, v.Tk.literal)
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
	Name       string
	Params     []*VariableDef
	Statements []Statement
}

func (v *FunctionDef) statementNode() {}
func (v *FunctionDef) String() string {
	return fmt.Sprintf("FunctionDef : Name=%s, Params=%s, Statements=%s", v.Name, v.Params, v.Statements)
}

type BlockStatement struct {
	Statements []Statement
}

func (v *BlockStatement) statementNode() {}
func (v *BlockStatement) String() string {
	s := fmt.Sprintf("BlockStatement={")
	sep := ""
	for _, v2 := range v.Statements {
		s += sep
		sep = ", "
		s += v2.String()
	}
	s += " }"
	return s
}

type AccessVar struct {
	Name string
}

func (v *AccessVar) statementNode() {}
func (v *AccessVar) String() string {
	return fmt.Sprintf("AccessVar : Name=%s", v.Name)
}

type Typedef struct {
	Name string
}

func (v *Typedef) statementNode() {}
func (v *Typedef) String() string {
	return fmt.Sprintf("Typedef : Name=%s", v.Name)
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
		if p.curToken().tokenType == comment {
			p.pos++
			continue
		}
		if s := p.parseStatement(); s != nil {
			ss = append(ss, s)
			if _, yes := s.(*InvalidStatement); yes {
				break
			}
		}
	}
	m := &Module{ss}
	return m
}

func (p *Parser) parseStatement() Statement {
	var s Statement = nil
	p.prevPos = p.pos
	switch p.curToken().tokenType {
	case keyTypedef:
		p.skipTypedef()
	case keyExtern:
		if s == nil {
			s = p.parsePrototypeDecl()
		}
		if s == nil {
			s = p.parseVariableDecl()
		}
		if s == nil {
			return &InvalidStatement{Contents: "could not parse", Tk: p.curToken()}
		}
	case keyStruct:
		p.skipStruct()
	case keyAttribute:
		p.pos++
		p.skipParen()
	default:
		if s == nil {
			s = p.parseFunctionDef()
		}
		if s == nil {
			s = p.parsePrototypeDecl()
		}
		if s == nil {
			s = p.parseVariableDef()
		}
		if s == nil {
			return &InvalidStatement{Contents: "could not parse", Tk: p.curToken()}
		}
	}
	return s
}

//func (p *Parser) parseBlockStatementSub() Statement {
//	var s Statement = nil
//	p.prevPos = p.pos
//	switch p.curToken().tokenType {
//	default:
//		if s == nil {
//			s = p.parseVariableDef()
//		}
//		if s == nil {
//			s = p.parseExpressionStatement()
//		}
//	}
//	return s
//}

func (p *Parser) parseVariableDef() Statement {
	// semicolon or assign or lbracket or eof の手前まで pos を進める
	typeCnt := 0
	n := p.peekToken()
	for n.tokenType != semicolon && n.tokenType != assign && n.tokenType != lbracket && n.tokenType != eof {
		// 現在トークンが識別子もしくは型に関するかチェック
		if !p.curToken().IsTypeToken() {
			p.posReset()
			return nil
		}
		typeCnt++
		p.pos++
		n = p.peekToken()
	}
	if n.tokenType == eof {
		p.posReset()
		return nil
	}

	if !p.curToken().IsTypeToken() {
		p.posReset()
		return nil
	}
	typeCnt++

	if typeCnt <= 1 {
		p.posReset()
		return nil
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
		p.posReset()
		return nil
	}
	return &VariableDef{Name: id}
}

func (p *Parser) parseVariableDecl() Statement {
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

func (p *Parser) parsePrototypeDecl() Statement {
	n := p.peekToken()
	for n.tokenType != lparen {
		// 現在トークンが識別子もしくは型に関するかチェック
		if !p.curToken().IsTypeToken() && p.curToken().tokenType != keyExtern {
			p.posReset()
			return nil
		}
		p.pos++
		n = p.peekToken()
	}

	if !p.curToken().IsTypeToken() {
		p.posReset()
		return nil
	}

	// Name
	id := p.curToken().literal

	// semicolon まで進める
	p.progUntil(semicolon)
	p.pos++
	// next
	return &PrototypeDecl{Name: id}

}

func (p *Parser) parseFunctionDef() Statement {
	// lparen or eof の手前まで pos を進める
	n := p.peekToken()
	for n.IsTypeToken() {
		p.pos++
		if p.curToken().tokenType == keyAttribute {
			p.skipParen()
		}
		n = p.peekToken()
	}
	if n.tokenType != lparen {
		p.posReset()
		return nil
	}

	// Name
	id := p.curToken().literal

	p.pos++

	// 引数のパース
	ps := p.parseParameter()

	// lbrace かチェック
	if p.curToken().tokenType != lbrace {
		p.posReset()
		return nil
	}

	ss := p.parseBlockStatement()

	return &FunctionDef{Name: id, Params: ps, Statements: ss}
}

func (p *Parser) parseBlockStatement() []Statement {
	p.pos++

	var ss []Statement = nil

	for p.curToken().tokenType != rbrace {
		var ts []Statement = nil
		p.prevPos = p.pos
		switch p.curToken().tokenType {
		case lbrace:
			ts = append(ts, p.parseBlockStatement()...)
		case lparen:
			fallthrough
		case word:
			if ts == nil {
				if s := p.parseVariableDef(); s != nil {
					ts = append(ts, s)
				}
			}
			if ts == nil {
				// other statement
				ts = p.parseExpressionStatement()
			}
		}
		ss = append(ss, ts...)
	}

	if ss == nil {
		// 何もない
		ss = []Statement{}
	}

	p.pos++

	return ss
}

func (p *Parser) parseExpressionStatement() []Statement {
	var ss []Statement = nil

	for p.curToken().tokenType != semicolon {
		if p.curToken().tokenType == word {
			id := p.curToken().literal
			ss = append(ss, &AccessVar{Name: id})
		}
		p.pos++
	}

	// semicolon
	p.pos++
	// next

	return ss
}

//func (p *Parser) parseBlockStatement______() Statement {
//	// lbrace の次へ
//	p.pos++
//
//	ss := []Statement{}
//	for p.curToken().tokenType != rbrace {
//		if p.curToken().IsTypeToken() {
//			s := p.parseBlockStatementSub()
//			if s != nil {
//				ss = append(ss, s)
//			} else {
//				return nil
//			}
//		} else if p.curToken().tokenType == lbrace {
//			s := p.parseBlockStatement()
//			if s != nil {
//				ss = append(ss, s)
//			} else {
//				return nil
//			}
//		} else {
//			// パース対象外のトークンの場合はスキップする
//			p.pos++
//		}
//	}
//	p.pos++
//	//next
//
//	return &BlockStatement{ss}
//}

func (p *Parser) parseAccessVar() Statement {
	if p.curToken().tokenType != word {
		p.posReset()
		return nil
	}
	n := p.fetchID()

	return &AccessVar{Name: n}
}

func (p *Parser) parseParameter() []*VariableDef {
	vs := []*VariableDef{}
	p.pos++

	if p.curToken().tokenType == rparen {
		// パラメータになにもなし
		p.pos++
		// next
		return vs
	}

	n := p.peekToken()
	for {
		for n.IsTypeToken() {
			p.pos++
			n = p.peekToken()
		}

		if p.curToken().tokenType == word {
			id := p.curToken().literal
			vs = append(vs, &VariableDef{Name: id})
		} else if p.curToken().tokenType == keyVoid {
			// 何もしない
		}

		if n.tokenType == rparen {
			break
		} else if n.tokenType == comma {
			p.pos++
			n = p.peekToken()
		} else if n.tokenType == lbracket {
			p.progUntil(rbracket)
			n = p.peekToken()
		}
	}

	p.pos++
	p.pos++
	// next

	return vs
}

//func (p *Parser) parseExpressionStatement() Statement {
//	if p.curToken().tokenType == word {
//		id := p.curToken().literal
//		p.pos++
//		// next
//		return &AccessVar{Name: id}
//	} else {
//		p.posReset()
//		return nil
//	}
//}

func (p *Parser) skipTypedef() {

	for p.curToken().tokenType != semicolon && p.curToken().tokenType != eof {
		if p.curToken().tokenType == lbrace {
			// { の場合は } まで進める
			p.progUntil(rbrace)
		} else {
			p.pos++
		}
	}
	p.pos++
	// next
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

func (p *Parser) skipStruct() {
	p.progUntil(rbrace)
	p.progUntil(semicolon)
	p.pos++
}

func (p *Parser) skipParen() {
	for {
		if p.curToken().tokenType == lparen {
			p.pos++
			p.skipParen()
			p.pos++
			return
		} else if p.curToken().tokenType == rparen {
			return
		} else {
			p.pos++
		}
	}
}

func (p *Parser) fetchID() string {

	id := ""
	for {
		id += p.curToken().literal
		p.pos++
		if p.curToken().tokenType != period && p.curToken().tokenType != arrow && p.curToken().tokenType != word {
			break
		}
	}
	return id
}
