package eventfilter_test

import (
	"strings"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnmarshalAction(t *testing.T) {
	Convey("Given an action without type", t, func() {
		bsonAction, err := bson.Marshal(bson.M{})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given an action with an invalid type", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type": "invalid_type",
		})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given a set_field action with an unexpected field", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":             "set_field",
			"name":             "Output",
			"value":            "output",
			"unexpected_field": "",
		})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given a set_field_from_template action with an unexpected field", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":             "set_field_from_template",
			"name":             "Output",
			"value":            "output",
			"unexpected_field": "",
		})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given a set_entity_info_from_template action with an unexpected field", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":             "set_entity_info_from_template",
			"name":             "info_name",
			"value":            "info_value",
			"unexpected_field": "",
		})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given a copy action with an unexpected field", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":             "copy",
			"from":             "ExternalData.Entity",
			"to":               "Entity",
			"unexpected_field": "",
		})
		So(err, ShouldBeNil)

		Convey("Decoding the action returns an error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldNotBeNil)
		})
	})

	Convey("Given a valid set_field action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":  "set_field",
			"name":  "component",
			"value": "test",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			Convey("The action type should be set correctly", func() {
				So(action.Type, ShouldEqual, eventfilter.SetField)
			})

			Convey("The processor should be set correctly", func() {
				processor, success := action.ActionProcessor.(eventfilter.SetFieldProcessor)
				So(success, ShouldBeTrue)
				So(processor.Name.Set, ShouldBeTrue)
				So(processor.Name.Value, ShouldEqual, "component")
				So(processor.Value.Set, ShouldBeTrue)
				So(processor.Value.Value, ShouldEqual, "test")
			})
		})
	})

	Convey("Given a valid set_field_from_template action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":  "set_field_from_template",
			"name":  "output",
			"value": "output by {{.author}}",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			Convey("The action type should be set correctly", func() {
				So(action.Type, ShouldEqual, eventfilter.SetFieldFromTemplate)
			})

			Convey("The processor should be set correctly", func() {
				processor, success := action.ActionProcessor.(eventfilter.SetFieldFromTemplateProcessor)
				So(success, ShouldBeTrue)
				So(processor.Name.Set, ShouldBeTrue)
				So(processor.Name.Value, ShouldEqual, "output")

				builder := strings.Builder{}
				So(processor.Value.Set, ShouldBeTrue)
				So(processor.Value.Value.Execute(&builder, map[string]string{"author": "Phileas"}), ShouldBeNil)
				So(builder.String(), ShouldEqual, "output by Phileas")
			})
		})
	})

	Convey("Given a valid set_entity_info_from_template action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":        "set_entity_info_from_template",
			"name":        "info_name",
			"description": "info_description",
			"value":       "info_value",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			Convey("The action type should be set correctly", func() {
				So(action.Type, ShouldEqual, eventfilter.SetEntityInfoFromTemplate)
			})

			Convey("The processor should be set correctly", func() {
				processor, success := action.ActionProcessor.(eventfilter.SetEntityInfoFromTemplateProcessor)
				So(success, ShouldBeTrue)
				So(processor.Name.Set, ShouldBeTrue)
				So(processor.Name.Value, ShouldEqual, "info_name")
				So(processor.Value.Set, ShouldBeTrue)
				So(processor.Description.Set, ShouldBeTrue)
				So(processor.Description.Value, ShouldEqual, "info_description")
			})
		})
	})

	Convey("Given a valid copy action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type": "copy",
			"from": "ExternalData.Entity",
			"to":   "Event.Entity",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			Convey("The action type should be set correctly", func() {
				So(action.Type, ShouldEqual, eventfilter.Copy)
			})

			Convey("The processor should be set correctly", func() {
				processor, success := action.ActionProcessor.(eventfilter.CopyProcessor)
				So(success, ShouldBeTrue)
				So(processor.From.Set, ShouldBeTrue)
				So(processor.From.Value, ShouldEqual, "ExternalData.Entity")
				So(processor.To.Set, ShouldBeTrue)
				So(processor.To.Value, ShouldEqual, "Event.Entity")
			})
		})
	})
}

func TestSetFieldProcessorApply(t *testing.T) {
	Convey("Given a SetFieldProcessor setting the state of an event to 3", t, func() {
		processor := eventfilter.SetFieldProcessor{
			Name:  utils.OptionalString{Set: true, Value: "State"},
			Value: utils.OptionalInterface{Set: true, Value: 3},
		}

		Convey("The processor can be applied to the event without error", func() {
			report := eventfilter.Report{}
			event := types.Event{
				State: 1,
			}
			event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
			So(err, ShouldBeNil)

			Convey("The state is set to 3", func() {
				So(event.State, ShouldEqual, 3)
			})

			Convey("The entity is not updated", func() {
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})

	Convey("Given a SetFieldProcessor setting the state of an event to a string", t, func() {
		processor := eventfilter.SetFieldProcessor{
			Name:  utils.OptionalString{Set: true, Value: "State"},
			Value: utils.OptionalInterface{Set: true, Value: "a string"},
		}

		Convey("Applying the processor returns an error", func() {
			report := eventfilter.Report{}
			event := types.Event{
				State: 1,
			}
			event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestSetFieldFromTemplateProcessorApply(t *testing.T) {
	Convey("Given a SetFieldFromTemplateProcessor setting the output of an event", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":  "set_field_from_template",
			"name":  "Output",
			"value": "{{.Event.Output}} (by {{.Event.Author}})",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			processor, success := action.ActionProcessor.(eventfilter.SetFieldFromTemplateProcessor)
			So(success, ShouldBeTrue)

			Convey("The processor can be applied to the event without error", func() {
				report := eventfilter.Report{}
				event := types.Event{
					Author: "Billy",
					Output: "test",
				}
				data := eventfilter.ActionParameters{
					DataSourceGetterParameters: eventfilter.DataSourceGetterParameters{
						Event: event,
					},
				}
				event, err := processor.Apply(event, data, &report)
				So(err, ShouldBeNil)

				Convey("The output is correctly set", func() {
					So(event.Output, ShouldEqual, "test (by Billy)")
				})

				Convey("The entity is not updated", func() {
					So(report.EntityUpdated, ShouldBeFalse)
				})
			})
		})
	})

	Convey("Given a SetFieldFromTemplateProcessor setting the state of an event to a string", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":  "set_field_from_template",
			"name":  "State",
			"value": "{{.Event.Output}} (by {{.Event.Author}})",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			processor, success := action.ActionProcessor.(eventfilter.SetFieldFromTemplateProcessor)
			So(success, ShouldBeTrue)

			Convey("Applying the processor returns an error", func() {
				report := eventfilter.Report{}
				event := types.Event{
					Author: "Billy",
					Output: "test",
				}
				data := eventfilter.ActionParameters{
					DataSourceGetterParameters: eventfilter.DataSourceGetterParameters{
						Event: event,
					},
				}
				event, err := processor.Apply(event, data, &report)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a SetFieldFromTemplateProcessor using unknown fields", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":  "set_field_from_template",
			"name":  "Output",
			"value": "{{.Event.NoSuchField}}",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			processor, success := action.ActionProcessor.(eventfilter.SetFieldFromTemplateProcessor)
			So(success, ShouldBeTrue)

			Convey("Applying the processor returns an error", func() {
				report := eventfilter.Report{}
				event := types.Event{
					Output: "test",
				}
				data := eventfilter.ActionParameters{
					DataSourceGetterParameters: eventfilter.DataSourceGetterParameters{
						Event: event,
					},
				}
				event, err := processor.Apply(event, data, &report)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSetEntityInfoFromTemplateProcessorApply(t *testing.T) {
	Convey("Given a valid set_entity_info_from_template action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type":        "set_entity_info_from_template",
			"name":        "info_name",
			"description": "info_description",
			"value":       "{{.Event.Output}}",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			processor, success := action.ActionProcessor.(eventfilter.SetEntityInfoFromTemplateProcessor)
			So(success, ShouldBeTrue)

			Convey("Applying the processor to an event without entity returns an error", func() {
				report := eventfilter.Report{}
				event := types.Event{
					Output: "test",
				}
				event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
				So(err, ShouldNotBeNil)
			})

			Convey("Applying the processor to an event with an entity updates the entity", func() {
				report := eventfilter.Report{}
				entity := types.NewEntity("r/c", "r/c", types.EntityTypeResource, nil, nil, nil)
				event := types.Event{
					Output: "test",
					Entity: &entity,
				}
				data := eventfilter.ActionParameters{
					DataSourceGetterParameters: eventfilter.DataSourceGetterParameters{
						Event: event,
					},
				}
				event, err := processor.Apply(event, data, &report)
				So(err, ShouldBeNil)
				So(event.Entity.Infos, ShouldContainKey, "info_name")
				So(event.Entity.Infos["info_name"].Name, ShouldEqual, "info_name")
				So(event.Entity.Infos["info_name"].Description, ShouldEqual, "info_description")
				So(event.Entity.Infos["info_name"].Value, ShouldEqual, "test")
				So(event.Entity.Infos["info_name"].RealValue, ShouldEqual, "test")
				So(report.EntityUpdated, ShouldBeTrue)
			})
		})
	})
}

func TestCopyProcessorApply(t *testing.T) {
	Convey("Given a valid copy action", t, func() {
		bsonAction, err := bson.Marshal(bson.M{
			"type": "copy",
			"from": "ExternalData.Entity",
			"to":   "Entity",
		})
		So(err, ShouldBeNil)

		Convey("The action should be decoded without error", func() {
			var action eventfilter.Action
			So(bson.Unmarshal(bsonAction, &action), ShouldBeNil)

			processor, success := action.ActionProcessor.(eventfilter.CopyProcessor)
			So(success, ShouldBeTrue)

			Convey("Applying the processor sets the event's entity", func() {
				report := eventfilter.Report{}
				event := types.Event{
					Output: "test",
				}
				data := eventfilter.ActionParameters{
					ExternalData: map[string]interface{}{
						"Entity": types.NewEntity("r/c", "r/c", types.EntityTypeResource, nil, nil, nil),
					},
				}
				event, err := processor.Apply(event, data, &report)
				So(err, ShouldBeNil)

				So(event.Entity.ID, ShouldEqual, "r/c")
				So(event.Entity.Name, ShouldEqual, "r/c")
				So(event.Entity.Type, ShouldEqual, types.EntityTypeResource)
				So(report.EntityUpdated, ShouldBeFalse)
			})
		})
	})
}
