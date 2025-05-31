package main

import (
	"fmt"

	"github.com/chp0304/gopl-exercises/ch7/Exercise7.14/eval"
)

func main() {
	expr := "sin(x) + 2!"
	env := eval.Env{
		eval.Var("x"): 3.1415926,
	}
	exp, err := eval.Parse(expr)
	fmt.Printf("expr = %s\n", exp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result:%f\n", exp.Eval(env))
}
