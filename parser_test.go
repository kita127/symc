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
const int *hoge;
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
extern const int *hoge;
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "hoge"},
				},
			},
		},
		{
			"test5",
			`
int hoge = 100;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"test6",
			`
int hoge[] = {0x00, 0x01, 0x02};
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"test7",
			`
int hoge[3] = {0x00, 0x01, 0x02};
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"test8",
			`
char hoge[] = "hello";
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"variable definition 9",
			`
int hoge = 0;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},

		{
			"prototype dec 1",
			`
void func_a( void );
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "func_a"},
				},
			},
		},
		{
			"prototype dec 2",
			`
extern void func_a( void );
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "func_a"},
				},
			},
		},
		{
			"typedef 1",
			`
typedef unsigned char __uint8_t;
`,
			&Module{
				[]Statement{
					&Typedef{Name: "__uint8_t"},
				},
			},
		},
		{
			"typedef 2",
			`
typedef union {
 char __mbstate8[128];
 long long _mbstateL;
} __mbstate_t;
`,
			&Module{
				[]Statement{
					&Typedef{Name: "__mbstate_t"},
				},
			},
		},
		{
			"struct 1",
			`
struct __darwin_pthread_handler_rec {
 void (*__routine)(void *);
 void *__arg;
 struct __darwin_pthread_handler_rec *__next;
};
`,
			&Module{
				[]Statement{
				},
			},
		},
		{
			"function def 1",
			`
void func_a( void ) {}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func_a", Block: &BlockStatement{Statements: []Statement{}}},
				},
			},
		},
		{
			"function def 2",
			`
int func(int a)
{
    int hoge = 0;
    hoge++;
    a = a + (10);
    return a;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func", Block: &BlockStatement{Statements: []Statement{
						&VariableDef{Name: "hoge"},
						&RefVar{Name: "hoge"},
						&AssignVar{Name: "a"},
						&RefVar{Name: "a"},
						&RefVar{Name: "a"},
					},
					},
					},
				},
			},
		},
		{
			"testx",
			`
# 1 "hoge.c"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 366 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "hoge.c" 2

int func(int a)
{
    int hoge = 0;
    hoge++;
    a = a + (10);
    return a;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func", Block: &BlockStatement{Statements: []Statement{
						&VariableDef{Name: "hoge"},
						&RefVar{Name: "hoge"},
						&AssignVar{Name: "a"},
						&RefVar{Name: "a"},
						&RefVar{Name: "a"},
					},
					},
					},
				},
			},
		},
		//		{
		//			"test err1",
		//			`int hoge`,
		//			&Module{
		//				[]Statement{
		//					&InvalidStatement{Contents: "parse, err parse function def, err parse prototype decl, err parse variable def"},
		//				},
		//			},
		//		},
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
