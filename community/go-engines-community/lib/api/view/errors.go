package view

type ValidationError struct {
	field string
	error error
}

func (v ValidationError) Error() string {
	return v.error.Error()
}
