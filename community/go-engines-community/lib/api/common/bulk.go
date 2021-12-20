package common

import (
	"github.com/valyala/fastjson"
)

func GetBulkResponseItem(ar *fastjson.Arena, id string, status int, rawUser, error *fastjson.Value) *fastjson.Value {
	item := ar.NewObject()
	item.Set("status", ar.NewNumberInt(status))
	item.Set("item", rawUser)

	if error == nil {
		item.Set("id", ar.NewString(id))
		return item
	}

	if error.Type() == fastjson.TypeString {
		item.Set("error", error)
	}

	if error.Type() == fastjson.TypeObject {
		item.Set("errors", error)
	}

	return item
}
