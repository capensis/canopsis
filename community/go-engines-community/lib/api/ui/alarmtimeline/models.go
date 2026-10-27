package alarmtimeline

const ExportType = "timeline"

type ExportRequest struct {
	ExportFetchParameters
	Separator  string `json:"separator" binding:"oneoforempty=comma semicolon tab space"`
	TimeFormat string `json:"time_format" binding:"time_format"`
}

type ExportFetchParameters struct {
	AlarmID string `json:"alarm" binding:"required"`
}

type ExportResponse struct {
	ID string `json:"_id"`

	//   * `0` - Running
	//   * `1` - Succeeded
	//   * `2` - Failed
	Status int64 `json:"status"`
}
