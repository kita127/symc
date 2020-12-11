package symc

import (
	"reflect"
	"testing"
)

func TestParseModule(t *testing.T) {
	testTbl := []struct {
		comment string
		src     string
		expect  *Module
	}{
		{
			"test1",
			`hoge;`,
			&Module{
				[]Statement{
					&VariableDef{"hoge"},
				},
			},
		},
	}

	for _, tt := range testTbl {
		got := ParseModule(tt.src)
		if !reflect.DeepEqual(got, tt.expect) {
			t.Errorf("got=%v, expect=%v", got.Statements[0], tt.expect.Statements[0])
		}
	}
}

// func TestConvertCaseLStoLC(t *testing.T) {
// 	testTbl := []struct {
// 		comment string
// 		src     string
// 		expect  string
// 	}{
// 		{"test convert 1",
// 			inputSrc,
// 			`package hoge
//
// var UPPER_SNAKE_VAR int
// var lowerSnakeVar int
// var UpperCamelVar int
// var lowerCamelVar int
//
// const UPPER_SNAKE_CONST int = 0
// const lowerSnakeConst int = 0
// const UpperCamelConst int = 0
// const lowerCamelConst int = 0
//
// func UPPER_SNAKE_FUNC() {
// 	LOCAL_VAR := 0
// }
//
// func lowerSnakeFunc() {
// 	localVar := 0
// }
//
// func UpperCamelFunc() {
// 	LocalVar := 0
// }
//
// func lowerCamelFunc() {
// 	localVar := 0
// }
// `,
// 		},
// 	}
//
// 	for _, tt := range testTbl {
// 		got, err := ConvertCase(tt.src, LowerSnake, LowerCamel)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		if got != tt.expect {
// 			t.Errorf("got=%v, expect=%v", got, tt.expect)
// 		}
// 	}
// }
