package msgpack

import (
	"bytes"
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"github.com/vmihailenco/msgpack/v4"
)

type msgpEncoder struct {
}

type msgpDecoder struct {
}

// NewEncoder for msgpack
func NewEncoder() encoding.Encoder {
	return &msgpEncoder{}
}

// NewDecoder for msgpack
func NewDecoder() encoding.Decoder {
	return &msgpDecoder{}
}

func (e *msgpEncoder) Encode(in interface{}) ([]byte, error) {
	var writer bytes.Buffer
	encoder := msgpack.NewEncoder(&writer)
	err := encoder.Encode(in)
	if err != nil {
		var empty []byte
		return empty, encoding.NewEncodingError(fmt.Errorf("msgpack encoder: %v", err))
	}
	return writer.Bytes(), nil
}

func (d *msgpDecoder) Decode(in []byte, out interface{}) error {
	decoder := msgpack.NewDecoder(bytes.NewReader(in))
	if err := decoder.Decode(out); err != nil {
		return encoding.NewDecodingError(fmt.Errorf("msgpack decoder: %v", err))
	}
	return nil
}
