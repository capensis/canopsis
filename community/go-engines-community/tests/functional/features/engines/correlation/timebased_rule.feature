Feature: correlation feature - timebased rule

  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-timebased-correlation-1",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        }
      },
      "patterns": {
        "entity_patterns": [
          {
            "component": "test-timebased-correlation-1"
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
      "connector": "test-timebased-1",
      "connector_name": "test-timebased-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-1",
      "resource": "test-timebased-correlation-resource-1",
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
      "connector": "test-timebased-1",
      "connector_name": "test-timebased-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-1",
      "resource": "test-timebased-correlation-resource-2",
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
                      "name": "test-timebased-correlation-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-1"
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
            "name": "test-timebased-correlation-1"
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

  Scenario: given meta alarm rule and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-timebased-correlation-2",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        }
      },
      "patterns": {
        "entity_patterns": [
          {
            "component": "test-timebased-correlation-2"
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
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-1",
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
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I wait 4s
    When I send an event:
    """
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-3",
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
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-timebased-2",
      "connector_name": "test-timebased-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-2",
      "resource": "test-timebased-correlation-resource-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=desc
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
                      "name": "test-timebased-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-2"
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
            "name": "test-timebased-correlation-2"
          }
        },
        {
          "consequences": {
            "data": [
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-2"
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
            "name": "test-timebased-correlation-2"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given meta alarm rule and events should create one single meta alarms because first group didn't reached default timebased threshold = 2 alarms
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-timebased-correlation-3",
      "type": "timebased",
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        }
      },
      "patterns": {
        "entity_patterns": [
          {
            "component": "test-timebased-correlation-3"
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
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 4s
    When I send an event:
    """
    {
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-2",
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
      "connector": "test-timebased-3",
      "connector_name": "test-timebased-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-timebased-correlation-3",
      "resource": "test-timebased-correlation-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=desc
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
                      "name": "test-timebased-correlation-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-timebased-correlation-3"
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
            "name": "test-timebased-correlation-3"
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
