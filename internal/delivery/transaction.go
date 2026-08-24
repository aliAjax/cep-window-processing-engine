package delivery

import "errors"

type Transaction interface {
	Commit() error
	Rollback() error
}

func Finish(tx Transaction, operationErr error) error {
	if tx == nil {
		return operationErr
	}
	if operationErr != nil {
		return errors.Join(operationErr, tx.Rollback())
	}
	return tx.Commit()
}
