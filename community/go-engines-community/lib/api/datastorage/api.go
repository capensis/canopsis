package datastorage

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/datastorage"
	"github.com/gin-gonic/gin"
	"net/http"
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

// Get conf
// @Summary Get conf
// @Description Get conf
// @Tags datastorage
// @ID datastorage-get
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 200 {object} DataStorage
// @Router /data-storage [get]
func (a *api) Get(c *gin.Context) {
	data, err := a.store.Get(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}

// Update conf
// @Summary Update conf
// @Description Update conf
// @Tags datastorage
// @ID datastorage-update
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body datastorage.Config true "body"
// @Success 200 {object} DataStorage
// @Router /data-storage [put]
func (a *api) Update(c *gin.Context) {
	conf := datastorage.Config{}
	if err := c.ShouldBind(&conf); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, conf))
		return
	}

	data, err := a.store.Update(c.Request.Context(), conf)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}
