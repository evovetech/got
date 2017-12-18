package file

type TypeCount map[Type]int

func (c TypeCount) add(t Type, i int) {
	c[t] = c[t] + i
}

func (c TypeCount) addAll(o TypeCount) {
	for t, n := range o {
		c.add(t, n)
	}
}
