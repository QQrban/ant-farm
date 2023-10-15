package types

import "fmt"

type Combination struct {
	len   int
	last  int
	paths [][]*Vertex
}

func (c *Combination) String() string {
	str := fmt.Sprintf("{%d: ", c.len)
	for _, path := range c.paths {
		for _, v := range path {
			str += v.Name + " "
		}
		str += "}\n"
	}
	return str
}
