Feature: correlation feature - complex rule with threshold rate

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-1-resource&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-1-connector",
                  "connector_name": "test-complex-rule-rate-1-connectorname",
                  "component":  "test-complex-rule-rate-1-component",
                  "resource": "test-complex-rule-rate-1-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-1-connector",
                  "connector_name": "test-complex-rule-rate-1-connectorname",
                  "component":  "test-complex-rule-rate-1-component",
                  "resource": "test-complex-rule-rate-1-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-1-connector",
                  "connector_name": "test-complex-rule-rate-1-connectorname",
                  "component":  "test-complex-rule-rate-1-component",
                  "resource": "test-complex-rule-rate-1-resource-3"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events shouldn't create meta alarm because rate was recomputed by the new entity event after, metaalarm should be create after reaching the new rate
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-2-resource&correlation=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-complex-rule-rate-2-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-2-resource-2"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-2-resource-new"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-2-resource&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 4,
          "meta_alarm_rule": {
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-2-connector",
                  "connector_name": "test-complex-rule-rate-2-connectorname",
                  "component":  "test-complex-rule-rate-2-component",
                  "resource": "test-complex-rule-rate-2-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-2-connector",
                  "connector_name": "test-complex-rule-rate-2-connectorname",
                  "component":  "test-complex-rule-rate-2-component",
                  "resource": "test-complex-rule-rate-2-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-2-connector",
                  "connector_name": "test-complex-rule-rate-2-connectorname",
                  "component":  "test-complex-rule-rate-2-component",
                  "resource": "test-complex-rule-rate-2-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-2-connector",
                  "connector_name": "test-complex-rule-rate-2-connectorname",
                  "component":  "test-complex-rule-rate-2-component",
                  "resource": "test-complex-rule-rate-2-resource-new"
                }
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

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events should create one single meta alarms because first group didn't reached threshold
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-3-resource&correlation=true&multi_sort[]=v.meta,desc&multi_sort[]=v.resource,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-threshold-rate-3"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-3-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-3-resource-2"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-3-connector",
                  "connector_name": "test-complex-rule-rate-3-connectorname",
                  "component":  "test-complex-rule-rate-3-component",
                  "resource": "test-complex-rule-rate-3-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-3-connector",
                  "connector_name": "test-complex-rule-rate-3-connectorname",
                  "component":  "test-complex-rule-rate-3-component",
                  "resource": "test-complex-rule-rate-3-resource-4"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-3-connector",
                  "connector_name": "test-complex-rule-rate-3-connectorname",
                  "component":  "test-complex-rule-rate-3-component",
                  "resource": "test-complex-rule-rate-3-resource-5"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-4-resource&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        }
      ]
    }
    """
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-4-resource&correlation=true&sort_by=t&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-threshold-rate-4"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ (index .lastResponse.data 1)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-4-connector",
                  "connector_name": "test-complex-rule-rate-4-connectorname",
                  "component":  "test-complex-rule-rate-4-component",
                  "resource": "test-complex-rule-rate-4-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-4-connector",
                  "connector_name": "test-complex-rule-rate-4-connectorname",
                  "component":  "test-complex-rule-rate-4-component",
                  "resource": "test-complex-rule-rate-4-resource-2"
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
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-4-connector",
                  "connector_name": "test-complex-rule-rate-4-connectorname",
                  "component":  "test-complex-rule-rate-4-component",
                  "resource": "test-complex-rule-rate-4-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-4-connector",
                  "connector_name": "test-complex-rule-rate-4-connectorname",
                  "component":  "test-complex-rule-rate-4-component",
                  "resource": "test-complex-rule-rate-4-resource-4"
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
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events should create one single meta alarm without first alarm, because interval shifting
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
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
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
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
    When I send an event and wait the end of event processing:
    """json
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
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-5-resource&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-threshold-rate-5"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-5-resource-1"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-5-connector",
                  "connector_name": "test-complex-rule-rate-5-connectorname",
                  "component":  "test-complex-rule-rate-5-component",
                  "resource": "test-complex-rule-rate-5-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-5-connector",
                  "connector_name": "test-complex-rule-rate-5-connectorname",
                  "component":  "test-complex-rule-rate-5-component",
                  "resource": "test-complex-rule-rate-5-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-5-connector",
                  "connector_name": "test-complex-rule-rate-5-connectorname",
                  "component":  "test-complex-rule-rate-5-component",
                  "resource": "test-complex-rule-rate-5-resource-4"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm rule with threshold rate and events should create meta alarm regarding total entity pattern
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-complex-correlation-threshold-rate-6",
      "type": "complex",
      "total_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-6-resource-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-complex-rule-rate-6-resource-2"
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
              "value": "test-complex-rule-rate-6-resource-2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-6-connector",
      "connector_name": "test-complex-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-6-component",
      "resource": "test-complex-rule-rate-6-resource-2-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-6-connector",
      "connector_name": "test-complex-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-6-component",
      "resource": "test-complex-rule-rate-6-resource-2-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-6-connector",
      "connector_name": "test-complex-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-6-component",
      "resource": "test-complex-rule-rate-6-resource-2-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-6-resource&correlation=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-complex-rule-rate-6-resource-2-1"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-6-resource-2-2"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-6-resource-2-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-6-connector",
      "connector_name": "test-complex-rule-rate-6-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-6-component",
      "resource": "test-complex-rule-rate-6-resource-2-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-6-resource&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 4,
          "meta_alarm_rule": {
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-6-connector",
                  "connector_name": "test-complex-rule-rate-6-connectorname",
                  "component":  "test-complex-rule-rate-6-component",
                  "resource": "test-complex-rule-rate-6-resource-2-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-6-connector",
                  "connector_name": "test-complex-rule-rate-6-connectorname",
                  "component":  "test-complex-rule-rate-6-component",
                  "resource": "test-complex-rule-rate-6-resource-2-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-6-connector",
                  "connector_name": "test-complex-rule-rate-6-connectorname",
                  "component":  "test-complex-rule-rate-6-component",
                  "resource": "test-complex-rule-rate-6-resource-2-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-6-connector",
                  "connector_name": "test-complex-rule-rate-6-connectorname",
                  "component":  "test-complex-rule-rate-6-component",
                  "resource": "test-complex-rule-rate-6-resource-2-4"
                }
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

  @concurrent
  Scenario: given meta alarm rule with threshold rate and old event patterns should create meta alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-1-component",
      "resource": "test-complex-rule-rate-backward-compatibility-1-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-1-component",
      "resource": "test-complex-rule-rate-backward-compatibility-1-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-1-component",
      "resource": "test-complex-rule-rate-backward-compatibility-1-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-backward-compatibility-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-complex-rule-rate-backward-compatibility-1-name"
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
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-1-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-1-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-1-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-1-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-1-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-1-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-1-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-1-resource-3"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm rule with threshold rate and old total event patterns should create meta alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-2-component",
      "resource": "test-complex-rule-rate-backward-compatibility-2-resource-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-2-component",
      "resource": "test-complex-rule-rate-backward-compatibility-2-resource-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-2-component",
      "resource": "test-complex-rule-rate-backward-compatibility-2-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-2-component",
      "resource": "test-complex-rule-rate-backward-compatibility-2-resource-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-backward-compatibility-2&correlation=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-complex-rule-rate-backward-compatibility-2-resource-1"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-backward-compatibility-2-resource-2"
          }
        },
        {
          "v": {
            "resource": "test-complex-rule-rate-backward-compatibility-2-resource-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
      "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-backward-compatibility-2-component",
      "resource": "test-complex-rule-rate-backward-compatibility-2-resource-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-backward-compatibility-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 4,
          "meta_alarm_rule": {
            "name": "test-complex-rule-rate-backward-compatibility-2-name"
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
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-2-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-2-resource-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-2-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-2-resource-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-2-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-2-resource-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-rule-rate-backward-compatibility-2-connector",
                  "connector_name": "test-complex-rule-rate-backward-compatibility-2-connectorname",
                  "component":  "test-complex-rule-rate-backward-compatibility-2-component",
                  "resource": "test-complex-rule-rate-backward-compatibility-2-resource-4"
                }
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

  @concurrent
  Scenario: given deleted meta alarm rule should delete meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-complex-correlation-threshold-rate-7",
      "type": "complex",
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
        "threshold_rate": 0.6
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-1",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-2",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-rule-rate-7-connector",
      "connector_name": "test-complex-rule-rate-7-connectorname",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-rule-rate-7-component",
      "resource": "test-complex-rule-rate-7-resource-3",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-7-resource&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-threshold-rate-7"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do DELETE /api/v4/cat/metaalarmrules/{{ .metaAlarmRuleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms/{{ .metaAlarmID }}
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-complex-rule-rate-7-component&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-complex-rule-rate-7-connector",
            "connector_name": "test-complex-rule-rate-7-connectorname",
            "component": "test-complex-rule-rate-7-component",
            "resource": "test-complex-rule-rate-7-resource-1",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-complex-rule-rate-7-connector",
            "connector_name": "test-complex-rule-rate-7-connectorname",
            "component": "test-complex-rule-rate-7-component",
            "resource": "test-complex-rule-rate-7-resource-2",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-complex-rule-rate-7-connector",
            "connector_name": "test-complex-rule-rate-7-connectorname",
            "component": "test-complex-rule-rate-7-component",
            "resource": "test-complex-rule-rate-7-resource-3",
            "parents": []
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  @concurrent
  Scenario: given meta alarm and removed child should update meta alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-8",
        "connector_name": "test-connector-name-correlation-complex-rate-8",
        "component":  "test-component-correlation-complex-rate-8",
        "resource": "test-resource-correlation-complex-rate-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-8",
        "connector_name": "test-connector-name-correlation-complex-rate-8",
        "component":  "test-component-correlation-complex-rate-8",
        "resource": "test-resource-correlation-complex-rate-8-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-8",
        "connector_name": "test-connector-name-correlation-complex-rate-8",
        "component":  "test-component-correlation-complex-rate-8",
        "resource": "test-resource-correlation-complex-rate-8-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-8",
        "connector_name": "test-connector-name-correlation-complex-rate-8",
        "component":  "test-component-correlation-complex-rate-8",
        "resource": "test-resource-correlation-complex-rate-8-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-8",
        "connector_name": "test-connector-name-correlation-complex-rate-8",
        "component":  "test-component-correlation-complex-rate-8",
        "resource": "test-resource-correlation-complex-rate-8-5",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-complex-rate-8",
      "type": "complex",
      "output_template": "{{ `{{ .Count }}` }}",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-correlation-complex-rate-8"
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
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-8",
      "connector_name": "test-connector-name-correlation-complex-rate-8",
      "component":  "test-component-correlation-complex-rate-8",
      "resource": "test-resource-correlation-complex-rate-8-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-complex-rate-8",
      "connector_name": "test-connector-name-correlation-complex-rate-8",
      "component":  "test-component-correlation-complex-rate-8",
      "resource": "test-resource-correlation-complex-rate-8-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-8",
      "connector_name": "test-connector-name-correlation-complex-rate-8",
      "component":  "test-component-correlation-complex-rate-8",
      "resource": "test-resource-correlation-complex-rate-8-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-8&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-8"
          },
          "v": {
            "output": "3",
            "state": {
              "val": 3
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
    When I save response metaAlarmId={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-correlation-complex-rate-8",
                  "connector_name": "test-connector-name-correlation-complex-rate-8",
                  "component":  "test-component-correlation-complex-rate-8",
                  "resource": "test-resource-correlation-complex-rate-8-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-rate-8",
                  "connector_name": "test-connector-name-correlation-complex-rate-8",
                  "component":  "test-component-correlation-complex-rate-8",
                  "resource": "test-resource-correlation-complex-rate-8-2"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-rate-8",
                  "connector_name": "test-connector-name-correlation-complex-rate-8",
                  "component":  "test-component-correlation-complex-rate-8",
                  "resource": "test-resource-correlation-complex-rate-8-3"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-complex-rate-8-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-8&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-8"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "connector": "test-connector-correlation-complex-rate-8",
            "connector_name": "test-connector-name-correlation-complex-rate-8",
            "component":  "test-component-correlation-complex-rate-8",
            "resource": "test-resource-correlation-complex-rate-8-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-connector-correlation-complex-rate-8",
                  "connector_name": "test-connector-name-correlation-complex-rate-8",
                  "component":  "test-component-correlation-complex-rate-8",
                  "resource": "test-resource-correlation-complex-rate-8-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-rate-8",
                  "connector_name": "test-connector-name-correlation-complex-rate-8",
                  "component":  "test-component-correlation-complex-rate-8",
                  "resource": "test-resource-correlation-complex-rate-8-3"
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
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and removed child should not add child to meta alarm again but add to new meta alarm on change
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-3",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-4",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-5",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 0,
        "connector": "test-connector-correlation-complex-rate-9",
        "connector_name": "test-connector-name-correlation-complex-rate-9",
        "component":  "test-component-correlation-complex-rate-9",
        "resource": "test-resource-correlation-complex-rate-9-6",
        "source_type": "resource"
      }
    ]
    """
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-complex-rate-9",
      "type": "complex",
      "output_template": "{{ `{{ .Count }}` }}",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "regexp",
              "value": "test-resource-correlation-complex-rate-9"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_rate": 0.5
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "3",
            "state": {
              "val": 3
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
    When I save response metaAlarmId1={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId1 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    Then the response code should be 207
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId1 }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-complex-rate-9-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-9&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-complex-rate-9-2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-9&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-complex-rate-9-2"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-9&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "3",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-complex-rate-9-2"
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
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-5",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-rate-9",
      "connector_name": "test-connector-name-correlation-complex-rate-9",
      "component":  "test-component-correlation-complex-rate-9",
      "resource": "test-resource-correlation-complex-rate-9-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-rate-9&correlation=true&sort_by=t&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "3",
            "state": {
              "val": 3
            }
          }
        },
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-rate-9"
          },
          "v": {
            "output": "3",
            "state": {
              "val": 2
            }
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
    When I save response metaAlarmId2={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmId2 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      },
      {
        "_id": "{{ .metaAlarmId1 }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
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
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-5"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-6"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-1"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-3"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-rate-9-4"
                }
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 3
            }
          }
        }
      }
    ]
    """
