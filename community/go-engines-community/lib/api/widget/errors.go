package widget

type ValidationErr struct {
	error error
}

func (v ValidationErr) Error() string {
	return v.error.Error()
}
