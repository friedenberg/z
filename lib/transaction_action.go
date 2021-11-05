package lib

type TransactionAction int

const (
	TransactionActionNone  = 0
	TransactionActionAdded = 1 << iota
	TransactionActionModified
	TransactionActionDeleted
)
