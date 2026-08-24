package delivery

type Transaction interface {
	Commit() error
	Rollback() error
}

func Finish(tx Transaction, operationErr error) error {
	if tx == nil {
		return operationErr
	}
	if operationErr != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}
		return nil
	}
	return tx.Commit()
}
