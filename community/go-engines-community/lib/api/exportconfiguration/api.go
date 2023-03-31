package exportconfiguration

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const exportFileName = "export_configurations.yml"

type API interface {
	Export(c *gin.Context)
}

type api struct {
	client          mongo.DbClient
	logger          zerolog.Logger
	collectionNames map[string]string
}

func NewApi(client mongo.DbClient, logger zerolog.Logger) API {
	return &api{
		client: client,
		logger: logger,
		collectionNames: map[string]string{
			"configuration":       mongo.ConfigurationMongoCollection,
			"acl":                 mongo.RightsMongoCollection,
			"pbehavior":           mongo.PbehaviorMongoCollection,
			"pbehavior_type":      mongo.PbehaviorTypeMongoCollection,
			"pbehavior_reason":    mongo.PbehaviorReasonMongoCollection,
			"pbehavior_exception": mongo.PbehaviorExceptionMongoCollection,
			"scenario":            mongo.ScenarioMongoCollection,
			"metaalarm":           mongo.MetaAlarmRulesMongoCollection,
			"idle_rule":           mongo.IdleRuleMongoCollection,
			"eventfilter":         mongo.EventFilterRulesMongoCollection,
			"dynamic_infos":       mongo.DynamicInfosRulesMongoCollection,
			"playlist":            mongo.PlaylistMongoCollection,
			"state_settings":      mongo.StateSettingsMongoCollection,
			"broadcast":           mongo.BroadcastMessageMongoCollection,
			"associative_table":   mongo.AssociativeTableCollection,
			"notification":        mongo.NotificationMongoCollection,
			"view":                mongo.ViewMongoCollection,
			"view_tab":            mongo.ViewTabMongoCollection,
			"widget":              mongo.WidgetMongoCollection,
			"widget_filter":       mongo.WidgetFiltersMongoCollection,
			"widget_template":     mongo.WidgetTemplateMongoCollection,
			"view_group":          mongo.ViewGroupMongoCollection,
			"instruction":         mongo.InstructionMongoCollection,
			"job_config":          mongo.JobConfigMongoCollection,
			"job":                 mongo.JobMongoCollection,
			"resolve_rule":        mongo.ResolveRuleMongoCollection,
			"flapping_rule":       mongo.FlappingRuleMongoCollection,
			"user_preferences":    mongo.UserPreferencesMongoCollection,
			"kpi_filter":          mongo.KpiFilterMongoCollection,
			"pattern":             mongo.PatternMongoCollection,
			"declare_ticket_rule": mongo.DeclareTicketRuleMongoCollection,
			"link_rule":           mongo.LinkRuleMongoCollection,
		},
	}
}

// Export
// @Param body body Request true "body"
func (a *api) Export(c *gin.Context) {
	var r Request

	if err := c.ShouldBindJSON(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	err := a.transformRequest(r)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	file, err := os.CreateTemp("", "cps_export_configurations_*.yml")
	if err != nil {
		panic(err)
	}

	filename := file.Name()

	defer func() {
		err := file.Close()
		if err != nil {
			a.logger.Err(err).Str("filename", filename).Msg("failed to close the file")
			return
		}

		err = os.Remove(file.Name())
		if err != nil {
			a.logger.Err(err).Str("filename", filename).Msg("failed to remove the file")
		}
	}()

	contents := make(map[string]ExportDocuments)

	for _, collection := range r.Exports {
		err = a.addContents(c, contents, collection)
		if err != nil {
			panic(err)
		}
	}

	b, err := yaml.Marshal(contents)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(b)
	if err != nil {
		panic(err)
	}

	c.FileAttachment(filename, exportFileName)
}

func (a *api) addContents(c *gin.Context, contents map[string]ExportDocuments, collectionName string) error {
	cursor, err := a.client.Collection(collectionName).Find(c, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(c)

	content := make(ExportDocuments)
	i := 0

	for cursor.Next(c) {
		var model map[string]interface{}
		err = cursor.Decode(&model)
		if err != nil {
			return err
		}

		content[i] = model
		i++
	}

	if i != 0 {
		contents[collectionName] = content
	}

	return nil
}

func (a *api) transformRequest(r Request) error {
	var invalid []string
	for idx, export := range r.Exports {
		collectionName, ok := a.collectionNames[export]
		if !ok {
			invalid = append(invalid, export)
			continue
		}

		r.Exports[idx] = collectionName
	}

	if len(invalid) != 0 {
		return fmt.Errorf("invalid export fields: [%s]", strings.Join(invalid, ","))
	}

	return nil
}
