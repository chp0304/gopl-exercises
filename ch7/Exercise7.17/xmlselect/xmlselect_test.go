package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestXmlSelect(t *testing.T) {
	input := `c b a a[id="2",id] [id="3",class=wide] b [id="3"] [id="4"]`
	if selectors, err := parseSelectors(input); err == nil {
		for _, s := range selectors {
			fmt.Println(s)
		}
	}
}

// func TestCompareAttr(t *testing.T) {
// 	tests := []struct {
// 		xAttrs []xml.Attr
// 		yAttrs []Attr
// 		want   string
// 	}{
// 		{
// 			xAttrs: []xml.Attr{
// 				{
// 					xml.Name:xml.Name{
// 						Local: "id",
// 					},
// 					xml.Value:"3",
// 				},

// 			},

// 		}
// 	}
// }

func TestXMLSelect(t *testing.T) {
	tests := []struct {
		selectors, xml, want string
	}{
		{`a[id="3"] [id="4"]`, `<a id="3"><p id="4">good</p></a>`, "good\n"},
		{`a[id="3"] [id="4"]`, `<a><p id="4">bad</p></a>`, ""},
		{`[id="3"] [id]`, `<a id="3"><p id="4">good</p></a><a><p id="4">bad</p></a>`, "good\n"},
		{`[id="3",class=big]`, `<a id="3" class="big">good</a><a id="3">bad</a>`, "good\n"},
		{`p a`, `<p><a>1</a><p><a>2</a></p></p><a>bad</a><p><a>3</a></p>`, "1\n2\n3\n"},
		{`div div h2`, `<div id="body"><div id="div1"><h2>good</h2></div></div>`, "good\n"},
	}
	for _, test := range tests {
		sels, err := parseSelectors(test.selectors)
		fmt.Println(sels)
		if err != nil {
			t.Error(test, err)
			continue
		}
		w := &bytes.Buffer{}
		xmlSelect(w, strings.NewReader(test.xml), sels)
		if w.String() != test.want {
			t.Errorf("%s: got %q, want %q", test, w.String(), test.want)
		}
	}
}
