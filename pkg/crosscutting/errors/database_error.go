package errors

const ErrorTypeDBTransactionNotFound ErrorType = "DBTransactionNotFound"

func NewTransactionNotFoundError() error {
	err := New("transaction not found")
	err = AddType(err, ErrorTypeDBTransactionNotFound)
	return err
}

func IsTransactionNotFoundErr(err error) bool {
	return IsType(err, ErrorTypeDBTransactionNotFound)
}
