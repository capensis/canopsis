Feature: correlation feature - valuegroup rule with threshold rate


  Scenario: given meta alarm rule with threshold rate and old event patterns should create meta alarm
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-1-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-1-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-1-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-1-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-1-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-1-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"test-valuegroup-rule-rate-backward-compatibility-1"}]}&with_steps=true&with_consequences=true&correlation=true
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
                      "name": "test-valuegroup-rule-rate-backward-compatibility-1-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-valuegroup-rule-rate-backward-compatibility-1-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-valuegroup-rule-rate-backward-compatibility-1-name"
                    }
                  ],
                  "total": 1
                }
              }
            ],
            "total": 3
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-rule-rate-backward-compatibility-1-name"
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

  Scenario: given meta alarm rule with threshold rate and old total event patterns should create meta alarm
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-2-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-2-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-2-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-2-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-2-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-2-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-1-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-1-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """
    {
      "connector": "test-valuegroup-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-valuegroup-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-valuegroup-rule-rate-backward-compatibility-2-component",
      "resource": "test-valuegroup-rule-rate-backward-compatibility-2-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """    
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"test-valuegroup-rule-rate-backward-compatibility-2"}]}&with_steps=true&with_consequences=true&correlation=true
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
                      "name": "test-valuegroup-rule-rate-backward-compatibility-2-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-valuegroup-rule-rate-backward-compatibility-2-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-valuegroup-rule-rate-backward-compatibility-2-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-valuegroup-rule-rate-backward-compatibility-2-name"
                    }
                  ],
                  "total": 1
                }
              }
            ],
            "total": 4
          },
          "metaalarm": true,
          "rule": {
            "name": "test-valuegroup-rule-rate-backward-compatibility-2-name"
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
