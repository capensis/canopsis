package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func NewCreateCmd(
	path, name string,
	logger zerolog.Logger,
) Cmd {
	return &createCmd{
		path:   path,
		name:   name,
		logger: logger,
	}
}

type createCmd struct {
	path, name string
	logger     zerolog.Logger
}

func (c *createCmd) Exec(_ context.Context) error {
	now := time.Now()
	filename := strings.Join([]string{now.Format(timeFormat), c.name}, fileNameDelimiter)
	fileUp := fmt.Sprintf("%s%s%s%s", filename, fileNameDelimiter, fileNameSuffixUp, fileExtJs)
	fileDown := fmt.Sprintf("%s%s%s%s", filename, fileNameDelimiter, fileNameSuffixDown, fileExtJs)
	err := ioutil.WriteFile(filepath.Join(c.path, fileUp), nil, filePerm)
	if err != nil {
		return fmt.Errorf("cannot create up migration file %q: %w", fileUp, err)
	}

	c.logger.Info().Str("filename", fileUp).Msg("up migration script created")

	err = ioutil.WriteFile(filepath.Join(c.path, fileDown), nil, filePerm)
	if err != nil {
		return fmt.Errorf("cannot create down migration file %q: %w", fileDown, err)
	}

	c.logger.Info().Str("filename", fileDown).Msg("down migration script created")

	return nil
}
