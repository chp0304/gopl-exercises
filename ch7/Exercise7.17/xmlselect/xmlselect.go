// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 214.
//!+

// Xmlselect prints the text of selected elements of an XML document.
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune // current lookahead token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

type lexPanic string

type Attr struct {
	Name  string
	Value string
}

// selector = tag
//
//	| tag attr ...
//	| attr ...
type selector struct {
	tag  string
	Attr []Attr
}

func (lex *lexer) eatWhiteSpace() {
	for lex.token == '\t' || lex.token == ' ' {
		lex.next()
	}
}

func (lex *lexer) describe() string {
	switch lex.token {
	case scanner.EOF:
		return "end of file"
	case scanner.Ident:
		return fmt.Sprintf("identifier %s", lex.text())
	case scanner.Int, scanner.Float:
		return fmt.Sprintf("number %s", lex.text())
	}
	return fmt.Sprintf("%q", rune(lex.token)) // any other rune
}

// a
// a[id="3",id="4"] [id="4"]
// [id="3"] [id="4"]
func parseSelectors(input string) (_ []selector, err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case lexPanic:
			err = fmt.Errorf("%s", x)
		default:
			// unexpected panic: resume state of panic.
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanStrings
	lex.scan.Whitespace = 0
	lex.next()
	selectors := make([]selector, 0)
	for lex.token != scanner.EOF {
		selector, err := parseSelector(lex)
		if err != nil {
			panic(lexPanic(fmt.Sprintf("err: %s", err.Error())))
		}
		selectors = append(selectors, selector)
	}
	return selectors, nil
}

func parseSelector(lex *lexer) (selector, error) {
	lex.eatWhiteSpace()
	s := selector{}
	if lex.token != '[' {
		s.tag = lex.text()
		lex.next()
	}
	if attrs, err := parseAttrs(lex); err == nil {
		s.Attr = attrs
	} else {
		panic(lexPanic(fmt.Sprintf("err %s", err.Error())))

	}
	if lex.token == ']' {
		lex.next()
	}
	return s, nil
}

func parseAttrs(lex *lexer) ([]Attr, error) {
	if lex.token != '[' {
		return nil, nil
	}
	attrs := make([]Attr, 0)
	lex.next()
	for lex.token != ']' {
		attr := Attr{}
		attr.Name = lex.text()
		lex.next()
		if lex.token != '=' {
			if lex.token == ']' {
				attrs = append(attrs, attr)
				break
			}
			panic(lexPanic(fmt.Sprintf("got %s, want ident", lex.describe())))
		}
		lex.next() // skip '='
		attr.Value = strings.Trim(lex.text(), `"`)
		attrs = append(attrs, attr)
		lex.next()
		if lex.token == ',' {
			lex.next()
		}
	}
	return attrs, nil
}

func StringElement(start xml.StartElement) string {
	str := fmt.Sprintf("%s[", start.Name.Local)
	flag := false
	for _, attr := range start.Attr {
		if flag {
			str = str + ","
		}
		flag = true
		str += attr.Value
	}
	str += "]"
	return str

}

func StringElements(elements []xml.StartElement) string {
	str := ""
	for _, e := range elements {
		str += StringElement(e)
		str += " "
	}
	return str

}

func isSelected(stack []xml.StartElement, selectors []selector) bool {
	// fmt.Println(stack)
	// for _, s := range stack {
	// 	if s.Name.Local == "h2" {
	// 		fmt.Println(StringElements(stack))
	// 	}
	// }
	// fmt.Println()
	for len(selectors) <= len(stack) {
		if len(selectors) == 0 {
			return true
		}
		if selectors[0].tag == "" || stack[0].Name.Local == selectors[0].tag {
			if compareAttr(stack[0].Attr, selectors[0].Attr) {
				selectors = selectors[1:]
			}
		}
		stack = stack[1:]
	}
	return false
}

// div[body] div[div1] h2[]
func compareAttr(xAttrs []xml.Attr, yAttrs []Attr) bool {
	if len(yAttrs) == 0 {
		return true
	}
	xMap := make(map[string]string)
	for _, x := range xAttrs {
		xMap[x.Name.Local] = x.Value
	}
	for _, y := range yAttrs {
		if value, ok := xMap[y.Name]; !ok || (value != y.Value && y.Value != "") {
			return false
		}
	}
	return true
}

func xmlSelect(w io.Writer, r io.Reader, selectors []selector) {
	dec := xml.NewDecoder(r)
	var stack []xml.StartElement // stack of element names
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok) // push
		case xml.EndElement:
			stack = stack[:len(stack)-1] // pop
		case xml.CharData:
			if isSelected(stack, selectors) {
				fmt.Fprintf(w, "%s\n", tok)
			}
		}
	}
}

// test cases:
// ./fetch http://www.w3.org/TR/2006/REC-xml11-20060816 | ./xmlselect div div h2
// ./fetch http://www.w3.org/TR/2006/REC-xml11-20060816 | ./xmlselect div 'div[class="div1"]' h2
func main() {
	selectors, err := parseSelectors(strings.Join(os.Args[1:], " "))
	if err != nil {
		panic(err)
	}
	xmlSelect(os.Stdout, os.Stdin, selectors)
}

//!-
