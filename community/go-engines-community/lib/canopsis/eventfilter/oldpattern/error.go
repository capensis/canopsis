package oldpattern

type UnexpectedFieldsError struct {
	Err error
}

func (e UnexpectedFieldsError) Error() string {
	if e.Err == nil {
		return "unknown"
	}
	return e.Err.Error()
}
