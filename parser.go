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
	Remain   []*Token
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
	return fmt.Sprintf("DEFINITION %s", v.Name)
}

type VariableDecl struct {
	Name string
}

func (v *VariableDecl) statementNode() {}
func (v *VariableDecl) String() string {
	return fmt.Sprintf("VariableDecl : Name=%s", v.Name)
}
func (v *VariableDecl) PrettyString() string {
	return fmt.Sprintf("DECLARE %s", v.Name)
}

type PrototypeDecl struct {
	Name string
}

func (v *PrototypeDecl) statementNode() {}
func (v *PrototypeDecl) String() string {
	return fmt.Sprintf("PrototypeDecl : Name=%s", v.Name)
}
func (v *PrototypeDecl) PrettyString() string {
	return fmt.Sprintf("PROTOTYPE %s", v.Name)
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
	txt := fmt.Sprintf("FUNC %s(", v.Name)

	sep := ""
	for _, p := range v.Params {
		txt += sep
		txt += p.PrettyString()
		sep = ", "
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
	leftVarInfo
}

// 代入先識別子情報
type leftVarInfo struct {
	idIndex int
	idName  string
}

// -----------------------------------------------------------
// 構文解析処理
// -----------------------------------------------------------

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
FOR:
	for p.curToken().tokenType != eof {
		if p.curToken().tokenType == comment {
			p.pos++
			continue
		}
		if ts := p.parseStatement(); ts != nil {
			ss = append(ss, ts...)
			for _, v := range ts {
				if _, yes := v.(*InvalidStatement); yes {
					break FOR
				}
			}
		}
		// エラーログを初期化
		p.errLog = ""
	}
	m := &Module{ss}
	return m
}

// parseStatement
func (p *Parser) parseStatement() []Statement {
	ss := []Statement{}
	switch p.curToken().tokenType {
	case keyTypedef:
		p.skipTypedef()
	case keyExtern:
		prePos := p.pos
		ss = p.parsePrototypeDecl()
		if ss == nil {
			p.pos = prePos
			ss = p.parseVariableDecl()
		}
		if ss == nil {
			return []Statement{&InvalidStatement{Contents: p.errLog, Tk: p.curToken(), Remain: p.tokens[p.pos:]}}
		}
	case keyUnion:
		if p.skipStructureLike() == nil {
			return nil
		}
	case keyStruct:
		if p.skipStructureLike() == nil {
			return nil
		}
	case keyAttribute:
		p.pos++
		p.skipParen()
	default:
		prePos := p.pos
		ss = p.parseFunctionDef()
		if ss == nil {
			p.pos = prePos
			ss = p.parsePrototypeDecl()
		}
		if ss == nil {
			p.pos = prePos
			ss = p.parseVariableDef()
		}
		if ss == nil {
			return []Statement{&InvalidStatement{Contents: p.errLog, Tk: p.curToken(), Remain: p.tokens[p.pos:]}}
		}
	}
	return ss
}

// parseVariableDef
func (p *Parser) parseVariableDef() []Statement {
	ss := []Statement{}

	prePos := p.pos
	ts := p.parseFuncPointerVarDef()
	if ts == nil {
		p.pos = prePos
		ts = p.parseNormalVarDef()
	}
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseVariableDef:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)
	return ss
}

// parseFuncPointerVarDef
func (p *Parser) parseFuncPointerVarDef() []Statement {

	ss := p.parseFuncPointerVarDefSub()
	if ss == nil {
		p.updateErrLog(fmt.Sprintf("parseFuncPointerVarDef:token[%s]", p.curToken().literal))
		return nil
	}
	if !p.curToken().isToken(semicolon) {
		p.updateErrLog(fmt.Sprintf("parseFuncPointerVarDef:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	return ss
}

// parseNormalVarDef
func (p *Parser) parseNormalVarDef() []Statement {
	ss := []Statement{}

	if !p.checkExists2word() {
		// 型名と変数名の合計が2以上なければ変数定義ではない
		p.updateErrLog(fmt.Sprintf("parseNormalVarDef:token[%s]", p.curToken().literal))
		return nil
	}

	for {

		if p.curToken().isToken(semicolon) {
			p.pos++
			break
		} else if p.curToken().isToken(comma) {
			p.pos++
		}

		for p.curToken().isTypeToken() {
			p.pos++
		}
		if !p.curToken().isToken(semicolon) &&
			!p.curToken().isToken(comma) &&
			!p.curToken().isToken(assign) &&
			!p.curToken().isToken(lbracket) {
			p.updateErrLog(fmt.Sprintf("parseNormalVarDef:token[%s]", p.curToken().literal))
			return nil
		}
		p.pos--
		id := p.curToken().literal
		p.pos++

		for p.curToken().isToken(lbracket) {
			// 配列
			p.pos++
			if !p.curToken().isToken(rbracket) {
				p.parseExpression()
				if !p.curToken().isToken(rbracket) {
					p.updateErrLog(fmt.Sprintf("parseNormalVarDef:token[%s]", p.curToken().literal))
					return nil
				}
			}
			// rbracket
			p.pos++
		}

		if p.curToken().isToken(assign) {
			// 初期化子あり
			p.pos++
			p.parseInitialValue()
		}

		ss = append(ss, &VariableDef{Name: id})
	}
	return ss
}

// checkExists2word
func (p *Parser) checkExists2word() bool {
	prePos := p.pos

	wordCnt := 0
	for p.curToken().isTypeToken() {
		if p.curToken().isToken(word) {
			wordCnt++
		}
		p.pos++
	}

	p.pos = prePos
	return wordCnt >= 2
}

// parseInitialValue
func (p *Parser) parseInitialValue() []Statement {

	if p.curToken().isToken(lbrace) {
		// 配列の初期化子
		xs := p.parseArrValue()
		if xs == nil {
			p.updateErrLog(fmt.Sprintf("parseVariableDef:token[%s]", p.curToken().literal))
			return nil
		}
	}

	xs := p.parseExpression()
	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parseVariableDef:token[%s]", p.curToken().literal))
		return nil
	}
	return []Statement{}
}

// parseArrValue
func (p *Parser) parseArrValue() []Statement {

	switch p.curToken().tokenType {
	case lbrace:
		p.pos++
		xs := p.parseArrValue()
		if xs == nil {
			p.updateErrLog(fmt.Sprintf("parseArrValue:token[%s]", p.curToken().literal))
			return nil
		}
		if !p.curToken().isToken(rbrace) {
			p.updateErrLog(fmt.Sprintf("parseArrValue:token[%s]", p.curToken().literal))
			return nil
		}
		p.pos++

		if p.curToken().isToken(comma) {
			p.pos++
			xs = p.parseArrValue()
			if xs == nil {
				p.updateErrLog(fmt.Sprintf("parseArrValue:token[%s]", p.curToken().literal))
				return nil
			}
		}
	default:
		for !p.curToken().isToken(rbrace) && !p.curToken().isToken(semicolon) && !p.curToken().isToken(eof) {
			p.pos++
		}
	}

	return []Statement{}
}

// parseVariableDefSub
func (p *Parser) parseVariableDefSub() []Statement {
	prePos := p.pos

	// はじめに関数ポインタか確認
	ss := p.parseFuncPointerVarDefSub()

	if ss == nil {
		// 関数ポインタ以外
		p.pos = prePos

		for p.curToken().isTypeToken() {
			p.pos++
		}
		p.pos--

		if p.curToken().tokenType != word {
			return nil
		}

		s := &VariableDef{Name: p.curToken().literal}

		p.pos++

		if p.curToken().tokenType == lbracket {
			// 配列の場合
			p.progUntil(rbracket)
			p.pos++
		}
		ss = append(ss, s)
	}

	return ss
}

// parseFuncPointerVarDefSub
func (p *Parser) parseFuncPointerVarDefSub() []Statement {
	for p.curToken().isTypeToken() {
		p.pos++
	}
	if p.curToken().tokenType != lparen {
		return nil
	}
	p.pos++
	ss := p.parseVariableDefSub()
	if ss == nil {
		return nil
	}
	if len(ss) != 1 {
		// 関数ポインタの識別子が一つだけあるはず
		return nil
	}
	// rparen
	p.pos++

	if p.curToken().tokenType != lparen {
		return nil
	}
	if xs := p.parsePrototypeParameter(); xs == nil {
		return nil
	}

	return ss
}

// parseVariableDecl
func (p *Parser) parseVariableDecl() []Statement {
	// extern
	p.pos++

	ss := p.parseVariableDef()
	if ss == nil {
		p.updateErrLog(fmt.Sprintf("parseVariableDecl:token[%s]", p.curToken().literal))
		return nil
	}

	// next

	if len(ss) != 1 {
		p.updateErrLog(fmt.Sprintf("parseVariableDecl:token[%s]", p.curToken().literal))
		return nil
	}
	defv, ok := ss[0].(*VariableDef)
	if !ok {
		p.updateErrLog(fmt.Sprintf("parseVariableDecl:token[%s]", p.curToken().literal))
		return nil
	}
	return []Statement{&VariableDecl{Name: defv.Name}}
}

// parsePrototypeDecl
func (p *Parser) parsePrototypeDecl() []Statement {

	if p.curToken().tokenType == keyExtern {
		p.pos++
	}

	xs := p.parsePrototypeDeclSub()

	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parsePrototypeDecl_3:token[%s]", p.curToken().literal))
		return nil
	}

	if p.curToken().tokenType == keyAttribute {
		// attribute の場合はセミコロンまでスキップ
		p.progUntil(semicolon)
	} else if p.curToken().isToken(keyAsm) {
		// __asm の場合はセミコロンまでスキップ
		p.progUntil(semicolon)
	} else if p.curToken().tokenType != semicolon {
		// セミコロン意外はプロトタイプ宣言ではない
		p.updateErrLog(fmt.Sprintf("parsePrototypeDecl_4:token[%s]", p.curToken().literal))
		return nil
	}

	p.pos++
	// next
	return xs

}

//parsePrototypeDeclSub
func (p *Parser) parsePrototypeDeclSub() []Statement {
	for p.curToken().isTypeToken() {
		p.pos++
	}

	if p.curToken().tokenType != lparen {
		// ( でなければプロトタイプ宣言ではない
		p.updateErrLog(fmt.Sprintf("parsePrototypeDeclSub:not lparen:token[%s]", p.curToken().literal))
		return nil
	}

	p.pos--

	if !p.curToken().isTypeToken() {
		p.updateErrLog(fmt.Sprintf("parsePrototypeDeclSub_2:token[%s]", p.curToken().literal))
		return nil
	}

	// Name
	// 仮の識別子名を取得
	id := p.curToken().literal
	p.pos++

	// 入れ子のプロトタイプ宣言のパターンを解析する
	// 以下の様な関数ポインタを返り値とする関数など
	// void(*signal(int, void (*)(int)))(int);
	// プロトタイプパラメータは lparen から解析開始するためこの位置で前回のポジションを記憶する
	prePos := p.pos
	xs := p.parsePrototypeParameter()
	if xs == nil {
		p.pos = prePos
		p.pos++
		xs = p.parsePrototypeDeclSub()
	}

	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parsePrototypeDeclSub_3:token[%s]", p.curToken().literal))
		return nil
	}

	// 入れ子のプロトタイプ宣言の場合識別名を更新
	if len(xs) == 1 {
		if v, ok := xs[0].(*PrototypeDecl); ok {
			id = v.Name
			if !p.curToken().isToken(rparen) {
				p.updateErrLog(fmt.Sprintf("parsePrototypeDeclSub_4:token[%s]", p.curToken().literal))
				return nil
			}
			p.pos++
			xs = p.parsePrototypeParameter()
			if xs == nil {
				p.updateErrLog(fmt.Sprintf("parsePrototypeDeclSub_5:token[%s]", p.curToken().literal))
				return nil
			}
		}
	}

	return []Statement{&PrototypeDecl{Name: id}}

}

// parsePrototypeParameter
// 構文解析のみ行い成功か失敗かを返すのみ
func (p *Parser) parsePrototypeParameter() []Statement {
	// lparen
	if !p.curToken().isToken(lparen) {
		return nil
	}
	p.pos++

	for {
		if p.curToken().isToken(rparen) {
			p.pos++
			return []Statement{}
		} else if p.curToken().isToken(comma) {
			p.pos++
		}

		prePos := p.pos
		xs := p.parseVariadicArgument()
		if xs == nil {
			p.pos = prePos
			xs = p.parsePrototypeParamVar()
		}
		if xs == nil {
			p.pos = prePos
			xs = p.parsePrototypeFPointerVar()
		}
		if xs == nil {
			p.updateErrLog(fmt.Sprintf("parsePrototypeParameter:token[%s]", p.curToken().literal))
			return nil
		}
	}
}

// parseVariadicArgument
func (p *Parser) parseVariadicArgument() []Statement {
	if !p.curToken().isToken(period) {
		p.updateErrLog(fmt.Sprintf("parseVariadicArgument:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	if !p.curToken().isToken(period) {
		p.updateErrLog(fmt.Sprintf("parseVariadicArgument:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	if !p.curToken().isToken(period) {
		p.updateErrLog(fmt.Sprintf("parseVariadicArgument:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	return []Statement{}
}

// parsePrototypeParamVar
func (p *Parser) parsePrototypeParamVar() []Statement {
	for p.curToken().isTypeToken() {
		p.pos++
	}

	if p.curToken().isToken(lbracket) {
		// 配列の場合
		p.pos++

		if p.curToken().isToken(rbracket) {
			// 空の配列
			p.pos++
		} else {
			xs := p.parseExpression()
			if xs == nil {
				p.updateErrLog(fmt.Sprintf("parsePrototypeParamVar_1:token[%s]", p.curToken().literal))
				return nil
			}
			if !p.curToken().isToken(rbracket) {
				p.updateErrLog(fmt.Sprintf("parsePrototypeParamVar_2:token[%s]", p.curToken().literal))
				return nil
			}
			p.pos++
		}
	}

	if p.curToken().isToken(comma) || p.curToken().isToken(rparen) {
		return []Statement{}
	} else if p.curToken().isToken(keyAttribute) {
		p.pos++
		p.skipParen()
		return []Statement{}
	} else {
		p.updateErrLog(fmt.Sprintf("parsePrototypeParamVar_3:token[%s]", p.curToken().literal))
		return nil
	}
}

// parsePrototypeFPointerVar
func (p *Parser) parsePrototypeFPointerVar() []Statement {
	for p.curToken().isTypeToken() {
		p.pos++
	}

	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_1:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	xs := p.parsePrototypeParamVar()
	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_2:token[%s]", p.curToken().literal))
		return nil
	}

	if !p.curToken().isToken(rparen) {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_3:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_4:token[%s]", p.curToken().literal))
		return nil
	}

	xs = p.parsePrototypeParameter()
	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_5:token[%s]", p.curToken().literal))
		return nil
	}

	if p.curToken().isToken(comma) || p.curToken().isToken(rparen) {
		return []Statement{}
	} else if p.curToken().isToken(keyAttribute) {
		return p.parseAttribute()
	} else {
		p.updateErrLog(fmt.Sprintf("parsePrototypeFPointerVar_6:token[%s]", p.curToken().literal))
		return nil
	}

}

// parseAttribute
func (p *Parser) parseAttribute() []Statement {
	if !p.curToken().isToken(keyAttribute) {
		p.updateErrLog(fmt.Sprintf("parseAttribute:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	p.skipParen()
	return []Statement{}
}

// parseFunctionDef
func (p *Parser) parseFunctionDef() []Statement {
	// lparen or eof の手前まで pos を進める
	for p.peekToken().isTypeToken() || p.peekToken().isToken(keyAttribute) {
		p.pos++
		if p.curToken().isToken(keyAttribute) {
			p.skipParen()
		}
	}
	if !p.peekToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parseFunctionDef:token[%s]", p.curToken().literal))
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
		return nil
	}

	ss := p.parseBlockStatement()
	if ss == nil {
		p.updateErrLog(fmt.Sprintf("parseFunctionDef:token[%s]", p.curToken().literal))
		return nil
	}

	return []Statement{&FunctionDef{Name: id, Params: ps, Statements: ss}}
}

// parseBlockStatement
func (p *Parser) parseBlockStatement() []Statement {
	ss := []Statement{}

	p.pos++

	for p.curToken().tokenType != rbrace {
		ts := p.parseInnerStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseBlockStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	}

	p.pos++

	return ss
}

// parseInnerStatement
func (p *Parser) parseInnerStatement() []Statement {
	ss := []Statement{}

	switch p.curToken().tokenType {
	case lbrace:
		ts := p.parseBlockStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyExtern:
		ts := p.parseVariableDecl()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyReturn:
		ts := p.parseReturn()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyIf:
		ts := p.parseIfStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyFor:
		ts := p.parseForStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyWhile:
		ts := p.parseWhileStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyDo:
		ts := p.parseDoWhileStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keySwitch:
		ts := p.parseSwitchStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyCase:
		ts := p.parseCaseStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyDefault:
		ts := p.parseDefaultStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	case keyBreak:
		p.pos++
		if !p.curToken().isToken(semicolon) {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		p.pos++
	default:
		prevPos := p.pos
		ts := p.parseVariableDef()
		if ts == nil {
			p.pos = prevPos
			// other statement
			ts = p.parseExpressionStatement()
		}
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseInnerStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
	}
	return ss
}

// parseDoWhileStatement
func (p *Parser) parseDoWhileStatement() []Statement {
	ss := []Statement{}
	// do
	p.pos++

	ts := p.parseBlockStatement()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}

	if !p.curToken().isToken(keyWhile) {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}

	p.pos++
	// lparen
	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}

	p.pos++

	us := p.parseExpression()
	if us == nil {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}

	if !p.curToken().isToken(rparen) {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}

	p.pos++
	if !p.curToken().isToken(semicolon) {
		p.updateErrLog(fmt.Sprintf("parseDoWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	ss = append(ss, ts...)
	ss = append(ss, us...)

	return ss
}

// parseCaseStatement
func (p *Parser) parseCaseStatement() []Statement {
	// case
	p.pos++

	xs := p.parseValue()
	if xs == nil {
		p.updateErrLog(fmt.Sprintf("parseCaseStatement:token[%s]", p.curToken().literal))
		return nil
	}

	if !p.curToken().isToken(colon) {
		p.updateErrLog(fmt.Sprintf("parseCaseStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	return []Statement{}
}

// parseValue
func (p *Parser) parseValue() []Statement {
	ss := []Statement{}
	switch p.curToken().tokenType {
	case float:
		fallthrough
	case letter:
		fallthrough
	case integer:
		p.pos++
		return ss
	default:
		return nil
	}
}

// parseDefaultStatement
func (p *Parser) parseDefaultStatement() []Statement {
	// default
	p.pos++

	if !p.curToken().isToken(colon) {
		p.updateErrLog(fmt.Sprintf("parseDefaultStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	return []Statement{}
}

// parseSwitchStatement
func (p *Parser) parseSwitchStatement() []Statement {
	ss := []Statement{}

	// switch
	p.pos++

	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parseSwitchStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	ts := p.parseExpression()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseSwitchStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	if !p.curToken().isToken(rparen) {
		p.updateErrLog(fmt.Sprintf("parseSwitchStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	if !p.curToken().isToken(lbrace) {
		p.updateErrLog(fmt.Sprintf("parseSwitchStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ts = p.parseBlockStatement()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseSwitchStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	return ss
}

// parseWhileStatement
func (p *Parser) parseWhileStatement() []Statement {
	ss := []Statement{}

	p.pos++

	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parseWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	ts := p.parseExpression()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	if !p.curToken().isToken(rparen) {
		p.updateErrLog(fmt.Sprintf("parseWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	ts = p.parseBlockStatement()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseWhileStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	return ss
}

// parseForStatement
func (p *Parser) parseForStatement() []Statement {
	ss := []Statement{}

	// for
	p.pos++

	if !p.curToken().isToken(lparen) {
		p.updateErrLog(fmt.Sprintf("parseForStatement:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++

	for {
		ts := p.parseExpression()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseForStatement:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)
		if p.curToken().isToken(semicolon) {
			p.pos++
		}
		if p.curToken().isToken(rparen) {
			p.pos++
			break
		}
	}

	ts := p.parseBlockStatement()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseForStatement:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	return ss
}

// parseIfStatement
func (p *Parser) parseIfStatement() []Statement {
	// if
	p.pos++
	// lparen
	p.pos++

	ss := []Statement{}

	// 条件式
	ts := p.parseExpression()
	if ts == nil {
		p.updateErrLog(fmt.Sprintf("parseIfStatement_1:token[%s]", p.curToken().literal))
		return nil
	}
	ss = append(ss, ts...)

	// rparen
	p.pos++

	if p.curToken().isToken(lbrace) {
		// ブロック文
		ts := p.parseBlockStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseIfStatement_2:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)

		if p.curToken().isToken(keyElse) {
			if p.peekToken().isToken(keyIf) {
				// else if 文
				p.pos++
				ts := p.parseIfStatement()
				if ts == nil {
					p.updateErrLog(fmt.Sprintf("parseIfStatement_3:token[%s]", p.curToken().literal))
					return nil
				}
				ss = append(ss, ts...)
			} else {
				// else 文
				p.pos++
				ts = p.parseBlockStatement()
				if ts == nil {
					p.updateErrLog(fmt.Sprintf("parseIfStatement_4:token[%s]", p.curToken().literal))
					return nil
				}
				ss = append(ss, ts...)
			}
		}
	} else {
		// １行命令
		ts := p.parseInnerStatement()
		if ts == nil {
			p.updateErrLog(fmt.Sprintf("parseIfStatement_5:token[%s]", p.curToken().literal))
			return nil
		}
		ss = append(ss, ts...)

		if p.curToken().isToken(keyElse) {
			// else 文あり
			p.pos++
			ts = p.parseInnerStatement()
			if ts == nil {
				p.updateErrLog(fmt.Sprintf("parseIfStatement_6:token[%s]", p.curToken().literal))
				return nil
			}
			ss = append(ss, ts...)
		}
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

// parseExpression
func (p *Parser) parseExpression() []Statement {
	ss := []Statement{}

	if p.curToken().isPrefixExpression() {
		// 前置式
		p.pos++
	}

	switch p.curToken().tokenType {
	case semicolon:
		// 空式
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
	case asterisk:
		p.pos++
		fallthrough
	case word:
		prePos := p.pos
		ls := p.parseCallFunc()
		if ls == nil {
			p.pos = prePos
			ls = p.parseRefVar()
			ss = append(ss, ls...)

			// RefVar の場合あとで assigne に変更する時のためにインデックスと変数名を記憶する
			if refv, ok := ss[len(ss)-1].(*RefVar); ok {
				p.leftVarInfo.idIndex = len(ss) - 1
				p.leftVarInfo.idName = refv.Name
			}

			for p.curToken().isToken(lbracket) {
				// 配列の場合
				p.pos++
				// leftVarInfo 上書き防止
				idIndex := p.leftVarInfo.idIndex
				idName := p.leftVarInfo.idName
				ts := p.parseExpression()
				p.leftVarInfo.idIndex = idIndex
				p.leftVarInfo.idName = idName
				p.pos++
				ss = append(ss, ts...)
			}

			if p.curToken().isPostExpression() {
				p.pos++
			}
		} else {
			ss = append(ss, ls...)
		}

		if ls == nil {
			p.updateErrLog(fmt.Sprintf("parseExpression:token[%s]", p.curToken().literal))
			return nil
		}

	case float:
		fallthrough
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
		if p.curToken().isToken(assign) || p.curToken().isCompoundOp() {
			// 代入式の場合は対象の識別子を Assigne 型に変更
			ss[p.leftVarInfo.idIndex] = &Assigne{p.leftVarInfo.idName}
		}
		p.pos++
		r := p.parseExpression()
		if r == nil {
			p.updateErrLog(fmt.Sprintf("parseExpression:token[%s]", p.curToken().literal))
			return nil
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

// parseCast
func (p *Parser) parseCast() []Statement {
	if p.curToken().tokenType != lparen {
		p.updateErrLog(fmt.Sprintf("parseCast:token[%s]", p.curToken().literal))
		return nil
	}
	p.pos++
	for p.curToken().tokenType != rparen {
		if !p.curToken().isToken(word) && !p.curToken().isToken(asterisk) {
			p.updateErrLog(fmt.Sprintf("parseCast:token[%s]", p.curToken().literal))
			return nil
		}
		p.pos++
	}

	p.pos++

	if p.curToken().isToken(semicolon) {
		p.updateErrLog(fmt.Sprintf("parseCast:token[%s]", p.curToken().literal))
		return nil
	}
	ss := p.parseExpression()

	return ss
}

// parseRefVar
func (p *Parser) parseRefVar() []Statement {
	if p.curToken().tokenType != word {
		p.updateErrLog(fmt.Sprintf("parseRefVar:token[%s]", p.curToken().literal))
		return nil
	}
	n := p.fetchID()

	return []Statement{&RefVar{Name: n}}
}

// parseParameter
func (p *Parser) parseParameter() []*VariableDef {
	ss := []*VariableDef{}
	// lparen
	p.pos++

	if p.curToken().tokenType == rparen {
		// パラメータになにもなし
		p.pos++
		// next
		return ss
	} else if p.curToken().tokenType == keyVoid {
		// void
		p.pos++
		// rparen
		p.pos++
		// next
		return ss
	}

	for {
		prePos := p.pos
		xs := p.parseVariadicArgument()
		if xs == nil {
			p.pos = prePos
			ts := p.parseVariableDefSub()
			if ts == nil {
				p.updateErrLog(fmt.Sprintf("parseParameter:token[%s]", p.curToken().literal))
				return nil
			}
			if len(ts) != 1 {
				p.updateErrLog(fmt.Sprintf("parseParameter:token[%s]", p.curToken().literal))
				return nil
			}
			v, ok := ts[0].(*VariableDef)
			if !ok {
				p.updateErrLog(fmt.Sprintf("parseParameter:token[%s]", p.curToken().literal))
				return nil
			}
			ss = append(ss, v)
		}

		switch p.curToken().tokenType {
		case rparen:
			p.pos++
			// next
			return ss
		case comma:
			p.pos++
			// next
		default:
			p.updateErrLog(fmt.Sprintf("parseParameter:token[%s]", p.curToken().literal))
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

// skipStructureLike
// 失敗した時は構文解析のパーサと同様 nil を返す
func (p *Parser) skipStructureLike() []Statement {

	for {
		if p.curToken().isToken(eof) ||
			p.curToken().isToken(semicolon) ||
			p.curToken().isToken(lbrace) {
			break
		}
		p.pos++
	}

	switch p.curToken().tokenType {
	case eof:
		return nil
	case semicolon:
		p.pos++
	case lbrace:
		if p.skipBrace() == nil {
			return nil
		}
		if !p.curToken().isToken(semicolon) {
			return nil
		}
		// semicolon
		p.pos++
	default:
		return nil
	}

	return []Statement{}
}

// skipBrace
func (p *Parser) skipBrace() []Statement {
	if !p.curToken().isToken(lbrace) {
		return nil
	}
	p.pos++

	for {
		if p.curToken().isToken(rbrace) {
			p.pos++
			return []Statement{}
		} else if p.curToken().isToken(lbrace) {
			xs := p.skipBrace()
			if xs == nil {
				return nil
			}
		} else if p.curToken().isToken(eof) {
			return nil
		}
		p.pos++
	}
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
