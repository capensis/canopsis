package api

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	pbehaviorapi "git.canopsis.net/canopsis/go-engines/lib/api/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorexception"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviorreason"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviortimespan"
	"git.canopsis.net/canopsis/go-engines/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/go-engines/lib/api/scenario"
	"git.canopsis.net/canopsis/go-engines/lib/api/watcherweather"
	libheartbeat "git.canopsis.net/canopsis/go-engines/lib/canopsis/heartbeat"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(client mongo.DbClient) {
	pbhRequestValidator := pbehaviorapi.NewValidator(client)
	reasonValidator := pbehaviorreason.NewValidator(client)
	typeValidator := pbehaviortype.NewValidator(client)
	exceptionValidator := pbehaviorexception.NewValidator(client)
	heartbeatUniqueNameValidator := common.NewUniqueFieldValidator(client, libheartbeat.HeartbeatCollectionName, "Name")
	heartbeatBulkUniqueNameValidator := common.NewUniqueBulkFieldValidator("Name")
	heartbeatValidator := heartbeat.NewValidator(client)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(common.ValidateFilteredQuery, pagination.FilteredQuery{})
		v.RegisterStructValidation(pbhRequestValidator.ValidateCreateRequest, pbehaviorapi.CreateRequest{})
		v.RegisterStructValidation(pbhRequestValidator.ValidateEditRequest, pbehaviorapi.EditRequest{})
		v.RegisterStructValidation(reasonValidator.ValidateReasonCreateRequest, pbehaviorreason.CreateRequest{})
		v.RegisterStructValidation(reasonValidator.ValidateReasonUpdateRequest, pbehaviorreason.UpdateRequest{})
		v.RegisterStructValidation(typeValidator.ValidateTypeCreateRequest, pbehaviortype.CreateRequest{})
		v.RegisterStructValidation(typeValidator.ValidateTypeUpdateRequest, pbehaviortype.UpdateRequest{})
		v.RegisterStructValidation(exceptionValidator.ValidateExceptionCreateRequest, pbehaviorexception.CreateRequest{})
		v.RegisterStructValidation(exceptionValidator.ValidateExceptionUpdateRequest, pbehaviorexception.UpdateRequest{})
		v.RegisterStructValidation(exceptionValidator.ValidateExdateRequest, pbehaviorexception.ExdateRequest{})
		v.RegisterStructValidation(pbehaviortimespan.ValidateTimespansRequest, pbehaviortimespan.TimespansRequest{})
		v.RegisterStructValidation(pbehaviortimespan.ValidateExdateRequest, pbehaviortimespan.ExdateRequest{})
		v.RegisterStructValidation(heartbeatValidator.ValidateCreateRequest, heartbeat.CreateRequest{})
		v.RegisterStructValidation(heartbeatUniqueNameValidator.Validate, heartbeat.UpdateRequest{})
		v.RegisterStructValidation(heartbeatUniqueNameValidator.Validate, heartbeat.BulkUpdateRequestItem{})
		v.RegisterStructValidation(heartbeatBulkUniqueNameValidator.Validate, heartbeat.BulkCreateRequest{})
		v.RegisterStructValidation(heartbeatBulkUniqueNameValidator.Validate, heartbeat.BulkUpdateRequest{})
		v.RegisterStructValidation(func(sl validator.StructLevel) {
			uniqueNameValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "Name")
			uniquePriorityValidator := common.NewUniqueFieldValidator(client, mongo.ScenarioMongoCollection, "Priority")

			scenario.ValidateEditRequest(sl)
			uniqueNameValidator.Validate(sl)
			uniquePriorityValidator.Validate(sl)
		}, scenario.EditRequest{})
		v.RegisterStructValidation(scenario.ValidateActionRequest, scenario.ActionRequest{})
		v.RegisterStructValidation(scenario.ValidateChangeStateParametersRequest, scenario.ChangeStateParametersRequest{})
		v.RegisterStructValidation(func(sl validator.StructLevel) {
			existReasonValidator := common.NewExistFieldValidator(client, mongo.PbehaviorReasonMongoCollection, "Reason")
			existTypeValidator := common.NewExistFieldValidator(client, mongo.PbehaviorTypeMongoCollection, "Type")

			scenario.ValidatePbehaviorParametersRequest(sl)
			existReasonValidator.Validate(sl)
			existTypeValidator.Validate(sl)
		}, scenario.PbehaviorParametersRequest{})
		v.RegisterStructValidation(scenario.ValidateWebhookRequest, scenario.WebhookRequest{})
		v.RegisterStructValidation(watcherweather.ValidateRequest, watcherweather.ListRequest{})
	}
}
