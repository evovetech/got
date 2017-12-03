package git

type Runner interface {
	Run() error
}

type runner func() error

func (r runner) Run() error {
	return r()
}

func ToRunner(f func() error) Runner {
	return (runner)(f)
}
