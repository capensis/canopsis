Feature: Get alarms
  I need to be able to get a alarms

  @concurrent
  Scenario: given get one request should return alarm
    When I am admin
    When I do GET /api/v4/alarms/test-alarm-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-alarm-to-get-1",
      "entity": {
        "_id": "test-resource-to-alarm-get-1/test-component-to-alarm-get",
        "category": {
          "_id": "test-category-to-alarm-get-1",
          "name": "test-category-to-alarm-get-1-name"
        },
        "connector": "test-connector-default/test-connector-default-name",
        "component": "test-component-to-alarm-get",
        "enabled": true,
        "impact_level": 1,
        "infos": {
          "test-resource-to-alarm-get-1-info-1": {
            "name": "test-resource-to-alarm-get-1-info-1-name",
            "description": "test-resource-to-alarm-get-1-info-1-description",
            "value": "test-resource-to-alarm-get-1-info-1-value"
          },
          "test-resource-to-alarm-get-1-info-2": {
            "name": "test-resource-to-alarm-get-1-info-2-name",
            "description": "test-resource-to-alarm-get-1-info-2-description",
            "value": false
          },
          "test-resource-to-alarm-get-1-info-3": {
            "name": "test-resource-to-alarm-get-1-info-3-name",
            "description": "test-resource-to-alarm-get-1-info-3-description",
            "value": 1022
          },
          "test-resource-to-alarm-get-1-info-4": {
            "name": "test-resource-to-alarm-get-1-info-4-name",
            "description": "test-resource-to-alarm-get-1-info-4-description",
            "value": 10.45
          },
          "test-resource-to-alarm-get-1-info-5": {
            "name": "test-resource-to-alarm-get-1-info-5-name",
            "description": "test-resource-to-alarm-get-1-info-5-description",
            "value": null
          },
          "test-resource-to-alarm-get-1-info-6": {
            "name": "test-resource-to-alarm-get-1-info-6-name",
            "description": "test-resource-to-alarm-get-1-info-6-description",
            "value": ["test-resource-to-alarm-get-1-info-6-value"]
          }
        },
        "name": "test-resource-to-alarm-get-1",
        "type": "resource",
        "depends_count": 0,
        "impacts_count": 0
      },
      "impact_state": 3,
      "infos": {},
      "t": 1597030219,
      "v": {
        "children": [],
        "component": "test-component-to-alarm-get",
        "connector": "test-connector-default",
        "connector_name": "test-connector-default-name",
        "creation_date": 1597030219,
        "display_name": "RC-KC_tW",
        "infos": {},
        "infos_rule_version": {},
        "initial_long_output": "",
        "initial_output": "test-alarm-to-get-1-output",
        "last_comment": {
          "_t": "comment",
          "a": "root",
          "user_id": "",
          "m": "test-alarm-to-get-1-comment-2",
          "t": 1597030221,
          "initiator": "user",
          "val": 0
        },
        "last_event_date": 1597030250,
        "last_update_date": 1597030219,
        "long_output": "",
        "long_output_history": [
          ""
        ],
        "output": "test-alarm-to-get-1-output",
        "parents": [],
        "resource": "test-resource-to-alarm-get-1",
        "pbh_inactive_duration": 0,
        "snooze_duration": 0,
        "state": {
          "_t": "stateinc",
          "a": "test-connector-default.test-connector-default-name",
          "user_id": "",
          "m": "test-alarm-to-get-1-output",
          "t": 1597030219,
          "initiator": "external",
          "val": 3
        },
        "status": {
          "_t": "statusinc",
          "a": "test-connector-default.test-connector-default-name",
          "user_id": "",
          "m": "test-alarm-to-get-1-output",
          "t": 1597030219,
          "initiator": "external",
          "val": 1
        },
        "total_state_changes": 1
      }
    }
    """
    When I do GET /api/v4/alarms/test-alarm-to-get-4
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-alarm-to-get-4",
      "entity": {
        "_id": "test-resource-to-alarm-get-3/test-component-to-alarm-get",
        "category": null,
        "connector": "test-connector-default/test-connector-default-name",
        "component": "test-component-to-alarm-get",
        "enabled": true,
        "old_entity_patterns": null,
        "impact_level": 1,
        "infos": {},
        "name": "test-resource-to-alarm-get-3",
        "type": "resource",
        "depends_count": 0,
        "impacts_count": 0
      },
      "impact_state": 0,
      "infos": {},
      "tags": [],
      "t": 1597030121,
      "v": {
        "children": [],
        "component": "test-component-to-alarm-get",
        "connector": "test-connector-default",
        "connector_name": "test-connector-default-name",
        "creation_date": 1597030121,
        "display_name": "PU-YA-QB",
        "infos": {},
        "infos_rule_version": {},
        "initial_long_output": "",
        "initial_output": "test-alarm-to-get-4-output",
        "last_event_date": 1597030121,
        "last_update_date": 1597030121,
        "resolved": 1597030141,
        "long_output": "",
        "long_output_history": [
          ""
        ],
        "output": "test-alarm-to-get-4-output",
        "parents": [],
        "resource": "test-resource-to-alarm-get-3",
        "duration": 20,
        "current_state_duration": 0,
        "active_duration": 20,
        "pbh_inactive_duration": 0,
        "snooze_duration": 0,
        "state": {
          "_t": "statedec",
          "a": "test-connector-default.test-connector-default-name",
          "user_id": "",
          "m": "test-alarm-to-get-4-output",
          "t": 1597030141,
          "initiator": "external",
          "val": 0
        },
        "status": {
          "_t": "statusdec",
          "a": "test-connector-default.test-connector-default-name",
          "user_id": "",
          "m": "test-alarm-to-get-4-output",
          "t": 1597030141,
          "initiator": "external",
          "val": 0
        },
        "total_state_changes": 1
      }
    }
    """

  @concurrent
  Scenario: given get one request should return not found error
    When I am admin
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given get one unauth request should not allow access
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 401

  @concurrent
  Scenario: given get one request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms/not-exist
    Then the response code should be 403
