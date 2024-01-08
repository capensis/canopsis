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

type api struct {
	reporter    StatusReporter
	filePattern string
	logger      zerolog.Logger
}

func NewApi(
	conf config.CanopsisConf,
	reporter StatusReporter,
	logger zerolog.Logger,
) API {
	a := &api{
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
	query := ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := ImportJob{
		Creation: time.Now(),
		Status:   StatusPending,
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
// @Param body body []importcontextgraph.EntityConfiguration true "body"
// @Success 200 {object} contextgraph.ImportResponse
func (a *api) ImportPartial(c *gin.Context) {
	query := ImportQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	job := ImportJob{
		Creation:  time.Now(),
		Status:    StatusPending,
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

	err = os.WriteFile(fmt.Sprintf(a.filePattern, job.ID), raw, os.ModePerm)
	if err != nil {
		return "", err
	}

	return job.ID, nil
}

// Status
// @Success 200 {object} contextgraph.ImportJob
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
