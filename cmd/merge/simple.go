package merge

import "github.com/evovetech/got/git"

type simple Merger

func (s *simple) Run() (RunStep, error) {
	if err := s.HeadRef.Checkout(); err != nil {
		return nil, err
	}

	m := s.MergeRef.Merge()
	m.IgnoreAllSpace()
	if err := m.Run(); err == nil {
		return nil, nil
	}

	s.reset()

	return (*multi)(s).Run()
}

func (s *simple) reset() {
	git.Merge().Abort()
	head := s.HeadRef
	if e := head.Checkout(); e == nil {
		head.Reset().Hard()
	}
}
