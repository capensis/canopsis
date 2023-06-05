Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  Scenario: given meta alarm rule and event should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-1
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-1",
      "connector_name": "test-connector-name-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-1",
      "resource": "test-resource-axe-correlation-1-1",
      "state": 2,
      "output": "test-output-axe-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-1",
      "connector_name": "test-connector-name-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-1",
      "resource": "test-resource-axe-correlation-1-2",
      "state": 3,
      "output": "test-output-axe-correlation-1"
    }
    """
    When I wait the end of 2 events processing
    When I save response createTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-1&correlation=true
    Then the response code should be 200
    Then the response body should contain:
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

  Scenario: given meta alarm and ack event should ack children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-2
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-2",
      "connector_name": "test-connector-name-axe-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-2",
      "resource": "test-resource-axe-correlation-2",
      "state": 2,
      "output": "test-output-axe-correlation-2"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-2&correlation=true
    Then the response code should be 200
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
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-2"
    }
    """
    When I wait the end of 2 events processing
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
              "a": "root",
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
                "a": "root",
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

  Scenario: given meta alarm child and change state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-3
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 1,
      "output": "test-output-axe-correlation-3"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "changestate",
      "component":  "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 2,
      "output": "test-output-axe-correlation-3"
    }
    """
    When I wait the end of event processing
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

  Scenario: given meta alarm and change state event should update children
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-4
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-4",
      "connector_name": "test-connector-name-axe-correlation-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-4",
      "resource": "test-resource-axe-correlation-4",
      "state": 1,
      "output": "test-output-axe-correlation-4"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-4&correlation=true
    Then the response code should be 200
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
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-4"
    }
    """
    When I wait the end of 2 events processing
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
              "a": "root",
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
                "a": "root",
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

  Scenario: given meta alarm child and resolve event should resolve parent
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-5
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "state": 1,
      "output": "test-output-axe-correlation-5"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "cancel",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5"
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-5&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1

  Scenario: given meta alarm child and check with inc state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-6
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 1,
      "output": "test-output-axe-correlation-6"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 2
    }
    """
    When I wait the end of event processing
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

  Scenario: given meta alarm child and check with dec state event should update parent state
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-7
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 2,
      "output": "test-output-axe-correlation-7"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 1,
      "output": "test-output-axe-correlation-7"
    }
    """
    When I wait the end of event processing
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

  Scenario: given meta alarm rule and event should update child alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-8
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-1",
      "state": 2,
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-8&correlation=true
    Then the response code should be 200
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
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8-2",
      "state": 2,
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of 3 events processing
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

  Scenario: given meta alarm rule and 2 simultaneous events should create meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-9
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-9",
        "connector_name": "test-connector-name-axe-correlation-9",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-9",
        "resource": "test-resource-axe-correlation-9-1",
        "state": 2,
        "output": "test-output-axe-correlation-9"
      },
      {
        "connector": "test-connector-axe-correlation-9",
        "connector_name": "test-connector-name-axe-correlation-9",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-9",
        "resource": "test-resource-axe-correlation-9-2",
        "state": 3,
        "output": "test-output-axe-correlation-9"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-9&correlation=true
    Then the response code should be 200
    Then the response body should contain:
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

  Scenario: given meta alarm rule and 2 simultaneous events should update meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-10
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-10",
      "connector_name": "test-connector-name-axe-correlation-10",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-10",
      "resource": "test-resource-axe-correlation-10-1",
      "state": 1,
      "output": "test-output-axe-correlation-10"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-10",
        "connector_name": "test-connector-name-axe-correlation-10",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-10",
        "resource": "test-resource-axe-correlation-10-2",
        "state": 2,
        "output": "test-output-axe-correlation-10"
      },
      {
        "connector": "test-connector-axe-correlation-10",
        "connector_name": "test-connector-name-axe-correlation-10",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-10",
        "resource": "test-resource-axe-correlation-10-3",
        "state": 3,
        "output": "test-output-axe-correlation-10"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-10&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-10; Count: 3; Children: test-component-axe-correlation-10",
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
                  "component": "test-component-axe-correlation-10",
                  "connector": "test-connector-axe-correlation-10",
                  "connector_name": "test-connector-name-axe-correlation-10",
                  "resource": "test-resource-axe-correlation-10-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-10",
                  "connector": "test-connector-axe-correlation-10",
                  "connector_name": "test-connector-name-axe-correlation-10",
                  "resource": "test-resource-axe-correlation-10-2"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-10",
                  "connector": "test-connector-axe-correlation-10",
                  "connector_name": "test-connector-name-axe-correlation-10",
                  "resource": "test-resource-axe-correlation-10-3"
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
      },
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 2)._id }}",
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

  Scenario: given meta alarm rule and 2 simultaneous events for the same entity should create meta alarm
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-11
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-11",
        "connector_name": "test-connector-name-axe-correlation-11",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-11",
        "resource": "test-resource-axe-correlation-11",
        "state": 2,
        "output": "test-output-axe-correlation-11"
      },
      {
        "connector": "test-connector-axe-correlation-11",
        "connector_name": "test-connector-name-axe-correlation-11",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-11",
        "resource": "test-resource-axe-correlation-11",
        "state": 3,
        "output": "test-output-axe-correlation-11"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-11&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-11; Count: 1; Children: test-component-axe-correlation-11",
            "children": [
              "test-resource-axe-correlation-11/test-component-axe-correlation-11"
            ],
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
                  "component": "test-component-axe-correlation-11",
                  "connector": "test-connector-axe-correlation-11",
                  "connector_name": "test-connector-name-axe-correlation-11",
                  "resource": "test-resource-axe-correlation-11"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "_t": "stateinc",
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

  Scenario: given 2 meta alarm rules and event should create 2 meta alarms
    Given I am admin
    When I save response metaAlarmRuleID1=test-metaalarmrule-axe-correlation-12-1
    When I save response metaAlarmRuleID2=test-metaalarmrule-axe-correlation-12-2
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-12",
      "connector_name": "test-connector-name-axe-correlation-12",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-12",
      "resource": "test-resource-axe-correlation-12",
      "state": 3,
      "output": "test-output-axe-correlation-12"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-12&correlation=true&sort_by=v.output&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-12-1; Count: 1; Children: test-component-axe-correlation-12",
            "children": [
              "test-resource-axe-correlation-12/test-component-axe-correlation-12"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
              "_t": "stateinc",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            }
          }
        },
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-12-2; Count: 1; Children: test-component-axe-correlation-12",
            "children": [
              "test-resource-axe-correlation-12/test-component-axe-correlation-12"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
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
        "total_count": 2
      }
    }
    """
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
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
                  "component": "test-component-axe-correlation-12",
                  "connector": "test-connector-axe-correlation-12",
                  "connector_name": "test-connector-name-axe-correlation-12",
                  "resource": "test-resource-axe-correlation-12"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID2 }}",
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
                  "component": "test-component-axe-correlation-12",
                  "connector": "test-connector-axe-correlation-12",
                  "connector_name": "test-connector-name-axe-correlation-12",
                  "resource": "test-resource-axe-correlation-12"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "a": "engine.correlation"
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
              "total_count": 4
            }
          }
        }
      }
    ]
    """

  Scenario: given 2 meta alarm rules and event should update 2 meta alarms
    Given I am admin
    When I save response metaAlarmRuleID1=test-metaalarmrule-axe-correlation-13-1
    When I save response metaAlarmRuleID2=test-metaalarmrule-axe-correlation-13-2
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-13",
      "connector_name": "test-connector-name-axe-correlation-13",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-13",
      "resource": "test-resource-axe-correlation-13-1",
      "state": 2,
      "output": "test-output-axe-correlation-13"
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-13",
      "connector_name": "test-connector-name-axe-correlation-13",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-13",
      "resource": "test-resource-axe-correlation-13-2",
      "state": 3,
      "output": "test-output-axe-correlation-13"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-13&correlation=true&sort_by=v.output&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-13-1; Count: 2; Children: test-component-axe-correlation-13",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
              "_t": "stateinc",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            }
          }
        },
        {
          "v": {
            "output": "Rule: test-metaalarmrule-axe-correlation-13-2; Count: 2; Children: test-component-axe-correlation-13",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
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
        "total_count": 2
      }
    }
    """
    When I save response metaAlarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmID2={{ (index .lastResponse.data 1)._id }}
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .metaAlarmID1 }}",
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
                  "component": "test-component-axe-correlation-13",
                  "connector": "test-connector-axe-correlation-13",
                  "connector_name": "test-connector-name-axe-correlation-13",
                  "resource": "test-resource-axe-correlation-13-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-13",
                  "connector": "test-connector-axe-correlation-13",
                  "connector_name": "test-connector-name-axe-correlation-13",
                  "resource": "test-resource-axe-correlation-13-2"
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
        "_id": "{{ .metaAlarmID2 }}",
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
                  "component": "test-component-axe-correlation-13",
                  "connector": "test-connector-axe-correlation-13",
                  "connector_name": "test-connector-name-axe-correlation-13",
                  "resource": "test-resource-axe-correlation-13-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-13",
                  "connector": "test-connector-axe-correlation-13",
                  "connector_name": "test-connector-name-axe-correlation-13",
                  "resource": "test-resource-axe-correlation-13-2"
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
                "val": 3
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
                "_t": "metaalarmattach",
                "a": "engine.correlation"
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

  Scenario: given meta alarm event and child event should update child alarm in correct order
    Given I am admin
    When I save response metaAlarmRuleID=test-metaalarmrule-axe-correlation-14
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-14",
      "connector_name": "test-connector-name-axe-correlation-14",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-14",
      "resource": "test-resource-axe-correlation-14",
      "state": 3,
      "output": "test-output-axe-correlation-14"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-14&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """json
    [
      {
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "source_type": "resource",
        "event_type": "snooze",
        "duration": 6000,
        "component":  "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "output": "test-output-axe-correlation-8"
      },
      {
        "connector": "test-connector-axe-correlation-14",
        "connector_name": "test-connector-name-axe-correlation-14",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "test-component-axe-correlation-14",
        "resource": "test-resource-axe-correlation-14",
        "output": "test-output-axe-correlation-14"
      }
    ]
    """
    When I wait the end of 3 events processing
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
                  "component": "test-component-axe-correlation-14",
                  "connector": "test-connector-axe-correlation-14",
                  "connector_name": "test-connector-name-axe-correlation-14",
                  "resource": "test-resource-axe-correlation-14"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index (index .lastResponse 0).data.children.data 0)._id }}",
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
                "a": "engine.correlation"
              },
              {
                "_t": "ack"
              },
              {
                "_t": "snooze"
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

  Scenario: given manual meta alarm and updated or removed child alarm should update meta alarm
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-15",
        "connector_name": "test-connector-name-axe-correlation-15",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-15",
        "resource": "test-resource-axe-correlation-15-1",
        "state": 1,
        "output": "test-output-axe-correlation-15"
      },
      {
        "connector": "test-connector-axe-correlation-15",
        "connector_name": "test-connector-name-axe-correlation-15",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-15",
        "resource": "test-resource-axe-correlation-15-2",
        "state": 2,
        "output": "test-output-axe-correlation-15"
      },
      {
        "connector": "test-connector-axe-correlation-15",
        "connector_name": "test-connector-name-axe-correlation-15",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-15",
        "resource": "test-resource-axe-correlation-15-3",
        "state": 3,
        "output": "test-output-axe-correlation-15"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-15",
      "display_name": "test-metaAlarm-axe-correlation-15",
      "ma_children": [
        "test-resource-axe-correlation-15-1/test-component-axe-correlation-15",
        "test-resource-axe-correlation-15-2/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-15-1
    Then the response code should be 200
    When I save response metaAlarmEntityID={{ (index (index .lastResponse.data 0).v.parents 0) }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-15&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "test-output-axe-correlation-15",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metaAlarm-axe-correlation-15",
            "state": {
              "_t": "stateinc",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-15-3"
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
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_update",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-15",
      "ma_parents": [ "{{ .metaAlarmEntityID }}" ],
      "ma_children": [
        "test-resource-axe-correlation-15-3/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
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
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-3",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_ungroup",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-15",
      "ma_parents": [ "{{ .metaAlarmEntityID }}" ],
      "ma_children": [
        "test-resource-axe-correlation-15-3/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
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
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metaAlarmEntityID }}"]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metaAlarmEntityID }}"]
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-15-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-15",
            "connector": "test-connector-axe-correlation-15",
            "connector_name": "test-connector-name-axe-correlation-15",
            "resource": "test-resource-axe-correlation-15-3",
            "parents": []
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

  Scenario: given meta alarm and ack event should double ack children if they're already acked
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-16",
      "type": "attribute",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-16"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-16",
      "connector_name": "test-connector-name-axe-correlation-16",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-16",
      "resource": "test-resource-axe-correlation-16-1",
      "state": 2,
      "output": "test-output-axe-correlation-16",
      "long_output": "test-long-output-axe-correlation-16",
      "author": "test-author-axe-correlation-16"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-16",
      "connector_name": "test-connector-name-axe-correlation-16",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-16",
      "resource": "test-resource-axe-correlation-16-2",
      "state": 2,
      "output": "test-output-axe-correlation-16",
      "long_output": "test-long-output-axe-correlation-16",
      "author": "test-author-axe-correlation-16"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-16",
      "connector_name": "test-connector-name-axe-correlation-16",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-axe-correlation-16",
      "resource": "test-resource-axe-correlation-16-2",
      "output": "previous ack",
      "author": "test-author-axe-correlation-16"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-16&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I send an event:
    """
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "metaalarm ack",
      "long_output": "test-long-output-axe-correlation-16",
      "author": "test-author-axe-correlation-16"
    }
    """
    When I wait the end of 3 events processing
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
                  "component": "test-component-axe-correlation-16",
                  "connector": "test-connector-axe-correlation-16",
                  "connector_name": "test-connector-name-axe-correlation-16",
                  "resource": "test-resource-axe-correlation-16-1"
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-16",
                  "connector": "test-connector-axe-correlation-16",
                  "connector_name": "test-connector-name-axe-correlation-16",
                  "resource": "test-resource-axe-correlation-16-2"
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
      },
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
                "_t": "ack",
                "a": "test-author-axe-correlation-16",
                "m": "metaalarm ack",
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
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-16",
                "m": "metaalarm ack",
                "val": 0
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
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-16",
                "m": "previous ack",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-16",
                "m": "metaalarm ack",
                "val": 0
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

  Scenario: given meta alarm child and cancel event should update parent state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-17-1",
      "type": "attribute",
      "auto_resolve": false,
      "output_template": "{{ `{{ .Rule.ID }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-17"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-17-2",
      "type": "attribute",
      "auto_resolve": true,
      "output_template": "{{ `{{ .Rule.ID }}` }}",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-17"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID2={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-17",
      "connector_name": "test-connector-name-axe-correlation-17",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-17",
      "resource": "test-resource-axe-correlation-17",
      "state": 2,
      "output": "test-output-axe-correlation-17"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-17&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-17/test-component-axe-correlation-17"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            }
          }
        },
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-17/test-component-axe-correlation-17"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
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
        "total_count": 2
      }
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-17",
      "connector_name": "test-connector-name-axe-correlation-17",
      "source_type": "resource",
      "event_type": "cancel",
      "component":  "test-component-axe-correlation-17",
      "resource": "test-resource-axe-correlation-17",
      "output": "test-output-axe-correlation-17"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-17",
      "connector_name": "test-connector-name-axe-correlation-17",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component":  "test-component-axe-correlation-17",
      "resource": "test-resource-axe-correlation-17",
      "output": "test-output-axe-correlation-17"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID1 }}&active_columns[]=v.output&correlation=true&opened=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-17/test-component-axe-correlation-17"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID1 }}",
            "state": {
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
    When I do GET /api/v4/alarms?search={{ .metaAlarmRuleID2 }}&active_columns[]=v.output&correlation=true&opened=false until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "is_meta_alarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-17/test-component-axe-correlation-17"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID2 }}",
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

  Scenario: given meta alarm and assoc ticket event should ticket to children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-18",
      "type": "attribute",
      "auto_resolve": true,
      "output_template": "some",
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-component-axe-correlation-18"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-18",
      "connector_name": "test-connector-name-axe-correlation-18",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-18",
      "resource": "test-resource-axe-correlation-18-1",
      "state": 2,
      "output": "test-output-axe-correlation-18"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-18",
      "connector_name": "test-connector-name-axe-correlation-18",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-18",
      "resource": "test-resource-axe-correlation-18-2",
      "state": 2,
      "output": "test-output-axe-correlation-18"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-18&correlation=true
    Then the response code should be 200
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
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
      "event_type": "assocticket",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "author": "test-author-axe-correlation-18",
      "initiator": "user",
      "ticket": "test-ticket-axe-correlation-18",
      "ticket_url": "test-url-axe-correlation-18",
      "ticket_system_name": "test-system-name-axe-correlation-18",
      "ticket_data": {
        "ticket_param_1": "ticket_value_1"
      },
      "ticket_comment": "test-comment-axe-correlation-18"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-18-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-author-axe-correlation-18",
                "user_id": "root",
                "initiator": "system",
                "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-18",
                "ticket_url": "test-url-axe-correlation-18",
                "ticket_system_name": "test-system-name-axe-correlation-18",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-18"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-author-axe-correlation-18",
              "user_id": "root",
              "initiator": "system",
              "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-ticket-axe-correlation-18",
              "ticket_url": "test-url-axe-correlation-18",
              "ticket_system_name": "test-system-name-axe-correlation-18",
              "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-comment-axe-correlation-18"
            },
            "children": [],
            "component": "test-component-axe-correlation-18",
            "connector": "test-connector-axe-correlation-18",
            "connector_name": "test-connector-name-axe-correlation-18",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-18-1",
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
                "_t": "assocticket",
                "a": "test-author-axe-correlation-18",
                "user_id": "root",
                "initiator": "system",
                "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-18",
                "ticket_url": "test-url-axe-correlation-18",
                "ticket_system_name": "test-system-name-axe-correlation-18",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-18"
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
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-18-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-author-axe-correlation-18",
                "user_id": "root",
                "initiator": "system",
                "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-18",
                "ticket_url": "test-url-axe-correlation-18",
                "ticket_system_name": "test-system-name-axe-correlation-18",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-18"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-author-axe-correlation-18",
              "user_id": "root",
              "initiator": "system",
              "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
              "ticket": "test-ticket-axe-correlation-18",
              "ticket_url": "test-url-axe-correlation-18",
              "ticket_system_name": "test-system-name-axe-correlation-18",
              "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-comment-axe-correlation-18"
            },
            "children": [],
            "component": "test-component-axe-correlation-18",
            "connector": "test-connector-axe-correlation-18",
            "connector_name": "test-connector-name-axe-correlation-18",
            "parents": [
              "{{ .metaAlarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-18-2",
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
                "_t": "assocticket",
                "a": "test-author-axe-correlation-18",
                "user_id": "root",
                "initiator": "system",
                "m": "Ticket ID: test-ticket-axe-correlation-18. Ticket URL: test-url-axe-correlation-18. Ticket ticket_param_1: ticket_value_1.",
                "ticket": "test-ticket-axe-correlation-18",
                "ticket_url": "test-url-axe-correlation-18",
                "ticket_system_name": "test-system-name-axe-correlation-18",
                "ticket_meta_alarm_id": "{{ .metaAlarmID }}",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-comment-axe-correlation-18"
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

  Scenario: given manual meta alarm and added child should inherit metaalarm actions
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-19",
        "connector_name": "test-connector-name-axe-correlation-19",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-19",
        "resource": "test-resource-axe-correlation-19-1",
        "state": 1,
        "output": "test-output-axe-correlation-19"
      },
      {
        "connector": "test-connector-axe-correlation-19",
        "connector_name": "test-connector-name-axe-correlation-19",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-19",
        "resource": "test-resource-axe-correlation-19-2",
        "state": 2,
        "output": "test-output-axe-correlation-19"
      },
      {
        "connector": "test-connector-axe-correlation-19",
        "connector_name": "test-connector-name-axe-correlation-19",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-correlation-19",
        "resource": "test-resource-axe-correlation-19-3",
        "state": 3,
        "output": "test-output-axe-correlation-19"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_group",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-19",
      "display_name": "test-metaAlarm-axe-correlation-19",
      "ma_children": [
        "test-resource-axe-correlation-19-1/test-component-axe-correlation-19",
        "test-resource-axe-correlation-19-2/test-component-axe-correlation-19"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-19-1&correlation=true
    Then the response code should be 200
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-19&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "test-output-axe-correlation-19",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metaAlarm-axe-correlation-19",
            "state": {
              "_t": "stateinc",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            }
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-19-3"
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
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-ack-19"
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "assocticket",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "ticket": "ticket-19"
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "comment-19"
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "snooze",
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "duration": 3600
    }
    """
    When I wait the end of 3 events processing
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_update",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-19",
      "ma_parents": [ "{{ .metaAlarmEntityID }}" ],
      "ma_children": [
        "test-resource-axe-correlation-19-3/test-component-axe-correlation-19"
      ]
    }
    """
    When I wait the end of 6 events processing
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
                  "component": "test-component-axe-correlation-19",
                  "connector": "test-connector-axe-correlation-19",
                  "connector_name": "test-connector-name-axe-correlation-19",
                  "resource": "test-resource-axe-correlation-19-1",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root"
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-19",
                  "connector": "test-connector-axe-correlation-19",
                  "connector_name": "test-connector-name-axe-correlation-19",
                  "resource": "test-resource-axe-correlation-19-2",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root"
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-19",
                  "connector": "test-connector-axe-correlation-19",
                  "connector_name": "test-connector-name-axe-correlation-19",
                  "resource": "test-resource-axe-correlation-19-3",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root"
                  }
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
