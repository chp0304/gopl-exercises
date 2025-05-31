package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/chp0304/gopl-exercises/ch7/Exercise7.14/eval"
)

func main() {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("please input your expr: ")
	if scanner.Scan() {
		input = scanner.Text()
	}
	fmt.Println("your expr:", input)
	vars := make(map[eval.Var]bool)
	expr, err := eval.ParseAndCheck(input, vars)
	if err != nil {
		panic(err)
	}
	env := make(eval.Env)
	var value float64
	for variable, _ := range vars {
		fmt.Printf("please input value of Var '%s' = ", variable)
		fmt.Scan(&value)
		env[variable] = value
	}
	fmt.Println()
	result := expr.Eval(env)
	fmt.Printf("your expr %s = %.2g\n", expr.String(), result)
}
