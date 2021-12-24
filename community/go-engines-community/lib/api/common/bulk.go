package common

import (
	"github.com/valyala/fastjson"
)

func GetBulkResponseItem(ar *fastjson.Arena, id string, status int, rawUser, err *fastjson.Value) *fastjson.Value {
	item := ar.NewObject()
	item.Set("status", ar.NewNumberInt(status))
	item.Set("item", rawUser)

	if err == nil {
		item.Set("id", ar.NewString(id))
		return item
	}

	if err.Type() == fastjson.TypeString {
		item.Set("error", err)
	}

	if err.Type() == fastjson.TypeObject {
		item.Set("errors", err)
	}

	return item
}
