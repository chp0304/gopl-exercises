package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chp0304/gopl-exercises/ch7/Exercise7.14/eval"
)

func calculate(w http.ResponseWriter, r *http.Request) {
	vars := make(map[eval.Var]bool)
	expr, err := eval.ParseAndCheck(r.URL.Query().Get("expr"), vars)
	fmt.Println(expr)
	if err != nil {
		http.Error(w, "bad expr: "+err.Error(), http.StatusBadRequest)
		return
	}
	env := make(eval.Env)
	xStr := r.URL.Query().Get("x")
	yStr := r.URL.Query().Get("y")
	if x, err := strconv.ParseFloat(xStr, 64); err == nil {
		env[eval.Var("x")] = x
	}
	if y, err := strconv.ParseFloat(yStr, 64); err == nil {
		env[eval.Var("y")] = y
	}
	fmt.Println(env)
	result := expr.Eval(env)
	fmt.Fprintf(w, "%s = %.6g\n", expr.String(), result)
}

func main() {
	http.HandleFunc("/calculate", calculate)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
