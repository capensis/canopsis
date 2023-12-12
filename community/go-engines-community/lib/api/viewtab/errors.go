package viewtab

type ValidationError struct {
	err error
}

func (v ValidationError) Error() string {
	return v.err.Error()
}

func (v ValidationError) Unwrap() error {
	return v.err
}
