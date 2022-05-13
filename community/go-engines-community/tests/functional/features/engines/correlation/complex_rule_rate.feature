Feature: correlation feature - complex rule with threshold rate

  Scenario: given meta alarm rule with threshold rate and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-1",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-1-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_rate": 0.6
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-1-connector",
      "connector_name": "test-complex-rule-rate-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-1-component",
      "resource": "test-complex-rule-rate-1-resource-1",
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
      "connector": "test-complex-rule-rate-1-connector",
      "connector_name": "test-complex-rule-rate-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-1-component",
      "resource": "test-complex-rule-rate-1-resource-2",
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
      "connector": "test-complex-rule-rate-1-connector",
      "connector_name": "test-complex-rule-rate-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-1-component",
      "resource": "test-complex-rule-rate-1-resource-3",
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
                      "name": "test-complex-correlation-threshold-rate-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-1"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-1"
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
            "name": "test-complex-correlation-threshold-rate-1"
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

  Scenario: given meta alarm rule with threshold rate and events shouldn't create meta alarm because rate was recomputed by the new entity event after, metaalarm should be create after reaching the new rate
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-2",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-2-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_rate": 0.6
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-2-connector",
      "connector_name": "test-complex-rule-rate-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-2-component",
      "resource": "test-complex-rule-rate-2-resource-1",
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
      "connector": "test-complex-rule-rate-2-connector",
      "connector_name": "test-complex-rule-rate-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-2-component",
      "resource": "test-complex-rule-rate-2-resource-2",
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
      "connector": "test-complex-rule-rate-2-connector",
      "connector_name": "test-complex-rule-rate-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-2-component",
      "resource": "test-complex-rule-rate-2-resource-new",
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
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-2-connector",
      "connector_name": "test-complex-rule-rate-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-2-component",
      "resource": "test-complex-rule-rate-2-resource-3",
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
                      "name": "test-complex-correlation-threshold-rate-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-2"
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
            "name": "test-complex-correlation-threshold-rate-2"
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

  Scenario: given meta alarm rule with threshold rate and events should create one single meta alarms because first group didn't reached threshold
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-3",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-3-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_rate": 0.6
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-3-connector",
      "connector_name": "test-complex-rule-rate-3-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-3-component",
      "resource": "test-complex-rule-rate-3-resource-1",
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
      "connector": "test-complex-rule-rate-3-connector",
      "connector_name": "test-complex-rule-rate-3-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-3-component",
      "resource": "test-complex-rule-rate-3-resource-2",
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
      "connector": "test-complex-rule-rate-3-connector",
      "connector_name": "test-complex-rule-rate-3-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-3-component",
      "resource": "test-complex-rule-rate-3-resource-3",
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
      "connector": "test-complex-rule-rate-3-connector",
      "connector_name": "test-complex-rule-rate-3-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-3-component",
      "resource": "test-complex-rule-rate-3-resource-4",
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
      "connector": "test-complex-rule-rate-3-connector",
      "connector_name": "test-complex-rule-rate-3-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-3-component",
      "resource": "test-complex-rule-rate-3-resource-5",
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
                      "name": "test-complex-correlation-threshold-rate-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-3"
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
            "name": "test-complex-correlation-threshold-rate-3"
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

  Scenario: given meta alarm rule with threshold rate and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-4",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-4-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_rate": 0.4
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-4-connector",
      "connector_name": "test-complex-rule-rate-4-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-4-component",
      "resource": "test-complex-rule-rate-4-resource-1",
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
      "connector": "test-complex-rule-rate-4-connector",
      "connector_name": "test-complex-rule-rate-4-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-4-component",
      "resource": "test-complex-rule-rate-4-resource-2",
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
      "connector": "test-complex-rule-rate-4-connector",
      "connector_name": "test-complex-rule-rate-4-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-4-component",
      "resource": "test-complex-rule-rate-4-resource-3",
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
      "connector": "test-complex-rule-rate-4-connector",
      "connector_name": "test-complex-rule-rate-4-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-4-component",
      "resource": "test-complex-rule-rate-4-resource-3",
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
      "connector": "test-complex-rule-rate-4-connector",
      "connector_name": "test-complex-rule-rate-4-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-4-component",
      "resource": "test-complex-rule-rate-4-resource-4",
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
                      "name": "test-complex-correlation-threshold-rate-4"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-4"
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
            "name": "test-complex-correlation-threshold-rate-4"
          }
        },
        {
          "consequences": {
            "data": [
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-4"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-4"
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
            "name": "test-complex-correlation-threshold-rate-4"
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

  Scenario: given meta alarm rule with threshold rate and events should create one single meta alarm without first alarm, because interval shifting
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-5",
      "type": "complex",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-5-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_rate": 0.6
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-5-connector",
      "connector_name": "test-complex-rule-rate-5-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-5-component",
      "resource": "test-complex-rule-rate-5-resource-1",
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
      "connector": "test-complex-rule-rate-5-connector",
      "connector_name": "test-complex-rule-rate-5-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-5-component",
      "resource": "test-complex-rule-rate-5-resource-2",
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
      "connector": "test-complex-rule-rate-5-connector",
      "connector_name": "test-complex-rule-rate-5-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-5-component",
      "resource": "test-complex-rule-rate-5-resource-3",
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
      "connector": "test-complex-rule-rate-5-connector",
      "connector_name": "test-complex-rule-rate-5-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-5-component",
      "resource": "test-complex-rule-rate-5-resource-4",
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
                      "name": "test-complex-correlation-threshold-rate-5"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-5"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-5"
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
            "name": "test-complex-correlation-threshold-rate-5"
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

  Scenario: given meta alarm rule with threshold rate and events should create meta alarm regarding total entity pattern
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-complex-correlation-threshold-rate-6",
      "type": "complex",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-6-resource"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-7-resource"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-7-resource"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        },
        "threshold_rate": 0.4
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-1",
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
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-2",
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
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-3",
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
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I send an event:
    """
    {
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-4",
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
                      "name": "test-complex-correlation-threshold-rate-6"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-6"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-6"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-complex-correlation-threshold-rate-6"
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
            "name": "test-complex-correlation-threshold-rate-6"
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
