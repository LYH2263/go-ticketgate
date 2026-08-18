package keyring

type Status byte

const (
	StatusActive   Status = 1
	StatusRetiring Status = 2
	StatusRetired  Status = 3
)

func (s Status) String() string {
	switch s {
	case StatusActive:
		return "active"
	case StatusRetiring:
		return "retiring"
	case StatusRetired:
		return "retired"
	default:
		return "unknown"
	}
}

func (s Status) CanIssue() bool {
	return s == StatusActive
}

func (s Status) CanVerify() bool {
	return s == StatusActive || s == StatusRetiring
}
