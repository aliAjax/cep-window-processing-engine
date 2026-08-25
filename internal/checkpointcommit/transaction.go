package checkpointcommit

type Transaction struct {
	commit   func() error
	rollback func() error
}

func NewTransaction(commit, rollback func() error) Transaction {
	return Transaction{commit: commit, rollback: rollback}
}

func (t Transaction) Finish(operationErr error) error {
	if operationErr != nil {
		if t.rollback != nil {
			_ = t.rollback()
		}
		return operationErr
	}
	if t.commit == nil {
		return nil
	}
	return t.commit()
}
