package encoding

type DecodingError interface {
	error
	IsDecoding()
}

type decodingError struct {
	error
}

func (e decodingError) IsDecoding() {
}

func NewDecodingError(err error) DecodingError {
	if err == nil {
		return nil
	}
	return decodingError{
		error: err,
	}
}

type EncodingError interface {
	error
	IsEncoding()
}

type encodingError struct {
	error
}

func (e encodingError) IsEncoding() {
}

func NewEncodingError(err error) EncodingError {
	if err == nil {
		return nil
	}
	return encodingError{
		error: err,
	}
}
