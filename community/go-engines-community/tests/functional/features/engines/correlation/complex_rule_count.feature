Feature: correlation feature - complex rule with threshold count

  Scenario: given meta alarm rule with threshold count and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-1",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-complex-correlation-1"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-1",
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
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-2",
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
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-3",
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
                      "name": "test-complex-correlation-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-1"
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
            "name": "test-complex-correlation-1"
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

  Scenario: given meta alarm rule with threshold count and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-2",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-complex-correlation-2"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-1",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-3",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-4",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-5",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-6",
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
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-7",
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
                      "name": "test-complex-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
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
            "name": "test-complex-correlation-2"
          }
        },
        {
          "consequences": {
            "data": [
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-2"
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
            "name": "test-complex-correlation-2"
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

  Scenario: given meta alarm rule with threshold count and events should create one single meta alarms because first group didn't reached threshold
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-3",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-complex-correlation-3"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-1",
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
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-2",
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
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3",
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
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-4",
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
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-5",
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
                      "name": "test-complex-correlation-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-3"
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
            "name": "test-complex-correlation-3"
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

  Scenario: given meta alarm rule with threshold count and events should create one single meta alarm without first alarm, because interval shifting
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-4",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-complex-correlation-4"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_count": 3
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """
    {
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """
    {
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-3",
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
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-4",
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
                      "name": "test-complex-correlation-4"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-4"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-4"
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
            "name": "test-complex-correlation-4"
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

  Scenario: given meta alarm rule with threshold count and events should create meta alarm
    Given I am admin
    When I send an event:
    """
    {
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-1",
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
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-2",
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
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"test-metaalarm-rule-backward-compatibility-1"}]}&with_steps=true&with_consequences=true&correlation=true
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
                      "name": "test-metaalarm-rule-backward-compatibility-1-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-metaalarm-rule-backward-compatibility-1-name"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-metaalarm-rule-backward-compatibility-1-name"
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
            "name": "test-metaalarm-rule-backward-compatibility-1-name"
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
