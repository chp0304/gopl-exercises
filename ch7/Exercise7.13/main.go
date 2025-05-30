package main

import (
	"fmt"

	"github.com/chp0304/gopl-exercises/ch7/Exercise7.13/eval"
)

func main() {
	expr := "sin(3.1415926)"
	env := eval.Env{
		eval.Var("F"): -40,
	}
	exp, err := eval.Parse(expr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result:%f\n", exp.Eval(env))
}
