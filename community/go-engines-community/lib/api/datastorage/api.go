package datastorage

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"github.com/gin-gonic/gin"
)

type API interface {
	Get(c *gin.Context)

	Update(c *gin.Context)
}

func NewApi(store Store) API {
	return &api{
		store: store,
	}
}

type api struct {
	store Store
}

// Get
// @Success 200 {object} DataStorage
func (a *api) Get(c *gin.Context) {
	data, err := a.store.Get(c)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}

// Update
// @Param body body datastorage.Config true "body"
// @Success 200 {object} DataStorage
func (a *api) Update(c *gin.Context) {
	conf := datastorage.Config{}
	if err := c.ShouldBind(&conf); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, conf))
		return
	}

	data, err := a.store.Update(c, conf)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}
