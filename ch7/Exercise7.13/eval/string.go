package eval

import (
	"fmt"
)

func (v Var) String() string {
	return fmt.Sprintf("%s", string(v))
}

func (l literal) String() string {
	return fmt.Sprintf("%g", l)
}

func (u unary) String() string {
	return fmt.Sprintf("%s%s", string(u.op), u.x.String())
}

func (b binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.x.String(), string(b.op), b.y.String())
}

func (c call) String() string {
	res := c.fn
	res += "("
	flag := false
	for _, arg := range c.args {
		if flag {
			res += ", "
		}
		flag = true
		res += arg.String()
	}
	res += ")"
	return res
}
