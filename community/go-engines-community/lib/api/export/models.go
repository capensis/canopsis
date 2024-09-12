package export

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

const (
	TaskStatusRunning = iota
	TaskStatusSucceeded
	TaskStatusFailed
)

type TaskParameters struct {
	Type           string
	Parameters     string
	Fields         Fields
	Separator      rune
	FilenamePrefix string
	UserID         string
}

type Task struct {
	ID         string         `bson:"_id"`
	Status     int64          `bson:"status"`
	Type       string         `bson:"type"`
	Parameters string         `bson:"parameters"`
	Fields     Fields         `bson:"fields"`
	Separator  rune           `bson:"separator"`
	File       string         `bson:"file,omitempty"`
	Filename   string         `bson:"filename"`
	FailReason string         `bson:"fail_reason,omitempty"`
	User       string         `bson:"user"`
	Created    types.CpsTime  `bson:"created"`
	Launched   *types.CpsTime `bson:"launched,omitempty"`
	Completed  *types.CpsTime `bson:"completed,omitempty"`
}

type Fields []Field

type Field struct {
	Name     string `bson:"name" json:"name"`
	Label    string `bson:"label" json:"label"`
	Template string `bson:"template" json:"template"`
}

func (f *Fields) Fields() []string {
	fields := make([]string, 0, len(*f))
	for _, field := range *f {
		if field.Name != "" {
			fields = append(fields, field.Name)
		}
	}

	return fields
}

func (f *Fields) Labels() []string {
	labels := make([]string, len(*f))
	for i, field := range *f {
		labels[i] = field.Label
	}

	return labels
}
