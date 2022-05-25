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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "id": "test-metaalarmrule-axe-correlation-1"
          },
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
            },
            "steps": [
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
    When I save response metaalarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I save response metaalarmCreationDate={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    Then the difference between metaalarmLastEventDate createTimestamp is in range -2,2
    Then the difference between metaalarmCreationDate createTimestamp is in range -2,2
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-1-2"}]}&with_steps=true
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
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-1-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I save response alarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I save response metaAlarmAttachStepDate={{ ( index (index .lastResponse.data 0).v.steps 2).t }}
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-2"}]}&with_steps=true
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
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc until the response code should be 200 and the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-3",
                  "connector": "test-connector-axe-correlation-3",
                  "connector_name": "test-connector-name-axe-correlation-3",
                  "resource": "test-resource-axe-correlation-3",
                  "steps": [
                    {
                      "_t": "stateinc",
                      "val": 1
                    },
                    {
                      "_t": "statusinc",
                      "val": 1
                    },
                    {
                      "_t": "metaalarmattach"
                    },
                    {
                      "_t": "changestate",
                      "val": 2
                    }
                  ]
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-3"
          },
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
            },
            "steps": [
              {
                "_t": "stateinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "statusinc",
                "val": 2
              }
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-4"}]}&with_steps=true&with_consequences=true
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
              "{{ .metalarmEntityID }}"
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
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc until the response code should be 200 and response key "data 0.v.resolved" is greater or equal than {{ .resolveTimestamp }}

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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc until the response code should be 200 and the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-6",
                  "connector": "test-connector-axe-correlation-6",
                  "connector_name": "test-connector-name-axe-correlation-6",
                  "resource": "test-resource-axe-correlation-6",
                  "steps": [
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
                      "val": 0
                    },
                    {
                      "_t": "stateinc",
                      "val": 2
                    }
                  ]
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-6"
          },
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
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc until the response code should be 200 and the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-7",
                  "connector": "test-connector-axe-correlation-7",
                  "connector_name": "test-connector-name-axe-correlation-7",
                  "resource": "test-resource-axe-correlation-7",
                  "steps": [
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
                      "val": 0
                    },
                    {
                      "_t": "statedec",
                      "val": 1
                    }
                  ]
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-7"
          },
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
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-8-2"}]}&with_steps=true
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
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-8-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID }}"
          },
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
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-9-1"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-9",
            "connector": "test-connector-axe-correlation-9",
            "connector_name": "test-connector-name-axe-correlation-9",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-9-1",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-9-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-9",
            "connector": "test-connector-axe-correlation-9",
            "connector_name": "test-connector-name-axe-correlation-9",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-9-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 3
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID }}"
          },
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
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-10-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-10",
            "connector": "test-connector-axe-correlation-10",
            "connector_name": "test-connector-name-axe-correlation-10",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-10-2",
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-10-3"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-10",
            "connector": "test-connector-axe-correlation-10",
            "connector_name": "test-connector-name-axe-correlation-10",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-10-3",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID }}"
          },
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
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-11"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-11",
            "connector": "test-connector-axe-correlation-11",
            "connector_name": "test-connector-name-axe-correlation-11",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-11",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID1 }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID1 }}"
          },
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
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "val": 1
              }
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
    When I save response metalarmEntityID1={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID2 }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID2 }}"
          },
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
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "val": 3
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "val": 1
              }
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
    When I save response metalarmEntityID2={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-12"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-12",
            "connector": "test-connector-axe-correlation-12",
            "connector_name": "test-connector-name-axe-correlation-12",
            "resource": "test-resource-axe-correlation-12",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID1 }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID1 }}"
          },
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
            },
            "steps": [
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
    When I save response metalarmEntityID1={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID2 }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
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
            "total": 2
          },
          "metaalarm": true,
          "rule": {
            "id": "{{ .metaAlarmRuleID2 }}"
          },
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
            },
            "steps": [
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
    When I save response metalarmEntityID2={{ (index .lastResponse.data 0).entity._id }}
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-13-2"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-13",
            "connector": "test-connector-axe-correlation-13",
            "connector_name": "test-connector-name-axe-correlation-13",
            "resource": "test-resource-axe-correlation-13-2",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-14"}]}&with_steps=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-correlation-14",
            "connector": "test-connector-axe-correlation-14",
            "connector_name": "test-connector-name-axe-correlation-14",
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-14",
            "state": {
              "val": 3
            },
            "status": {
              "val": 1
            },
            "steps": [
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
      "display_name": "test-metalarm-axe-correlation-15",
      "ma_children": [
        "test-resource-axe-correlation-15-1/test-component-axe-correlation-15",
        "test-resource-axe-correlation-15-2/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-15-1"}]}
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index (index .lastResponse.data 0).v.parents 0) }}
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"{{ .metalarmEntityID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              }
            ],
            "total": 2
          },
          "metaalarm": true,
          "v": {
            "output": "",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metalarm-axe-correlation-15",
            "state": {
              "_t": "stateinc",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "val": 1
              }
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
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_update",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-15",
      "ma_parents": [ "{{ .metalarmEntityID }}" ],
      "ma_children": [
        "test-resource-axe-correlation-15-3/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"{{ .metalarmEntityID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-3",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              }
            ],
            "total": 3
          },
          "metaalarm": true,
          "v": {
            "output": "",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metalarm-axe-correlation-15",
            "state": {
              "_t": "stateinc",
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            },
            "steps": [
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
    When I send an event:
    """json
    {
      "connector": "engine",
      "connector_name": "correlation",
      "source_type": "component",
      "event_type": "manual_metaalarm_ungroup",
      "component":  "metaalarm",
      "output": "test-output-axe-correlation-15",
      "ma_parents": [ "{{ .metalarmEntityID }}" ],
      "ma_children": [
        "test-resource-axe-correlation-15-3/test-component-axe-correlation-15"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"d":"{{ .metalarmEntityID }}"}]}&with_steps=true&with_consequences=true&correlation=true&sort_key=v.resource&sort_dir=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "consequences": {
            "data": [
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-1",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-15",
                  "connector": "test-connector-axe-correlation-15",
                  "connector_name": "test-connector-name-axe-correlation-15",
                  "resource": "test-resource-axe-correlation-15-2",
                  "parents": ["{{ .metalarmEntityID }}"],
                  "steps": [
                    {
                      "_t": "stateinc"
                    },
                    {
                      "_t": "statusinc"
                    },
                    {
                      "_t": "metaalarmattach"
                    }
                  ]
                }
              }
            ],
            "total": 2
          },
          "metaalarm": true,
          "v": {
            "output": "",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metalarm-axe-correlation-15",
            "state": {
              "_t": "statedec",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "val": 1
            },
            "steps": [
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
              },
              {
                "_t": "statedec",
                "a": "engine.correlation",
                "val": 2
              }
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-15-3"}]}&with_steps=true
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
            "parents": [],
            "steps": [
              {
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "metaalarmattach"
              }
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
