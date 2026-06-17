package encoding

type DecodingError struct {
	Err error
}

func (e DecodingError) Unwrap() error {
	return e.Err
}

func (e DecodingError) Error() string {
	return e.Err.Error()
}

func NewDecodingError(err error) DecodingError {
	return DecodingError{
		Err: err,
	}
}

type EncodingError struct {
	Err error
}

func (e EncodingError) Unwrap() error {
	return e.Err
}

func (e EncodingError) Error() string {
	return e.Err.Error()
}

func NewEncodingError(err error) EncodingError {
	return EncodingError{
		Err: err,
	}
}
