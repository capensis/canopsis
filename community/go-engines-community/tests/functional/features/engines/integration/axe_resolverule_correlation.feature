Feature: resolve meta alarm
  I need to be able to resolve meta alarm

  Scenario: given meta alarm and resolved child should auto resolve meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-resolverule-correlation-1",
      "type": "attribute",
      "auto_resolve": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-resolverule-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-axe-resolverule-correlation-1",
      "name": "test-resolve-rule-axe-resolverule-correlation-1-name",
      "description": "test-resolve-rule-axe-resolverule-correlation-1-desc",
      "entity_pattern":[
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-resolverule-correlation-1"
            }
          }
        ]
      ],
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "priority": 1
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-resolverule-correlation-1",
      "connector_name": "test-connector-name-axe-resolverule-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-resolverule-correlation-1",
      "resource": "test-resource-axe-resolverule-correlation-1",
      "state": 2,
      "output": "test-output-axe-resolverule-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-resolverule-correlation-1",
      "connector_name": "test-connector-name-axe-resolverule-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-resolverule-correlation-1",
      "resource": "test-resource-axe-resolverule-correlation-1",
      "state": 0,
      "output": "test-output-axe-resolverule-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-resolverule-correlation-1&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "state": {
              "val": 0
            },
            "status": {
              "val": 0
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
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
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
