package symc

import "fmt"

func ParseModule(src string) *Module {
	l := NewLexer(src)
	p := NewParser(l)
	ast := p.Parse()
	fmt.Println(ast)
	return &Module{[]Statement{&VariableDef{"hoge"}}}
}
