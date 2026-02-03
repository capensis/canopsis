package externaldata

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/externaldata/externaldata.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata Getter

import (
	"context"
	"fmt"
)

type GetterContainer struct {
	getters map[string]Getter
}

func NewGetterContainer() *GetterContainer {
	return &GetterContainer{
		getters: make(map[string]Getter),
	}
}

func (c *GetterContainer) Get(dataType string) (Getter, bool) {
	g, ok := c.getters[dataType]

	return g, ok
}

func (c *GetterContainer) Set(dataType string, service Getter) {
	if c.Has(dataType) {
		panic(fmt.Errorf("data getter %q already exists", dataType))
	}

	c.getters[dataType] = service
}

func (c *GetterContainer) Has(dataType string) bool {
	_, ok := c.getters[dataType]

	return ok
}

type Getter interface {
	Get(ctx context.Context, parameters ParsedRefParameters, templateParameters any) (any, error)
}

type GetterError struct {
	failReason      string
	isParamsInvalid bool
	err             error
}

func NewGetterError(err error, failReason string, isParamsInvalid bool) error {
	return &GetterError{err: err, failReason: failReason, isParamsInvalid: isParamsInvalid}
}

func (e *GetterError) Error() string {
	return e.err.Error()
}

func (e *GetterError) FailReason() string {
	return e.failReason
}

func (e *GetterError) IsParamsInvalid() bool {
	return e.isParamsInvalid
}

func (e *GetterError) Unwrap() error {
	return e.err
}

type GetterTplError struct {
	GetterError
}

func NewGetterTplError(err error, failReason string, isParamsInvalid bool) error {
	return &GetterTplError{GetterError{err: err, failReason: failReason, isParamsInvalid: isParamsInvalid}}
}
