package widget

type ValidationError struct {
	error error
}

func (v ValidationError) Error() string {
	return v.error.Error()
}
