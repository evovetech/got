// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package collect

type StringSet = StringBooleanMap

func (set *StringSet) Add(t String) bool {
	return set.Put(t, true)
}

type StringCounterSet = StringIntegerMap

func (set *StringCounterSet) Add(t String, num int) int {
	s := set.Init()
	next := s[t] + num
	s[t] = next
	return next
}

func (set *StringCounterSet) Subtract(t String, num int) int {
	return set.Add(t, -1)
}

func (set *StringCounterSet) Increment(t String) int {
	return set.Add(t, 1)
}

func (set *StringCounterSet) Decrement(t String) int {
	return set.Subtract(t, 1)
}

type ShaSet = ShaBooleanMap

func (set *ShaSet) Add(t Sha) bool {
	return set.Put(t, true)
}

type ShaCounterSet = ShaIntegerMap

func (set *ShaCounterSet) Add(t Sha, num int) int {
	s := set.Init()
	next := s[t] + num
	s[t] = next
	return next
}

func (set *ShaCounterSet) Subtract(t Sha, num int) int {
	return set.Add(t, -1)
}

func (set *ShaCounterSet) Increment(t Sha) int {
	return set.Add(t, 1)
}

func (set *ShaCounterSet) Decrement(t Sha) int {
	return set.Subtract(t, 1)
}
