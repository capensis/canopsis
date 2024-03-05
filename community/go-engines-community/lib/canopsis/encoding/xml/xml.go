package xml

import (
	"encoding/xml"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
)

type xmlEncoder struct {
}

type xmlDecoder struct {
}

func NewEncoder() encoding.Encoder {
	return &xmlEncoder{}
}

func (e *xmlEncoder) Encode(in interface{}) ([]byte, error) {
	b, err := xml.Marshal(in)
	if err != nil {
		return []byte{}, encoding.NewEncodingError(fmt.Errorf("xml encoder: %w", err))
	}
	return b, nil
}

func NewDecoder() encoding.Decoder {
	return &xmlDecoder{}
}

func (e *xmlDecoder) Decode(in []byte, out interface{}) error {
	if err := xml.Unmarshal(in, out); err != nil {
		return encoding.NewDecodingError(fmt.Errorf("xml decoder: %w", err))
	}
	return nil
}
