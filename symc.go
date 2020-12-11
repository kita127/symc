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

func ParseModule(src string) *Module {
	//lexicalize(src)
	return &Module{[]Statement{&VariableDef{"hoge"}}}
}
