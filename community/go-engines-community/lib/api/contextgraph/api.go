package contextgraph

import (
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type API interface {
	Import(c *gin.Context)
	Status(c *gin.Context)
}

type Graph struct {
	Impact  *[]string
	Depends *[]string
}

type api struct {
	jobQueue    JobQueue
	reporter    StatusReporter
	filePattern string
	logger      zerolog.Logger
}

func NewApi(
	conf config.CanopsisConf,
	jobQueue JobQueue,
	reporter StatusReporter,
	logger zerolog.Logger,
) API {
	a := &api{
		jobQueue:    jobQueue,
		filePattern: conf.ImportCtx.FilePattern,
		reporter:    reporter,
		logger:      logger,
	}

	return a
}

// Import
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
func (a *api) Import(c *gin.Context) {
	query := ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := ImportJob{
		Creation: time.Now(),
		Status:   statusPending,
		Source:   query.Source,
	}

	err := a.reporter.ReportCreate(c.Request.Context(), &job)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf(a.filePattern, job.ID), raw, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	a.jobQueue.Push(job)

	c.JSON(http.StatusOK, ImportResponse{ID: job.ID})
}

// Status
// @Success 200 {object} ImportJob
func (a *api) Status(c *gin.Context) {
	status, err := a.reporter.GetStatus(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, status)
}
