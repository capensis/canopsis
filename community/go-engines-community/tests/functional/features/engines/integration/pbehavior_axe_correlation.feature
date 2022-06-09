Feature: update meta alarm on pbehavior
  I need to be able to update meta alarm on pbehavior

  Scenario: given meta alarm and pbehavior should update meta alarm and not update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-pbehavior-axe-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-pbehavior-axe-correlation-1"
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
      "connector": "test-connector-pbehavior-axe-correlation-1",
      "connector_name": "test-connector-name-pbehavior-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-pbehavior-axe-correlation-1",
      "resource": "test-resource-pbehavior-axe-correlation-1",
      "state": 2,
      "output": "test-output-pbehavior-axe-correlation-1",
      "long_output": "test-long-output-pbehavior-axe-correlation-1",
      "author": "test-author-pbehavior-axe-correlation-1",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "name": "test-pbehavior-axe-correlation-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "_id": "{{ .metalarmEntityID }}"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-pbehavior-axe-correlation-1",
                  "connector": "test-connector-pbehavior-axe-correlation-1",
                  "connector_name": "test-connector-name-pbehavior-axe-correlation-1",
                  "resource": "test-resource-pbehavior-axe-correlation-1",
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
                      "val": 0
                    }
                  ]
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-pbehavior-axe-correlation-1"
          },
          "v": {
            "children": [
              "test-resource-pbehavior-axe-correlation-1/test-component-pbehavior-axe-correlation-1"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "pbehavior_info": {
              "canonical_type": "maintenance",
              "name": "test-pbehavior-axe-correlation-1"
            },
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
                "_t": "pbhenter",
                "m": "Pbehavior test-pbehavior-axe-correlation-1. Type: Engine maintenance. Reason: Test Engine.",
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
