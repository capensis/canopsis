Feature: update alarm on idle rule
  I need to be able to update alarm on idle rule

  @concurrent
  Scenario: given idle rule and no events for alarm and entity should apply most priority entity rule
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-second-1-1-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 53,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-second-1"
            }
          }
        ]
      ],
      "operation": {
        "type": "ack",
        "parameters": {
          "output": "test-idlerule-axe-idlerule-second-1-1-output"
        }
      }
    }
    """
    Then the response code should be 201
    Then I save response rule1ID={{ .lastResponse._id }}
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-second-1-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 52,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-axe-idlerule-second-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response rule2ID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-second-1",
      "connector_name": "test-connector-name-axe-idlerule-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-second-1",
      "resource": "test-resource-axe-idlerule-second-1",
      "state": 2,
      "output": "test-output-axe-idlerule-second-1"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-second-1",
      "connector_name": "test-connector-name-axe-idlerule-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-second-1",
      "resource": "test-resource-axe-idlerule-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-1",
            "connector": "test-connector-axe-idlerule-second-1",
            "connector_name": "test-connector-name-axe-idlerule-second-1",
            "resource": "test-resource-axe-idlerule-second-1",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I wait 5s
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-1
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-second-1",
      "connector_name": "test-connector-name-axe-idlerule-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-second-1",
      "resource": "test-resource-axe-idlerule-second-1",
      "state": 2,
      "output": "test-output-axe-idlerule-second-1"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "noevents",
      "connector": "test-connector-axe-idlerule-second-1",
      "connector_name": "test-connector-name-axe-idlerule-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-second-1",
      "resource": "test-resource-axe-idlerule-second-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-1
    Then the response code should be 200
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              },
              {
                "_t": "statedec",
                "val": 2
              },
              {
                "_t": "statusdec",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 8
            }
          }
        }
      }
    ]
    """
    When I do DELETE /api/v4/idle-rules/{{ .rule1ID }}
    Then the response code should be 204
    When I do DELETE /api/v4/idle-rules/{{ .rule2ID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule and no events for component which is created by resource event should create alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-second-2-name",
      "type": "entity",
      "enabled": true,
      "priority": 54,
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-idlerule-second-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response ruleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-idlerule-second-2",
      "connector_name": "test-connector-name-axe-idlerule-second-2",
      "source_type": "resource",
      "component":  "test-component-axe-idlerule-second-2",
      "resource":  "test-resource-axe-idlerule-second-2",
      "state": 0,
      "output": "test-output-axe-idlerule-second-2"
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "activate",
      "connector": "test-connector-axe-idlerule-second-2",
      "connector_name": "test-connector-name-axe-idlerule-second-2",
      "source_type": "component",
      "component":  "test-component-axe-idlerule-second-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-axe-idlerule-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-2",
            "connector": "test-connector-axe-idlerule-second-2",
            "connector_name": "test-connector-name-axe-idlerule-second-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 5
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 5
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
    When I do DELETE /api/v4/idle-rules/{{ .ruleID }}
    Then the response code should be 204

  @concurrent
  Scenario: given idle rule with ok changestate operation should update next alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-second-3-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
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
              "value": "test-resource-axe-idlerule-second-3"
            }
          }
        ]
      ],
      "operation": {
        "type": "changestate",
        "parameters": {
          "state": 0
        }
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-idlerule-second-3",
      "connector": "test-connector-axe-idlerule-second-3",
      "connector_name": "test-connector-name-axe-idlerule-second-3",
      "component":  "test-component-axe-idlerule-second-3",
      "resource": "test-resource-axe-idlerule-second-3",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "changestate",
      "connector": "test-connector-axe-idlerule-second-3",
      "connector_name": "test-connector-name-axe-idlerule-second-3",
      "component":  "test-component-axe-idlerule-second-3",
      "resource": "test-resource-axe-idlerule-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-3",
            "connector": "test-connector-axe-idlerule-second-3",
            "connector_name": "test-connector-name-axe-idlerule-second-3",
            "resource": "test-resource-axe-idlerule-second-3",
            "state": {
              "_t": "changestate",
              "val": 0
            },
            "status": {
              "val": 0
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_close",
      "output": "test-output-axe-idlerule-second-3",
      "connector": "test-connector-axe-idlerule-second-3",
      "connector_name": "test-connector-name-axe-idlerule-second-3",
      "component":  "test-component-axe-idlerule-second-3",
      "resource": "test-resource-axe-idlerule-second-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-idlerule-second-3",
      "connector": "test-connector-axe-idlerule-second-3",
      "connector_name": "test-connector-name-axe-idlerule-second-3",
      "component":  "test-component-axe-idlerule-second-3",
      "resource": "test-resource-axe-idlerule-second-3",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "changestate",
      "connector": "test-connector-axe-idlerule-second-3",
      "connector_name": "test-connector-name-axe-idlerule-second-3",
      "component":  "test-component-axe-idlerule-second-3",
      "resource": "test-resource-axe-idlerule-second-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-3",
            "connector": "test-connector-axe-idlerule-second-3",
            "connector_name": "test-connector-name-axe-idlerule-second-3",
            "resource": "test-resource-axe-idlerule-second-3",
            "state": {
              "_t": "changestate",
              "val": 0
            },
            "status": {
              "val": 0
            }
          }
        },
        {
          "v": {
            "component": "test-component-axe-idlerule-second-3",
            "connector": "test-connector-axe-idlerule-second-3",
            "connector_name": "test-connector-name-axe-idlerule-second-3",
            "resource": "test-resource-axe-idlerule-second-3",
            "state": {
              "_t": "changestate",
              "val": 0
            },
            "status": {
              "val": 0
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

  @concurrent
  Scenario: given idle rule with cancel operation should update next alarm
    Given I am admin
    When I do POST /api/v4/idle-rules:
    """json
    {
      "name": "test-idlerule-axe-idlerule-second-4-name",
      "type": "alarm",
      "alarm_condition": "last_event",
      "enabled": true,
      "priority": 40,
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
              "value": "test-resource-axe-idlerule-second-4"
            }
          }
        ]
      ],
      "operation": {
        "type": "cancel"
      }
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-idlerule-second-4",
      "connector": "test-connector-axe-idlerule-second-4",
      "connector_name": "test-connector-name-axe-idlerule-second-4",
      "component":  "test-component-axe-idlerule-second-4",
      "resource": "test-resource-axe-idlerule-second-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-idlerule-second-4",
      "connector_name": "test-connector-name-axe-idlerule-second-4",
      "component":  "test-component-axe-idlerule-second-4",
      "resource": "test-resource-axe-idlerule-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-4",
            "connector": "test-connector-axe-idlerule-second-4",
            "connector_name": "test-connector-name-axe-idlerule-second-4",
            "resource": "test-resource-axe-idlerule-second-4",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "output": "test-output-axe-idlerule-second-4",
      "connector": "test-connector-axe-idlerule-second-4",
      "connector_name": "test-connector-name-axe-idlerule-second-4",
      "component":  "test-component-axe-idlerule-second-4",
      "resource": "test-resource-axe-idlerule-second-4",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-idlerule-second-4",
      "connector": "test-connector-axe-idlerule-second-4",
      "connector_name": "test-connector-name-axe-idlerule-second-4",
      "component":  "test-component-axe-idlerule-second-4",
      "resource": "test-resource-axe-idlerule-second-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-idlerule-second-4",
      "connector_name": "test-connector-name-axe-idlerule-second-4",
      "component":  "test-component-axe-idlerule-second-4",
      "resource": "test-resource-axe-idlerule-second-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-idlerule-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-idlerule-second-4",
            "connector": "test-connector-axe-idlerule-second-4",
            "connector_name": "test-connector-name-axe-idlerule-second-4",
            "resource": "test-resource-axe-idlerule-second-4",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
            }
          }
        },
        {
          "v": {
            "component": "test-component-axe-idlerule-second-4",
            "connector": "test-connector-axe-idlerule-second-4",
            "connector_name": "test-connector-name-axe-idlerule-second-4",
            "resource": "test-resource-axe-idlerule-second-4",
            "canceled": {
              "_t": "cancel"
            },
            "state": {
              "val": 2
            },
            "status": {
              "val": 4
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
