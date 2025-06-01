package main

import (
	"strings"
	"testing"
)

func TestXmlParse(t *testing.T) {
	tests := []struct {
		xml, want string
	}{
		{
			xml: `<doc><a id="b"><b/>hi<b>rah</b></a></doc>`,
			want: `doc []
  a [{{ id} b}]
    b []
    "hi"
    b []
      "rah"
`,
		},
	}

	for _, test := range tests {
		node, err := xmlParse(strings.NewReader(test.xml))
		if err != nil {
			t.Error(test, err)
		}
		if n, ok := node.(*Element); ok {
			if n.String() != test.want {
				t.Errorf("%s got: %s,want: %s\n", test, n.String(), test.want)
			}
		}
	}
}
