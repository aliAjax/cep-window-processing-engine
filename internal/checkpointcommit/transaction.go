package checkpointcommit

type Transaction struct {
	commit   func() error
	rollback func() error
}

func NewTransaction(commit, rollback func() error) Transaction {
	return Transaction{commit: commit, rollback: rollback}
}

func (t Transaction) Finish(operationErr error) error {
	if t.commit != nil {
		return t.commit()
	}
	if operationErr != nil && t.rollback != nil {
		return t.rollback()
	}
	return operationErr
}
