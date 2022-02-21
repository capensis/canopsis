package cli

import "context"

const (
	collectionName     = "migration"
	timeFormat         = "20060102150405"
	fileNameDelimiter  = "_"
	fileExtJs          = ".js"
	fileNameSuffixUp   = fileNameDelimiter + "up" + fileExtJs
	fileNameSuffixDown = fileNameDelimiter + "down" + fileExtJs
	filePerm           = 0644
)

type Cmd interface {
	Exec(ctx context.Context) error
}
