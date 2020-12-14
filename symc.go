package symc

func ParseModule(src string) *Module {
	l := NewLexer(src)
	p := NewParser(l)
	return p.Parse()
}
