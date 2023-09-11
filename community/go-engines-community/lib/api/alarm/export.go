package alarm

import (
	"context"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	liblink "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
)

var stateTitles = map[int]string{
	types.AlarmStateOK:       types.AlarmStateTitleOK,
	types.AlarmStateMinor:    types.AlarmStateTitleMinor,
	types.AlarmStateMajor:    types.AlarmStateTitleMajor,
	types.AlarmStateCritical: types.AlarmStateTitleCritical,
}

var statusTitles = map[int]string{
	types.AlarmStatusOff:       types.AlarmStatusTitleOff,
	types.AlarmStatusOngoing:   types.AlarmStatusTitleOngoing,
	types.AlarmStatusStealthy:  types.AlarmStatusTitleStealthy,
	types.AlarmStatusFlapping:  types.AlarmStatusTitleFlapping,
	types.AlarmStatusCancelled: types.AlarmStatusTitleCancelled,
}

func newExportCursor(
	cursor mongo.Cursor,
	fields []string,
	timeFormat string,
	location *time.Location,
	instructions []Instruction,
	linkGenerator liblink.Generator,
	logger zerolog.Logger,
) export.DataCursor {
	return &mongoCursor{
		cursor:        cursor,
		fields:        fields,
		timeFormat:    timeFormat,
		location:      location,
		instructions:  instructions,
		linkGenerator: linkGenerator,
		logger:        logger,
	}
}

type mongoCursor struct {
	cursor        mongo.Cursor
	fields        []string
	timeFormat    string
	location      *time.Location
	instructions  []Instruction
	linkGenerator liblink.Generator
	logger        zerolog.Logger
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *mongoCursor) Scan(m *map[string]any) error {
	err := c.cursor.Decode(m)
	if err != nil {
		return err
	}

	var model types.AlarmWithEntity
	if len(c.instructions) > 0 || c.linkGenerator != nil {
		delete(*m, "model")
		item := struct {
			Model types.AlarmWithEntity `bson:"model"`
		}{}
		err := c.cursor.Decode(&item)
		if err != nil {
			return err
		}

		model = item.Model
	}

	*m = c.filterFields(context.Background(), *m, c.fields, model)
	return nil
}

func (c *mongoCursor) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

func (c *mongoCursor) filterFields(
	ctx context.Context,
	m map[string]any,
	fields []string,
	model types.AlarmWithEntity,
) map[string]any {
	var links liblink.LinksByCategory
	var err error
	if c.linkGenerator != nil {
		links, err = c.linkGenerator.GenerateForAlarm(ctx, model.Alarm, model.Entity)
		if err != nil {
			c.logger.Err(err).Str("alarm", model.Alarm.ID).Msg("cannot generate links")
		}
	}

	res := make(map[string]any, len(fields))
	for _, field := range fields {
		v, _ := c.getNestedVal(m, strings.Split(field, "."))
		res[field] = c.transformField(field, v, model, links)
	}

	return res
}

func (c *mongoCursor) transformField(k string, v any, model types.AlarmWithEntity, linksByCategory liblink.LinksByCategory) any {
	switch k {
	case "v.state.val":
		if i, ok := c.getInt64(v); ok {
			return stateTitles[int(i)]
		}
	case "v.status.val":
		if i, ok := c.getInt64(v); ok {
			return statusTitles[int(i)]
		}
	case "t",
		"v.creation_date",
		"v.activation_date",
		"v.last_update_date",
		"v.last_event_date",
		"v.resolved":
		if i, ok := c.getInt64(v); ok {
			return types.NewCpsTime(i).In(c.location).Time.Format(c.timeFormat)
		}
	case "v.infos":
		values := make([]string, 0)
		if m, ok := v.(map[string]any); ok {
			for _, mv := range m {
				if infos, ok := mv.(map[string]any); ok {
					for ik, iv := range infos {
						if s, ok := iv.(string); ok {
							values = append(values, ik+": "+s)
						}
					}
				}
			}
		}

		return strings.Join(values, ",")
	case "entity.infos",
		"entity.component_infos":
		values := make([]string, 0)
		if m, ok := v.(map[string]any); ok {
			for mk, mv := range m {
				if info, ok := mv.(map[string]any); ok {
					if s, ok := info["value"].(string); ok {
						values = append(values, mk+": "+s)
					}
				}
			}
		}

		return strings.Join(values, ",")
	case "assigned_instructions":
		names := c.matchInstructions(model)
		return strings.Join(names, ",")
	case "links":
		values := make([]string, 0)
		for _, links := range linksByCategory {
			for _, link := range links {
				values = append(values, link.Label+": "+link.Url)
			}
		}

		return strings.Join(values, ",")
	default:
		if strings.HasSuffix(k, ".t") {
			if i, ok := c.getInt64(v); ok {
				return types.NewCpsTime(i).In(c.location).Time.Format(c.timeFormat)
			}
		}

		if strings.HasPrefix(k, "links.") {
			category := strings.TrimPrefix(k, "links.")
			values := make([]string, 0)
			for _, link := range linksByCategory[category] {
				values = append(values, link.Label+": "+link.Url)
			}

			return strings.Join(values, ",")
		}
	}

	return v
}

func (c *mongoCursor) getNestedVal(m map[string]any, keys []string) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	if v, ok := m[keys[0]]; ok {
		if len(keys) == 1 {
			return v, true
		}

		if mv, ok := v.(map[string]any); ok {
			return c.getNestedVal(mv, keys[1:])
		}
	}

	return nil, false
}

func (c *mongoCursor) getInt64(v any) (int64, bool) {
	switch i := v.(type) {
	case int64:
		return i, true
	case int32:
		return int64(i), true
	case int:
		return int64(i), true
	default:
		return 0, false
	}
}

func (c *mongoCursor) matchInstructions(model types.AlarmWithEntity) []string {
	names := make([]string, 0)
	alarmPbhType := model.Alarm.Value.PbehaviorInfo.TypeID
	for _, instruction := range c.instructions {
		match, err := pattern.Match(model.Entity, model.Alarm, instruction.EntityPattern, instruction.AlarmPattern,
			instruction.OldEntityPatterns, instruction.OldAlarmPatterns)
		if err != nil || !match {
			continue
		}

		found := false
		for _, pbhType := range instruction.DisabledOnPbh {
			if alarmPbhType == pbhType {
				found = true
				break
			}
		}

		if found {
			continue
		}

		if len(instruction.ActiveOnPbh) > 0 {
			found := false
			for _, pbhType := range instruction.ActiveOnPbh {
				if alarmPbhType == pbhType {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		names = append(names, instruction.Name)
	}

	return names
}
