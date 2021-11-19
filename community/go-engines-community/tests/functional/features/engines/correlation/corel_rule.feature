Feature: correlation feature - corel rule

  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-1",
      "type": "corel",
      "config": {
        "entity_patterns": [
          {
            "component": "child"
          },
          {
            "component": "parent"
          }
        ],
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
    """
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
    """
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
                      "name": "test-corel-1"
                    }
                  ],
                  "total": 1
                },
                "v": {
                  "component": "child",
                  "resource": "test-2"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
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

  Scenario: given meta alarm rule and events shouldn't create meta alarm without a parent, after parent event metaalarm should contain all children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-2",
      "type": "corel",
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "entity_patterns": [
          {
            "component": "child-2"
          },
          {
            "component": "parent-2"
          }
        ],
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
    """
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
    """
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
    """
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
                      "name": "test-corel-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-corel-2"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-corel-2"
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

  Scenario: given meta alarm rule and events shouldn't create meta alarm without children, after children events metaalarm should be based only on first parent
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-3",
      "type": "corel",
      "config": {
        "time_interval": {
          "value": 20,
          "unit": "s"
        },
        "threshold_count": 2,
        "entity_patterns": [
          {
            "component": "child-3"
          },
          {
            "component": "parent-3"
          }
        ],
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
    """
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
    """
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
    """
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
    When I wait 1s
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
    """
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
                        "name": "test-corel-3"
                    }
                  ],
                  "total": 1
                }
              },
              {
                "causes": {
                  "rules": [
                    {
                        "name": "test-corel-3"
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
            "name": "test-corel-3"
          },
          "v": {
            "resource": "test-1",
            "component": "parent-3"
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

  Scenario: given meta alarm rule and events after time interval should shift child interval
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-4",
      "type": "corel",
      "config": {
        "entity_patterns": [
          {
            "component": "child-4"
          },
          {
            "component": "parent-4"
          }
        ],
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
    """
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
    """
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
    """
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
                  "component": "child-4",
                  "resource": "test-2"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-corel-4"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-4",
            "children": [
              "test-2/child-4"
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

  Scenario: given meta alarm rule and events after time interval should shift parent interval
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-5",
      "type": "corel",
      "config": {
        "entity_patterns": [
          {
            "component": "child-5"
          },
          {
            "component": "parent-5"
          }
        ],
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
    """
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
    """
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
    """
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
                  "component": "child-5",
                  "resource": "test-3"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-corel-5"
          },
          "v": {
            "resource": "test-2",
            "component": "parent-5",
            "children": [
              "test-3/child-5"
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

  Scenario: given meta alarm rule and events should create 2 meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-6",
      "type": "corel",
      "config": {
        "entity_patterns": [
          {
            "component": "child-6"
          },
          {
            "component": "parent-6"
          }
        ],
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
    """
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
    """
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
    """
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
    """
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=asc
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
                      "name": "test-corel-6"
                    }
                  ],
                  "total": 1
                },
                "v": {
                  "component": "child-6",
                  "resource": "test-2"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
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
          "consequences": {
            "data": [
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-corel-6"
                    }
                  ],
                  "total": 1
                },
                "v": {
                  "component": "child-6",
                  "resource": "test-4"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
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

  Scenario: given meta alarm rule and events with different corel_id, should create 2 meta alarm without mixing
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-corel-7",
      "type": "corel",
      "config": {
        "entity_patterns": [
          {
            "component": "child-7"
          },
          {
            "component": "parent-7"
          }
        ],
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
    """
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
    """
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
    """
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
    """
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
                      "name": "test-corel-7"
                    }
                  ],
                  "total": 1
                },
                "v": {
                  "component": "child-7",
                  "resource": "test-2"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
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
          "consequences": {
            "data": [
              {
                "causes": {
                  "rules": [
                    {
                      "name": "test-corel-7"
                    }
                  ],
                  "total": 1
                },
                "v": {
                  "component": "child-7",
                  "resource": "test-4"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
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
    """
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
    """
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=t&sort_dir=desc
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "name": "test-corel-7"
          },
          "v": {
            "resource": "test-3",
            "component": "parent-7",
            "connector": "test-corel-7-2"
          }
        },
        {
          "consequences": {
            "total": 2
          },
          "metaalarm": true,
          "rule": {
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
