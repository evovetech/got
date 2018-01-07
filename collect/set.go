package collect

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "T_=String,Sha"

type T_Set = T_BooleanMap

func (set *T_Set) Add(t T_) bool {
	return set.Put(t, true)
}

type T_CounterSet = T_IntegerMap

func (set *T_CounterSet) Add(t T_, num int) int {
	s := set.Init()
	next := s[t] + num
	s[t] = next
	return next
}

func (set *T_CounterSet) Subtract(t T_, num int) int {
	return set.Add(t, -1)
}

func (set *T_CounterSet) Increment(t T_) int {
	return set.Add(t, 1)
}

func (set *T_CounterSet) Decrement(t T_) int {
	return set.Subtract(t, 1)
}
