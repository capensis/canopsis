package binding

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Bind(c *gin.Context, v any) error {
	if err := c.ShouldBind(v); err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}

func BindQuery(c *gin.Context, v any) error {
	if err := c.ShouldBindQuery(v); err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}

func ValidateStruct(v any) error {
	if err := binding.Validator.ValidateStruct(v); err != nil {
		return validation.ConvertValErr(err, v)
	}

	return nil
}
