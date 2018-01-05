package types

type Sha String

func (sha Sha) String() string {
	return string(sha)
}
