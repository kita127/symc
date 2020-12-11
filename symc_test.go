package symc

//import (
//	"reflect"
//	"testing"
//)
//
//func TestParseModule(t *testing.T) {
//	testTbl := []struct {
//		comment string
//		src     string
//		expect  *Module
//	}{
//		{
//			"test1",
//			`char hoge;`,
//			&Module{
//				[]Statement{
//					&VariableDef{"hoge"},
//				},
//			},
//		},
//		{
//			"test2",
//			`char hoge; int fuga`,
//			&Module{
//				[]Statement{
//					&VariableDef{"hoge"},
//					&VariableDef{"fuga"},
//				},
//			},
//		},
//	}
//
//	for _, tt := range testTbl {
//		t.Logf("%s", tt.comment)
//		got := ParseModule(tt.src)
//		if !reflect.DeepEqual(got, tt.expect) {
//			t.Errorf("got=%v, expect=%v", got.Statements[0], tt.expect.Statements[0])
//		}
//	}
//}
