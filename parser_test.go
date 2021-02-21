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
							// %
							&Assigne{Name: "var10"},
							&RefVar{Name: "a10"},
							&RefVar{Name: "b10"},
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
