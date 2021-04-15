package eval

import (
	"encoding/json"
	"testing"
)

type input struct {
	cond string
	env map[string]interface{}
	res bool
}

func TestParse(t *testing.T) {
	// antlr4 -Dlanguage=Go -o parser cond_parser.g4 -visitor -no-listener

	ins := []input{
		{
			cond: `b==null`,
			env:  map[string]interface{}{},
			res:  true,
		},
		{
			cond: `a=="123"`,
			env:  map[string]interface{}{"a": "123"},
			res:  true,
		},
		{
			cond: `a=="123" && b == 456`,
			env:  map[string]interface{}{"a": "123", "b": json.Number("456")},
			res:  true,
		},
		{
			cond: `a=="123" || b==456 && c =="789"`,
			env:  map[string]interface{}{"b": json.Number("456"), "c": "789"},
			res:  true,
		},
		{
			cond: `a=="123" || (b==456 && c =="789")`,
			env:  map[string]interface{}{"a": "", "b": json.Number("456"), "c": "789"},
			res:  true,
		},
		{
			cond: `a=="123" || (b==456 && c =="789") && d=="d"`,
			env:  map[string]interface{}{"a": "", "b": json.Number("456"), "c": "789", "d":"d"},
			res:  true,
		},
		{
			cond: `a=="123" || ((b==456 || bb==777) && c =="789") && d=="d"`,
			env:  map[string]interface{}{"a": "", "bb": json.Number("777"), "c": "789", "d":"d"},
			res:  true,
		},
	}

	for _, in := range ins {
		ast, err := GetConditionAst(in.cond)
		if err != nil {
			t.Fatal(err)
		}
		if ast.Match(in.env) != in.res {
			t.Fatal(in.cond, "match fail")
		}
	}
}


