package healthcheck

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
)

type API interface {
	Get(c *gin.Context)
	IsLive(c *gin.Context)
	GetStatus(c *gin.Context)
	GetEnginesOrder(c *gin.Context)
	GetParameters(c *gin.Context)
	UpdateParameters(c *gin.Context)
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
// @Success 200 {object} Info
func (a *api) Get(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.GetInfo())
}

// IsLive
// @Success 200 {object} LiveResponse
func (a *api) IsLive(c *gin.Context) {
	s := a.store.GetStatus()
	c.JSON(http.StatusOK, LiveResponse{
		Ok: len(s.Services) == 0 && len(s.Engines) == 0 && !s.HasInvalidEnginesOrder,
	})
}

// GetStatus
// @Success 200 {object} Status
func (a *api) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.GetStatus())
}

// GetEnginesOrder
// @Success 200 {object} Graph
func (a *api) GetEnginesOrder(c *gin.Context) {
	c.JSON(http.StatusOK, a.store.GetEnginesOrder())
}

// GetParameters
// @Success 200 {object} config.HealthCheckParameters
func (a *api) GetParameters(c *gin.Context) {
	data, err := a.store.GetParameters(c)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}

// UpdateParameters
// @Param body body config.HealthCheckParameters true "body"
// @Success 200 {object} config.HealthCheckParameters
func (a *api) UpdateParameters(c *gin.Context) {
	conf := config.HealthCheckParameters{}
	if err := c.ShouldBind(&conf); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, conf))
		return
	}

	data, err := a.store.UpdateParameters(c, conf)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, data)
}
