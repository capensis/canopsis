Feature: update meta alarm on idle rule
  I need to be able to update meta alarm on idle rule

  @concurrent
  Scenario: given meta alarm and entity idle rule should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-idlerule-correlation-1
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-idlerule-correlation-1",
      "connector_name": "test-connector-name-axe-idlerule-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-idlerule-correlation-1",
      "resource": "test-resource-axe-idlerule-correlation-1",
      "state": 2,
      "output": "test-output-axe-idlerule-correlation-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-correlation-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-correlation-1-name",
      "type": "entity",
      "enabled": true,
      "priority": 10,
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-correlation-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-correlation-1",
      "connector_name": "test-connector-name-axe-idlerule-correlation-1",
      "component":  "test-component-axe-idlerule-correlation-1",
      "resource": "test-resource-axe-idlerule-correlation-1",
      "source_type": "resource"
    }
    """
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
                  "component": "test-component-axe-idlerule-correlation-1",
                  "connector": "test-connector-axe-idlerule-correlation-1",
                  "connector_name": "test-connector-name-axe-idlerule-correlation-1",
                  "resource": "test-resource-axe-idlerule-correlation-1"
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
    When I save response childAlarmID={{ (index (index .lastResponse 0).data.children.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .childAlarmID }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 3
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "stateinc",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "val": 5
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 5
            }
          }
        }
      }
    ]
    """
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-idlerule-correlation-1",
      "connector_name": "test-connector-name-axe-idlerule-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-idlerule-correlation-1",
      "resource": "test-resource-axe-idlerule-correlation-1",
      "state": 2,
      "output": "test-output-axe-idlerule-correlation-1"
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .childAlarmID }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 3
              },
              {
                "_t": "statedec",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 2
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "stateinc",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 2
              },
              {
                "_t": "statusdec",
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 7
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm and alarm idle rule should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-idlerule-correlation-2
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-idlerule-correlation-2",
      "connector_name": "test-connector-name-axe-idlerule-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-idlerule-correlation-2",
      "resource": "test-resource-axe-idlerule-correlation-2",
      "state": 2,
      "output": "test-output-axe-idlerule-correlation-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-correlation-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-correlation-2-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 10,
      "duration": {
        "value": 2,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-correlation-2"
            }
          }
        ]
      ],
      "operation": {
        "type": "changestate",
        "parameters": {
          "state": 3
        }
      }
    }
    """
    Then the response code should be 201
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "changestate",
      "connector": "test-connector-axe-idlerule-correlation-2",
      "connector_name": "test-connector-name-axe-idlerule-correlation-2",
      "component":  "test-component-axe-idlerule-correlation-2",
      "resource": "test-resource-axe-idlerule-correlation-2",
      "source_type": "resource"
    }
    """
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
                  "component": "test-component-axe-idlerule-correlation-2",
                  "connector": "test-connector-axe-idlerule-correlation-2",
                  "connector_name": "test-connector-name-axe-idlerule-correlation-2",
                  "resource": "test-resource-axe-idlerule-correlation-2"
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
    When I save response childAlarmID={{ (index (index .lastResponse 0).data.children.data 0)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .childAlarmID }}",
        "steps": {
          "page": 1
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system",
                "val": 3
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
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "user_id": "",
                "initiator": "system"
              },
              {
                "_t": "changestate",
                "a": "system",
                "user_id": "",
                "initiator": "system",
                "val": 3
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
