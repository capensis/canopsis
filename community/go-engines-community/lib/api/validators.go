package api

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/appinfo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgettemplate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libdatastorage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libidlerule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators(
	client mongo.DbClient,
	secConfig libsecurity.Config,
	enforcer libsecurity.Enforcer,
	tplExecutor libtemplate.Executor,
) error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("unknown validator engine")
	}

	libvalidator.RegisterTranslations(v)
	err := registerTagValidations(v, tplExecutor)
	if err != nil {
		return err
	}

	v.RegisterCustomTypeFunc(validation.ValidateCpsTimeType, datetime.CpsTime{})
	registerStructValidations(v, client, secConfig, enforcer, tplExecutor)

	return nil
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
	client mongo.DbClient,
	secConfig libsecurity.Config,
	enforcer libsecurity.Enforcer,
	tplExecutor libtemplate.Executor,
) {
	v.RegisterStructValidation(validation.ValidateFilteredQuery, pagination.FilteredQuery{})

	pbhValidator := pbehavior.NewValidator(client)
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhValidator.ValidateCreateRequest(sl)
	}, pbehavior.CreateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.UpdateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateEditRequest, pbehavior.EditRequest{})
	v.RegisterStructValidationCtx(pbhValidator.ValidatePatchRequest, pbehavior.PatchRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateCalendarRequest, pbehavior.CalendarByEntityIDRequest{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateEntityCreateRequest, pbehavior.BulkEntityCreateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateConnectorCreateRequest, pbehavior.BulkConnectorCreateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateConnectorEditRequest, pbehavior.BulkConnectorEditRequestItem{})

	v.RegisterStructValidation(exdate.ValidateExdateRequest, exdate.Request{})

	v.RegisterStructValidation(pbehaviortimespan.ValidateTimespansRequest, pbehaviortimespan.TimespansRequest{})

	scenarioExistReasonValidator := validation.NewExistFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Reason")
	scenarioExistTypeValidator := validation.NewExistFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Type")

	scenarioValidator := scenario.NewValidator(tplExecutor)
	v.RegisterStructValidation(scenarioValidator.ValidateActionRequest, scenario.ActionRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, action.Parameters{})

	entitybasicValidator := entitybasic.NewValidator(client)
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entitybasicValidator.ValidateEditRequest(ctx, sl)
	}, entitybasic.EditRequest{})

	entityserviceValidator := entityservice.NewValidator(client)
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceValidator.ValidateCreateRequest(sl)
	}, entityservice.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceValidator.ValidateUpdateRequest(sl)
	}, entityservice.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceValidator.ValidateUpdateRequest(sl)
	}, entityservice.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(entityserviceValidator.ValidateEditRequest, entityservice.EditRequest{})

	roleValidator := role.NewValidator(client)
	v.RegisterStructValidationCtx(roleValidator.ValidateEditRequest, role.EditRequest{})
	v.RegisterStructValidationCtx(roleValidator.ValidateBulkUpdatePermissionsRequestItem, role.BulkUpdatePermissionsRequestItem{})

	userValidator := user.NewValidator(client, secConfig)
	v.RegisterStructValidationCtx(userValidator.ValidateUpdateRequest, user.UpdateRequest{})
	v.RegisterStructValidationCtx(userValidator.ValidateCreateRequest, user.CreateRequest{})
	v.RegisterStructValidationCtx(userValidator.ValidatePatchRequest, user.PatchRequest{})
	v.RegisterStructValidationCtx(userValidator.ValidateBulkUpdateRequestItem, user.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(userValidator.ValidateBulkPatchRequestItem, user.BulkPatchRequestItem{})

	accountValidator := account.NewValidator(client)
	v.RegisterStructValidationCtx(accountValidator.ValidateEditRequest, account.EditRequest{})

	v.RegisterStructValidation(view.ValidateEditPositionRequest, view.EditPositionRequest{})

	viewGroupValidator := viewgroup.NewValidator(client)
	v.RegisterStructValidationCtx(viewGroupValidator.ValidateEditRequest, viewgroup.EditRequest{})

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

	eventfilterValidator := eventfilter.NewValidator(client, tplExecutor)
	v.RegisterStructValidationCtx(eventfilterValidator.ValidateEditRequest, eventfilter.EditRequest{})
	v.RegisterStructValidation(eventfilterValidator.ValidateTemplateRuleRequest, eventfilter.TemplateRuleRequest{})

	broadcastmessageValidator := broadcastmessage.NewValidator(client)
	v.RegisterStructValidationCtx(broadcastmessageValidator.Validate, broadcastmessage.CreateRequest{})

	v.RegisterStructValidation(messageratestats.ValidateListRequest, messageratestats.ListRequest{})

	idleRuleValidator := idlerule.NewValidator()
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleValidator.ValidateCreateRequest(sl)
	}, idlerule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleValidator.ValidateUpdateRequest(sl)
	}, idlerule.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, libidlerule.Parameters{})
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

	appInfoValidator := appinfo.NewValidator(client)
	v.RegisterStructValidationCtx(appInfoValidator.ValidateRequest, appinfo.UserInterfaceConf{})

	tplValidator := template.NewValidator(enforcer)
	v.RegisterStructValidation(tplValidator.ValidateEditDataRequest, template.EditDataRequest{})
	v.RegisterStructValidation(tplValidator.ValidateEditTestRequest, template.EditTestRequest{})
}
