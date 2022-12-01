Feature: correlation feature - corel rule

  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-1",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent",
        "corel_child":  "child"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-1",
      "connector_name": "test-corel-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-1",
      "connector_name": "test-corel-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "child",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-1"
          },
          "v": {
            "resource": "test-1",
            "component": "parent",
            "children": [
              "test-2/child"
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
                  "connector": "test-corel-1",
                  "connector_name": "test-corel-1-name",
                  "component": "child",
                  "resource": "test-2"
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
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events shouldn't create meta alarm without a parent, after parent event metaalarm should contain all children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-2",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-2"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-2"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-2",
        "corel_child":  "child-2"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-2",
      "connector_name": "test-corel-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-2",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-2",
      "connector_name": "test-corel-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-2",
      "resource": "test-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-2",
      "connector_name": "test-corel-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-2",
      "resource": "test-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search=test-corel-2&correlation=true&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-3"
          }
        },
        {
          "v": {
            "resource": "test-4"
          }
        },
        {
          "v": {
            "resource": "test-5"
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
    When I send an event:
    """json
    {
      "connector": "test-corel-2",
      "connector_name": "test-corel-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-2",
      "resource": "test-6",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-2&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-2"
          },
          "v": {
            "resource": "test-6",
            "component": "parent-2"
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
                  "connector": "test-corel-2",
                  "connector_name": "test-corel-2-name",
                  "component": "child-2",
                  "resource": "test-3"
                }
              },
              {
                "v": {
                  "connector": "test-corel-2",
                  "connector_name": "test-corel-2-name",
                  "component": "child-2",
                  "resource": "test-4"
                }
              },
              {
                "v": {
                  "connector": "test-corel-2",
                  "connector_name": "test-corel-2-name",
                  "component": "child-2",
                  "resource": "test-5"
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

  Scenario: given meta alarm rule and events shouldn't create meta alarm without children, after children events metaalarm should be based only on first parent
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-3",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-3"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-3"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-3",
        "corel_child":  "child-3"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-3",
      "connector_name": "test-corel-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-3",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector": "test-corel-3",
      "connector_name": "test-corel-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-3",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 1s
    When I send an event:
    """json
    {
      "connector": "test-corel-3",
      "connector_name": "test-corel-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-3",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I do GET /api/v4/alarms?search=test-corel-3&correlation=true&sort_by=v.resource&sort=asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-1"
          }
        },
        {
          "v": {
            "resource": "test-2"
          }
        },
        {
          "v": {
            "resource": "test-3"
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
    When I send an event:
    """json
    {
      "connector": "test-corel-3",
      "connector_name": "test-corel-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-3",
      "resource": "test-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-3",
      "connector_name": "test-corel-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-3",
      "resource": "test-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-3&correlation=true&multi_sort[]=v.meta,desc&multi_sort[]=v.resource,asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-3"
          },
          "v": {
            "resource": "test-1",
            "component": "parent-3"
          }
        },
        {
          "v": {
            "resource": "test-2"
          }
        },
        {
          "v": {
            "resource": "test-3"
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
                  "connector": "test-corel-3",
                  "connector_name": "test-corel-3-name",
                  "component": "child-3",
                  "resource": "test-4"
                }
              },
              {
                "v": {
                  "connector": "test-corel-3",
                  "connector_name": "test-corel-3-name",
                  "component": "child-3",
                  "resource": "test-5"
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

  Scenario: given meta alarm rule and events after time interval should shift child interval
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-4",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-4"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-4"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-4",
        "corel_child":  "child-4"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-4",
      "connector_name": "test-corel-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-4",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector": "test-corel-4",
      "connector_name": "test-corel-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-4",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector": "test-corel-4",
      "connector_name": "test-corel-4-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-4",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-4&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-4"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-4",
            "children": [
              "test-2/child-4"
            ]
          }
        },
        {
          "v": {
            "resource": "test-1"
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
                  "connector": "test-corel-4",
                  "connector_name": "test-corel-4-name",
                  "component": "child-4",
                  "resource": "test-2"
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
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events after time interval should shift parent interval
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-5",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-5"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-5"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-5",
        "corel_child":  "child-5"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-5",
      "connector_name": "test-corel-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-5",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector": "test-corel-5",
      "connector_name": "test-corel-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-5",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector": "test-corel-5",
      "connector_name": "test-corel-5-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-5",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-5&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-5"
          },
          "v": {
            "resource": "test-2",
            "component": "parent-5",
            "children": [
              "test-3/child-5"
            ]
          }
        },
        {
          "v": {
            "resource": "test-1"
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
                  "connector": "test-corel-5",
                  "connector_name": "test-corel-5-name",
                  "component": "child-5",
                  "resource": "test-3"
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
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events should create 2 meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-6",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-6"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-6"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 5,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-6",
        "corel_child":  "child-6"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-6",
      "connector_name": "test-corel-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-6",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-6",
      "connector_name": "test-corel-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "child-6",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I wait 5s
    When I send an event:
    """json
    {
      "connector": "test-corel-6",
      "connector_name": "test-corel-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-6",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-6",
      "connector_name": "test-corel-6-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "child-6",
      "resource": "test-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-6&correlation=true&sort_by=t&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-6"
          },
          "v": {
            "resource": "test-1",
            "component": "parent-6",
            "children": [
              "test-2/child-6"
            ]
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-6"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-6",
            "children": [
              "test-4/child-6"
            ]
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
                  "connector": "test-corel-6",
                  "connector_name": "test-corel-6-name",
                  "component": "child-6",
                  "resource": "test-2"
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
        }
      },
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-corel-6",
                  "connector_name": "test-corel-6-name",
                  "component": "child-6",
                  "resource": "test-4"
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
        }
      }
    ]
    """

  Scenario: given meta alarm rule and events with different corel_id, should create 2 meta alarm without mixing
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-7",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "child-7"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "parent-7"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "parent-7",
        "corel_child":  "child-7"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-7-1",
      "connector_name": "test-corel-7-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-7",
      "resource": "test-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-7-2",
      "connector_name": "test-corel-7-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-7",
      "resource": "test-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 1 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "connector": "test-corel-7-2",
      "connector_name": "test-corel-7-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "parent-7",
      "resource": "test-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-7-1",
      "connector_name": "test-corel-7-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-7",
      "resource": "test-4",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-7&correlation=true&sort_by=t&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-7"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-7",
            "connector": "test-corel-7-2",
            "children": [
              "test-2/child-7"
            ]
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-7"
          },
          "v": {
            "connector": "test-corel-7-1",
            "resource": "test-1",
            "component": "parent-7",
            "children": [
              "test-4/child-7"
            ]
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
    When I send an event:
    """json
    {
      "connector": "test-corel-7-2",
      "connector_name": "test-corel-7-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-7",
      "resource": "test-5",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-7-1",
      "connector_name": "test-corel-7-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "child-7",
      "resource": "test-6",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-7&correlation=true&sort_by=t&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-7"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-7",
            "connector": "test-corel-7-2"
          }
        },
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-corel-7"
          },
          "v": {
            "connector": "test-corel-7-1",
            "resource": "test-1",
            "component": "parent-7"
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
                  "connector": "test-corel-7-2",
                  "connector_name": "test-corel-7-2-name",
                  "component": "child-7",
                  "resource": "test-2"
                }
              },
              {
                "v": {
                  "connector": "test-corel-7-2",
                  "connector_name": "test-corel-7-2-name",
                  "component":  "child-7",
                  "resource": "test-5"
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
                  "connector": "test-corel-7-1",
                  "connector_name": "test-corel-7-1-name",
                  "component": "child-7",
                  "resource": "test-4"
                }
              },
              {
                "v": {
                  "connector": "test-corel-7-1",
                  "connector_name": "test-corel-7-1-name",
                  "component":  "child-7",
                  "resource": "test-6"
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

  Scenario: given deleted meta alarm rule should delete meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-corel-8",
      "type": "corel",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-corel-8-child"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-corel-8-parent"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "test-corel-8-parent",
        "corel_child":  "test-corel-8-child"
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-corel-1",
      "connector_name": "test-corel-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-corel-8-parent",
      "resource": "test-corel-8-resource-1",
      "state": 2
    }
    """
    When I wait the end of 1 events processing
    When I send an event:
    """json
    {
      "connector": "test-corel-1",
      "connector_name": "test-corel-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-corel-8-child",
      "resource": "test-corel-8-resource-2",
      "state": 2
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-corel-8&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-corel-8"
          },
          "v": {
            "resource": "test-corel-8-resource-1",
            "component": "test-corel-8-parent",
            "children": [
              "test-corel-8-resource-2/test-corel-8-child"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do DELETE /api/v4/cat/metaalarmrules/{{ .metaAlarmRuleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms/{{ .metaAlarmID }}
    Then the response code should be 200
    When I do GET /api/v4/alarms?search=test-corel-8&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-corel-8-resource-1/test-corel-8-parent"
          },
          "v": {
            "connector": "test-corel-1",
            "connector_name": "test-corel-1-name",
            "component": "test-corel-8-parent",
            "resource": "test-corel-8-resource-1",
            "children": [],
            "parents": []
          }
        },
        {
          "entity": {
            "_id": "test-corel-8-resource-2/test-corel-8-child"
          },
          "v": {
            "connector": "test-corel-1",
            "connector_name": "test-corel-1-name",
            "component": "test-corel-8-child",
            "resource": "test-corel-8-resource-2",
            "children": [],
            "parents": []
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
