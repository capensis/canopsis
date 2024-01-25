package techmetrics

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/gin-gonic/gin"
)

type API interface {
	StartExport(c *gin.Context)
	GetExport(c *gin.Context)
	DownloadExport(c *gin.Context)
}

type api struct {
	taskExecutor TaskExecutor

	timezoneConfigProvider config.TimezoneConfigProvider
}

func NewApi(
	taskExecutor TaskExecutor,
	timezoneConfigProvider config.TimezoneConfigProvider,
) API {
	return &api{
		taskExecutor: taskExecutor,

		timezoneConfigProvider: timezoneConfigProvider,
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

	created := datetime.NewCpsTime(task.Created.Unix())
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
			Status: task.Status,
		})
		return
	}

	created := datetime.NewCpsTime(task.Created.Unix())
	response := ExportResponse{
		Status:  task.Status,
		Created: &created,
	}
	if task.Status == TaskStatusSucceeded && task.Completed != nil {
		d := int(task.Completed.Sub(task.Created).Seconds())
		response.Duration = &d
	}
	c.JSON(http.StatusOK, response)
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

	location := a.timezoneConfigProvider.Get().Location
	filename := "cps_tech_metrics_" + task.Started.In(location).Format("2006-01-02T15-04-05-MST") + ".bak"
	c.FileAttachment(task.Filepath, filename)
}
