package api

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/account"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/playlist"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/resolverule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/role"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators(client mongo.DbClient) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	// Translations
	libvalidator.RegisterTranslations(v)

	// Common validation rules
	err := v.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("oneoforempty", common.ValidateOneOfOrEmpty)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("iscolororempty", common.ValidateColorOrEmpty)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("id", common.ValidateID)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("time_format", common.ValidateTimeFormat)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("event_pattern", common.ValidateEventPattern)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("alarm_pattern", common.ValidateAlarmPattern)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("entity_pattern", common.ValidateEntityPattern)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("pbehavior_pattern", common.ValidatePbehaviorPattern)
	if err != nil {
		panic(err)
	}
	err = v.RegisterValidation("info_value", common.ValidateInfoValue)
	if err != nil {
		panic(err)
	}
	v.RegisterCustomTypeFunc(common.ValidateCpsTimeType, datetime.CpsTime{})

	// Request validators
	v.RegisterStructValidation(common.ValidateFilteredQuery, pagination.FilteredQuery{})

	pbhValidator := pbehavior.NewValidator(client)
	pbhUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorMongoCollection, "ID")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhUniqueIDValidator.Validate(ctx, sl)
		pbhValidator.ValidateCreateRequest(sl)
	}, pbehavior.CreateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.UpdateRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateUpdateRequest, pbehavior.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateEditRequest, pbehavior.EditRequest{})
	v.RegisterStructValidationCtx(pbhValidator.ValidatePatchRequest, pbehavior.PatchRequest{})
	v.RegisterStructValidation(pbhValidator.ValidateCalendarRequest, pbehavior.CalendarByEntityIDRequest{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateEntityCreateRequest, pbehavior.BulkEntityCreateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateConnectorCreateRequest, pbehavior.BulkConnectorCreateRequestItem{})

	pbhReasonUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "ID")
	pbhReasonUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhReasonUniqueIDValidator.Validate(ctx, sl)
		pbhReasonUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviorreason.CreateRequest{})
	v.RegisterStructValidationCtx(pbhReasonUniqueNameValidator.Validate, pbehaviorreason.UpdateRequest{})

	pbhTypeUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "ID")
	pbhTypeUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhTypeUniqueIDValidator.Validate(ctx, sl)
		pbhTypeUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviortype.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhTypeUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviortype.UpdateRequest{})

	pbhExceptionUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorExceptionMongoCollection, "ID")
	pbhExceptionUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorExceptionMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhExceptionUniqueIDValidator.Validate(ctx, sl)
		pbhExceptionUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviorexception.CreateRequest{})
	v.RegisterStructValidationCtx(pbhExceptionUniqueNameValidator.Validate, pbehaviorexception.UpdateRequest{})
	v.RegisterStructValidation(exdate.ValidateExdateRequest, exdate.Request{})

	v.RegisterStructValidation(pbehaviortimespan.ValidateTimespansRequest, pbehaviortimespan.TimespansRequest{})

	scenarioUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "Name")
	scenarioExistReasonValidator := common.NewExistFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Reason")
	scenarioExistTypeValidator := common.NewExistFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Type")
	scenarioExistIdValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "ID")

	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
		scenarioExistIdValidator.Validate(ctx, sl)
	}, scenario.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
	}, scenario.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
	}, scenario.BulkUpdateRequestItem{})

	v.RegisterStructValidation(scenario.ValidateActionRequest, scenario.ActionRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, action.Parameters{})

	entitybasicValidator := entitybasic.NewValidator(client)
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entitybasicValidator.ValidateEditRequest(ctx, sl)
	}, entitybasic.EditRequest{})

	entityserviceValidator := entityservice.NewValidator(client)
	entityserviceUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.EntityMongoCollection, "ID")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceUniqueIDValidator.Validate(ctx, sl)
		entityserviceValidator.ValidateCreateRequest(sl)
	}, entityservice.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceValidator.ValidateUpdateRequest(sl)
	}, entityservice.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceValidator.ValidateUpdateRequest(sl)
	}, entityservice.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(entityserviceValidator.ValidateEditRequest, entityservice.EditRequest{})

	entityCategoryUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.EntityCategoryMongoCollection, "Name")
	v.RegisterStructValidationCtx(entityCategoryUniqueNameValidator.Validate, entitycategory.EditRequest{})

	roleValidator := role.NewValidator(client)
	v.RegisterStructValidationCtx(roleValidator.ValidateCreateRequest, role.CreateRequest{})
	v.RegisterStructValidationCtx(roleValidator.ValidateEditRequest, role.EditRequest{})

	userValidator := user.NewValidator(client)
	v.RegisterStructValidationCtx(userValidator.ValidateUpdateRequest, user.UpdateRequest{})
	v.RegisterStructValidationCtx(userValidator.ValidateCreateRequest, user.CreateRequest{})
	v.RegisterStructValidationCtx(userValidator.ValidateBulkUpdateRequestItem, user.BulkUpdateRequestItem{})

	accountValidator := account.NewValidator(client)
	v.RegisterStructValidationCtx(accountValidator.ValidateEditRequest, account.EditRequest{})

	v.RegisterStructValidation(view.ValidateEditPositionRequest, view.EditPositionRequest{})

	viewGroupValidator := viewgroup.NewValidator(client)
	v.RegisterStructValidationCtx(viewGroupValidator.ValidateEditRequest, viewgroup.EditRequest{})

	widgetValidator := widget.NewValidator()
	v.RegisterStructValidation(widgetValidator.ValidateEditRequest, widget.EditRequest{})
	v.RegisterStructValidation(widgetValidator.ValidateFilterRequest, widget.FilterRequest{})

	v.RegisterStructValidation(widgetfilter.NewValidator().ValidateCreateRequest, widgetfilter.CreateRequest{})
	v.RegisterStructValidation(widgetfilter.NewValidator().ValidateUpdateRequest, widgetfilter.UpdateRequest{})

	v.RegisterStructValidation(widgettemplate.ValidateEditRequest, widgettemplate.EditRequest{})

	playlistUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PlaylistMongoCollection, "Name")
	v.RegisterStructValidationCtx(playlistUniqueNameValidator.Validate, playlist.EditRequest{})

	stateSettingsValidator := statesettings.NewValidator()
	v.RegisterStructValidation(stateSettingsValidator.ValidateEditRequest, statesettings.EditRequest{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateJUnitThresholds, statesettings.JUnitThreshold{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateThreshold, statesettings.StateThreshold{})

	eventfilterValidator := eventfilter.NewValidator(client)
	eventfilterExistIdValidator := common.NewUniqueFieldValidator(client, mongo.EventFilterRuleCollection, "ID")
	v.RegisterStructValidationCtx(eventfilterValidator.ValidateUpdateRequest, eventfilter.UpdateRequest{})
	v.RegisterStructValidationCtx(eventfilterValidator.ValidateBulkUpdateRequestItem, eventfilter.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		eventfilterValidator.ValidateCreateRequest(ctx, sl)
		eventfilterExistIdValidator.Validate(ctx, sl)
	}, eventfilter.CreateRequest{})

	broadcastmessageValidator := broadcastmessage.NewValidator(client)
	v.RegisterStructValidationCtx(broadcastmessageValidator.Validate, broadcastmessage.BroadcastMessage{})

	v.RegisterStructValidation(messageratestats.ValidateListRequest, messageratestats.ListRequest{})

	idleRuleValidator := idlerule.NewValidator()
	idleRuleUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "Name")
	idleRuleExistIdValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "ID")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleExistIdValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateCreateRequest(sl)
	}, idlerule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateUpdateRequest(sl)
	}, idlerule.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, libidlerule.Parameters{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateBulkUpdateRequestItem(sl)
	}, idlerule.BulkUpdateRequestItem{})

	v.RegisterStructValidation(alarm.ValidateListRequest, alarm.ListRequest{})
	v.RegisterStructValidation(alarm.ValidateDetailsRequest, alarm.DetailsRequest{})

	v.RegisterStructValidation(datastorage.ValidateConfig, libdatastorage.Config{})

	resolveRuleValidator := resolverule.NewValidator()
	resolveRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "ID")
	resolveRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleIdUniqueValidator.Validate(ctx, sl)
		resolveRuleNameUniqueValidator.Validate(ctx, sl)
		resolveRuleValidator.ValidateCreateRequest(sl)
	}, resolverule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleNameUniqueValidator.Validate(ctx, sl)
		resolveRuleValidator.ValidateUpdateRequest(sl)
	}, resolverule.UpdateRequest{})

	flappingRuleValidator := flappingrule.NewValidator()
	flappingRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "ID")
	flappingRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleIdUniqueValidator.Validate(ctx, sl)
		flappingRuleNameUniqueValidator.Validate(ctx, sl)
		flappingRuleValidator.ValidateCreateRequest(sl)
	}, flappingrule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleNameUniqueValidator.Validate(ctx, sl)
		flappingRuleValidator.ValidateUpdateRequest(sl)
	}, flappingrule.UpdateRequest{})

	v.RegisterStructValidation(pattern.ValidateEditRequest, pattern.EditRequest{})

	linkRuleUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.LinkRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		linkrule.ValidateEditRequest(sl)
		linkRuleUniqueNameValidator.Validate(ctx, sl)
	}, linkrule.EditRequest{})

	v.RegisterStructValidation(alarmtag.ValidateCreateRequest, alarmtag.CreateRequest{})
	v.RegisterStructValidation(alarmtag.ValidateUpdateRequest, alarmtag.UpdateRequest{})

	v.RegisterStructValidation(healthcheck.ValidateEngineParameters, config.EngineParameters{})
	v.RegisterStructValidation(healthcheck.ValidateLimitParameters, config.LimitParameters{})
}
