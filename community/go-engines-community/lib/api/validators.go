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
	libdatastorage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	err = v.RegisterValidation("alarmpatterns", common.ValidateAlarmPatterns)
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
	v.RegisterCustomTypeFunc(common.ValidateCpsTimeType, types.CpsTime{})

	// Request validators
	v.RegisterStructValidation(common.ValidateFilteredQuery, pagination.FilteredQuery{})

	pbhValidator := pbehavior.NewValidator(client)
	pbhUniqueIDValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorMongoCollection, "ID")
	pbhUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PbehaviorMongoCollection, "Name")
	pbhExistReasonValidator := common.NewExistFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Reason")
	pbhExistTypeValidator := common.NewExistFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Type")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhUniqueIDValidator.Validate(ctx, sl)
		pbhUniqueNameValidator.Validate(ctx, sl)
	}, pbehavior.CreateRequest{})
	v.RegisterStructValidationCtx(pbhUniqueNameValidator.Validate, pbehavior.UpdateRequest{})
	v.RegisterStructValidationCtx(pbhUniqueNameValidator.Validate, pbehavior.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhValidator.ValidateEditRequest(ctx, sl)
		pbhExistReasonValidator.Validate(ctx, sl)
		pbhExistTypeValidator.Validate(ctx, sl)
	}, pbehavior.EditRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		pbhValidator.ValidateEditRequest(ctx, sl)
	}, pbehavior.PatchRequest{})

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
	v.RegisterStructValidation(pbehaviortype.ValidateEditRequest, pbehaviortype.EditRequest{})

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

	v.RegisterStructValidation(scenario.ValidateEditRequest, scenario.EditRequest{})
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
	v.RegisterStructValidation(scenario.ValidateChangeStateParametersRequest, scenario.ChangeStateParametersRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		scenario.ValidatePbehaviorParametersRequest(sl)
		scenarioExistReasonValidator.Validate(ctx, sl)
		scenarioExistTypeValidator.Validate(ctx, sl)
	}, scenario.PbehaviorParametersRequest{})
	v.RegisterStructValidation(scenario.ValidateWebhookRequest, scenario.WebhookRequest{})

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
		entityserviceUniqueNameValidator.Validate(ctx, sl)
	}, entityservice.CreateRequest{})
	v.RegisterStructValidationCtx(entityserviceUniqueNameValidator.Validate, entityservice.UpdateRequest{})
	v.RegisterStructValidationCtx(entityserviceUniqueNameValidator.Validate, entityservice.BulkUpdateRequestItem{})
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
	viewBulkUniqueIDValidator := common.NewUniqueBulkFieldValidator("ID")
	viewBulkUniqueTitleValidator := common.NewUniqueBulkFieldValidator("Title")
	v.RegisterStructValidationCtx(viewValidator.ValidateEditRequest, view.EditRequest{})
	v.RegisterStructValidation(view.ValidateWidgetParametersJunitRequest, view.WidgetParametersJunitRequest{})
	v.RegisterStructValidation(view.ValidateEditPositionRequest, view.EditPositionRequest{})
	v.RegisterStructValidationCtx(viewBulkUniqueTitleValidator.Validate, view.BulkCreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		viewBulkUniqueIDValidator.Validate(ctx, sl)
		viewBulkUniqueTitleValidator.Validate(ctx, sl)
	}, view.BulkUpdateRequest{})

	viewGroupUniqueTitleValidator := common.NewUniqueFieldValidator(client, mongo.ViewGroupMongoCollection, "Title")
	viewGroupBulkUniqueIDValidator := common.NewUniqueBulkFieldValidator("ID")
	viewGroupBulkUniqueTitleValidator := common.NewUniqueBulkFieldValidator("Title")
	v.RegisterStructValidationCtx(viewGroupUniqueTitleValidator.Validate, viewgroup.EditRequest{}, viewgroup.BulkUpdateRequestItem{})
	v.RegisterStructValidationCtx(viewGroupBulkUniqueTitleValidator.Validate, viewgroup.BulkCreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		viewGroupBulkUniqueIDValidator.Validate(ctx, sl)
		viewGroupBulkUniqueTitleValidator.Validate(ctx, sl)
	}, viewgroup.BulkUpdateRequest{})

	playlistValidator := playlist.NewPlaylistValidator(client)
	playlistUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.PlaylistMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		playlistUniqueNameValidator.Validate(ctx, sl)
		playlistValidator.ValidateEditRequest(ctx, sl)
	}, playlist.EditRequest{})

	stateSettingsValidator := statesettings.NewValidator()
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateSettingRequest, statesettings.StateSettingRequest{})
	v.RegisterStructValidation(stateSettingsValidator.ValidateStateThresholds, statesettings.StateThresholds{})

	eventfilterValidator := eventfilter.NewValidator(client)
	eventfilterExistIdValidator := common.NewUniqueFieldValidator(client, mongo.EventFilterRulesMongoCollection, "ID")
	v.RegisterStructValidationCtx(eventfilterValidator.Validate, eventfilter.EventFilterPayload{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		eventfilterExistIdValidator.Validate(ctx, sl)
	}, eventfilter.EventFilter{})

	broadcastmessageValidator := broadcastmessage.NewValidator(client)
	v.RegisterStructValidationCtx(broadcastmessageValidator.Validate, broadcastmessage.BroadcastMessage{})

	v.RegisterStructValidation(messageratestats.ValidateListRequest, messageratestats.ListRequest{})

	idleRuleUniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "Name")
	idleRuleExistIdValidator := common.NewUniqueFieldValidator(client, mongo.IdleRuleMongoCollection, "ID")
	v.RegisterStructValidation(idlerule.ValidateEditRequest, idlerule.EditRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
		idleRuleExistIdValidator.Validate(ctx, sl)
	}, idlerule.CreateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
	}, idlerule.UpdateRequest{})
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		idleRuleUniqueNameValidator.Validate(ctx, sl)
	}, idlerule.BulkUpdateRequestItem{})
	v.RegisterStructValidation(idlerule.ValidateCountPatternRequest, idlerule.CountByPatternRequest{})

	v.RegisterStructValidation(alarm.ValidateListRequest, alarm.ListRequest{})
	v.RegisterStructValidation(datastorage.ValidateConfig, libdatastorage.Config{})

	resolveRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "ID")
	resolveRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.ResolveRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		resolveRuleIdUniqueValidator.Validate(ctx, sl)
		resolveRuleNameUniqueValidator.Validate(ctx, sl)
	}, resolverule.CreateRequest{})
	v.RegisterStructValidationCtx(resolveRuleNameUniqueValidator.Validate, resolverule.UpdateRequest{})
	v.RegisterStructValidation(resolverule.ValidateEditRequest, resolverule.EditRequest{})

	flappingRuleIdUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "ID")
	flappingRuleNameUniqueValidator := common.NewUniqueFieldValidator(client, mongo.FlappingRuleMongoCollection, "Name")
	v.RegisterStructValidationCtx(func(ctx context.Context, sl validator.StructLevel) {
		flappingRuleIdUniqueValidator.Validate(ctx, sl)
		flappingRuleNameUniqueValidator.Validate(ctx, sl)
	}, flappingrule.CreateRequest{})
	v.RegisterStructValidationCtx(flappingRuleNameUniqueValidator.Validate, flappingrule.UpdateRequest{})
	v.RegisterStructValidation(flappingrule.ValidateEditRequest, flappingrule.EditRequest{})
}
