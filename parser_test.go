package symc

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"test1",
			`
int hoge;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"test2",
			`
extern char fuga;`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "fuga"},
				},
			},
		},
		{
			"test3",
			`
const int hoge;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"test4",
			`
extern const int hoge;
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "hoge"},
				},
			},
		},
		{
			"testx",
			`
int hoge;
char fuga;
extern long piyo;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
					&VariableDef{Name: "fuga"},
					&VariableDecl{Name: "piyo"},
				},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("got=%v, expect=%v", got, tt.expect)
		}
	}
}
