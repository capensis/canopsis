package json

import (
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	jsoniter "github.com/json-iterator/go"
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
	b, err := e.jsonAPI.Marshal(in)
	if err != nil {
		return []byte{}, encoding.NewEncodingError(fmt.Errorf("json encoder: %v", err))
	}
	return b, nil
}

func NewDecoder() encoding.Decoder {
	return &jsonDecoder{
		jsonAPI: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}

func (e *jsonDecoder) Decode(in []byte, out interface{}) error {
	if err := e.jsonAPI.Unmarshal(in, out); err != nil {
		return encoding.NewDecodingError(fmt.Errorf("json decoder: %v", err))
	}
	return nil
}
