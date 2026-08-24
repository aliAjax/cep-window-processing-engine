package checkpointcommit

import "errors"

type Transaction struct {
	commit   func() error
	rollback func() error
}

func NewTransaction(commit, rollback func() error) Transaction {
	return Transaction{commit: commit, rollback: rollback}
}

func (t Transaction) Finish(operationErr error) error {
	if operationErr != nil {
		if t.rollback == nil {
			return operationErr
		}
		return errors.Join(operationErr, t.rollback())
	}
	if t.commit == nil {
		return nil
	}
	return t.commit()
}
