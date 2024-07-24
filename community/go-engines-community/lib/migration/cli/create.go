package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
	fileUp := filename + fileNameSuffixUp
	fileDown := filename + fileNameSuffixDown
	err := os.WriteFile(filepath.Join(c.path, fileUp), []byte(mongo.MigrationHelperComment), filePerm)
	if err != nil {
		return fmt.Errorf("cannot create up migration file %q: %w", fileUp, err)
	}

	c.logger.Info().Str("filename", fileUp).Msg("up migration script created")

	err = os.WriteFile(filepath.Join(c.path, fileDown), []byte(mongo.MigrationHelperComment), filePerm)
	if err != nil {
		return fmt.Errorf("cannot create down migration file %q: %w", fileDown, err)
	}

	c.logger.Info().Str("filename", fileDown).Msg("down migration script created")

	return nil
}
