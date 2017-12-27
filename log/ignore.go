package log

type ignore uint8

const devNull = ignore(0)

func (i ignore) Write(p []byte) (int, error) {
	// write nothing
	return len(p), nil
}
