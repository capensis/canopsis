Feature: correlation feature - relation rule

  @concurrent
  Scenario: given meta alarm rule and events should trigger metaalarm by component event
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-relation-correlation-1",
      "type": "relation",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-relation-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-1",
      "connector_name": "test-relation-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-1",
      "resource": "test-relation-correlation-resource-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-1",
      "connector_name": "test-relation-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-1",
      "resource": "test-relation-correlation-resource-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-1",
      "connector_name": "test-relation-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-1",
      "resource": "test-relation-correlation-resource-1-3",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-1",
      "connector_name": "test-relation-1-name",
      "source_type": "component",
      "event_type": "check",
      "component": "test-relation-correlation-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-relation-correlation-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "component",
            "name": "test-relation-correlation-1"
          },
          "v": {
            "component": "test-relation-correlation-1"
          },
          "is_meta_alarm": true,
          "children": 3,
          "meta_alarm_rule": {
            "name": "test-relation-correlation-1"
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
                  "connector": "test-relation-1",
                  "connector_name": "test-relation-1-name",
                  "component": "test-relation-correlation-1",
                  "resource": "test-relation-correlation-resource-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-relation-1",
                  "connector_name": "test-relation-1-name",
                  "component": "test-relation-correlation-1",
                  "resource": "test-relation-correlation-resource-1-2"
                }
              },
              {
                "v": {
                  "connector": "test-relation-1",
                  "connector_name": "test-relation-1-name",
                  "component": "test-relation-correlation-1",
                  "resource": "test-relation-correlation-resource-1-3"
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
  Scenario: given meta alarm rule and events should trigger metaalarm by first resource event
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-relation-correlation-2",
      "type": "relation",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-relation-correlation-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-2",
      "connector_name": "test-relation-2-name",
      "source_type": "component",
      "event_type": "check",
      "component": "test-relation-correlation-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-2",
      "connector_name": "test-relation-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-2",
      "resource": "test-relation-correlation-resource-2-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-2",
      "connector_name": "test-relation-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-2",
      "resource": "test-relation-correlation-resource-2-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-relation-correlation-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "component",
            "name": "test-relation-correlation-2"
          },
          "v": {
            "component": "test-relation-correlation-2"
          },
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-relation-correlation-2"
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
                  "connector": "test-relation-2",
                  "connector_name": "test-relation-2-name",
                  "component": "test-relation-correlation-2",
                  "resource": "test-relation-correlation-resource-2-1"
                }
              },
              {
                "v": {
                  "connector": "test-relation-2",
                  "connector_name": "test-relation-2-name",
                  "component": "test-relation-correlation-2",
                  "resource": "test-relation-correlation-resource-2-2"
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
  Scenario: given deleted meta alarm rule should delete meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-relation-correlation-3",
      "type": "relation",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-relation-correlation-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-3",
      "connector_name": "test-relation-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-3",
      "resource": "test-relation-correlation-resource-3-1",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-3",
      "connector_name": "test-relation-3-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-relation-correlation-3",
      "resource": "test-relation-correlation-resource-3-2",
      "state": 2
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-relation-3",
      "connector_name": "test-relation-3-name",
      "source_type": "component",
      "event_type": "check",
      "component": "test-relation-correlation-3",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-relation-correlation-3&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "entity": {
            "type": "component",
            "name": "test-relation-correlation-3"
          },
          "v": {
            "component": "test-relation-correlation-3"
          },
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-relation-correlation-3"
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
    When I do GET /api/v4/alarms?search=test-relation-correlation-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-relation-correlation-3"
          },
          "v": {
            "connector": "test-relation-3",
            "connector_name": "test-relation-3-name",
            "component": "test-relation-correlation-3",
            "children": [],
            "parents": []
          }
        },
        {
          "entity": {
            "_id": "test-relation-correlation-resource-3-1/test-relation-correlation-3"
          },
          "v": {
            "connector": "test-relation-3",
            "connector_name": "test-relation-3-name",
            "component": "test-relation-correlation-3",
            "resource": "test-relation-correlation-resource-3-1",
            "children": [],
            "parents": []
          }
        },
        {
          "entity": {
            "_id": "test-relation-correlation-resource-3-2/test-relation-correlation-3"
          },
          "v": {
            "connector": "test-relation-3",
            "connector_name": "test-relation-3-name",
            "component": "test-relation-correlation-3",
            "resource": "test-relation-correlation-resource-3-2",
            "children": [],
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
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-relation-4",
      "type": "relation",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-relation-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "component": "test-component-correlation-relation-4",
        "connector": "test-connector-correlation-relation-4",
        "connector_name": "test-connector-name-correlation-relation-4",
        "resource": "test-resource-correlation-relation-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "component": "test-component-correlation-relation-4",
        "connector": "test-connector-correlation-relation-4",
        "connector_name": "test-connector-name-correlation-relation-4",
        "resource": "test-resource-correlation-relation-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "component": "test-component-correlation-relation-4",
      "connector": "test-connector-correlation-relation-4",
      "connector_name": "test-connector-name-correlation-relation-4",
      "source_type": "component"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-relation-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-relation-4"
          },
          "v": {
            "output": "2",
            "state": {
              "val": 3
            },
            "component": "test-component-correlation-relation-4",
            "connector": "test-connector-correlation-relation-4",
            "connector_name": "test-connector-name-correlation-relation-4"
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
                  "component": "test-component-correlation-relation-4",
                  "connector": "test-connector-correlation-relation-4",
                  "connector_name": "test-connector-name-correlation-relation-4",
                  "resource": "test-resource-correlation-relation-4-1"
                }
              },
              {
                "v": {
                  "component": "test-component-correlation-relation-4",
                  "connector": "test-connector-correlation-relation-4",
                  "connector_name": "test-connector-name-correlation-relation-4",
                  "resource": "test-resource-correlation-relation-4-2"
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
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-relation-4-remove-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-relation-4&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-relation-4"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-relation-4-1/test-component-correlation-relation-4"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-component-correlation-relation-4",
            "connector": "test-connector-correlation-relation-4",
            "connector_name": "test-connector-name-correlation-relation-4",
            "resource": "test-resource-correlation-relation-4-2"
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
                  "component": "test-component-correlation-relation-4",
                  "connector": "test-connector-correlation-relation-4",
                  "connector_name": "test-connector-name-correlation-relation-4",
                  "resource": "test-resource-correlation-relation-4-1"
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

  @concurrent
  Scenario: given meta alarm and removed child should not add child to meta alarm again
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-relation-5",
      "type": "relation",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-correlation-relation-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "component": "test-component-correlation-relation-5",
        "connector": "test-connector-correlation-relation-5",
        "connector_name": "test-connector-name-correlation-relation-5",
        "resource": "test-resource-correlation-relation-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "component": "test-component-correlation-relation-5",
        "connector": "test-connector-correlation-relation-5",
        "connector_name": "test-connector-name-correlation-relation-5",
        "resource": "test-resource-correlation-relation-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 3,
      "component": "test-component-correlation-relation-5",
      "connector": "test-connector-correlation-relation-5",
      "connector_name": "test-connector-name-correlation-relation-5",
      "source_type": "component"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-relation-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-relation-5"
          },
          "v": {
            "output": "2"
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
    When I save response childAlarmId2={{ (index (index .lastResponse 0).data.children.data 1)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "comment": "test-metaalarmrule-correlation-relation-5-remove-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-relation-5&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-relation-5"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-relation-5-1/test-component-correlation-relation-5"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-component-correlation-relation-5",
            "connector": "test-connector-correlation-relation-5",
            "connector_name": "test-connector-name-correlation-relation-5",
            "resource": "test-resource-correlation-relation-5-2"
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
      "component": "test-component-correlation-relation-5",
      "connector": "test-connector-correlation-relation-5",
      "connector_name": "test-connector-name-correlation-relation-5",
      "resource": "test-resource-correlation-relation-5-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-relation-5&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-relation-5"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-relation-5-1/test-component-correlation-relation-5"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-component-correlation-relation-5",
            "connector": "test-connector-correlation-relation-5",
            "connector_name": "test-connector-name-correlation-relation-5",
            "resource": "test-resource-correlation-relation-5-2"
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
