package symc

import "fmt"

func ParseModule(src string) *Module {
	l := NewLexer(src)
	res := l.lexicalize()
	fmt.Println(res)
	return &Module{[]Statement{&VariableDef{"hoge"}}}
}
