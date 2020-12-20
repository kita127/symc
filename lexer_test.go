package symc

import (
	_ "reflect"
	"testing"
)

func TestLexicalize(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  []*Token
	}{
		{
			"test0",
			``,
			[]*Token{
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test1",
			`   char   `,
			[]*Token{
				{
					word,
					"char",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test2",
			`   char		hoge   `,
			[]*Token{
				{
					word,
					"char",
				},
				{
					word,
					"hoge",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test3",
			`=`,
			[]*Token{
				{
					assign,
					"=",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test4",
			`=+-!*/<>;(),{}[]`,
			[]*Token{
				{
					assign,
					"=",
				},
				{
					plus,
					"+",
				},
				{
					minus,
					"-",
				},
				{
					bang,
					"!",
				},
				{
					asterisk,
					"*",
				},
				{
					slash,
					"/",
				},
				{
					lt,
					"<",
				},
				{
					gt,
					">",
				},
				{
					semicolon,
					";",
				},
				{
					lparen,
					"(",
				},
				{
					rparen,
					")",
				},
				{
					comma,
					",",
				},
				{
					lbrace,
					"{",
				},
				{
					rbrace,
					"}",
				},
				{
					lbracket,
					"[",
				},
				{
					rbracket,
					"]",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test5",
			` = + - ! * / < > ; ( ) , { } [ ] `,
			[]*Token{
				{
					assign,
					"=",
				},
				{
					plus,
					"+",
				},
				{
					minus,
					"-",
				},
				{
					bang,
					"!",
				},
				{
					asterisk,
					"*",
				},
				{
					slash,
					"/",
				},
				{
					lt,
					"<",
				},
				{
					gt,
					">",
				},
				{
					semicolon,
					";",
				},
				{
					lparen,
					"(",
				},
				{
					rparen,
					")",
				},
				{
					comma,
					",",
				},
				{
					lbrace,
					"{",
				},
				{
					rbrace,
					"}",
				},
				{
					lbracket,
					"[",
				},
				{
					rbracket,
					"]",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test6",
			`&~^|:?.\`,
			[]*Token{
				{
					ampersand,
					"&",
				},
				{
					tilde,
					"~",
				},
				{
					caret,
					"^",
				},
				{
					vertical,
					"|",
				},
				{
					colon,
					":",
				},
				{
					question,
					"?",
				},
				{
					period,
					".",
				},
				{
					backslash,
					"\\",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test8",
			`   ident00+123;   `,
			[]*Token{
				{
					word,
					"ident00",
				},
				{
					plus,
					"+",
				},
				{
					integer,
					"123",
				},
				{
					semicolon,
					";",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test9",
			`# 1 "hoge.c"`,
			[]*Token{
				{
					comment,
					" 1 \"hoge.c\"",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test10",
			`# 1 "hoge.c"
# 1 "<built-in>" 1`,
			[]*Token{
				{
					comment,
					" 1 \"hoge.c\"",
				},
				{
					comment,
					" 1 \"<built-in>\" 1",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test11",
			`0 0U 123 0xA1c 0765 0b0110 567u 567U 567l 567L 567lu 567UL`,
			[]*Token{
				{
					integer,
					"0",
				},
				{
					integer,
					"0U",
				},
				{
					integer,
					"123",
				},
				{
					integer,
					"0xA1c",
				},
				{
					integer,
					"0765",
				},
				{
					integer,
					"0b0110",
				},
				{
					integer,
					"567u",
				},
				{
					integer,
					"567U",
				},
				{
					integer,
					"567l",
				},
				{
					integer,
					"567L",
				},
				{
					integer,
					"567lu",
				},
				{
					integer,
					"567UL",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test12",
			`0.123 987.123 123.`,
			[]*Token{
				{
					float,
					"0.123",
				},
				{
					float,
					"987.123",
				},
				{
					float,
					"123.",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test13",
			`return if else while do goto for break continue switch case default
 extern volatile const typedef union struct enum`,
			[]*Token{
				{
					keyReturn,
					"return",
				},
				{
					keyIf,
					"if",
				},
				{
					keyElse,
					"else",
				},
				{
					keyWhile,
					"while",
				},
				{
					keyDo,
					"do",
				},
				{
					keyGoto,
					"goto",
				},
				{
					keyFor,
					"for",
				},
				{
					keyBreak,
					"break",
				},
				{
					keyContinue,
					"continue",
				},
				{
					keySwitch,
					"switch",
				},
				{
					keyCase,
					"case",
				},
				{
					keyDefault,
					"default",
				},
				{
					keyExtern,
					"extern",
				},
				{
					keyVolatile,
					"volatile",
				},
				{
					keyConst,
					"const",
				},
				{
					keyTypedef,
					"typedef",
				},
				{
					keyUnion,
					"union",
				},
				{
					keyStruct,
					"struct",
				},
				{
					keyEnum,
					"enum",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test13",
			`char hoge[] = "hello";`,
			[]*Token{
				{
					word,
					"char",
				},
				{
					word,
					"hoge",
				},
				{
					lbracket,
					"[",
				},
				{
					rbracket,
					"]",
				},
				{
					assign,
					"=",
				},
				{
					str,
					"\"hello\"",
				},
				{
					semicolon,
					";",
				},
				{
					eof,
					"eof",
				},
			},
		},
		{
			"test14",
			`int hoge = 0;`,
			[]*Token{
				{
					word,
					"int",
				},
				{
					word,
					"hoge",
				},
				{
					assign,
					"=",
				},
				{
					integer,
					"0",
				},
				{
					semicolon,
					";",
				},
				{
					eof,
					"eof",
				},
			},
		},

		{
			"testx",
			`
# 1 "hoge.c"

int func(int a) {
      a = a + (10);
        return a;
    }
`,
			[]*Token{
				{
					comment,
					` 1 "hoge.c"`,
				},
				{
					word,
					`int`,
				},
				{
					word,
					`func`,
				},
				{
					lparen,
					`(`,
				},
				{
					word,
					`int`,
				},
				{
					word,
					`a`,
				},
				{
					rparen,
					`)`,
				},
				{
					lbrace,
					`{`,
				},
				{
					word,
					`a`,
				},
				{
					assign,
					`=`,
				},
				{
					word,
					`a`,
				},
				{
					plus,
					`+`,
				},
				{
					lparen,
					`(`,
				},
				{
					integer,
					`10`,
				},
				{
					rparen,
					`)`,
				},
				{
					semicolon,
					`;`,
				},
				{
					keyReturn,
					`return`,
				},
				{
					word,
					`a`,
				},
				{
					semicolon,
					`;`,
				},
				{
					rbrace,
					`}`,
				},
				{
					eof,
					"eof",
				},
			},
		},
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		got := l.lexicalize()
		if len(got) != len(tt.expect) {
			t.Fatalf("got len=%v, expect len=%v", len(got), len(tt.expect))
		}
		for i, v := range got {
			e := tt.expect[i]
			if v.tokenType != e.tokenType {
				t.Errorf("got type=%v, expect type=%v", v.tokenType, tt.expect[i].tokenType)
			}
			if v.literal != e.literal {
				t.Errorf("got literal=%v, expect literal=%v", v.literal, tt.expect[i].literal)
			}
		}
	}
}
