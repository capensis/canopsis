package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/workers"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const (
	dirPerm     = 0o770
	filePerm    = 0o640
	filePattern = "import_%s.json"
)

type api struct {
	reporter       StatusReporter
	jobPublisher   workers.JobPublisher
	dir            string
	filePattern    string
	maxImportSize  uint64
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

func NewApi(
	conf config.CanopsisConf,
	reporter StatusReporter,
	jobPublisher workers.JobPublisher,
	maxImportSize uint64,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		dir:            filepath.Join(conf.File.Dir, canopsis.SubDirImport),
		filePattern:    filePattern,
		reporter:       reporter,
		jobPublisher:   jobPublisher,
		maxImportSize:  maxImportSize,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

// ImportAll
// @Param body body []importcontextgraph.EntityConfiguration true "body"
// @Success 200 {object} contextgraph.ImportResponse
func (a *api) ImportAll(c *gin.Context) {
	query := ImportQuery{}
	if err := validation.BindQuery(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	job := ImportJob{
		Creation: time.Now(),
		Status:   StatusPending,
		Source:   query.Source,
	}

	raw, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	if a.maxImportSize > 0 && uint64(len(raw)) > a.maxImportSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("request size %d exceeds limit %d", len(raw), a.maxImportSize)))
		return
	}

	jobID, err := a.createImportJob(c, job, raw)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, ImportResponse{ID: jobID})
}

// ImportPartial
// @Param body body []importcontextgraph.EntityConfiguration true "body"
// @Success 200 {object} contextgraph.ImportResponse
func (a *api) ImportPartial(c *gin.Context) {
	query := ImportQuery{}
	if err := validation.BindQuery(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

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
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	if a.maxImportSize > 0 && uint64(len(raw)) > a.maxImportSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("request size %d exceeds limit %d", len(raw), a.maxImportSize)))
		return
	}

	jobID, err := a.createImportJob(c, job, raw)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, ImportResponse{ID: jobID})
}

func (a *api) createImportJob(ctx context.Context, job ImportJob, raw []byte) (string, error) {
	err := a.reporter.ReportCreate(ctx, &job)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(a.dir, os.ModeDir|dirPerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filepath.Join(a.dir, fmt.Sprintf(a.filePattern, job.ID)), raw, filePerm)
	if err != nil {
		return "", err
	}

	err = a.jobPublisher.Publish(ctx, "")
	if err != nil {
		return "", err
	}

	return job.ID, nil
}

// Status
// @Success 200 {object} contextgraph.ImportJob
func (a *api) Status(c *gin.Context) {
	status, err := a.reporter.GetStatus(c, c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, status)
}
