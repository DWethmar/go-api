package handler

type ErrValidation struct {
	errors []string
}

func (e *ErrValidation) Error() string {
	return ":D"
}
