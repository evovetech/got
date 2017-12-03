package merge

import (
	"github.com/evovetech/got/git"
)

type simple Merger

func (s *simple) Run() (RunStep, error) {
	if err := s.HeadRef.Checkout(); err != nil {
		return nil, err
	}

	m := s.MergeRef.Merge()
	m.IgnoreAllSpace()
	if err := m.Run(); err != nil {
		if err = git.AbortMerge(); err != nil {
			return nil, err
		}
		return (*multi)(s).Run()
	}

	return nil, nil
}
