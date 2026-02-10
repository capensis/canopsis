package api

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/healthcheck"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/linkrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/messageratestats"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/webhook"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgettemplate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libdatastorage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	libwebhook "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators(secConfig libsecurity.Config, tplExecutor libtemplate.Executor) (*ut.UniversalTranslator, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("unknown validator engine")
	}

	trans, err := libvalidator.NewTranslator(v, validation.InvalidIDChars)
	if err != nil {
		return nil, err
	}

	err = registerTagValidations(v, tplExecutor)
	if err != nil {
		return nil, err
	}

	v.RegisterCustomTypeFunc(validation.ValidateCpsTimeType, datetime.CpsTime{})
	registerStructValidations(v, secConfig, tplExecutor)

	return trans, nil
}

func registerTagValidations(v *validator.Validate, tplExecutor libtemplate.Executor) error {
	var tagValidations = map[string]validator.Func{
		"notblank":                validators.NotBlank,
		"oneoforempty":            validation.ValidateOneOfOrEmpty,
		"iscolororempty":          validation.ValidateColorOrEmpty,
		"id":                      validation.ValidateID,
		"time_format":             validation.ValidateTimeFormat,
		"info_value":              validation.ValidateInfoValue,
		"table_name":              validation.ValidateTableName,
		"event_pattern":           patternfields.ValidateEventPattern,
		"alarm_pattern":           patternfields.ValidateAlarmPattern,
		"entity_pattern":          patternfields.ValidateEntityPattern,
		"pbehavior_pattern":       patternfields.ValidatePbehaviorPattern,
		"weather_service_pattern": patternfields.ValidateWeatherServicePattern,
		"template":                validation.ValidateTemplate(tplExecutor),
	}
	for tag, f := range tagValidations {
		err := v.RegisterValidation(tag, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerStructValidations(
	v *validator.Validate,
	secConfig libsecurity.Config,
	tplExecutor libtemplate.Executor,
) {
	v.RegisterStructValidation(validation.ValidateFilteredQuery, pagination.FilteredQuery{})

	pbhValidator := pbehavior.NewValidator()
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhValidator.ValidateCreateRequest(sl)
	}, pbehavior.CreateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.UpdateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.BulkUpdateRequestItem{})
	v.RegisterStructValidation(pbhValidator.ValidateEditRequest, pbehavior.EditRequest{})
	v.RegisterStructValidation(pbhValidator.ValidatePatchRequest, pbehavior.PatchRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateCalendarRequest, pbehavior.CalendarByEntityIDRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateEntityCreateRequest, pbehavior.BulkEntityCreateRequestItem{})
	v.RegisterStructValidation(pbhValidator.ValidateConnectorCreateRequest, pbehavior.BulkConnectorCreateRequestItem{})
	v.RegisterStructValidation(pbhValidator.ValidateConnectorEditRequest, pbehavior.BulkConnectorEditRequestItem{})

	v.RegisterStructValidation(exdate.ValidateExdateRequest, exdate.Request{})

	v.RegisterStructValidation(pbehaviortimespan.ValidateTimespansRequest, pbehaviortimespan.TimespansRequest{})

	scenarioValidator := scenario.NewValidator(tplExecutor)
	v.RegisterStructValidation(scenarioValidator.ValidateActionRequest, scenario.ActionRequest{})

	v.RegisterStructValidation(entityservice.ValidateEditRequest, entityservice.EditRequest{})

	userValidator := user.NewValidator(secConfig)
	v.RegisterStructValidation(userValidator.ValidateUpdateRequest, user.UpdateRequest{})
	v.RegisterStructValidation(userValidator.ValidateCreateRequest, user.CreateRequest{})
	v.RegisterStructValidation(userValidator.ValidatePatchRequest, user.PatchRequest{})
	v.RegisterStructValidation(userValidator.ValidateBulkUpdateRequestItem, user.BulkUpdateRequestItem{})
	v.RegisterStructValidation(userValidator.ValidateBulkPatchRequestItem, user.BulkPatchRequestItem{})

	v.RegisterStructValidation(account.ValidateEditRequest, account.EditRequest{})

	v.RegisterStructValidation(view.ValidateEditPositionRequest, view.EditPositionRequest{})

	widgetValidator := widget.NewValidator(tplExecutor)
	v.RegisterStructValidation(widgetValidator.ValidateEditRequest, widget.EditRequest{})
	v.RegisterStructValidation(widgetValidator.ValidateFilterRequest, widget.FilterRequest{})

	v.RegisterStructValidation(widgetfilter.NewValidator().ValidateCreateRequest, widgetfilter.CreateRequest{})
	v.RegisterStructValidation(widgetfilter.NewValidator().ValidateUpdateRequest, widgetfilter.UpdateRequest{})

	v.RegisterStructValidation(widgettemplate.ValidateEditRequest, widgettemplate.EditRequest{})

	stateSettingsValidator := statesettings.NewValidator()
	v.RegisterStructValidation(stateSettingsValidator.ValidateEditRequest, statesettings.EditRequest{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateJUnitThresholds, statesettings.JUnitThreshold{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateThreshold, statesettings.StateThreshold{})

	eventfilterValidator := eventfilter.NewValidator(tplExecutor)
	v.RegisterStructValidationCtx(eventfilterValidator.ValidateEditRequest, eventfilter.EditRequest{})
	v.RegisterStructValidation(eventfilterValidator.ValidateTemplateRuleRequest, eventfilter.TemplateRuleRequest{})

	v.RegisterStructValidation(messageratestats.ValidateListRequest, messageratestats.ListRequest{})

	idleRuleValidator := idlerule.NewValidator()
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleValidator.ValidateCreateRequest(sl)
	}, idlerule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleValidator.ValidateUpdateRequest(sl)
	}, idlerule.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleValidator.ValidateBulkUpdateRequestItem(sl)
	}, idlerule.BulkUpdateRequestItem{})

	v.RegisterStructValidation(alarm.ValidateListRequest, alarm.ListRequest{})
	v.RegisterStructValidation(alarm.ValidateDetailsRequest, alarm.DetailsRequest{})

	v.RegisterStructValidation(datastorage.ValidateConfig, libdatastorage.Config{})

	resolveRuleValidator := resolverule.NewValidator()
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleValidator.ValidateCreateRequest(sl)
	}, resolverule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleValidator.ValidateUpdateRequest(sl)
	}, resolverule.UpdateRequest{})

	flappingRuleValidator := flappingrule.NewValidator()
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleValidator.ValidateCreateRequest(sl)
	}, flappingrule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleValidator.ValidateUpdateRequest(sl)
	}, flappingrule.UpdateRequest{})

	v.RegisterStructValidation(pattern.ValidateEditRequest, pattern.EditRequest{})

	linkRuleValidator := linkrule.NewValidator(tplExecutor)
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		linkRuleValidator.ValidateEditRequest(sl)
	}, linkrule.EditRequest{})
	v.RegisterStructValidation(linkRuleValidator.ValidateTemplateRequest, linkrule.TemplateRequest{})

	v.RegisterStructValidation(alarmtag.ValidateCreateRequest, alarmtag.CreateRequest{})
	v.RegisterStructValidation(alarmtag.ValidateUpdateRequest, alarmtag.UpdateRequest{})

	v.RegisterStructValidation(healthcheck.ValidateEngineParameters, config.EngineParameters{})
	v.RegisterStructValidation(healthcheck.ValidateLimitParameters, config.LimitParameters{})

	v.RegisterStructValidation(template.ValidateEditDataRequest, template.EditDataRequest{})
	v.RegisterStructValidation(template.ValidateEditTestRequest, template.EditTestRequest{})

	checkTicketStatusValidator := webhook.NewCheckTicketStatusValidator(tplExecutor)
	v.RegisterStructValidation(checkTicketStatusValidator.ValidateCheckTicketStatus, libwebhook.CheckTicketStatus{})
}
