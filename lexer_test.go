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
	}

	for _, tt := range testTbl {
		t.Logf("%s", tt.comment)
		l := NewLexer(tt.src)
		got := l.lexicalize(tt.src)
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
