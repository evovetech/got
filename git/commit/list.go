package commit

type List struct {
	value Commit
	next  *List
}

func (list *List) Append(c Commit) bool {
	if l := list.Last(); l != nil {
		if l.value == nil {
			l.value = c
		} else {
			l.next = &List{value: c}
		}
		return true
	}
	return false
}

func (list *List) Value() Commit {
	if list == nil {
		return nil
	}
	return list.value
}

func (list *List) Next() *List {
	if list == nil {
		return nil
	}
	return list.next
}

func (list *List) Last() *List {
	if list == nil {
		return nil
	}

	l := list
	for l.next != nil {
		l = l.next
	}
	return l
}
