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

type Nothing struct {
}

func (v *Nothing) statementNode() {}
func (v *Nothing) String() string {
	return fmt.Sprintf("Nothing")
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

func (p *Parser) extractVarName() (string, error) {
	e := fmt.Errorf("fail extractVarName token=%s", p.curToken().literal)

	for {
		if p.curToken().tokenType == lparen {
			p.pos++
			s, err := p.extractVarName()
			if err != nil {
				return "", e
			}
			// rparen
			p.pos++
			return s, nil
		} else if !p.curToken().IsTypeToken() {
			break
		}
		p.pos++
	}

	p.pos--

	if p.curToken().tokenType != word {
		return "", e
	}

	s := p.curToken().literal
	p.pos++
	return s, nil

}

func (p *Parser) parseVariableDef() Statement {

	// 変数定義か確認する
	if !p.isVariabeDef() {
		p.posReset()
		return nil
	}

	s, err := p.extractVarName()
	if err != nil {
		p.posReset()
		return nil
	}
	// semicolon or assign or lbracket
	switch p.curToken().tokenType {
	case semicolon:
		fallthrough
	case assign:
		fallthrough
	case lparen:
		fallthrough
	case lbracket:
		// semicolon まで進める
		p.progUntil(semicolon)
		p.pos++
		// next
	default:
		p.posReset()
		return nil
	}
	return &VariableDef{Name: s}
}

func (p *Parser) isVariabeDef() bool {
	// セミコロンもしくはイコールまでの間に型を表すトークンが2つ以上なければ変数定義ではない
	wordCnt := 0
	pPrev := p.pos
	t := p.curToken()
	for t.tokenType != semicolon && t.tokenType != assign && t.tokenType != eof {
		if t.IsTypeToken() {
			wordCnt++
		}
		p.pos++
		t = p.curToken()
	}
	p.pos = pPrev
	return wordCnt >= 2
}

func (p *Parser) parseVariableDecl() Statement {
	// extern
	p.pos++

	id, err := p.extractVarName()
	if err != nil {
		p.posReset()
		return nil
	}

    // セミコロンまでスキップ
	p.progUntil(semicolon)

	p.pos++
    // next

	return &VariableDecl{Name: id}
}

func (p *Parser) parsePrototypeDecl() Statement {

	if p.curToken().tokenType == keyExtern {
		p.pos++
	}

	for p.curToken().IsTypeToken() {
		p.pos++
	}

	if p.curToken().tokenType != lparen {
		// ( でなければプロトタイプ宣言ではない
		p.posReset()
		return nil
	}

	p.pos--

	if p.curToken().tokenType != word {
		p.posReset()
		return nil
	}

	// Name
	id := p.curToken().literal

	p.pos++

	// 関数の引数の括弧は飛ばす
	p.skipParen()

	if p.curToken().tokenType == keyAttribute {
		// attribute の場合はセミコロンまでスキップ
		p.progUntil(semicolon)
	} else if p.curToken().tokenType != semicolon {
		// セミコロン意外はプロトタイプ宣言ではない
		p.posReset()
		return nil
	}

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
		prevPos := p.pos
		switch p.curToken().tokenType {
		case lbrace:
			ts = append(ts, p.parseBlockStatement()...)
		case keyReturn:
			ts = p.parseReturn()
		case lparen:
			fallthrough
		case word:
			if ts == nil {
				if s := p.parseVariableDef(); s != nil {
					ts = append(ts, s)
				}
			}
			if ts == nil {
				p.pos = prevPos
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

func (p *Parser) parseReturn() []Statement {
	var ss []Statement = nil
	if p.curToken().tokenType == keyReturn {
		ss = []Statement{}
		p.pos++
		if p.curToken().tokenType != semicolon {
			// 何らかの式がある
			ts := p.parseExpression()
			ss = append(ss, ts...)
		}
		p.pos++
		// next
	}
	return ss
}

func (p *Parser) parseExpressionStatement() []Statement {
	ss := p.parseExpression()
	// semicolon
	p.pos++
	return ss
}

func (p *Parser) parseExpression() []Statement {
	var ss []Statement = nil

	switch p.curToken().tokenType {
	case lparen:
		prePos := p.pos
		ts := p.parseCast()
		if ts != nil {
			ss = append(ss, ts...)
		} else {
			p.pos = prePos
			p.pos++
			ts := p.parseExpression()
			ss = append(ss, ts...)
			// rparen
			p.pos++
		}
	case word:
		l := p.parseAccessVar()
		ss = append(ss, l)

	case letter:
		fallthrough
	case integer:
		p.pos++
		ss = []Statement{}
	default:
		return nil
	}

	if p.curToken().isOperator() {
		for p.curToken().isOperator() {
			p.pos++
		}
		r := p.parseExpression()
		ss = append(ss, r...)
	}

	return ss
}

func (p *Parser) parseCast() []Statement {
	if p.curToken().tokenType != lparen {
		p.posReset()
		return nil
	}
	p.pos++
	for p.curToken().tokenType != rparen {
		if p.curToken().tokenType != word {
			p.posReset()
			return nil
		}
		p.pos++
	}

	p.pos++
	ss := p.parseExpression()

	return ss
}

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
