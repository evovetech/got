package git

type CommitList struct {
	value Commit
	next  *CommitList
}

func (list *CommitList) Append(c Commit) bool {
	if l := list.Last(); l != nil {
		if l.value == nil {
			l.value = c
		} else {
			l.next = &CommitList{value: c}
		}
		return true
	}
	return false
}

func (list *CommitList) Value() Commit {
	if list == nil {
		return nil
	}
	return list.value
}

func (list *CommitList) Next() *CommitList {
	if list == nil {
		return nil
	}
	return list.next
}

func (list *CommitList) Last() *CommitList {
	if list == nil {
		return nil
	}

	l := list
	for l.next != nil {
		l = l.next
	}
	return l
}
