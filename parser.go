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
	lexer *Lexer
}

func NewParser(l *Lexer) *Parser {
	return &Parser{lexer: l}
}

func (p *Parser) Parse() *Module {
	return &Module{Statements: []Statement{&VariableDef{Name: "hoge"}}}
}
