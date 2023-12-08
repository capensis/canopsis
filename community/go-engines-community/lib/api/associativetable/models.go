package associativetable

import (
	"encoding/json"

	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

const bsonKey = "val"

type GetRequest struct {
	Name string `form:"name" binding:"required"`
}

type AssociativeTable struct {
	Name    string  `json:"name" bson:"name" binding:"required"`
	Content Content `json:"content" bson:"content"  binding:"required"`
}

type Content struct {
	value interface{}
}

func (c *Content) UnmarshalJSON(b []byte) error {
	var v interface{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	c.value = v
	return nil
}

func (c Content) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

// MarshalBSONValue stores value to map because it's impossible to decode struct and array of struct
// to interface without bson.D and bson.D cannot be encoded to JSON properly.
// Try to use interface{} in mongo-driver > 1.3.7
func (c Content) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return mongobson.MarshalValue(map[string]interface{}{
		bsonKey: c.value,
	})
}

func (c *Content) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	var v map[string]interface{}
	err := mongobson.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	c.value = v[bsonKey]

	return nil
}
