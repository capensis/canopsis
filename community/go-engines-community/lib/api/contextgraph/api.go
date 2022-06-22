package contextgraph

import (
	"context"
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
	ImportAll(c *gin.Context)
	ImportPartial(c *gin.Context)
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
// @Summary Create import task
// @Description Create import task
// @Tags contextgraph-import
// @ID contextgraph-import-create-import
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param source query string true "source"
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /contextgraph/import [put]
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

// ImportPartial
// @Summary Create import task
// @Description Create import task
// @Tags contextgraph-import
// @ID contextgraph-import-create-import-partial
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param source query string true "source"
// @Param body body ImportRequest true "body"
// @Success 200 {object} ImportResponse
// @Failure 400 {object} common.ErrorResponse
// @Router /contextgraph/import-partial [put]
func (a *api) ImportPartial(c *gin.Context) {
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

	err = ioutil.WriteFile(fmt.Sprintf(a.filePattern, job.ID), raw, os.ModePerm)
	if err != nil {
		return "", err
	}

	a.jobQueue.Push(job)
	return job.ID, nil
}

// Get import status by id
// @Summary Get import status by id
// @Description Get import status by id
// @Tags contextgraph-import
// @ID contextgraph-import-get-status-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "import id"
// @Success 200 {object} ImportJob
// @Failure 404 {object} common.ErrorResponse
// @Router /contextgraph/import/status/{id} [get]
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
