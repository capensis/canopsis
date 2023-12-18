package json

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	jsoniter "github.com/json-iterator/go"
	"github.com/mailru/easyjson"
)

type jsonEncoder struct {
	jsonAPI jsoniter.API
}

type jsonDecoder struct {
	jsonAPI jsoniter.API
}

func NewEncoder() encoding.Encoder {
	return &jsonEncoder{
		jsonAPI: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (e *jsonEncoder) Encode(in interface{}) ([]byte, error) {
	var b []byte
	var err error

	if m, ok := in.(easyjson.Marshaler); ok {
		b, err = easyjson.Marshal(m)
	} else {
		b, err = e.jsonAPI.Marshal(in)
	}

	if err != nil {
		return []byte{}, encoding.NewEncodingError(fmt.Errorf("json encoder: %w", err))
	}

	return b, nil
}

func NewDecoder() encoding.Decoder {
	return &jsonDecoder{
		jsonAPI: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (e *jsonDecoder) Decode(in []byte, out interface{}) error {
	var err error

	if m, ok := out.(easyjson.Unmarshaler); ok {
		err = easyjson.Unmarshal(in, m)
	} else {
		err = e.jsonAPI.Unmarshal(in, out)
	}

	if err != nil {
		return encoding.NewDecodingError(fmt.Errorf("json decoder: %w", err))
	}

	return nil
}
