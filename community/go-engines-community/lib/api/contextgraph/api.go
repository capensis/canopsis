package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	ImportAll(c *gin.Context)
	ImportOldAll(c *gin.Context)
	ImportOldPartial(c *gin.Context)
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

// ImportAll
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
func (a *api) ImportAll(c *gin.Context) {
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

	raw, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	jobID, err := a.createImportJob(c.Request.Context(), job, raw)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, ImportResponse{ID: jobID})
}

// ImportOldAll
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
func (a *api) ImportOldAll(c *gin.Context) {
	query := ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := ImportJob{
		Creation: time.Now(),
		Status:   statusPending,
		Source:   query.Source,
		IsOld:    true,
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

	c.JSON(http.StatusOK, ImportResponse{ID: jobID})
}

// ImportOldPartial
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
func (a *api) ImportOldPartial(c *gin.Context) {
	query := ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := ImportJob{
		Creation:  time.Now(),
		Status:    statusPending,
		Source:    query.Source,
		IsPartial: true,
		IsOld:     true,
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

	c.JSON(http.StatusOK, ImportResponse{ID: jobID})
}

func (a *api) createImportJob(ctx context.Context, job ImportJob, raw []byte) (string, error) {
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
