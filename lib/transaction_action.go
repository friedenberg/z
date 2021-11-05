package lib

type TransactionAction int

const (
	TransactionActionNone  = 0
	TransactionActionAdded = iota
	TransactionActionModified
	TransactionActionDeleted
)
