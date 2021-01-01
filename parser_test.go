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
						Statements: []Statement{&AccessVar{Name: "hoge"}},
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
						Statements: []Statement{&AccessVar{Name: "hoge"}},
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
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}},
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
						Statements: []Statement{&AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}},
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
    hoge = 10;
    fuga = 20;
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params:     []*VariableDef{{Name: "a"}},
						Statements: []Statement{&VariableDef{Name: "aaa"}, &AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"function def 11",
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
						Statements: []Statement{&AccessVar{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 12",
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
						Statements: []Statement{&AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"function def 13",
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
						Statements: []Statement{&AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}},
					},
				},
			},
		},
		{
			"function def 14",
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
						Statements: []Statement{&AccessVar{Name: "hoge"}, &AccessVar{Name: "fuga"}, &AccessVar{Name: "piyo"}},
					},
				},
			},
		},
		{
			"function def 15",
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
						Statements: []Statement{&AccessVar{Name: "hoge"}},
					},
				},
			},
		},
		{
			"function def 16",
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
			"function def 17",
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
			"function def 18",
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
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// +
							&AccessVar{Name: "var1"},
							&AccessVar{Name: "a1"},
							&AccessVar{Name: "b1"},
							// -
							&AccessVar{Name: "var2"},
							&AccessVar{Name: "a2"},
							&AccessVar{Name: "b2"},
							// *
							&AccessVar{Name: "var3"},
							&AccessVar{Name: "a3"},
							&AccessVar{Name: "b3"},
							// /
							&AccessVar{Name: "var4"},
							&AccessVar{Name: "a4"},
							&AccessVar{Name: "b4"},
							// |
							&AccessVar{Name: "var5"},
							&AccessVar{Name: "a5"},
							&AccessVar{Name: "b5"},
							// &
							&AccessVar{Name: "var6"},
							&AccessVar{Name: "a6"},
							&AccessVar{Name: "b6"},
							// ||
							&AccessVar{Name: "var7"},
							&AccessVar{Name: "a7"},
							&AccessVar{Name: "b7"},
							// &&
							&AccessVar{Name: "var8"},
							&AccessVar{Name: "a8"},
							&AccessVar{Name: "b8"},
						}},
				},
			},
		},
		{
			"function def 19",
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
							&AccessVar{Name: "var1"},
							&AccessVar{Name: "a1"},
							&AccessVar{Name: "b1"},
							// >>
							&AccessVar{Name: "var2"},
							&AccessVar{Name: "a2"},
							&AccessVar{Name: "b2"},
							// ^
							&AccessVar{Name: "var3"},
							&AccessVar{Name: "a3"},
							&AccessVar{Name: "b3"},
						}},
				},
			},
		},
		{
			"function def 20",
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
}
`,
			&Module{
				[]Statement{
					&FunctionDef{Name: "func",
						Params: []*VariableDef{},
						Statements: []Statement{
							// +=
							&AccessVar{Name: "var1"},
							&AccessVar{Name: "a1"},
							// -=
							&AccessVar{Name: "var2"},
							&AccessVar{Name: "a2"},
							// /=
							&AccessVar{Name: "var3"},
							&AccessVar{Name: "a3"},
							// *=
							&AccessVar{Name: "var4"},
							&AccessVar{Name: "a4"},
							// |=
							&AccessVar{Name: "var5"},
							&AccessVar{Name: "a5"},
							// &=
							&AccessVar{Name: "var6"},
							&AccessVar{Name: "a6"},
							// <<=
							&AccessVar{Name: "var7"},
							&AccessVar{Name: "a7"},
							// >>=
							&AccessVar{Name: "var8"},
							&AccessVar{Name: "a8"},
						}},
				},
			},
		},
		//		{
		//			"function def 4",
		//			`
		//typedef struct {
		//  char xxx;
		//} Gt;
		//
		//typedef struct {
		//  int aaa;
		//  Gt bbb;
		//} St;
		//
		//int muruchi_piyomi(char *s) {
		//  St purin;
		//  St *p;
		//  purin.aaa = 100;
		//  purin.bbb.xxx = 'A';
		//  p = &purin;
		//  p->aaa = purin.aaa + 200;
		//  return (0);
		//}
		//`,
		//			&Module{
		//				[]Statement{
		//					&FunctionDef{Name: "muruchi_piyomi",
		//						Params: []*VariableDef{{Name: "s"}},
		//						Statements: []Statement{
		//							&VariableDef{Name: "purin"},
		//							&VariableDef{Name: "p"},
		//							&AccessVar{Name: "purin.aaa"},
		//							&AccessVar{Name: "purin.bbb.xxx"},
		//							&AccessVar{Name: "p"},
		//							&AccessVar{Name: "purin"},
		//							&AccessVar{Name: "p->aaa"},
		//							&AccessVar{Name: "purin.aaa"},
		//						}},
		//				},
		//			},
		//		},

		//		{
		//			"function def 5",
		//			`
		//void whsxks(int a) {
		//}
		//`,
		//			&Module{
		//				[]Statement{
		//					&FunctionDef{Name: "whsxks",
		//						Params: []*VariableDef{{Name: "a"}},
		//						Block:  &BlockStatement{Statements: []Statement{}}},
		//				},
		//			},
		//		},
		//		{
		//			"function def 6",
		//			`
		//void haraheri(int a, char *b, unsigned char s[]) {
		//}
		//`,
		//			&Module{
		//				[]Statement{
		//					&FunctionDef{Name: "haraheri",
		//						Params: []*VariableDef{{Name: "a"}, {Name: "b"}, {Name: "s"}},
		//						Block:  &BlockStatement{Statements: []Statement{}}},
		//				},
		//			},
		//		},
		//		{
		//			"function def 7",
		//			`
		//void hoge(void){
		//  _p->_p++ = _c;
		//  _p->_p-- = _c;
		//}
		//`,
		//			&Module{
		//				[]Statement{
		//					&FunctionDef{Name: "hoge",
		//						Params: []*VariableDef{},
		//						Block: &BlockStatement{Statements: []Statement{
		//							&AccessVar{Name: "_p->_p"},
		//							&AccessVar{Name: "_c"},
		//							&AccessVar{Name: "_p->_p"},
		//							&AccessVar{Name: "_c"},
		//						}}},
		//				},
		//			},
		//		},
		//		{
		//			"testx",
		//			`
		//# 1 "hoge.c"
		//# 1 "<built-in>" 1
		//# 1 "<built-in>" 3
		//# 366 "<built-in>" 3
		//# 1 "<command line>" 1
		//# 1 "<built-in>" 2
		//# 1 "hoge.c" 2
		//
		//int func(int a)
		//{
		//    int hoge = 0;
		//    hoge++;
		//    a = a + (10);
		//    return a;
		//}
		//`,
		//			&Module{
		//				[]Statement{
		//					&FunctionDef{Name: "func",
		//						Params: []*VariableDef{{Name: "a"}},
		//						Block: &BlockStatement{Statements: []Statement{
		//							&VariableDef{Name: "hoge"},
		//							&AccessVar{Name: "hoge"},
		//							&AccessVar{Name: "a"},
		//							&AccessVar{Name: "a"},
		//							&AccessVar{Name: "a"},
		//						},
		//						},
		//					},
		//				},
		//			},
		//		},
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
