package encoding

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
)

type DecodingError interface {
	errt.ErrT
	IsDecoding()
}

type decodingError struct {
	errt.ErrT
}

func (e decodingError) IsDecoding() {
}

func NewDecodingError(err error) DecodingError {
	if err == nil {
		return nil
	}
	return decodingError{
		ErrT: errt.NewErrT(err),
	}
}

type EncodingError interface {
	errt.ErrT
	IsEncoding()
}

type encodingError struct {
	errt.ErrT
}

func (e encodingError) IsEncoding() {
}

func NewEncodingError(err error) EncodingError {
	if err == nil {
		return nil
	}
	return encodingError{
		ErrT: errt.NewErrT(err),
	}
}
