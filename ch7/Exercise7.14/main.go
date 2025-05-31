package main

import (
	"fmt"

	"github.com/chp0304/gopl-exercises/ch7/Exercise7.14/eval"
)

func main() {
	expr := "sin(x)! + 3"
	env := eval.Env{
		eval.Var("x"): 0,
	}
	exp, err := eval.Parse(expr)
	fmt.Printf("expr = %s\n", exp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result:%f\n", exp.Eval(env))
}
