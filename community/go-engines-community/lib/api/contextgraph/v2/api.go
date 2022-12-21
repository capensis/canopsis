package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/contextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type api struct {
	jobQueue    contextgraph.JobQueue
	reporter    contextgraph.StatusReporter
	filePattern string
	logger      zerolog.Logger
}

func NewApi(
	conf config.CanopsisConf,
	jobQueue contextgraph.JobQueue,
	reporter contextgraph.StatusReporter,
	logger zerolog.Logger,
) contextgraph.API {
	a := &api{
		jobQueue:    jobQueue,
		filePattern: conf.ImportCtx.FilePattern,
		reporter:    reporter,
		logger:      logger,
	}

	return a
}

// ImportAll
// @Param body body []importcontextgraph.EntityConfiguration true "body"
// @Success 200 {object} contextgraph.ImportResponse
func (a *api) ImportAll(c *gin.Context) {
	query := contextgraph.ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := contextgraph.ImportJob{
		Creation: time.Now(),
		Status:   contextgraph.StatusPending,
		Source:   query.Source,
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	jobID, err := a.createImportJob(c.Request.Context(), job, raw)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, contextgraph.ImportResponse{ID: jobID})
}

// ImportPartial
// @Param body body []importcontextgraph.EntityConfiguration true "body"
// @Success 200 {object} contextgraph.ImportResponse
func (a *api) ImportPartial(c *gin.Context) {
	query := contextgraph.ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := contextgraph.ImportJob{
		Creation:  time.Now(),
		Status:    contextgraph.StatusPending,
		Source:    query.Source,
		IsPartial: true,
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	jobID, err := a.createImportJob(c.Request.Context(), job, raw)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, contextgraph.ImportResponse{ID: jobID})
}

func (a *api) createImportJob(ctx context.Context, job contextgraph.ImportJob, raw []byte) (string, error) {
	err := a.reporter.ReportCreate(ctx, &job)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(fmt.Sprintf(a.filePattern, job.ID), raw, os.ModePerm)
	if err != nil {
		return "", err
	}

	a.jobQueue.Push(job)
	return job.ID, nil
}

// Status
// @Success 200 {object} contextgraph.ImportJob
func (a *api) Status(c *gin.Context) {
	status, err := a.reporter.GetStatus(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, contextgraph.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, status)
}
