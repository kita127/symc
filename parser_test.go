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
			"struct 1",
			`
struct __darwin_pthread_handler_rec {
 void (*__routine)(void *);
 void *__arg;
 struct __darwin_pthread_handler_rec *__next;
};
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"struct 2",
			`
struct __sFILEX;
# 126 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/_stdio.h" 3 4
typedef struct __sFILE {
 unsigned char *_p;
 int _r;
 int _w;
 short _flags;
 short _file;
 struct __sbuf _bf;
 int _lbfsize;


 void *_cookie;
 int (* _Nullable _close)(void *);
 int (* _Nullable _read) (void *, char *, int);
 fpos_t (* _Nullable _seek) (void *, fpos_t, int);
 int (* _Nullable _write)(void *, const char *, int);


 struct __sbuf _ub;
 struct __sFILEX *_extra;
 int _ur;


 unsigned char _ubuf[3];
 unsigned char _nbuf[1];


 struct __sbuf _lb;


 int _blksize;
 fpos_t _offset;
} FILE;
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"attribute 1",
			`
__attribute__()
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"attribute 2",
			`
__attribute__((()))
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"attribute 3",
			`
__attribute__((__availability__(swift, unavailable, message="Use mkstemp(3) instead.")))
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"attribute 4",
			`
__attribute__ ((__always_inline__))
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"attribute 5",
			`
__attribute__ ((__always_inline__)) int hoge;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"function parameter 1",
			`
void f_hoge(int a){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "f_hoge",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function parameter 2",
			`
void f_fuga(int a, char b){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "f_fuga",
						Params:     []*VariableDef{{Name: "a"}, {Name: "b"}},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function parameter 3",
			`
void f_piyo(int a, char b, AnyType c[]){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "f_piyo",
						Params:     []*VariableDef{{Name: "a"}, {Name: "b"}, {Name: "c"}},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function parameter 4",
			`
void f_ice(int a, char b, AnyType c[100]){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "f_ice",
						Params:     []*VariableDef{{Name: "a"}, {Name: "b"}, {Name: "c"}},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function parameter 5",
			`
void func(void){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function parameter 6",
			`
void func(int a, ...){}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{
							{Name: "a"},
						},
						Statements: []Statement{}},
				},
			},
		},
		{
			"character 1",
			`
void func(void)
{
    char c = 'A';
    c = 'B';
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "c"},
							&Assigne{Name: "c"},
						},
					},
				},
			},
		},
		{
			"access struct var 1",
			`
void func(void)
{
    hoge.a1 = 'A';
    fuga->b1 = 100;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "hoge.a1"},
							&Assigne{Name: "fuga->b1"},
						},
					},
				},
			},
		},
		{
			"return 1",
			`
void f_xxx(void)
{
    return;
}
void f_yyy(void)
{
    return 10;
}
void f_zzz(void)
{
    return a;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "f_xxx",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
					&FunctionDef{Name: "f_yyy",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
					&FunctionDef{Name: "f_zzz",
						Params:     []*VariableDef{},
						Statements: []Statement{&RefVar{Name: "a"}},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestParseVardef
func TestParseVardef(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"variable definition 1",
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
			"variable definition 2",
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
			"variable definition 3",
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
			"variable definition 4",
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
			"variable definition 5",
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
			"variable definition 6",
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
			"variable definition 7",
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
			"variable definition 8",
			`
int mx, my;
char ma, mb, mc;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "mx"},
					&VariableDef{Name: "my"},
					&VariableDef{Name: "ma"},
					&VariableDef{Name: "mb"},
					&VariableDef{Name: "mc"},
				},
			},
		},
		{
			"variable definition 9",
			`
int len = 0, len_buf;
int var_a = 0, var_b = 100;
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "len"},
					&VariableDef{Name: "len_buf"},
					&VariableDef{Name: "var_a"},
					&VariableDef{Name: "var_b"},
				},
			},
		},
		{
			"variable definition 10",
			`
BOOTINFO *binfo = (BOOTINFO *)(0x00000ff0);
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "binfo"},
				},
			},
		},
		{
			"variable definition 11",
			`
int hoge[3][3] = { {0x0a, 0x0b, 0x0c}, {0x00, 0x01, 0x02}, {0x00, 0x01, 0x09} };
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"variable definition 12",
			`
int hoge[2][3][4] = {
                      {
                        {0x00 ,0x00 ,0x00, 0x00},
                        {0x00 ,0x00 ,0x00, 0x00},
                        {0x00 ,0x00 ,0x00, 0x00}
                      },
                      {
                        {0x00 ,0x00 ,0x00, 0x00},
                        {0x00 ,0x00 ,0x00, 0x00},
                        {0x00 ,0x00 ,0x00, 0x00}
                      },
                    };
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
				},
			},
		},
		{
			"variable definition 13",
			`
char func(void) {
    Buffer *b = make_buffer();
    int len = vec_len(inits);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "b"},
							&CallFunc{Name: "make_buffer", Args: []Statement{}},
							&VariableDef{Name: "len"},
							&CallFunc{Name: "vec_len",
								Args: []Statement{
									&RefVar{Name: "inits"},
								},
							},
						},
					},
				},
			},
		},
		{
			"variable definition 14",
			`
int hoge = sizeof(int);
int fuga = sizeof(arr)/sizeof(&arr[0]);
Node **buf = malloc(len * sizeof(Node *));
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "hoge"},
					&VariableDef{Name: "fuga"},
					&VariableDef{Name: "buf"},
					&CallFunc{
						Name: "malloc",
						Args: []Statement{
							&RefVar{Name: "len"},
						},
					},
				},
			},
		},
		{
			"variable definition 15",
			`
void **v = calloc(newsize, sizeof(void *));
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "v"},
					&CallFunc{
						Name: "calloc",
						Args: []Statement{
							&RefVar{Name: "newsize"},
						},
					},
				},
			},
		},
		{
			"function pointer def 1",
			`
void (* p_f)();
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "p_f"},
				},
			},
		},
		{
			"function pointer def 3",
			`
int (* p_f)(void);
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "p_f"},
				},
			},
		},
		{
			"function pointer def 4",
			`
const * AnyType (* p_f)(int a, char b[]);
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "p_f"},
				},
			},
		},
		{
			"function pointer def 5",
			`
int (*_Nullable _read)(void *, char *, int);
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "_read"},
				},
			},
		},
		{
			"function pointer def 6",
			`
int (* fp)(void *, char *, int, ...);
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "fp"},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestApp
func TestApp(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"app 1",
			`
# 1 "hoge.c"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 366 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "hoge.c" 2
# 1 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/stdio.h" 1 3 4
# 64 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/stdio.h" 3 4
# 1 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/_stdio.h" 1 3 4
# 68 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/_stdio.h" 3 4
# 1 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/sys/cdefs.h" 1 3 4
# 647 "/Library/Developer/CommandLineTools/SDKs/MacOSX.sdk/usr/include/sys/cdefs.h" 3 4

typedef struct {
  char xxx;
} Gt;

typedef struct {
  int aaa;
  Gt bbb;
} St;

int muruchi_piyomi(char *s) {
  St purin;
  St *p;
  purin.aaa = 100;
  purin.bbb.xxx = 'A';
  p = &purin;
  p->aaa = purin.aaa + 200;
  return (0);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "muruchi_piyomi",
						Params: []*VariableDef{{Name: "s"}},
						Statements: []Statement{
							&VariableDef{Name: "purin"},
							&VariableDef{Name: "p"},
							&Assigne{Name: "purin.aaa"},
							&Assigne{Name: "purin.bbb.xxx"},
							&Assigne{Name: "p"},
							&RefVar{Name: "purin"},
							&Assigne{Name: "p->aaa"},
							&RefVar{Name: "purin.aaa"},
						}},
				},
			},
		},
		{
			"app 2",
			`
void hoge(void){
 if (--_p->_w >= 0 || (_p->_w >= _p->_lbfsize && (char)_c != '\n'))
  return (*_p->_p++ = _c);
 else
  return (__swbuf(_c, _p));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "hoge",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "_p->_w"},
							&RefVar{Name: "_p->_w"},
							&RefVar{Name: "_p->_lbfsize"},
							&RefVar{Name: "_c"},
							&Assigne{Name: "_p->_p"},
							&RefVar{Name: "_c"},
							&CallFunc{
								Name: "__swbuf",
								Args: []Statement{
									&RefVar{Name: "_c"},
									&RefVar{Name: "_p"},
								},
							},
						}},
				},
			},
		},
		{
			"app 3",
			`
FILE *funopen(const void *,
                 int (* _Nullable)(void *, char *, int),
                 int (* _Nullable)(void *, const char *, int),
                 fpos_t (* _Nullable)(void *, fpos_t, int),
                 int (* _Nullable)(void *));
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "funopen"},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestStatements
func TestStatements(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"break 1",
			`
void func(void)
{
    break;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestWhile
func TestWhile(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"while 1",
			`
void func(void)
{
    while(1){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"while 2",
			`
void func(void)
{
    while(condition1){
        var = 100;
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "condition1"},
							&Assigne{Name: "var"},
						},
					},
				},
			},
		},
		{
			"while 3",
			`
void func(void)
{
    while (*p)
        print(b, *p++);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "p"},
							&CallFunc{Name: "print", Args: []Statement{
								&RefVar{Name: "b"},
								&RefVar{Name: "p"},
							},
							},
						},
					},
				},
			},
		},
		{
			"while 4",
			`
void func(void)
{
    while (*p)
        while(x)
            print(b, *p++);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "p"},
							&RefVar{Name: "x"},
							&CallFunc{Name: "print", Args: []Statement{
								&RefVar{Name: "b"},
								&RefVar{Name: "p"},
							},
							},
						},
					},
				},
			},
		},
		{
			"while 5",
			`
void func(void)
{
    while (*p)
        while(x)
            ;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "p"},
							&RefVar{Name: "x"},
						},
					},
				},
			},
		},
		{
			"do while 1",
			`
void func(void)
{
    do {
        var = 100;
    } while(0);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "var"},
						},
					},
				},
			},
		},
		{
			"do while 2",
			`
void func(void)
{
    do {
        var = 100;
    } while(a = read(buf));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "var"},
							&Assigne{Name: "a"},
							&CallFunc{Name: "read", Args: []Statement{&RefVar{Name: "buf"}}},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestSwitch
func TestSwitch(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"switch 1",
			`
void func(void)
{
    switch (c) {
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "c"},
						},
					},
				},
			},
		},
		{
			"switch 2",
			`
void func(void)
{
    switch (c) {
        default:
            break;
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "c"},
						},
					},
				},
			},
		},
		{
			"switch 3",
			`
void func(void)
{
    switch (c) {
        case 0x01:
            var1++;
            break;
        case 'a':
        case 0.01:
        default:
            break;
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "c"},
							&RefVar{Name: "var1"},
						},
					},
				},
			},
		},
		{
			"switch 4",
			`
void func(void)
{
    switch (macro->kind) {
    case MACRO_OBJ:
        Set *hideset = set_add(tok->hideset, name);
        Vector *tokens = subst(macro, ((void *)0), hideset);
        return read_expand();
    case MACRO_FUNC:
        Vector *args = read_args(tok, macro);
        expect(')');
    case MACRO_SPECIAL:
        macro->fn(tok);
        return read_expand();
    default:
        errorf("cpp.c" ":" "360", ((void *)0), "internal error");
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "macro->kind"},
							&VariableDef{Name: "hideset"},
							&CallFunc{
								Name: "set_add",
								Args: []Statement{
									&RefVar{Name: "tok->hideset"},
									&RefVar{Name: "name"},
								},
							},
							&VariableDef{Name: "tokens"},
							&CallFunc{
								Name: "subst",
								Args: []Statement{
									&RefVar{Name: "macro"},
									&RefVar{Name: "hideset"},
								},
							},
							&CallFunc{
								Name: "read_expand",
								Args: []Statement{},
							},
							&VariableDef{Name: "args"},
							&CallFunc{
								Name: "read_args",
								Args: []Statement{
									&RefVar{Name: "tok"},
									&RefVar{Name: "macro"},
								},
							},
							&CallFunc{
								Name: "expect",
								Args: []Statement{},
							},
							&CallFunc{
								Name: "macro->fn",
								Args: []Statement{
									&RefVar{Name: "tok"},
								},
							},
							&CallFunc{
								Name: "read_expand",
								Args: []Statement{},
							},
							&CallFunc{
								Name: "errorf",
								Args: []Statement{},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestAssigne
func TestAssigne(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"assigne 1",
			`
void func(void)
{
    hoge = 1;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{"hoge"},
						},
					},
				},
			},
		},
		{
			"assigne 2",
			`
void func(void)
{
    hoge = a;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{"hoge"},
							&RefVar{"a"},
						},
					},
				},
			},
		},
		{
			"assigne 3",
			`
void func(void)
{
    arrVar[i] = 0xAA;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{"arrVar"},
							&RefVar{"i"},
						},
					},
				},
			},
		},
		{
			"assigne 4",
			`
void func(void)
{
    arrVar2[i][j] = 0xAA;
    arrVar3[i][j][k] = 0xAA;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{"arrVar2"},
							&RefVar{"i"},
							&RefVar{"j"},
							&Assigne{"arrVar3"},
							&RefVar{"i"},
							&RefVar{"j"},
							&RefVar{"k"},
						},
					},
				},
			},
		},
		{
			"assigne 5",
			`
void func(void)
{
    *_p->_p++ = _c;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{"_p->_p"},
							&RefVar{"_c"},
						},
					},
				},
			},
		},
		{
			"assigne 6",
			`
char global_arr[(char)1000];

void func(void)
{
    char local_arr[(99)+(1)];
}
`,
			&Module{
				[]Statement{
					&VariableDef{Name: "global_arr"},
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "local_arr"},
						},
					},
				},
			},
		},
		{
			"assigne 7",
			`
void func(void)
{
    v = *p;
    v = ***p;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "v"},
							&RefVar{Name: "p"},
							&Assigne{Name: "v"},
							&RefVar{Name: "p"},
						},
					},
				},
			},
		},
		{
			"assigne 8",
			`
void func(void)
{
    p = &v;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "p"},
							&RefVar{Name: "v"},
						},
					},
				},
			},
		},
		{
			"assigne 9",
			`
void func(void)
{
    k[j] = m->key[i];
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "k"},
							&RefVar{Name: "j"},
							&RefVar{Name: "m->key"},
							&RefVar{Name: "i"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestRef
func TestRef(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestVarDecl
func TestVarDecl(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"variable decl 1",
			`
extern char fuga;`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "fuga"},
				},
			},
		},
		{
			"variable decl 2",
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
			"variable decl 3",
			`
extern int (* p_f)(void);
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "p_f"},
				},
			},
		},
		{
			"variable decl 4",
			`
extern struct StType st_var;
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "st_var"},
				},
			},
		},
		{
			"variable decl 5",
			`
extern int optind, opterr, optopt;
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "optind"},
					&VariableDecl{Name: "opterr"},
					&VariableDecl{Name: "optopt"},
				},
			},
		},
		{
			"local variable decl 1",
			`
void func(void)
{
    extern char hankaku;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDecl{Name: "hankaku"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestTypedef
func TestTypedef(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"typedef 1",
			`
typedef unsigned char __uint8_t;
`,
			&Module{
				[]Statement{},
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
				[]Statement{},
			},
		},
		{
			"typedef 3",
			`
typedef union {
    t v;
    struct {
        t2 x;
        t2 y;
    }

} HOGE;
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"typedef 4",
			`
typedef struct {
    int kind;
    File *file;
    int line;
    int column;
    _Bool space;
    _Bool bol;
    int count;
    Set *hideset;
    union {

        int id;

        struct {
            char *sval;
            int slen;
            int c;
            int enc;
        };

        struct {
            _Bool is_vararg;
            int position;
        };
    };
} Token;
`,
			&Module{
				[]Statement{},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestUnion
func TestUnion(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"union 1",
			`
union unionType {

 int int_var;
 void *void_ptr;
};
`,
			&Module{
				[]Statement{},
			},
		},
		{
			"union 2",
			`
union wait {
 int w_status;
 struct {
  unsigned int w_Termsig:7,
      w_Coredump:1,
      w_Retcode:8,
      w_Filler:16;
 } w_T;

 struct {
  unsigned int w_Stopval:8,
      w_Stopsig:8,
      w_Filler:16;
 } w_S;
};
`,
			&Module{
				[]Statement{},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestCast
func TestCast(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"cast 1",
			`
void func(int a)
{
    hoge = (char)10;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}},
					},
				},
			},
		},
		{
			"cast 2",
			`
void func(int a)
{
    hoge = (char)fuga;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}, &RefVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"cast 3",
			`
void func(int a)
{
    hoge = (char)(fuga);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}, &RefVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"cast 4",
			`
void func(int a)
{
    hoge = (char)(fuga + piyo);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}, &RefVar{Name: "fuga"}, &RefVar{Name: "piyo"}},
					},
				},
			},
		},
		{
			"cast 5",
			`
void func(int a)
{
    hoge = (unsigned char)10;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}},
					},
				},
			},
		},
		{
			"cast 6",
			`
void func(void)
{
    hoge = (char *)"strings";
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{&Assigne{Name: "hoge"}},
					},
				},
			},
		},
		{
			"cast 6",
			`
void func(void)
{
    (const void *)xxx;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{&RefVar{Name: "xxx"}},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestFunctionDef
func TestFunctionDef(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"function def 1",
			`
void func_name( void ) {}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func_name",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"function def 2",
			`
void func(int a)
{
    {
        {}
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"function def 4",
			`
void func(int a)
{
    int hoge;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&VariableDef{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 5",
			`
void func(int a)
{
    {
        int hoge;
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&VariableDef{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 6",
			`
void func(int a)
{
    hoge;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&RefVar{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 7",
			`
void func(int a)
{
    hoge = 10;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 8",
			`
void func(int a)
{
    hoge = fuga;
    p_var = &address;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{{Name: "a"}},
						Statements: []Statement{
							&Assigne{Name: "hoge"},
							&RefVar{Name: "fuga"},
							&Assigne{Name: "p_var"},
							&RefVar{Name: "address"},
						},
					},
				},
			},
		},
		{
			"function def 9",
			`
void func(int a)
{
    (hoge) = (fuga);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&Assigne{Name: "hoge"}, &RefVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"function def 10",
			`
void func(int a)
{
    unsigned char aaa;
    a = 10;
    b = 0x00;
    c = 10U;
    d = 100UL;
    e = 0.001;
    f = 0.002f;
    g = 0.01F;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{{Name: "a"}},
						Statements: []Statement{
							&VariableDef{Name: "aaa"},
							&Assigne{Name: "a"},
							&Assigne{Name: "b"},
							&Assigne{Name: "c"},
							&Assigne{Name: "d"},
							&Assigne{Name: "e"},
							&Assigne{Name: "f"},
							&Assigne{Name: "g"}},
					},
				},
			},
		},
		{
			"function def 11",
			`
inline __attribute__ ((__always_inline__)) int __sputc(int _c, FILE *_p) {
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "__sputc",
						Params:     []*VariableDef{{Name: "_c"}, {Name: "_p"}},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function def 12",
			`
void func()
{
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{}},
				},
			},
		},
		{
			"function def 13",
			`
void func()
{
    (hoge);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{"hoge"},
						}},
				},
			},
		},
		{
			"function def 14",
			`
static void pop_function(void *ignore) {
    if (dumpstack)
        vec_pop(functions);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "pop_function",
						Params: []*VariableDef{
							{Name: "ignore"},
						},
						Statements: []Statement{
							&RefVar{"dumpstack"},
							&CallFunc{
								Name: "vec_pop",
								Args: []Statement{
									&RefVar{Name: "functions"},
								},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestPrototypeDecl
func TestPrototypeDecl(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
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
			"prototype dec 3",
			`
int renameat(int, const char *, int, const char *) __attribute__((availability(macosx,introduced=10.10)));
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "renameat"},
				},
			},
		},
		{
			"prototype dec 4",
			`
extern int __vsnprintf_chk (char * restrict, size_t, int, size_t,
       const char * restrict, va_list);
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "__vsnprintf_chk"},
				},
			},
		},
		{
			"prototype dec 5",
			`
int  _read(void *, char *, int);
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "_read"},
				},
			},
		},
		{
			"prototype dec 6",
			`
void(*signal(int, void (*)(int)))(int);

`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "signal"},
				},
			},
		},
		{
			"prototype dec 7",
			`
int yakitori(void (* _Nonnull)(void));

`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "yakitori"},
				},
			},
		},
		{
			"prototype dec 8",
			`
int nanosleep(const struct timespec *__rqtp, struct timespec *__rmtp) __asm("_" "nanosleep" );
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "nanosleep"},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestPragma
func TestPragma(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"pragma 1",
			`
#pragma hoge xxxxxxxxxx
`,
			&Module{
				[]Statement{},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestGcc
func TestGcc(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"test gcc 1",
			`
FILE *fopen(const char * restrict __filename, const char * restrict __mode) __asm("_" "fopen" );
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "fopen"},
				},
			},
		},
		{
			"test gcc 2",
			`
inline __attribute__ ((__always_inline__)) int __sputc(int _c, FILE *_p) {
}
`,
			&Module{
				[]Statement{
					&FunctionDef{
						Name: "__sputc",
						Params: []*VariableDef{
							{Name: "_c"},
							{Name: "_p"},
						},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"test gcc 3",
			`
FILE *fdopen(int, const char *) __asm("_" "fdopen" );
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "fdopen"},
				},
			},
		},
		{
			"test gcc 4",
			`
int fprintf(FILE * restrict, const char * restrict, ...);
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "fprintf"},
				},
			},
		},
		{
			"test gcc 5",
			`
int oden(void (^ _Nonnull)(void)) __attribute__((availability(macosx,introduced=11.2)));
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "oden"},
				},
			},
		},
		{
			"test gcc 6",
			`
int heapsort_b(void *__base, size_t __nel, size_t __width __attribute__((__noescape__)));
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "heapsort_b"},
				},
			},
		},
		{
			"test gcc 7",
			`
int heapsort_b(void *__base, size_t __nel, size_t __width,
     int (^ _Nonnull __compar)(const void *, const void *) __attribute__((__noescape__)))
     __attribute__((availability(macosx,introduced=31.7)));
`,
			&Module{
				[]Statement{
					&PrototypeDecl{Name: "heapsort_b"},
				},
			},
		},
		{
			"test gcc 8",
			`
extern long timezone __asm("_" "timezone" );
`,
			&Module{
				[]Statement{
					&VariableDecl{Name: "timezone"},
				},
			},
		},
		{
			"test gcc 9",
			`
static void push_xmm(int reg) {
    int save_hook __attribute__((unused, cleanup(pop_function)));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "push_xmm",
						Params: []*VariableDef{
							{Name: "reg"},
						},
						Statements: []Statement{
							&VariableDef{Name: "save_hook"},
						},
					},
				},
			},
		},
		{
			"test gcc 10",
			`
static void push_xmm(int reg) {
    int save_hook __attribute__((unused, cleanup(pop_function))); if (dumpstack) vec_push(functions, (void *)__func__);;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "push_xmm",
						Params: []*VariableDef{
							{Name: "reg"},
						},
						Statements: []Statement{
							&VariableDef{Name: "save_hook"},
							&RefVar{Name: "dumpstack"},
							&CallFunc{Name: "vec_push",
								Args: []Statement{
									&RefVar{Name: "functions"},
									&RefVar{Name: "__func__"},
								},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestIf
func TestIf(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"if 1",
			`
void func(void)
{
    if (1){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"if 2",
			`
void func(void)
{
    if (hoge == 0){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "hoge"},
						},
					},
				},
			},
		},
		{
			"if 3",
			`
void func(void)
{
    if ( hoge == 0){
        fuga = a + (1 - 2);
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "hoge"},
							&Assigne{Name: "fuga"},
							&RefVar{Name: "a"},
						},
					},
				},
			},
		},
		{
			"if 4",
			`
void func(void)
{
    if (1)
        hoge;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "hoge"},
						},
					},
				},
			},
		},
		{
			"if 5",
			`
void func(void)
{
 if (flag)
  return 1;
 else
  return 0;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "flag"},
						},
					},
				},
			},
		},
		{
			"if 6",
			`
void func(void)
{
   if (a >= 0 || (b >= c && (char)d != '\n')){}
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "a"},
							&RefVar{Name: "b"},
							&RefVar{Name: "c"},
							&RefVar{Name: "d"},
						},
					},
				},
			},
		},
		{
			"if 7",
			`
void func(void)
{
   if (condition1){

   }else if (condition2){

   }else{

   }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "condition1"},
							&RefVar{Name: "condition2"},
						},
					},
				},
			},
		},
		{
			"if 8",
			`
void func(void)
{
   if (condition1)
        if (condition2)
            return;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "condition1"},
							&RefVar{Name: "condition2"},
						},
					},
				},
			},
		},
		{
			"if 9",
			`
void func(void)
{
   if (condition1)
        if (condition2)
            ;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "condition1"},
							&RefVar{Name: "condition2"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestCallExpression
func TestCallExpression(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"call expression 1",
			`
void func()
{
    hoge();
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{
								Name: "hoge",
								Args: []Statement{},
							},
						}},
				},
			},
		},
		{
			"call expression 2",
			`
void func()
{
    hoge(100);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "hoge",
								Args: []Statement{},
							},
						}},
				},
			},
		},
		{
			"call expression 3",
			`
void func()
{
    hoge(a);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "hoge",
								Args: []Statement{&RefVar{Name: "a"}},
							},
						}},
				},
			},
		},
		{
			"call expression 4",
			`
void func()
{
    hoge(a, b, 100);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "hoge",
								Args: []Statement{
									&RefVar{Name: "a"},
									&RefVar{Name: "b"},
								},
							},
						}},
				},
			},
		},
		{
			"call expression 5",
			`
void func()
{
    hoge(a, fuga(piyo(1, b), 2));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "hoge",
								Args: []Statement{
									&RefVar{Name: "a"},
									&CallFunc{Name: "fuga",
										Args: []Statement{
											&CallFunc{Name: "piyo",
												Args: []Statement{&RefVar{Name: "b"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"call expression 6",
			`
void func()
{
    _read(0, "A", 1);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "_read",
								Args: []Statement{},
							},
						},
					},
				},
			},
		},
		{
			"call expression 7",
			`
void func()
{
    __darwin_check_fd_set(_fd, (const void *) _p);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "__darwin_check_fd_set",
								Args: []Statement{&RefVar{Name: "_fd"}, &RefVar{Name: "_p"}}},
						},
					},
				},
			},
		},
		{
			"call expression 8",
			`
void func()
{
    fprintf(exitcode ? __stderrp : __stdoutp);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "fprintf",
								Args: []Statement{
									&RefVar{Name: "exitcode"},
									&RefVar{Name: "__stderrp"},
									&RefVar{Name: "__stdoutp"},
								},
							},
						},
					},
				},
			},
		},
		{
			"call expression 9",
			`
void func()
{
    fprintf( "Usage: hoge [ -E ][ -a ] [ -h ] <file>\n\n"
            "\n"
            "  -I<path>          kds;la ksksks kkasfnf\n"
            "  -E                source code\n"
            "  -S                x,ssi kskgng owoiw ,c,dd(default)\n"
            "  -c                Do not run linker (default)\n"
            "  -U name           shdsh kskalal\n"
            "  -o filename       shdha lllal ,c,c,c, ieiri\n"
            "  -g                Do nothing at this moment\n"
            "  -O<number>        Does nothing at this moment\n"
            "  -w                Disable all warnings\n"
            "  -h                print this help\n"
            "\n"
            "One of -a, -c, -E or -S must be specified.\n\n");
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "fprintf",
								Args: []Statement{},
							},
						},
					},
				},
			},
		},
		{
			"call expression 10",
			`
void func()
{
    macro->fn(tok);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{Name: "macro->fn",
								Args: []Statement{
									&RefVar{Name: "tok"},
								},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestExpression
func TestExpression(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"expression 1",
			`
void func()
{
    var1 = a1 +  b1;
    var2 = a2 -  b2;
    var3 = a3 *  b3;
    var4 = a4 /  b4;
    var5 = a5 |  b5;
    var6 = a6 &  b6;
    var7 = a7 || b7;
    var8 = a8 && b8;
    var9 = a9 > b9;
    var10 = a10 & b10;
    var11 = a11 % b11;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// +
							&Assigne{Name: "var1"},
							&RefVar{Name: "a1"},
							&RefVar{Name: "b1"},
							// -
							&Assigne{Name: "var2"},
							&RefVar{Name: "a2"},
							&RefVar{Name: "b2"},
							// *
							&Assigne{Name: "var3"},
							&RefVar{Name: "a3"},
							&RefVar{Name: "b3"},
							// /
							&Assigne{Name: "var4"},
							&RefVar{Name: "a4"},
							&RefVar{Name: "b4"},
							// |
							&Assigne{Name: "var5"},
							&RefVar{Name: "a5"},
							&RefVar{Name: "b5"},
							// &
							&Assigne{Name: "var6"},
							&RefVar{Name: "a6"},
							&RefVar{Name: "b6"},
							// ||
							&Assigne{Name: "var7"},
							&RefVar{Name: "a7"},
							&RefVar{Name: "b7"},
							// &&
							&Assigne{Name: "var8"},
							&RefVar{Name: "a8"},
							&RefVar{Name: "b8"},
							// >
							&Assigne{Name: "var9"},
							&RefVar{Name: "a9"},
							&RefVar{Name: "b9"},
							// &
							&Assigne{Name: "var10"},
							&RefVar{Name: "a10"},
							&RefVar{Name: "b10"},
							// %
							&Assigne{Name: "var11"},
							&RefVar{Name: "a11"},
							&RefVar{Name: "b11"},
						}},
				},
			},
		},
		{
			"expression 2",
			`
void func()
{
    var1 = a1 <<  b1;
    var2 = a2 >>  b2;
    var3 = a3 ^ b3;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// <<
							&Assigne{Name: "var1"},
							&RefVar{Name: "a1"},
							&RefVar{Name: "b1"},
							// >>
							&Assigne{Name: "var2"},
							&RefVar{Name: "a2"},
							&RefVar{Name: "b2"},
							// ^
							&Assigne{Name: "var3"},
							&RefVar{Name: "a3"},
							&RefVar{Name: "b3"},
						}},
				},
			},
		},
		{
			"expression 3",
			`
void func()
{
    var1 += a1;
    var2 -= a2;
    var3 /= a3;
    var4 *= a4;
    var5 |= a5;
    var6 &= a6;
    var7 <<= a7;
    var8 >>= a8;
    var9 ~= a9;
    var10 ^= a10;
    var11 %= a11;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// +=
							&Assigne{Name: "var1"},
							&RefVar{Name: "a1"},
							// -=
							&Assigne{Name: "var2"},
							&RefVar{Name: "a2"},
							// /=
							&Assigne{Name: "var3"},
							&RefVar{Name: "a3"},
							// *=
							&Assigne{Name: "var4"},
							&RefVar{Name: "a4"},
							// |=
							&Assigne{Name: "var5"},
							&RefVar{Name: "a5"},
							// &=
							&Assigne{Name: "var6"},
							&RefVar{Name: "a6"},
							// <<=
							&Assigne{Name: "var7"},
							&RefVar{Name: "a7"},
							// >>=
							&Assigne{Name: "var8"},
							&RefVar{Name: "a8"},
							// ~=
							&Assigne{Name: "var9"},
							&RefVar{Name: "a9"},
							// ^=
							&Assigne{Name: "var10"},
							&RefVar{Name: "a10"},
							// %=
							&Assigne{Name: "var11"},
							&RefVar{Name: "a11"},
						}},
				},
			},
		},
		{
			"expression 4",
			`
void func()
{
    var1 = a1 + b1 + c1;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// +=
							&Assigne{Name: "var1"},
							&RefVar{Name: "a1"},
							&RefVar{Name: "b1"},
							&RefVar{Name: "c1"},
						}},
				},
			},
		},
		{
			"expression 5",
			`
void func()
{
    var1 = (a1);
    var2 = (a2 + b2);
    var3 = ((a3 + b3) + c3);
    var4 = (a4 + (b4 + c4));
    var5 = ((a5 + b5) + (c5 + d5));
    var6 = (((a6 + b6) - c6) + (d6 + e6));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// 1
							&Assigne{Name: "var1"},
							&RefVar{Name: "a1"},
							// 2
							&Assigne{Name: "var2"},
							&RefVar{Name: "a2"},
							&RefVar{Name: "b2"},
							// 3
							&Assigne{Name: "var3"},
							&RefVar{Name: "a3"},
							&RefVar{Name: "b3"},
							&RefVar{Name: "c3"},
							// 4
							&Assigne{Name: "var4"},
							&RefVar{Name: "a4"},
							&RefVar{Name: "b4"},
							&RefVar{Name: "c4"},
							// 5
							&Assigne{Name: "var5"},
							&RefVar{Name: "a5"},
							&RefVar{Name: "b5"},
							&RefVar{Name: "c5"},
							&RefVar{Name: "d5"},
							// 6
							&Assigne{Name: "var6"},
							&RefVar{Name: "a6"},
							&RefVar{Name: "b6"},
							&RefVar{Name: "c6"},
							&RefVar{Name: "d6"},
							&RefVar{Name: "e6"},
						}},
				},
			},
		},
		{
			"expression 6",
			`
void func()
{
    --a;
    ++b;
    c--;
    d++;
    e++ + --f;
    ++g - h--;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "a"},
							&RefVar{Name: "b"},
							&RefVar{Name: "c"},
							&RefVar{Name: "d"},
							&RefVar{Name: "e"},
							&RefVar{Name: "f"},
							&RefVar{Name: "g"},
							&RefVar{Name: "h"},
						},
					},
				},
			},
		},
		{
			"expression 7",
			`
void func()
{
    &hoge;
    (unsigned char)fuga;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "hoge"},
							&RefVar{Name: "fuga"},
						},
					},
				},
			},
		},
		{
			"expression 8",
			`
void func()
{
    ;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"expression 9",
			`
void func()
{
    a == b;
    c != d;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "a"},
							&RefVar{Name: "b"},
							&RefVar{Name: "c"},
							&RefVar{Name: "d"},
						},
					},
				},
			},
		},
		{
			"expression 10",
			`
void func()
{
    ((unsigned long)_fd % (sizeof(__int32_t) * 8));
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "_fd"},
						},
					},
				},
			},
		},
		{
			"expression 11",
			`
void func()
{
    a = ~a;
    b = !b;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "a"},
							&RefVar{Name: "a"},
							&Assigne{Name: "b"},
							&RefVar{Name: "b"},
						},
					},
				},
			},
		},
		{
			"expression 12",
			`
void func()
{
    &((Vector){});
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"expression 13",
			`
void func()
{
    a = -1;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "a"},
						},
					},
				},
			},
		},
		{
			"expression 14",
			`
void func()
{
    ++++++a;
    b--------;
    --++c++--;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "a"},
							&RefVar{Name: "b"},
							&RefVar{Name: "c"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestEnum
func TestEnum(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"enum 1",
			`
enum {
    TIME,
    PLACE,
    NUMBER,
    FUGA,
    INU,
    NEKO,
    HIYOKO,

    PAN_DA,
    CHIHUAHUA,
    TORI,
    TMACRO_PARAM,
};
`,
			&Module{
				[]Statement{},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		p := NewParser(l)
		got := p.Parse()
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestTernaryOp
func TestTernaryOp(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"ternary op 1",
			`
void func()
{
    hoge = exitcode ? __stderrp : __stdoutp;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "hoge"},
							&RefVar{Name: "exitcode"},
							&RefVar{Name: "__stderrp"},
							&RefVar{Name: "__stdoutp"},
						},
					},
				},
			},
		},
		{
			"ternary op 2",
			`
void func()
{
    char *dir = file->name ? dirname(strdup(file->name)) : ".";
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "dir"},
							&RefVar{Name: "file->name"},
							&CallFunc{Name: "dirname", Args: []Statement{
								&CallFunc{
									Name: "strdup",
									Args: []Statement{
										&RefVar{Name: "file->name"},
									},
								},
							},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestContinue
func TestContinue(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"continue 1",
			`
void func()
{
    continue;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestForStatement
func TestForStatement(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"for 1",
			`
void func(void)
{
    for (i = 0; i < 10; i++){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "i"},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
						},
					},
				},
			},
		},
		{
			"for 2",
			`
void func(void)
{
    for (;;){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{},
						Statements: []Statement{},
					},
				},
			},
		},
		{
			"for 3",
			`
void func(void)
{
    for (i = 0; i < 10; i++){
        arrVar[i] = i * x;
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "i"},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
							&Assigne{Name: "arrVar"},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
							&RefVar{Name: "x"},
						},
					},
				},
			},
		},
		{
			"for 4",
			`
void func(void)
{
    for (i = 0; i < len; i++)
        buf_write(b, s[i]);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "i"},
							&RefVar{Name: "i"},
							&RefVar{Name: "len"},
							&RefVar{Name: "i"},
							&CallFunc{Name: "buf_write", Args: []Statement{
								&RefVar{Name: "b"},
								&RefVar{Name: "s"},
								&RefVar{Name: "i"},
							},
							},
						},
					},
				},
			},
		},
		{
			"for 5",
			`
void func(void)
{
    for (int i = vec_len(tokens) - 1; i >= 0; i--){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "i"},
							&CallFunc{Name: "vec_len", Args: []Statement{
								&RefVar{Name: "tokens"},
							},
							},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
						},
					},
				},
			},
		},
		{
			"for 6",
			`
void func(void)
{
    for (int i = vec_len(tokens) - 1; i >= 0; i--)
        for (int j = 0; j < i; j++)
            a += j;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "i"},
							&CallFunc{
								Name: "vec_len",
								Args: []Statement{
									&RefVar{Name: "tokens"},
								},
							},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
							&VariableDef{Name: "j"},
							&RefVar{Name: "j"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
							&Assigne{Name: "a"},
							&RefVar{Name: "j"},
						},
					},
				},
			},
		},
		{
			"for 7",
			`
void func(void)
{
    for (int i = vec_len(tokens) - 1; i >= 0; i--)
        for (int j = 0; j < i; j++)
            ;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&VariableDef{Name: "i"},
							&CallFunc{
								Name: "vec_len",
								Args: []Statement{
									&RefVar{Name: "tokens"},
								},
							},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
							&VariableDef{Name: "j"},
							&RefVar{Name: "j"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestGoto
func TestGoto(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"goto 1",
			`
void func(void)
{
    goto err;

    err:
        errorf("cpp.c" ":" "711", token_pos(hash), "cannot find header file: %s", filename);
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&CallFunc{
								Name: "errorf",
								Args: []Statement{
									&CallFunc{
										Name: "token_pos",
										Args: []Statement{
											&RefVar{Name: "hash"},
										},
									},
									&RefVar{Name: "filename"},
								},
							},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestComment
func TestComment(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"comment 1",
			`
void func(void)
{
    switch (tok->kind) {
        case HOGE:
            return x;
        case FUGA:
            switch (tok->id) {

# 1 "./piyopiyo.inc" 1

case OP_ARROW: return "->";
case OP_A_ADD: return "+=";
case OP_A_AND: return "&=";
            }
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "tok->kind"},
							&RefVar{Name: "x"},
							&RefVar{Name: "tok->id"},
						},
					},
				},
			},
		},
		{
			"comment 2",
			`
void func(void)
{
    for
# 1 "./piyopiyo.inc" 1
    (i = 0; i < 10; i++){
    }
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "i"},
							&RefVar{Name: "i"},
							&RefVar{Name: "i"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestStVar
func TestStVar(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"st var 1",
			`
void func(void)
{
    s.v;
    p->v;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "s"},
							&RefVar{Name: "p"},
						},
					},
				},
			},
		},
		{
			"st var 2",
			`
void func(void)
{
    s.v = 100;
    p->v = 100;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "s"},
							&Assigne{Name: "p"},
						},
					},
				},
			},
		},
		{
			"st var 3",
			`
void func(void)
{
    s.v.v2.v3;
    p->v->v2->v3;
    s.v.v2.v3 = 100;
    p->v->v2->v3 = 100;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "s"},
							&RefVar{Name: "p"},
							&Assigne{Name: "s"},
							&Assigne{Name: "p"},
						},
					},
				},
			},
		},
		{
			"st var 4",
			`
void func(void)
{
    (*p).v;
    (&s)->v;
    (*p).v = 100;
    (&s)->v = 100;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "p"},
							&RefVar{Name: "s"},
							&Assigne{Name: "p"},
							&Assigne{Name: "s"},
						},
					},
				},
			},
		},
		{
			"st var 5",
			`
void func(void)
{
    s.x[i];
    s.y[i][j];
    s.z[i][j][k];
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&RefVar{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
							&RefVar{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
							&RefVar{Name: "k"},
						},
					},
				},
			},
		},
		{
			"st var 6",
			`
void func(void)
{
    h[i] = 0;
    s.x[i] = a;
    s.y[i][j] = b;
    s.z[i][j][k] = c;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							&Assigne{Name: "h"},
							&RefVar{Name: "i"},
							&Assigne{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "a"},
							&Assigne{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
							&RefVar{Name: "b"},
							&Assigne{Name: "s"},
							&RefVar{Name: "i"},
							&RefVar{Name: "j"},
							&RefVar{Name: "k"},
							&RefVar{Name: "c"},
						},
					},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}

// TestTmp
func TestTmp(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"tmp 1",
			`
int hoge[] = {100, 200};
`,
			&Module{
				[]Statement{
					&VariableDef{"hoge"},
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
			t.Errorf("\ngot=   %v\nexpect=%v\n", got, tt.expect)
		}
	}
}
