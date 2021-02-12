package symc

// parse モジュール
// 字句解析の結果を構文解析する
// 構文解析後、次のトークンに必ず位置を合わせること

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
func (v *Module) PrettyString() string {
	txt := ""

	for _, s := range v.Statements {
		txt += s.PrettyString()
		txt += "\n"
	}
	return txt
}

type PrettyStringer interface {
	PrettyString() string
}

type Statement interface {
	statementNode()
	fmt.Stringer
	PrettyStringer
}

type InvalidStatement struct {
	Contents string
	Tk       *Token
}

func (v *InvalidStatement) statementNode() {}
func (v *InvalidStatement) String() string {
	return fmt.Sprintf("InvalidStatement : %s", v.Contents)
}
func (v *InvalidStatement) PrettyString() string {
	return v.String()
}

type VariableDef struct {
	Name string
}

func (v *VariableDef) statementNode() {}
func (v *VariableDef) String() string {
	return fmt.Sprintf("VariableDef : Name=%s", v.Name)
}
func (v *VariableDef) PrettyString() string {
	return fmt.Sprintf("VAR %s", v.Name)
}

type VariableDecl struct {
	Name string
}

func (v *VariableDecl) statementNode() {}
func (v *VariableDecl) String() string {
	return fmt.Sprintf("VariableDecl : Name=%s", v.Name)
}
func (v *VariableDecl) PrettyString() string {
	return fmt.Sprintf("dec %s", v.Name)
}

type PrototypeDecl struct {
	Name string
}

func (v *PrototypeDecl) statementNode() {}
func (v *PrototypeDecl) String() string {
	return fmt.Sprintf("PrototypeDecl : Name=%s", v.Name)
}
func (v *PrototypeDecl) PrettyString() string {
	return fmt.Sprintf("prototype %s", v.Name)
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
func (v *FunctionDef) PrettyString() string {
	txt := fmt.Sprintf("func %s (", v.Name)

	sep := ""
	for _, p := range v.Params {
		txt += sep
		txt += p.PrettyString()
		sep = " ,"
	}

	txt += ") {\n"

	for _, t := range v.Statements {
		txt += "    "
		txt += t.PrettyString()
		txt += "\n"
	}
	txt += "}\n"

	return txt
}

type RefVar struct {
	Name string
}

func (v *RefVar) statementNode() {}
func (v *RefVar) String() string {
	return fmt.Sprintf("RefVar : Name=%s", v.Name)
}
func (v *RefVar) PrettyString() string {
	return fmt.Sprintf("%s", v.Name)
}

type Assigne struct {
	Name string
}

func (v *Assigne) statementNode() {}
func (v *Assigne) String() string {
	return fmt.Sprintf("Assigne : Name=%s", v.Name)
}
func (v *Assigne) PrettyString() string {
	return fmt.Sprintf("ASSIGNE %s", v.Name)
}

type CallFunc struct {
	Name string
	Args []Statement
}

func (v *CallFunc) statementNode() {}
func (v *CallFunc) String() string {
	return fmt.Sprintf("CallFunc : Name=%s, Args=%v", v.Name, v.Args)
}
func (v *CallFunc) PrettyString() string {
	txt := fmt.Sprintf("%s(", v.Name)
	sep := ""
	for _, a := range v.Args {
		txt += sep
		txt += fmt.Sprintf("%s", a.PrettyString())
		sep = ", "
	}
	txt += ")"
	return txt
}

type Typedef struct {
	Name string
}

func (v *Typedef) statementNode() {}
func (v *Typedef) String() string {
	return fmt.Sprintf("Typedef : Name=%s", v.Name)
}
func (v *Typedef) PrettyString() string {
	return fmt.Sprintf("%s", v.Name)
}

// 構文解析器
type Parser struct {
	lexer   *Lexer
	tokens  []*Token
	pos     int
	prevPos int
	errLog  string
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
		// エラーログを初期化
		p.errLog = ""
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
			return &InvalidStatement{Contents: p.errLog, Tk: p.curToken()}
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
			return &InvalidStatement{Contents: p.errLog, Tk: p.curToken()}
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
		} else if !p.curToken().isTypeToken() {
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

	s := p.parseVariableDefSub()
	if s == nil {
		p.posReset()
		return nil
	}

	switch p.curToken().tokenType {
	case semicolon:
		p.pos++
	case assign:
		p.progUntil(semicolon)
		p.pos++
	default:
		p.posReset()
		return nil
	}

	return s
}

func (p *Parser) parseVariableDefSub() Statement {
	prePos := p.pos

	// はじめに関数ポインタか確認
	s := p.parseFuncPointerVarDef()

	if s == nil {
		// 関数ポインタ以外
		p.pos = prePos

		for p.curToken().isTypeToken() {
			p.pos++
		}
		p.pos--

		if p.curToken().tokenType != word {
			return nil
		}

		s = &VariableDef{Name: p.curToken().literal}

		p.pos++

		if p.curToken().tokenType == lbracket {
			// 配列の場合
			p.progUntil(rbracket)
			p.pos++
		}
	}

	return s
}

func (p *Parser) parseFuncPointerVarDef() Statement {
	for p.curToken().isTypeToken() {
		p.pos++
	}
	if p.curToken().tokenType != lparen {
		return nil
	}
	p.pos++
	s := p.parseVariableDefSub()
	if s == nil {
		return nil
	}
	// rparen
	p.pos++

	if p.curToken().tokenType != lparen {
		return nil
	}
	if x := p.parsePrototypeParameter(); x == nil {
		return nil
	}

	return s
}

// parsePrototypeParameter
// プロトタイプ宣言、関数ポインタの定義で使用するパラメータのパース
// 識別子がないパターンのパラメータ定義に対応する
// パラメータの識別子は返さず構文チェックのみ行う
func (p *Parser) parsePrototypeParameter() []Statement {
	p.pos++

FOR:
	for {
		t := p.curToken()
		switch t.tokenType {
		case keyVoid:
		case asterisk:
		case comma:
		case rbracket:
		case lbracket:
		case word:
		case rparen:
			break FOR
		default:
			return nil
		}
		p.pos++
	}

	p.pos++

	if !p.curToken().isToken(semicolon) {
		return nil
	}

	return []Statement{}
}

func (p *Parser) isVariabeDef() bool {
	// セミコロンもしくはイコールまでの間に型を表すトークンが2つ以上なければ変数定義ではない
	wordCnt := 0
	pPrev := p.pos
	t := p.curToken()
	for t.tokenType != semicolon && t.tokenType != assign && t.tokenType != eof {
		if t.isTypeToken() {
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

	for p.curToken().isTypeToken() {
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
	} else if p.curToken().isToken(keyAsm) {
		// __asm の場合はセミコロンまでスキップ
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

// parseFunctionDef
func (p *Parser) parseFunctionDef() Statement {
	// lparen or eof の手前まで pos を進める
	n := p.peekToken()
	for n.isTypeToken() {
		p.pos++
		if p.curToken().tokenType == keyAttribute {
			p.skipParen()
		}
		n = p.peekToken()
	}
	if n.tokenType != lparen {
		p.updateErrLog(fmt.Sprintf("parseFunctionDef:token[%s]", p.curToken().literal))
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
		p.updateErrLog(fmt.Sprintf("parseFunctionDef:token[%s]", p.curToken().literal))
		p.posReset()
		return nil
	}

	ss := p.parseBlockStatement()
	if ss == nil {
		p.updateErrLog(fmt.Sprintf("parseFunctionDef:token[%s]", p.curToken().literal))
		p.posReset()
		return nil
	}

	return &FunctionDef{Name: id, Params: ps, Statements: ss}
}

// parseBlockStatement
func (p *Parser) parseBlockStatement() []Statement {
	p.pos++

	ss := []Statement{}

	for p.curToken().tokenType != rbrace {
		var ts []Statement = nil
		prevPos := p.pos
		switch p.curToken().tokenType {
		case lbrace:
			ts = append(ts, p.parseBlockStatement()...)
		case keyReturn:
			ts = p.parseReturn()
		case keyIf:
			ts = p.parseIfStatement()
		case keyFor:
			ts = p.parseForStatement()
		default:
			if ts == nil {
				p.pos = prevPos
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

	p.pos++

	return ss
}

// parseAssigne
func (p *Parser) parseAssigne() []Statement {
	ss := []Statement{}

	if !p.curToken().isToken(word) {
		return nil
	}
	id := p.fetchID()

	if p.curToken().isPreAssigneOperator() {
		// 代入の前に置くことができる演算子 e.g. +=
		p.pos++
	}

	if !p.curToken().isToken(assign) {
		return nil
	}
	p.pos++

	exps := p.parseExpression()
	if exps == nil {
		return nil
	}

	ss = append(ss, &Assigne{id})
	ss = append(ss, exps...)

	return ss

}

func (p *Parser) parseForStatement() []Statement {
	ss := []Statement{}

	// for
	p.pos++

	if !p.curToken().isToken(lparen) {
		return nil
	}
	p.pos++

	for {
		ts := p.parseExpression()
		if ts == nil {
			return nil
		}
		ss = append(ss, ts...)
		if p.curToken().isToken(rparen) {
			p.pos++
			break
		} else if p.curToken().isToken(semicolon) {
			p.pos++
		} else {
			p.updateErrLog(fmt.Sprintf("parseForStatement:token[%s]", p.curToken().literal))
			return nil
		}
	}

	p.progUntil(rbrace)
	p.pos++
	return ss
}

func (p *Parser) parseIfStatement() []Statement {
	// if
	p.pos++
	// lparen
	p.pos++

	ss := []Statement{}

	// 条件式
	ts := p.parseExpression()
	if ts == nil {
		return nil
	}
	ss = append(ss, ts...)

	// rparen
	p.pos++

	if p.curToken().isToken(lbrace) {
		// ブロック文
		ts := p.parseBlockStatement()
		if ts == nil {
			return nil
		}
		ss = append(ss, ts...)
	} else {
		// １行命令
		ts := p.parseExpressionStatement()
		if ts == nil {
			return nil
		}
		ss = append(ss, ts...)
	}

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
	if ss == nil {
		p.updateErrLog(fmt.Sprintf("parseExpressionStatement:token[%s]", p.curToken().literal))
		return nil
	}
	// semicolon
	p.pos++
	return ss
}

func (p *Parser) parseExpression() []Statement {
	var varName string
	ss := []Statement{}

	if p.curToken().isPrefixExpression() {
		// 前置式
		p.pos++
	}

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
	case ampersand:
		p.pos++
		ls := p.parseExpression()
		if ls == nil {
			p.updateErrLog(fmt.Sprintf("parseExpression:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ls...)
	case word:
		prePos := p.pos
		ls := p.parseCallFunc()
		if ls == nil {
			p.pos = prePos
			ls = p.parseRefVar()
			if p.curToken().isInfixExpression() {
				p.pos++
			}
		}
		ss = append(ss, ls...)
	case str:
		fallthrough
	case letter:
		fallthrough
	case integer:
		p.pos++
	default:
		p.updateErrLog(fmt.Sprintf("parseExpression:token[%s]", p.curToken().literal))
		return nil
	}

	if p.curToken().isOperator() {
		op := p.curToken()
		p.pos++
		r := p.parseExpression()
		if r == nil {
			p.updateErrLog(fmt.Sprintf("parseExpression:token[%s]", p.curToken().literal))
			return nil
		}
		if op.isToken(assign) || op.isCompoundOp() {
			// 代入は最後の Statement を Assigne に変更する

			// 代入式の場合、直前の値の型は *RefVar であることが保証されている
			// Assigne に変更するため変数名を取得しておく
			refv, _ := ss[len(ss)-1].(*RefVar)
			varName = refv.Name
			ss[len(ss)-1] = &Assigne{varName}
		}
		ss = append(ss, r...)
	}

	return ss
}

func (p *Parser) parseCallFunc() []Statement {
	id := p.curToken().literal
	p.pos++
	if p.curToken().tokenType != lparen {
		return nil
	}
	p.pos++

	ss := []Statement{}

	if p.curToken().tokenType == rparen {
		// 引数なし
		p.pos++
		return []Statement{&CallFunc{Name: id, Args: ss}}
	}

	for {
		ts := p.parseExpression()
		if ts == nil {
			return nil
		}
		ss = append(ss, ts...)

		if p.curToken().tokenType == comma {
			p.pos++
		} else if p.curToken().tokenType == rparen {
			p.pos++
			break
		} else {
			return nil
		}
	}

	return []Statement{&CallFunc{Name: id, Args: ss}}
}

func (p *Parser) parseCast() []Statement {
	if p.curToken().tokenType != lparen {
		return nil
	}
	p.pos++
	for p.curToken().tokenType != rparen {
		if p.curToken().tokenType != word {
			return nil
		}
		p.pos++
	}

	p.pos++
	ss := p.parseExpression()

	return ss
}

func (p *Parser) parseRefVar() []Statement {
	if p.curToken().tokenType != word {
		p.posReset()
		return nil
	}
	n := p.fetchID()

	return []Statement{&RefVar{Name: n}}
}

func (p *Parser) parseParameter() []*VariableDef {
	vs := []*VariableDef{}
	// lparen
	p.pos++

	if p.curToken().tokenType == rparen {
		// パラメータになにもなし
		p.pos++
		// next
		return vs
	} else if p.curToken().tokenType == keyVoid {
		// void
		p.pos++
		// rparen
		p.pos++
		// next
		return vs
	}

	for {
		s := p.parseVariableDefSub()
		if s == nil {
			return nil
		}
		v, ok := s.(*VariableDef)
		if !ok {
			return nil
		}
		vs = append(vs, v)

		switch p.curToken().tokenType {
		case rparen:
			p.pos++
			// next
			return vs
		case comma:
			p.pos++
			// next
		default:
			return nil
		}
	}
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

func (p *Parser) updateErrLog(msg string) {
	delimiter := ";"
	p.errLog += msg
	p.errLog += delimiter
}
