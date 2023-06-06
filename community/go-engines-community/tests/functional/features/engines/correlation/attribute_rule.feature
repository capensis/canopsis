Feature: correlation feature - attribute rule

  @concurrent
  Scenario: given meta alarm rule and events should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-attribute-correlation-1",
      "type": "attribute",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-attribute-correlation-1"
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
      "connector": "test-attribute-1",
      "connector_name": "test-attribute-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-attribute-correlation-1",
      "resource": "test-attribute-correlation-resource-1-1",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-attribute-correlation-resource-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "meta_alarm_rule": {
            "name": "test-attribute-correlation-1"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                  "connector": "test-attribute-1",
                  "connector_name": "test-attribute-1-name",
                  "component": "test-attribute-correlation-1",
                  "resource": "test-attribute-correlation-resource-1-1"
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
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-attribute-1",
      "connector_name": "test-attribute-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-attribute-correlation-1",
      "resource": "test-attribute-correlation-resource-1-2",
      "state": 2,
      "output": "test",
      "long_output": "test",
      "author": "test-author"
    }
    """
    When I save request:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "children": {
          "page": 1,
          "sort_by": "v.resource",
          "sort": "asc"
        }
      }
    ]
    """
    When I do POST /api/v4/alarm-details until response code is 207 and body contains:
    """json
    [
      {
        "status": 200,
        "data": {
          "children": {
            "data": [
              {
                "v": {
                  "connector": "test-attribute-1",
                  "connector_name": "test-attribute-1-name",
                  "component": "test-attribute-correlation-1",
                  "resource": "test-attribute-correlation-resource-1-1"
                }
              },
              {
                "v": {
                  "connector": "test-attribute-1",
                  "connector_name": "test-attribute-1-name",
                  "component": "test-attribute-correlation-1",
                  "resource": "test-attribute-correlation-resource-1-2"
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
      "name": "test-attribute-correlation-2",
      "type": "attribute",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-attribute-correlation-2"
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
      "connector": "test-attribute-2",
      "connector_name": "test-attribute-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-attribute-correlation-2",
      "resource": "test-attribute-correlation-resource-2-1",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-attribute-correlation-resource-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-attribute-2",
      "connector_name": "test-attribute-2-name",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-attribute-correlation-2",
      "resource": "test-attribute-correlation-resource-2-2",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-attribute-correlation-resource-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-attribute-correlation-2"
          }
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-attribute-2",
      "connector_name": "test-attribute-2-name",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-attribute-correlation-2",
      "resource": "test-attribute-correlation-resource-2-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-attribute-2",
      "connector_name": "test-attribute-2-name",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-attribute-correlation-2",
      "resource": "test-attribute-correlation-resource-2-2"
    }
    """
    When I do DELETE /api/v4/cat/metaalarmrules/{{ .metaAlarmRuleID }}
    Then the response code should be 204
    When I do GET /api/v4/alarms/{{ .metaAlarmID }}
    Then the response code should be 404
    When I do GET /api/v4/entitybasics?_id={{ .metaAlarmEntityID }}
    Then the response code should be 404
    When I do GET /api/v4/alarms?search=test-attribute-correlation-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-attribute-2",
            "connector_name": "test-attribute-2-name",
            "component": "test-attribute-correlation-2",
            "resource": "test-attribute-correlation-resource-2-1",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-attribute-2",
            "connector_name": "test-attribute-2-name",
            "component": "test-attribute-correlation-2",
            "resource": "test-attribute-correlation-resource-2-2",
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

  @concurrent
  Scenario: given deleted resolved meta alarm rule should delete meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-attribute-correlation-3",
      "type": "attribute",
      "auto_resolve": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-attribute-correlation-3"
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
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-1",
        "state": 2
      },
      {
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-2",
        "state": 2
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-attribute-correlation-resource-3&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "cancel",
        "component":  "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-1"
      },
      {
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "cancel",
        "component": "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component":  "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-1"
      },
      {
        "connector": "test-attribute-3",
        "connector_name": "test-attribute-3-name",
        "source_type": "resource",
        "event_type": "resolve_cancel",
        "component": "test-attribute-correlation-3",
        "resource": "test-attribute-correlation-resource-3-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-attribute-correlation-resource-3&opened=false&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-attribute-correlation-3"
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
    When I do GET /api/v4/alarms?search=test-attribute-correlation-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-attribute-3",
            "connector_name": "test-attribute-3-name",
            "component": "test-attribute-correlation-3",
            "resource": "test-attribute-correlation-resource-3-1",
            "parents": []
          }
        },
        {
          "v": {
            "connector": "test-attribute-3",
            "connector_name": "test-attribute-3-name",
            "component": "test-attribute-correlation-3",
            "resource": "test-attribute-correlation-resource-3-2",
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

  @concurrent
  Scenario: given meta alarm and removed child should update meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-correlation-attribute-4",
      "type": "attribute",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-attribute-correlation-attribute-4"
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
        "component": "test-attribute-correlation-attribute-4",
        "connector": "test-connector-correlation-attribute-4",
        "connector_name": "test-connector-name-correlation-attribute-4",
        "resource": "test-resource-correlation-attribute-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "component": "test-attribute-correlation-attribute-4",
        "connector": "test-connector-correlation-attribute-4",
        "connector_name": "test-connector-name-correlation-attribute-4",
        "resource": "test-resource-correlation-attribute-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-attribute-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-attribute-4"
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
                  "component": "test-attribute-correlation-attribute-4",
                  "connector": "test-connector-correlation-attribute-4",
                  "connector_name": "test-connector-name-correlation-attribute-4",
                  "resource": "test-resource-correlation-attribute-4-1"
                }
              },
              {
                "v": {
                  "component": "test-attribute-correlation-attribute-4",
                  "connector": "test-connector-correlation-attribute-4",
                  "connector_name": "test-connector-name-correlation-attribute-4",
                  "resource": "test-resource-correlation-attribute-4-2"
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
      "comment": "test-metaalarmrule-correlation-attribute-4-remove-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-attribute-4&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-attribute-4"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-attribute-4-1/test-attribute-correlation-attribute-4"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-attribute-correlation-attribute-4",
            "connector": "test-connector-correlation-attribute-4",
            "connector_name": "test-connector-name-correlation-attribute-4",
            "resource": "test-resource-correlation-attribute-4-2"
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
                  "component": "test-attribute-correlation-attribute-4",
                  "connector": "test-connector-correlation-attribute-4",
                  "connector_name": "test-connector-name-correlation-attribute-4",
                  "resource": "test-resource-correlation-attribute-4-1"
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
      "name": "test-metaalarmrule-correlation-attribute-5",
      "type": "attribute",
      "output_template": "{{ `{{ .Count }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-attribute-correlation-attribute-5"
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
        "component": "test-attribute-correlation-attribute-5",
        "connector": "test-connector-correlation-attribute-5",
        "connector_name": "test-connector-name-correlation-attribute-5",
        "resource": "test-resource-correlation-attribute-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 3,
        "component": "test-attribute-correlation-attribute-5",
        "connector": "test-connector-correlation-attribute-5",
        "connector_name": "test-connector-name-correlation-attribute-5",
        "resource": "test-resource-correlation-attribute-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-correlation-attribute-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 2,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-attribute-5"
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
      "comment": "test-metaalarmrule-correlation-attribute-5-remove-comment",
      "alarms": ["{{ .childAlarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-correlation-attribute-5&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-attribute-5"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-attribute-5-1/test-attribute-correlation-attribute-5"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-attribute-correlation-attribute-5",
            "connector": "test-connector-correlation-attribute-5",
            "connector_name": "test-connector-name-correlation-attribute-5",
            "resource": "test-resource-correlation-attribute-5-2"
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
      "component": "test-attribute-correlation-attribute-5",
      "connector": "test-connector-correlation-attribute-5",
      "connector_name": "test-connector-name-correlation-attribute-5",
      "resource": "test-resource-correlation-attribute-5-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-correlation-attribute-5&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "children": 1,
          "meta_alarm_rule": {
            "name": "test-metaalarmrule-correlation-attribute-5"
          },
          "v": {
            "output": "1",
            "children": [
              "test-resource-correlation-attribute-5-1/test-attribute-correlation-attribute-5"
            ],
            "state": {
              "val": 2
            }
          }
        },
        {
          "v": {
            "component": "test-attribute-correlation-attribute-5",
            "connector": "test-connector-correlation-attribute-5",
            "connector_name": "test-connector-name-correlation-attribute-5",
            "resource": "test-resource-correlation-attribute-5-2"
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
