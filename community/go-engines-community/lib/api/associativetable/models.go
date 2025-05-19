package associativetable

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
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
func (c Content) MarshalBSONValue() (byte, []byte, error) {
	bsonType, bsonBytes, err := bson.MarshalValue(map[string]any{
		bsonKey: c.value,
	})
	return byte(bsonType), bsonBytes, err
}

func (c *Content) UnmarshalBSONValue(_ byte, b []byte) error {
	var v map[string]interface{}
	err := bson.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	c.value = v[bsonKey]

	return nil
}
