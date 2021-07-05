Feature: correlation feature - attribute rule

  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-attribute-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-attribute-correlation-1"
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
      "connector": "test-attribute-1",
      "connector_name": "test-attribute-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-attribute-correlation-1",
      "resource": "test-attribute-correlation-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
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
                "causes": {
                  "rules": [
                    {
                      "name": "test-attribute-correlation-1"
                    }
                  ],
                  "total": 1
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-attribute-correlation-1"
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
    When I send an event:
    """
    {
      "connector": "test-attribute-1",
      "connector_name": "test-attribute-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-attribute-correlation-1",
      "resource": "test-attribute-correlation-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
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
                "causes": {
                  "rules": [
                    {
                      "name": "test-attribute-correlation-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-attribute-correlation-1"
                    }
                  ],
                  "total": 1
                }
              }
            ],
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-attribute-correlation-1"
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
