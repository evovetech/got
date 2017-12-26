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

func (c TypeCount) All(t Type) (n int) {
	for k, v := range c {
		if k.HasFlag(t) {
			n += v
		}
	}
	return
}

func (c TypeCount) AllAdd() int {
	return c.All(Add)
}

func (c TypeCount) AllDel() int {
	return c.All(Del)
}
