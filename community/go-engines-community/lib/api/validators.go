package api

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/broadcastmessage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitybasic"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/idlerule"
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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/serviceweather"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/user"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	libdatastorage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	libidlerule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libvalidator "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/validator"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func RegisterValidators(client mongo.DbClient, enableSameServiceNames bool) {
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
	v.RegisterCustomTypeFunc(common.ValidateCpsTimeType, types.CpsTime{})

	// Request validators
	v.RegisterStructValidation(common.ValidateFilteredQuery, pagination.FilteredQuery{})
	v.RegisterStructValidation(common.ValidateAlarmPatternFieldsRequest, common.AlarmPatternFieldsRequest{})
	v.RegisterStructValidation(common.ValidateEntityPatternFieldsRequest, common.EntityPatternFieldsRequest{})
	v.RegisterStructValidation(common.ValidatePbehaviorPatternFieldsRequest, common.PbehaviorPatternFieldsRequest{})

	pbhValidator := pbehavior.NewValidator(client)
	pbhUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorMongoCollection, "ID")
	pbhUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhUniqueIDValidator.Validate(ctx, sl)
		pbhUniqueNameValidator.Validate(ctx, sl)
		pbhValidator.ValidateCreateRequest(sl)
	}, pbehavior.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhUniqueNameValidator.Validate(ctx, sl)
		pbhValidator.ValidateUpdateRequest(ctx, sl)
	}, pbehavior.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhUniqueNameValidator.Validate(ctx, sl)
		pbhValidator.ValidateUpdateRequest(ctx, sl)
	}, pbehavior.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(pbhValidator.ValidateEditRequest, pbehavior.EditRequest{})
	v.RegisterStructValidationCtx(pbhValidator.ValidatePatchRequest, pbehavior.PatchRequest{})

	pbhReasonUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "ID")
	pbhReasonUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhReasonUniqueIDValidator.Validate(ctx, sl)
		pbhReasonUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviorreason.CreateRequest{})
	v.RegisterStructValidationCtx(pbhReasonUniqueNameValidator.Validate, pbehaviorreason.UpdateRequest{})

	pbhTypeUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "ID")
	pbhTypeUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Name")
	pbhTypeUniquePriorityValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Priority")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhTypeUniqueIDValidator.Validate(ctx, sl)
		pbhTypeUniqueNameValidator.Validate(ctx, sl)
		pbhTypeUniquePriorityValidator.Validate(ctx, sl)
	}, pbehaviortype.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhTypeUniqueNameValidator.Validate(ctx, sl)
		pbhTypeUniquePriorityValidator.Validate(ctx, sl)
	}, pbehaviortype.UpdateRequest{})

	pbhExceptionUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorExceptionMongoCollection, "ID")
	pbhExceptionUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorExceptionMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhExceptionUniqueIDValidator.Validate(ctx, sl)
		pbhExceptionUniqueNameValidator.Validate(ctx, sl)
	}, pbehaviorexception.CreateRequest{})
	v.RegisterStructValidationCtx(pbhExceptionUniqueNameValidator.Validate, pbehaviorexception.UpdateRequest{})
	v.RegisterStructValidation(pbehaviorexception.ValidateExdateRequest, pbehaviorexception.ExdateRequest{})

	v.RegisterStructValidation(pbehaviortimespan.ValidateTimespansRequest, pbehaviortimespan.TimespansRequest{})
	v.RegisterStructValidation(pbehaviortimespan.ValidateExdateRequest, pbehaviortimespan.ExdateRequest{})

	scenarioUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "Name")
	scenarioUniquePriorityValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "Priority")
	scenarioExistReasonValidator := common.NewExistFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Reason")
	scenarioExistTypeValidator := common.NewExistFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Type")
	scenarioExistIdValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "ID")

	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
		scenarioUniquePriorityValidator.Validate(ctx, sl)
		scenarioExistIdValidator.Validate(ctx, sl)
	}, scenario.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
		scenarioUniquePriorityValidator.Validate(ctx, sl)
	}, scenario.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioUniqueNameValidator.Validate(ctx, sl)
		scenarioUniquePriorityValidator.Validate(ctx, sl)
	}, scenario.BulkUpdateRequestItem{})

	v.RegisterStructValidation(scenario.ValidateActionRequest, scenario.ActionRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, action.Parameters{})

	v.RegisterStructValidation(serviceweather.ValidateRequest, serviceweather.ListRequest{})

	entitybasicValidator := entitybasic.NewValidator(client)
	v.RegisterStructValidation(entity.ValidateListRequest, entity.ListRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entitybasicValidator.ValidateEditRequest(ctx, sl)
	}, entitybasic.EditRequest{})

	entityserviceValidator := entityservice.NewValidator(client)
	entityserviceUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.EntityMongoCollection, "ID")
	entityserviceUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.EntityMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		entityserviceUniqueIDValidator.Validate(ctx, sl)
		if !enableSameServiceNames {
			entityserviceUniqueNameValidator.Validate(ctx, sl)
		}
		entityserviceValidator.ValidateCreateRequest(sl)
	}, entityservice.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		if !enableSameServiceNames {
			entityserviceUniqueNameValidator.Validate(ctx, sl)
		}
		entityserviceValidator.ValidateUpdateRequest(ctx, sl)
	}, entityservice.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		if !enableSameServiceNames {
			entityserviceUniqueNameValidator.Validate(ctx, sl)
		}
		entityserviceValidator.ValidateUpdateRequest(ctx, sl)
	}, entityservice.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(entityserviceValidator.ValidateEditRequest, entityservice.EditRequest{})

	entityCategoryUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.EntityCategoryMongoCollection, "Name")
	v.RegisterStructValidationCtx(entityCategoryUniqueNameValidator.Validate, entitycategory.EditRequest{})

	roleValidator := role.NewValidator(client)
	v.RegisterStructValidationCtx(roleValidator.ValidateCreateRequest, role.CreateRequest{})
	v.RegisterStructValidationCtx(roleValidator.ValidateEditRequest, role.EditRequest{})

	userValidator := user.NewValidator(client)
	v.RegisterStructValidationCtx(userValidator.ValidateRequest, user.Request{})
	v.RegisterStructValidationCtx(userValidator.ValidateBulkUpdateRequestItem, user.BulkUpdateRequestItem{})

	viewValidator := view.NewValidator(client)
	v.RegisterStructValidationCtx(viewValidator.ValidateEditRequest, view.EditRequest{})
	v.RegisterStructValidation(view.ValidateEditPositionRequest, view.EditPositionRequest{})

	viewGroupUniqueTitleValidator := common.NewUniqueFieldValidator(client, mongo.ViewGroupMongoCollection, "Title")
	v.RegisterStructValidationCtx(viewGroupUniqueTitleValidator.Validate, viewgroup.EditRequest{})

	widgetValidator := widget.NewValidator(client)
	v.RegisterStructValidation(widgetValidator.ValidateEditRequest, widget.EditRequest{})
	v.RegisterStructValidationCtx(widgetValidator.ValidateFilterRequest, widget.FilterRequest{})

	v.RegisterStructValidationCtx(widgetfilter.NewValidator(client).ValidateEditRequest, widgetfilter.EditRequest{})

	playlistUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PlaylistMongoCollection, "Name")
	v.RegisterStructValidationCtx(playlistUniqueNameValidator.Validate, playlist.EditRequest{})

	stateSettingsValidator := statesettings.NewValidator()
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateSettingRequest, statesettings.StateSettingRequest{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateThresholds, statesettings.StateThresholds{})

	eventfilterValidator := eventfilter.NewValidator(client)
	eventfilterExistIdValidator := common.NewUniqueFieldValidator(client, mongo.EventFilterRulesMongoCollection, "ID")
	v.RegisterStructValidation(eventfilterValidator.ValidateEventFilter, eventfilter.EditRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		eventfilterExistIdValidator.Validate(ctx, sl)
	}, eventfilter.CreateRequest{})

	broadcastmessageValidator := broadcastmessage.NewValidator(client)
	v.RegisterStructValidationCtx(broadcastmessageValidator.Validate, broadcastmessage.BroadcastMessage{})

	v.RegisterStructValidation(messageratestats.ValidateListRequest, messageratestats.ListRequest{})

	idleRuleValidator := idlerule.NewValidator(client)
	idleRuleUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "Name")
	idleRuleExistIdValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "ID")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleExistIdValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateCreateRequest(ctx, sl)
	}, idlerule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateUpdateRequest(ctx, sl)
	}, idlerule.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, libidlerule.Parameters{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleValidator.ValidateBulkUpdateRequestItem(ctx, sl)
	}, idlerule.BulkUpdateRequestItem{})

	v.RegisterStructValidation(alarm.ValidateListRequest, alarm.ListRequest{})
	v.RegisterStructValidation(alarm.ValidateDetailsRequest, alarm.DetailsRequest{})

	v.RegisterStructValidation(datastorage.ValidateConfig, libdatastorage.Config{})

	resolveRuleValidator := resolverule.NewValidator(client)
	resolveRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "ID")
	resolveRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleIdUniqueValidator.Validate(ctx, sl)
		resolveRuleNameUniqueValidator.Validate(ctx, sl)
		resolveRuleValidator.ValidateCreateRequest(ctx, sl)
	}, resolverule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleNameUniqueValidator.Validate(ctx, sl)
		resolveRuleValidator.ValidateUpdateRequest(ctx, sl)
	}, resolverule.UpdateRequest{})

	flappingRuleValidator := flappingrule.NewValidator(client)
	flappingRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "ID")
	flappingRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleIdUniqueValidator.Validate(ctx, sl)
		flappingRuleNameUniqueValidator.Validate(ctx, sl)
		flappingRuleValidator.ValidateCreateRequest(ctx, sl)
	}, flappingrule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleNameUniqueValidator.Validate(ctx, sl)
		flappingRuleValidator.ValidateUpdateRequest(ctx, sl)
	}, flappingrule.UpdateRequest{})

	v.RegisterStructValidation(pattern.ValidateEditRequest, pattern.EditRequest{})
}
