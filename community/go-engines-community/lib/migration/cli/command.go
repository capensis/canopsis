package cli

import "context"

const (
	collectionName     = "migration"
	timeFormat         = "20060102150405"
	fileNameDelimiter  = "_"
	fileNameSuffixUp   = "up"
	fileNameSuffixDown = "down"
	fileExtJs          = ".js"
	filePerm           = 0644
)

type Cmd interface {
	Exec(ctx context.Context) error
}
