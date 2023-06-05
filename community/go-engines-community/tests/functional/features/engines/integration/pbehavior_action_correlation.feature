Feature: update meta alarm on action
  I need to be able to update meta alarm on action

  Scenario: given meta alarm and pbehavior action should update meta alarm and not update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-pbehavior-action-correlation-1",
      "type": "attribute",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-action-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-action-correlation-1",
      "connector_name": "test-connector-name-pbehavior-action-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-action-correlation-1",
      "resource": "test-resource-pbehavior-action-correlation-1",
      "state": 2,
      "output": "test-output-pbehavior-action-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-action-correlation-1&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do POST /api/v4/scenarios:
    """json
    {
      "name": "test-scenario-pbehavior-action-correlation-1-name",
      "priority": 10066,
      "enabled": true,
      "triggers": ["comment"],
      "actions": [
        {
          "_id": "test-action-pbehavior-action-correlation-1",
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "{{ .metaAlarmResource }}"
                }
              }
            ]
          ],
          "parameters": {
            "name": "test-pbehavior-action-correlation-1",
            "tstart": {{ now }},
            "tstop": {{ nowAdd "1h" }},
            "type": "test-maintenance-type-to-engine",
            "reason": "test-reason-to-engine"
          },
          "type": "pbehavior",
          "drop_scenario_if_not_matched": false,
          "emit_trigger": false
        }
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-pbehavior-action-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-action-correlation-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-pbehavior-action-correlation-1"
          },
          "v": {
            "children": [
              "test-resource-pbehavior-action-correlation-1/test-component-pbehavior-action-correlation-1"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-action-correlation-1"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "comment"
              },
              {
                "_t": "pbhenter",
                "m": "Pbehavior test-pbehavior-action-correlation-1. Type: Engine maintenance. Reason: Test Engine.",
                "val": 0
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-action-correlation-1
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "val": 0
              },
              {
                "_t": "comment"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 4
            }
          }
        }
      }
    ]
    """
