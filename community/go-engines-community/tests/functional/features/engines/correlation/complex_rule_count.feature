Feature: correlation feature - complex rule with threshold count

  @concurrent
  Scenario: given meta alarm rule with threshold count and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-1",
      "connector_name": "test-complex-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-1",
      "resource": "test-complex-correlation-resource-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
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
                  "connector": "test-complex-1",
                  "connector_name": "test-complex-1-name",
                  "component": "test-complex-correlation-1",
                  "resource": "test-complex-correlation-resource-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-1",
                  "connector_name": "test-complex-1-name",
                  "component": "test-complex-correlation-1",
                  "resource": "test-complex-correlation-resource-1-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-1",
                  "connector_name": "test-complex-1-name",
                  "component": "test-complex-correlation-1",
                  "resource": "test-complex-correlation-resource-1-3"
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
  Scenario: given meta alarm rule with threshold count and events should create 2 meta alarms because of 2 separate time intervals
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        }
      ]
    }
    """
    When I wait 4s
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-6",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3
        },
        {
          "children": 3
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-2",
      "connector_name": "test-complex-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-2",
      "resource": "test-complex-correlation-resource-2-7",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-2&correlation=true&sort_by=t&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 4,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-2"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
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
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-4"
                }
              },
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-5"
                }
              },
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-6"
                }
              },
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-7"
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
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-1"
                }
              },
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-2",
                  "connector_name": "test-complex-2-name",
                  "component": "test-complex-correlation-2",
                  "resource": "test-complex-correlation-resource-2-3"
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
  Scenario: given meta alarm rule with threshold count and events should create one single meta alarms because first group didn't reached threshold
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3-2",
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
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-3",
      "connector_name": "test-complex-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-3",
      "resource": "test-complex-correlation-resource-3-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-3&correlation=true&multi_sort[]=v.meta,desc&multi_sort[]=v.resource,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-3"
          }
        },
        {
          "v": {
            "resource": "test-complex-correlation-resource-3-1"
          }
        },
        {
          "v": {
            "resource": "test-complex-correlation-resource-3-2"
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
                  "connector": "test-complex-3",
                  "connector_name": "test-complex-3-name",
                  "component": "test-complex-correlation-3",
                  "resource": "test-complex-correlation-resource-3-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-3",
                  "connector_name": "test-complex-3-name",
                  "component": "test-complex-correlation-3",
                  "resource": "test-complex-correlation-resource-3-4"
                }
              },
              {
                "v": {
                  "connector": "test-complex-3",
                  "connector_name": "test-complex-3-name",
                  "component": "test-complex-correlation-3",
                  "resource": "test-complex-correlation-resource-3-5"
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
  Scenario: given meta alarm rule with threshold count and events should create one single meta alarm without first alarm, because interval shifting
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-4-1",
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
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-4-2",
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
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-4-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-4",
      "connector_name": "test-complex-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-4",
      "resource": "test-complex-correlation-resource-4-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-4&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-4"
          }
        },
        {
          "v": {
            "resource": "test-complex-correlation-resource-4-1"
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
                  "connector": "test-complex-4",
                  "connector_name": "test-complex-4-name",
                  "component": "test-complex-correlation-4",
                  "resource": "test-complex-correlation-resource-4-2"
                }
              },
              {
                "v": {
                  "connector": "test-complex-4",
                  "connector_name": "test-complex-4-name",
                  "component": "test-complex-correlation-4",
                  "resource": "test-complex-correlation-resource-4-3"
                }
              },
              {
                "v": {
                  "connector": "test-complex-4",
                  "connector_name": "test-complex-4-name",
                  "component": "test-complex-correlation-4",
                  "resource": "test-complex-correlation-resource-4-4"
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
  Scenario: given meta alarm rule with threshold count and events should create meta alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
      "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-metaalarm-rule-backward-compatibility-component-1",
      "resource": "test-metaalarm-rule-backward-compatibility-resource-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-metaalarm-rule-backward-compatibility-resource-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
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
                  "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
                  "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
                  "component":  "test-metaalarm-rule-backward-compatibility-component-1",
                  "resource": "test-metaalarm-rule-backward-compatibility-resource-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
                  "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
                  "component":  "test-metaalarm-rule-backward-compatibility-component-1",
                  "resource": "test-metaalarm-rule-backward-compatibility-resource-1-2"
                }
              },
              {
                "v": {
                  "connector": "test-metaalarm-rule-backward-compatibility-connector-1",
                  "connector_name": "test-metaalarm-rule-backward-compatibility-connector-name-1",
                  "component":  "test-metaalarm-rule-backward-compatibility-component-1",
                  "resource": "test-metaalarm-rule-backward-compatibility-resource-1-3"
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
  Scenario: given deleted meta alarm rule should delete meta alarms
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-complex-correlation-5",
      "type": "complex",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-complex-correlation-5"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 2
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-5",
      "connector_name": "test-complex-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-5",
      "resource": "test-complex-correlation-resource-5-1",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-5",
      "connector_name": "test-complex-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-5",
      "resource": "test-complex-correlation-resource-5-2",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-5&correlation=true until response code is 200 and body contains:
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
      "connector": "test-complex-5",
      "connector_name": "test-complex-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-5",
      "resource": "test-complex-correlation-resource-5-3",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-complex-5",
      "connector_name": "test-complex-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-complex-correlation-5",
      "resource": "test-complex-correlation-resource-5-4",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-complex-correlation-resource-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-5"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-complex-correlation-5"
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
    When I save response metaAlarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmId2={{ (index .lastResponse.data 1)._id }}
    When I do DELETE /api/v4/cat/metaalarmrules/{{ .metaAlarmRuleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms/{{ .metaAlarmId1 }}
    Then the response code should be 404
    When I do GET /api/v4/alarms/{{ .metaAlarmId2 }}
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-complex-correlation-5&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-complex-5",
            "connector_name": "test-complex-5-name",
            "component":  "test-complex-correlation-5",
            "resource": "test-complex-correlation-resource-5-1",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-complex-5",
            "connector_name": "test-complex-5-name",
            "component":  "test-complex-correlation-5",
            "resource": "test-complex-correlation-resource-5-2",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-complex-5",
            "connector_name": "test-complex-5-name",
            "component":  "test-complex-correlation-5",
            "resource": "test-complex-correlation-resource-5-3",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-complex-5",
            "connector_name": "test-complex-5-name",
            "component":  "test-complex-correlation-5",
            "resource": "test-complex-correlation-resource-5-4",
            "parents": []
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
    """

  @concurrent
  Scenario: given meta alarm and removed child should update meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-complex-count-6",
      "type": "complex",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-complex-count-6"
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
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-6",
      "connector_name": "test-connector-name-correlation-complex-count-6",
      "component":  "test-component-correlation-complex-count-6",
      "resource": "test-resource-correlation-complex-count-6-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-complex-count-6",
      "connector_name": "test-connector-name-correlation-complex-count-6",
      "component":  "test-component-correlation-complex-count-6",
      "resource": "test-resource-correlation-complex-count-6-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-6",
      "connector_name": "test-connector-name-correlation-complex-count-6",
      "component":  "test-component-correlation-complex-count-6",
      "resource": "test-resource-correlation-complex-count-6-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-6&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-6"
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
                  "connector": "test-connector-correlation-complex-count-6",
                  "connector_name": "test-connector-name-correlation-complex-count-6",
                  "component":  "test-component-correlation-complex-count-6",
                  "resource": "test-resource-correlation-complex-count-6-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-count-6",
                  "connector_name": "test-connector-name-correlation-complex-count-6",
                  "component":  "test-component-correlation-complex-count-6",
                  "resource": "test-resource-correlation-complex-count-6-2"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-count-6",
                  "connector_name": "test-connector-name-correlation-complex-count-6",
                  "component":  "test-component-correlation-complex-count-6",
                  "resource": "test-resource-correlation-complex-count-6-3"
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
      "comment": "test-metaalarmrule-correlation-complex-count-6-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-6&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-6"
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
            "connector": "test-connector-correlation-complex-count-6",
            "connector_name": "test-connector-name-correlation-complex-count-6",
            "component":  "test-component-correlation-complex-count-6",
            "resource": "test-resource-correlation-complex-count-6-2"
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
                  "connector": "test-connector-correlation-complex-count-6",
                  "connector_name": "test-connector-name-correlation-complex-count-6",
                  "component":  "test-component-correlation-complex-count-6",
                  "resource": "test-resource-correlation-complex-count-6-1"
                }
              },
              {
                "v": {
                  "connector": "test-connector-correlation-complex-count-6",
                  "connector_name": "test-connector-name-correlation-complex-count-6",
                  "component":  "test-component-correlation-complex-count-6",
                  "resource": "test-resource-correlation-complex-count-6-3"
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
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-complex-count-7",
      "type": "complex",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-complex-count-7"
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
    Then I save response metaAlarmRuleId={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-7&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
      "comment": "test-metaalarmrule-correlation-complex-count-7-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-7&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
            "resource": "test-resource-correlation-complex-count-7-2"
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
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component": "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-7&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
            "resource": "test-resource-correlation-complex-count-7-2"
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
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-7&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
            "resource": "test-resource-correlation-complex-count-7-2"
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
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-5",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-complex-count-7",
      "connector_name": "test-connector-name-correlation-complex-count-7",
      "component":  "test-component-correlation-complex-count-7",
      "resource": "test-resource-correlation-complex-count-7-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-complex-count-7&correlation=true&sort_by=t&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
            "name": "test-metaalarmrule-correlation-complex-count-7"
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
                  "resource": "test-resource-correlation-complex-count-7-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-count-7-5"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-count-7-6"
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
                  "resource": "test-resource-correlation-complex-count-7-1"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-count-7-3"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-complex-count-7-4"
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
