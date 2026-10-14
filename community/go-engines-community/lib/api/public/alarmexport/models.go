package alarmexport

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/dataexport"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
)

type Job struct {
	asyncjob.Base `bson:",inline"`
	Params        Request `json:"params"`
}

func (j Job) DownloadName() string {
	return "alarms_" + j.Created.Time.Format(dataexport.DownloadNameTimeLayout) + "." + j.Params.Format
}

func (j Job) MimeType() (string, string) {
	switch j.Params.Format {
	case "csv":
		return response.MIMECSV, j.Params.Format
	case "json":
		return response.MIMEJSON, j.Params.Format
	default:
		return "", ""
	}
}

type Request struct {
	alarm.SearchRequest `bson:",inline"`
	Format              string             `json:"format" bson:"format" binding:"required,oneof=csv json"`
	Fields              []dataexport.Field `json:"fields" bson:"fields" binding:"required"`
	TimeFormat          string             `json:"time_format" bson:"time_format" binding:"time_format"`
	Author              string             `json:"author" bson:"-" swaggerignore:"true"`

	// csv
	Delimiter string `json:"delimiter" bson:"delimiter" binding:"oneoforempty=comma semicolon tab space"`
}
