package eval

import "testing"

func TestString(t *testing.T) {
	tests := []struct {
		expr Expr
		want string
	}{
		{
			Var("X"),
			"X",
		},
		{
			literal(1),
			"1",
		},
		{
			unary{
				'-',
				literal(1),
			},
			"-1",
		},
		{
			binary{
				x:  literal(1),
				op: '*',
				y:  Var("x"),
			},
			"(1 * x)",
		},
		{
			call{
				fn: "pow",
				args: []Expr{
					Var("x"),
					Var("y"),
				},
			},
			"pow(x, y)",
		},
	}

	for _, test := range tests {
		result := test.expr.String()
		if result != test.want {
			t.Errorf("expr.String() = %s, want = %s\n", result, test.want)
		}
	}
}
