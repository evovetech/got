package file

type Type uint32

func (t Type) String() string {
	var str string
	if t.HasFlag(Rn) {
		str += "R"
	}
	switch {
	case t.HasFlag(Add):
		str += "A"
	case t.HasFlag(Del):
		str += "D"
	default:
		str += "?"
	}
	return str
}

func (t Type) HasFlag(flag Type) bool {
	return t&flag != 0
}

const (
	Add Type = 1 << iota
	Del
	Rn
)
