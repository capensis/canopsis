package gob

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
)

type gobEncoder struct {
}

type gobDecoder struct {
}

// NewEncoder for GOB
func NewEncoder() encoding.Encoder {
	return &gobEncoder{}
}

func (e *gobEncoder) Encode(in interface{}) ([]byte, error) {
	var writer bytes.Buffer
	encoder := gob.NewEncoder(&writer)

	if err := encoder.Encode(in); err != nil {
		var empty []byte
		return empty, encoding.NewEncodingError(fmt.Errorf("gob encoder: %v", err))
	}

	return writer.Bytes(), nil
}

// NewDecoder for GOB
func NewDecoder() encoding.Decoder {
	return &gobDecoder{}
}

func (d *gobDecoder) Decode(in []byte, out interface{}) error {
	reader := bytes.NewReader(in)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(out)

	if err != nil {
		return encoding.NewDecodingError(fmt.Errorf("gob decoder: %v", err))
	}

	return nil
}
