package encoding

//go:generate mockgen -destination=../../../mocks/lib/canopsis/encoding/encoding.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding Encoder,Decoder

// Encoder interface ease the use of encoders like encoding/gob...
type Encoder interface {
	// Encode in to []byte.
	// Example:
	//
	// var in MyTime
	// bytes, err := encoder.Encode(in)
	Encode(in interface{}) ([]byte, error)
}

// Decoder interface ease the use of decoders like encoding/gob...
type Decoder interface {
	// Decode the given byte in to the given out type.
	// Example:
	//
	// var out MyType
	// err := decoder.Decode(data, &out)
	Decode(in []byte, out interface{}) error
}
