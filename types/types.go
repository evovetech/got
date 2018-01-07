package types

type (
	String  = string
	Boolean = bool
	Integer = int
)

type Sha String

func (sha Sha) String() string {
	return string(sha)
}
