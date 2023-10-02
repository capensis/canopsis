Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  @concurrent
  Scenario: given meta alarm rule and event should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-1
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-1",
      "connector_name": "test-connector-name-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-1",
      "resource": "test-resource-axe-correlation-1-1",
      "state": 2,
      "output": "test-output-axe-correlation-1"
    }
    """
    When I save response createTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-1; Count: 1; Children: test-component-axe-correlation-1",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            }
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-1",
      "connector_name": "test-connector-name-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-1",
      "resource": "test-resource-axe-correlation-1-2",
      "state": 3,
      "output": "test-output-axe-correlation-1"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-1; Count: 2; Children: test-component-axe-correlation-1",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "val": 1
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I save response metaalarmCreationDate={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    Then the difference between metaalarmLastEventDate createTimestamp is in range -2,2
    Then the difference between metaalarmCreationDate createTimestamp is in range -2,2
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
                  "component": "test-component-axe-correlation-1",
                  "connector": "test-connector-axe-correlation-1",
                  "connector_name": "test-connector-name-axe-correlation-1",
                  "resource": "test-resource-axe-correlation-1-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-1",
                  "connector": "test-connector-axe-correlation-1",
                  "connector_name": "test-connector-name-axe-correlation-1",
                  "resource": "test-resource-axe-correlation-1-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID }}",
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "engine.correlation",
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
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-1-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "children": [],
            "component": "test-component-axe-correlation-1",
            "connector": "test-connector-axe-correlation-1",
            "connector_name": "test-connector-name-axe-correlation-1",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-1-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
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
    When I save response alarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
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
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "val": 0
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
    When I save response metaAlarmAttachStepDate={{ (index (index .lastResponse 0).data.steps.data 2).t }}
    Then the difference between alarmLastEventDate createTimestamp is in range -2,2
    Then the difference between metaAlarmAttachStepDate createTimestamp is in range -2,2

  @concurrent
  Scenario: given meta alarm and ack event should ack children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-2
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-2",
      "connector_name": "test-connector-name-axe-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-2",
      "resource": "test-resource-axe-correlation-2",
      "state": 2,
      "output": "test-output-axe-correlation-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-2&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "ack",
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-2"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-2",
        "connector_name": "test-connector-name-axe-correlation-2",
        "component": "test-component-axe-correlation-2",
        "resource": "test-resource-axe-correlation-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "root John Doe admin@canopsis.net",
              "m": "test-output-axe-correlation-2"
            },
            "children": [],
            "component": "test-component-axe-correlation-2",
            "connector": "test-connector-axe-correlation-2",
            "connector_name": "test-connector-name-axe-correlation-2",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
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
                "_t": "metaalarmattach"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "m": "test-output-axe-correlation-2"
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
  Scenario: given meta alarm child and change state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-3
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 1,
      "output": "test-output-axe-correlation-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-3&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "changestate",
      "component": "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 2,
      "output": "test-output-axe-correlation-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-3&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-3; Count: 1; Children: test-component-axe-correlation-3",
            "children": [
              "test-resource-axe-correlation-3/test-component-axe-correlation-3"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-3&correlation=true
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
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
  Scenario: given meta alarm and change state event should update children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-4
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-4",
      "connector_name": "test-connector-name-axe-correlation-4",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-4",
      "resource": "test-resource-axe-correlation-4",
      "state": 1,
      "output": "test-output-axe-correlation-4"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "changestate",
      "state": 2,
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-4"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "changestate",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "changestate",
        "connector": "test-connector-axe-correlation-4",
        "connector_name": "test-connector-name-axe-correlation-4",
        "component": "test-component-axe-correlation-4",
        "resource": "test-resource-axe-correlation-4",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "children": [],
            "component": "test-component-axe-correlation-4",
            "connector": "test-connector-axe-correlation-4",
            "connector_name": "test-connector-name-axe-correlation-4",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-4",
            "state": {
              "_t": "changestate",
              "a": "root John Doe admin@canopsis.net",
              "m": "test-output-axe-correlation-4",
              "val": 2
            },
            "status": {
              "val": 1
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "changestate",
                "a": "root John Doe admin@canopsis.net",
                "m": "test-output-axe-correlation-4",
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
      }
    ]
    """

  @concurrent
  Scenario: given meta alarm child and resolve event should resolve parent
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-5
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "state": 1,
      "output": "test-output-axe-correlation-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-5&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1

  @concurrent
  Scenario: given meta alarm child and check with inc state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-6
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 1,
      "output": "test-output-axe-correlation-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-6&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 2
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-6&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-6; Count: 1; Children: test-component-axe-correlation-6",
            "children": [
              "test-resource-axe-correlation-6/test-component-axe-correlation-6"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-6&correlation=true
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
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "stateinc",
                "val": 2
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
  Scenario: given meta alarm child and check with dec state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-7
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 2,
      "output": "test-output-axe-correlation-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-7&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 1,
      "output": "test-output-axe-correlation-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-7&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-7; Count: 1; Children: test-component-axe-correlation-7",
            "children": [
              "test-resource-axe-correlation-7/test-component-axe-correlation-7"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "val": 1
            },
            "status": {
              "val": 1
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-7&correlation=true
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
                "_t": "statedec",
                "val": 1
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
  Scenario: given meta alarm rule and event should update child alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-8
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-1",
      "state": 2,
      "output": "test-output-axe-correlation-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-8&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "meta": "{{ .metaAlarmRuleID }}"
          }
        }
      ]
    }
    """
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "snooze",
      "duration": 6000,
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-8"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "snooze",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "test-connector-axe-correlation-8",
        "connector_name": "test-connector-name-axe-correlation-8",
        "component": "test-component-axe-correlation-8",
        "resource": "test-resource-axe-correlation-8-1",
        "source_type": "resource"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-2",
      "state": 2,
      "output": "test-output-axe-correlation-8"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "activate",
        "connector": "test-connector-axe-correlation-8",
        "connector_name": "test-connector-name-axe-correlation-8",
        "component": "test-component-axe-correlation-8",
        "resource": "test-resource-axe-correlation-8-2",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "test-connector-axe-correlation-8",
        "connector_name": "test-connector-name-axe-correlation-8",
        "component": "test-component-axe-correlation-8",
        "resource": "test-resource-axe-correlation-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-8-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-8",
            "connector": "test-connector-axe-correlation-8",
            "connector_name": "test-connector-name-axe-correlation-8",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-8-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
              },
              {
                "_t": "snooze"
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
  Scenario: given meta alarm rule and 2 simultaneous events should create meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-9
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-9",
        "connector_name": "test-connector-name-axe-correlation-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-9",
        "resource": "test-resource-axe-correlation-9-1",
        "state": 2,
        "output": "test-output-axe-correlation-9"
      },
      {
        "connector": "test-connector-axe-correlation-9",
        "connector_name": "test-connector-name-axe-correlation-9",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-9",
        "resource": "test-resource-axe-correlation-9-2",
        "state": 3,
        "output": "test-output-axe-correlation-9"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-9&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-9; Count: 2; Children: test-component-axe-correlation-9",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "val": 1
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
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
                  "component": "test-component-axe-correlation-9",
                  "connector": "test-connector-axe-correlation-9",
                  "connector_name": "test-connector-name-axe-correlation-9",
                  "resource": "test-resource-axe-correlation-9-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-9",
                  "connector": "test-connector-axe-correlation-9",
                  "connector_name": "test-connector-name-axe-correlation-9",
                  "resource": "test-resource-axe-correlation-9-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 1)._id }}",
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
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
                "val": 3
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "metaalarmattach",
                "a": "engine.correlation"
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
