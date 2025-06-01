package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Node interface{}

type CharData string

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (n *Element) String() string {
	b := &bytes.Buffer{}
	visit(n, b, 0)
	return b.String()
}

func visit(n Node, w io.Writer, depth int) {
	switch n := n.(type) {
	case *Element:
		fmt.Fprintf(w, "%*s%s %s\n", depth*2, "", n.Type.Local, n.Attr)
		for _, c := range n.Children {
			visit(c, w, depth+1)
		}
	case CharData:
		fmt.Fprintf(w, "%*s%q\n", depth*2, "", n)
	default:
		panic(fmt.Sprintf("got %T", n))
	}
}

func xmlParse(r io.Reader) (Node, error) {
	var root Node
	dec := xml.NewDecoder(r)
	var stack []Node // stack of element names
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
			node := &Element{
				Type:     tok.Name,
				Attr:     tok.Attr,
				Children: make([]Node, 0),
			}
			if len(stack) == 0 {
				root = node
			}
			stack = append(stack, node) // push
		case xml.EndElement:
			child := stack[len(stack)-1]
			// fmt.Println(child)
			stack = stack[:len(stack)-1] // pop
			if len(stack) > 0 {
				if e, ok := stack[len(stack)-1].(*Element); ok {
					e.Children = append(e.Children, child)
				}
			}
		case xml.CharData:
			char := CharData(tok)
			if len(stack) > 0 {
				if e, ok := stack[len(stack)-1].(*Element); ok {
					e.Children = append(e.Children, char)
				}
			}
		}
	}
	return root, nil
}

func main() {
	node, err := xmlParse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	fmt.Println(node)
}
