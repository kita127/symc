package symc

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
	p.pos++
	id := p.tokens[p.pos].literal
	p.pos++ // semicolon
	p.pos++ // next
	m := &Module{[]Statement{&VariableDef{Name: id}}}
	return m
}
