Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  Scenario: given meta alarm and ack event should double ack children if they're already acked
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-8",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-8"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-1",
      "state": 2,
      "output": "test-output-axe-correlation-8",
      "long_output": "test-long-output-axe-correlation-8",
      "author": "test-author-axe-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-2",
      "state": 2,
      "output": "test-output-axe-correlation-8",
      "long_output": "test-long-output-axe-correlation-8",
      "author": "test-author-axe-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-2",
      "output": "previous ack",
      "author": "test-author-axe-correlation-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "metaalarm ack",
      "long_output": "test-long-output-axe-correlation-8",
      "author": "test-author-axe-correlation-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-8"
          },
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [
              "test-resource-axe-correlation-8-1/test-component-axe-correlation-8",
              "test-resource-axe-correlation-8-2/test-component-axe-correlation-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-8",
                "m": "metaalarm ack",
                "val": 0
              }
            ]
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-8-1"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [],
            "component": "test-component-axe-correlation-8",
            "connector": "test-connector-axe-correlation-8",
            "connector_name": "test-connector-name-axe-correlation-8",
            "initial_long_output": "test-long-output-axe-correlation-8",
            "initial_output": "test-output-axe-correlation-8",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-8-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-8",
                "m": "metaalarm ack",
                "val": 0
              }
            ]
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-8-2"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-correlation-8",
              "m": "metaalarm ack",
              "val": 0
            },
            "children": [],
            "component": "test-component-axe-correlation-8",
            "connector": "test-connector-axe-correlation-8",
            "connector_name": "test-connector-name-axe-correlation-8",
            "initial_long_output": "test-long-output-axe-correlation-8",
            "initial_output": "test-output-axe-correlation-8",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-8-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-8",
                "m": "previous ack",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-8",
                "m": "metaalarm ack",
                "val": 0
              }
            ]
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
