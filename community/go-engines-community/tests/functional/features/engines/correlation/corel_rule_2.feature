Feature: correlation feature - corel rule

  @concurrent
  Scenario: given meta alarm and removed child should not add child to meta alarm again but add to new meta alarm on change
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-corel-second-1",
      "type": "corel",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-corel-second-1-child"
            }
          }
        ],
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-corel-second-1-parent"
            }
          }
        ]
      ],
      "config": {
        "time_interval": {
          "value": 3,
          "unit": "s"
        },
        "threshold_count": 2,
        "corel_id": "{{ `{{ .Alarm.Value.Connector }}` }}",
        "corel_status": "{{ `{{ .Entity.Component }}` }}",
        "corel_parent": "test-component-correlation-corel-second-1-parent",
        "corel_child":  "test-component-correlation-corel-second-1-child"
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
      "state": 3,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-parent",
      "resource": "test-resource-correlation-corel-second-1-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-corel-second-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
          },
          "v": {
            "output": "2",
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
      "comment": "test-metaalarmrule-correlation-corel-second-1-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-corel-second-1&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
          },
          "v": {
            "output": "1",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-corel-second-1-3"
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
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-3",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-corel-second-1&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
          },
          "v": {
            "output": "1",
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-correlation-corel-second-1-3"
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
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-corel-second-1&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
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
            "resource": "test-resource-correlation-corel-second-1-3"
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
    When I wait 3s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-parent",
      "resource": "test-resource-correlation-corel-second-1-5",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "connector": "test-connector-correlation-corel-second-1",
      "connector_name": "test-connector-name-correlation-corel-second-1",
      "component":  "test-component-correlation-corel-second-1-child",
      "resource": "test-resource-correlation-corel-second-1-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-corel-second-1&correlation=true&sort_by=t&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 3
            },
            "component":  "test-component-correlation-corel-second-1-parent",
            "resource": "test-resource-correlation-corel-second-1-5"
          }
        },
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-corel-second-1"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 2
            },
            "component":  "test-component-correlation-corel-second-1-parent",
            "resource": "test-resource-correlation-corel-second-1-1"
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
                  "resource": "test-resource-correlation-corel-second-1-3"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-corel-second-1-6"
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
                  "resource": "test-resource-correlation-corel-second-1-2"
                }
              },
              {
                "v": {
                  "resource": "test-resource-correlation-corel-second-1-4"
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
