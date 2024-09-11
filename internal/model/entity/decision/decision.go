package decision

import "github.com/google/uuid"

type Verdict string

const (
	Approved Verdict = "Approved"
	Rejected Verdict = "Rejected"
)

func IsDecisionVerdict(string2 string) bool {
	mapped := Verdict(string2)
	return Approved == mapped || Rejected == mapped
}

type Decision struct {
	Id       uuid.UUID
	Verdict  Verdict
	Username string
	BidId    uuid.UUID
}
