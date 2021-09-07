package types

import (
	"strings"
	"text/template"
)

// Operation represents alarm modification operation.
type Operation struct {
	Type       string      `bson:"type"`
	Parameters interface{} `bson:"parameters,omitempty"`
}

type Templater interface {
	Template(data interface{}) error
}

// OperationParameters represents default operation parameters.
type OperationParameters struct {
	Output string `bson:"output" json:"output"`
	Author string `bson:"author" json:"author"`
}

func (p *OperationParameters) Template(data interface{}) error {
	output, err := renderTemplate(p.Output, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Output = output
	author, err := renderTemplate(p.Author, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Author = author

	return nil
}

type OperationAssocTicketParameters struct {
	Ticket string `bson:"ticket" json:"ticket"`
	Output string `bson:"output" json:"output"`
	Author string `bson:"author" json:"author"`
}

func (p *OperationAssocTicketParameters) Template(data interface{}) error {
	output, err := renderTemplate(p.Output, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Output = output
	author, err := renderTemplate(p.Author, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Author = author

	return nil
}

type OperationSnoozeParameters struct {
	Duration DurationWithUnit `bson:"duration" json:"duration"`
	Output   string           `bson:"output" json:"output"`
	Author   string           `bson:"author" json:"author"`
}

func (p *OperationSnoozeParameters) Template(data interface{}) error {
	output, err := renderTemplate(p.Output, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Output = output
	author, err := renderTemplate(p.Author, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Author = author

	return nil
}

type OperationChangeStateParameters struct {
	State  CpsNumber `bson:"state" json:"state"`
	Output string    `bson:"output" json:"output"`
	Author string    `bson:"author" json:"author"`
}

func (p *OperationChangeStateParameters) Template(data interface{}) error {
	output, err := renderTemplate(p.Output, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Output = output
	author, err := renderTemplate(p.Author, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Author = author

	return nil
}

type OperationDeclareTicketParameters struct {
	Ticket string            `bson:"ticket" json:"ticket"`
	Data   map[string]string `bson:"data" json:"data"`
	Output string            `bson:"output" json:"output"`
	Author string            `bson:"author" json:"author"`
}

func (p *OperationDeclareTicketParameters) Template(data interface{}) error {
	output, err := renderTemplate(p.Output, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Output = output
	author, err := renderTemplate(p.Author, data, GetTemplateFunc())
	if err != nil {
		return err
	}
	p.Author = author

	return nil
}

type OperationPbhParameters struct {
	PbehaviorInfo PbehaviorInfo `json:"pbehavior_info"`
	Output        string        `json:"output"`
	Author        string        `json:"author"`
}

type OperationInstructionParameters struct {
	Execution string `bson:"execution" json:"execution"`
	Output    string `bson:"output" json:"output"`
	Author    string `bson:"author" json:"author"`
}

func renderTemplate(templateStr string, data interface{}, f template.FuncMap) (string, error) {
	t, err := template.New("template").Funcs(f).Parse(templateStr)
	if err != nil {
		return "", err
	}
	b := strings.Builder{}
	err = t.Execute(&b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
