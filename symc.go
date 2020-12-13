package symc

import "fmt"

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
	l := NewLexer(src)
	res := l.lexicalize(src)
	fmt.Println(res)
	return &Module{[]Statement{&VariableDef{"hoge"}}}
}
