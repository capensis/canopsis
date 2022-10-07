package techmetrics

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
)

type API interface {
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
}

type api struct {
	taskExecutor TaskExecutor
}

func NewApi(
	taskExecutor TaskExecutor,
) API {
	return &api{
		taskExecutor: taskExecutor,
	}
}

// StartExport
// @Success 200 {object} ExportResponse
func (a *api) StartExport(c *gin.Context) {
	task, err := a.taskExecutor.StartExecute(c)
	if err != nil {
		panic(err)
	}

	if task.ID == 0 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Already in progress"})
		return
	}

	created := types.NewCpsTime(task.Created.Unix())
	c.JSON(http.StatusCreated, ExportResponse{
		Status:  task.Status,
		Created: &created,
	})
}

// GetExport
// @Success 200 {object} ExportResponse
func (a *api) GetExport(c *gin.Context) {
	task, err := a.taskExecutor.GetStatus(c)
	if err != nil {
		panic(err)
	}

	if task.ID == 0 {
		c.JSON(http.StatusOK, ExportResponse{
			Status: TaskStatusNone,
		})
		return
	}

	created := types.NewCpsTime(task.Created.Unix())
	c.JSON(http.StatusOK, ExportResponse{
		Status:  task.Status,
		Created: &created,
	})
}

func (a *api) DownloadExport(c *gin.Context) {
	task, err := a.taskExecutor.GetStatus(c)
	if err != nil {
		panic(err)
	}

	if task.ID == 0 || task.Status != TaskStatusSucceeded || task.Filepath == "" {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.FileAttachment(task.Filepath, "cps_tech_metrics_"+task.Started.Format("2006-01-02T15-04-05")+".bak")
}
