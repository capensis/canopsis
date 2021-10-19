package file

type ValidationError struct {
	field string
	error string
}

func (e ValidationError) Error() string {
	return e.error
}
