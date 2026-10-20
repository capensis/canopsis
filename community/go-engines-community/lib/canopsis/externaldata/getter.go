package externaldata

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/externaldata/externaldata.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata Getter

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
)

const (
	OnMismatchError = iota
	OnMismatchSkipRule
	OnMismatchSkipActions
)

type Getter interface {
	Get(ctx context.Context, rule Rule, tplData any) (GetterResult, error)
}

type GetterResult struct {
	ExternalData  map[string]any
	CorrelationID string
	WebhookEvent  *rpc.WebhookEvent
	RequestCount  map[string]int64
	SkipRule      bool
	SkippedRefs   map[string]bool
}

type Rule struct {
	ID           string
	Name         string
	ExternalData []ParsedRefParameters
	// OnMismatch overrides the default error outcome for each entry in ExternalData, in the same order.
	OnMismatch    []int
	TriggerUserID string
	TriggerAuthor string
}

func NewNullGetter() Getter {
	return &nullGetter{}
}

type nullGetter struct {
}

func (*nullGetter) Get(_ context.Context, _ Rule, _ any) (GetterResult, error) {
	return GetterResult{}, errors.New("null getter")
}

type GetterError struct {
	t               string
	failReason      string
	isParamsInvalid bool
	err             error
	notFound        bool
}

func NewGetterError(err error, t string, failReason string, isParamsInvalid, notFound bool) error {
	return &GetterError{
		err:             err,
		t:               t,
		failReason:      failReason,
		isParamsInvalid: isParamsInvalid,
		notFound:        notFound,
	}
}

func (e *GetterError) Error() string {
	return e.err.Error()
}

func (e *GetterError) Type() string {
	return e.t
}

func (e *GetterError) FailReason() string {
	return e.failReason
}

func (e *GetterError) IsParamsInvalid() bool {
	return e.isParamsInvalid
}

func (e *GetterError) IsNotFound() bool {
	return e.notFound
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
