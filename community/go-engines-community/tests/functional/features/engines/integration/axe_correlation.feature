Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  Scenario: given meta alarm rule and event should create meta alarm
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-1",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-1"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-1",
      "connector_name": "test-connector-name-axe-correlation-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-1",
      "resource": "test-resource-axe-correlation-1",
      "state": 2,
      "output": "test-output-axe-correlation-1",
      "long_output": "test-long-output-axe-correlation-1",
      "author": "test-author-axe-correlation-1",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I save response createTimestamp={{ now }}
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
                  "component": "test-component-axe-correlation-1",
                  "connector": "test-connector-axe-correlation-1",
                  "connector_name": "test-connector-name-axe-correlation-1",
                  "resource": "test-resource-axe-correlation-1"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-1"
          },
          "v": {
            "children": [
              "test-resource-axe-correlation-1/test-component-axe-correlation-1"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "last_update_date": {{ .checkEventTimestamp }},
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "stateinc",
              "a": "engine.correlation",
              "m": "",
              "t": {{ .checkEventTimestamp }},
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "m": "",
              "t": {{ .checkEventTimestamp }},
              "val": 1
            },
            "steps": [
              {
                "_t": "stateinc",
                "a": "engine.correlation",
                "m": "",
                "t": {{ .checkEventTimestamp }},
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "engine.correlation",
                "m": "",
                "t": {{ .checkEventTimestamp }},
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
    When I save response metaalarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    When I save response metaalarmCreationDate={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    Then the difference between metaalarmLastEventDate createTimestamp is in range -2,2
    Then the difference between metaalarmCreationDate createTimestamp is in range -2,2
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-1"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "children": [],
            "component": "test-component-axe-correlation-1",
            "connector": "test-connector-axe-correlation-1",
            "connector_name": "test-connector-name-axe-correlation-1",
            "initial_long_output": "test-long-output-axe-correlation-1",
            "initial_output": "test-output-axe-correlation-1",
            "last_update_date": {{ .checkEventTimestamp }},
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-1",
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
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-2",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-2"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-2",
      "connector_name": "test-connector-name-axe-correlation-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-2",
      "resource": "test-resource-axe-correlation-2",
      "state": 2,
      "output": "test-output-axe-correlation-2",
      "long_output": "test-long-output-axe-correlation-2",
      "author": "test-author-axe-correlation-2",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
      "output": "test-output-axe-correlation-2",
      "long_output": "test-long-output-axe-correlation-2",
      "author": "test-author-axe-correlation-2",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
                  "component": "test-component-axe-correlation-2",
                  "connector": "test-connector-axe-correlation-2",
                  "connector_name": "test-connector-name-axe-correlation-2",
                  "resource": "test-resource-axe-correlation-2"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-2"
          },
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-correlation-2",
              "m": "test-output-axe-correlation-2",
              "t": {{ .ackEventTimestamp }},
              "val": 0
            },
            "children": [
              "test-resource-axe-correlation-2/test-component-axe-correlation-2"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "last_update_date": {{ .checkEventTimestamp }},
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
                "val": 2
              },
              {
                "_t": "statusinc",
                "val": 1
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-2",
                "m": "test-output-axe-correlation-2",
                "t": {{ .ackEventTimestamp }},
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-2"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-correlation-2",
              "m": "test-output-axe-correlation-2",
              "t": {{ .ackEventTimestamp }},
              "val": 0
            },
            "children": [],
            "component": "test-component-axe-correlation-2",
            "connector": "test-connector-axe-correlation-2",
            "connector_name": "test-connector-name-axe-correlation-2",
            "initial_long_output": "test-long-output-axe-correlation-2",
            "initial_output": "test-output-axe-correlation-2",
            "last_update_date": {{ .checkEventTimestamp }},
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
                "_t": "metaalarmattach",
                "a": "engine.correlation",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-correlation-2",
                "m": "test-output-axe-correlation-2",
                "t": {{ .ackEventTimestamp }},
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

  Scenario: given meta alarm child and change state event should not update parent state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-3",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-3"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 1,
      "output": "test-output-axe-correlation-3",
      "long_output": "test-long-output-axe-correlation-3",
      "author": "test-author-axe-correlation-3",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-3",
      "connector_name": "test-connector-name-axe-correlation-3",
      "source_type": "resource",
      "event_type": "changestate",
      "component":  "test-component-axe-correlation-3",
      "resource": "test-resource-axe-correlation-3",
      "state": 2,
      "output": "test-output-axe-correlation-3",
      "long_output": "test-long-output-axe-correlation-3",
      "author": "test-author-axe-correlation-3",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response changeStateEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
                      "_t": "metaalarmattach",
                      "val": 0
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
            "children": [
              "test-resource-axe-correlation-3/test-component-axe-correlation-3"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "last_update_date": {{ .checkEventTimestamp }},
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
                "val": 1
              },
              {
                "_t": "statusinc",
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

  Scenario: given meta alarm and change state event should update children
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-4",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-4"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-4",
      "connector_name": "test-connector-name-axe-correlation-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-4",
      "resource": "test-resource-axe-correlation-4",
      "state": 1,
      "output": "test-output-axe-correlation-4",
      "long_output": "test-long-output-axe-correlation-4",
      "author": "test-author-axe-correlation-4",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    When I save response metalarmEntityID={{ (index .lastResponse.data 0).entity._id }}
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
      "event_type": "changestate",
      "state": 2,
      "component":  "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-output-axe-correlation-4",
      "long_output": "test-long-output-axe-correlation-4",
      "author": "test-author-axe-correlation-4",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response changeStateEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
                  "component": "test-component-axe-correlation-4",
                  "connector": "test-connector-axe-correlation-4",
                  "connector_name": "test-connector-name-axe-correlation-4",
                  "resource": "test-resource-axe-correlation-4"
                }
              }
            ],
            "total": 1
          },
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-4"
          },
          "v": {
            "children": [
              "test-resource-axe-correlation-4/test-component-axe-correlation-4"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "last_update_date": {{ .checkEventTimestamp }},
            "meta": "{{ .metaAlarmRuleID }}",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-correlation-4",
              "m": "test-output-axe-correlation-4",
              "t": {{ .changeStateEventTimestamp }},
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
                "_t": "changestate",
                "a": "test-author-axe-correlation-4",
                "m": "test-output-axe-correlation-4",
                "t": {{ .changeStateEventTimestamp }},
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
    When I do GET /api/v4/alarms?filter={"$and":[{"entity.name":"test-resource-axe-correlation-4"}]}&with_steps=true&with_consequences=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "children": [],
            "component": "test-component-axe-correlation-4",
            "connector": "test-connector-axe-correlation-4",
            "connector_name": "test-connector-name-axe-correlation-4",
            "initial_long_output": "test-long-output-axe-correlation-4",
            "initial_output": "test-output-axe-correlation-4",
            "last_update_date": {{ .checkEventTimestamp }},
            "parents": [
              "{{ .metalarmEntityID }}"
            ],
            "resource": "test-resource-axe-correlation-4",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-correlation-4",
              "m": "test-output-axe-correlation-4",
              "t": {{ .changeStateEventTimestamp }},
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
                "a": "test-author-axe-correlation-4",
                "m": "test-output-axe-correlation-4",
                "t": {{ .changeStateEventTimestamp }},
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
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-5",
      "type": "attribute",
      "auto_resolve": true,
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-5"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "state": 1,
      "output": "test-output-axe-correlation-5",
      "long_output": "test-long-output-axe-correlation-5",
      "author": "test-author-axe-correlation-5",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "done",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5",
      "long_output": "test-long-output-axe-correlation-5",
      "author": "test-author-axe-correlation-5",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-5",
      "connector_name": "test-connector-name-axe-correlation-5",
      "source_type": "resource",
      "event_type": "resolve_done",
      "component":  "test-component-axe-correlation-5",
      "resource": "test-resource-axe-correlation-5",
      "output": "test-output-axe-correlation-5",
      "long_output": "test-long-output-axe-correlation-5",
      "author": "test-author-axe-correlation-5",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.meta":"{{ .metaAlarmRuleID }}"}]}&with_steps=true&with_consequences=true&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "metaalarm": true,
          "rule": {
            "name": "test-metaalarmrule-axe-correlation-5"
          },
          "v": {
            "children": [
              "test-resource-axe-correlation-5/test-component-axe-correlation-5"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "last_update_date": {{ .checkEventTimestamp }},
            "meta": "{{ .metaAlarmRuleID }}",
            "resolved": {{ .resolveTimestamp }},
            "state": {
              "val": 1
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

  Scenario: given meta alarm child and check with inc state event should update parent state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-6",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-6"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 1,
      "output": "test-output-axe-correlation-6",
      "long_output": "test-long-output-axe-correlation-6",
      "author": "test-author-axe-correlation-6",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-6",
      "connector_name": "test-connector-name-axe-correlation-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-6",
      "resource": "test-resource-axe-correlation-6",
      "state": 2,
      "output": "test-output-axe-correlation-6",
      "long_output": "test-long-output-axe-correlation-6",
      "author": "test-author-axe-correlation-6",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
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
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "name": "test-metaalarmrule-axe-correlation-7",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-7"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 2,
      "output": "test-output-axe-correlation-7",
      "long_output": "test-long-output-axe-correlation-7",
      "author": "test-author-axe-correlation-7",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector": "test-connector-axe-correlation-7",
      "connector_name": "test-connector-name-axe-correlation-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-7",
      "resource": "test-resource-axe-correlation-7",
      "state": 1,
      "output": "test-output-axe-correlation-7",
      "long_output": "test-long-output-axe-correlation-7",
      "author": "test-author-axe-correlation-7",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
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

  Scenario: given meta alarm child and cancel event should update parent state
    Given I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """json
    {
      "name": "test-metaalarmrule-axe-correlation-8",
      "type": "attribute",
      "config": {
        "alarm_patterns": [
          {
            "v": {
              "component": "test-component-axe-correlation-8"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then I save response metaAlarmRuleID={{ .lastResponse._id }}
    When I wait the next periodical process
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8",
      "state": 2,
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-8&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "metaalarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-8/test-component-axe-correlation-8"
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
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "cancel",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8",
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-axe-correlation-8",
      "connector_name": "test-connector-name-axe-correlation-8",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component":  "test-component-axe-correlation-8",
      "resource": "test-resource-axe-correlation-8",
      "output": "test-output-axe-correlation-8"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-8&correlation=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "metaalarm": true,
          "v": {
            "children": [
              "test-resource-axe-correlation-8/test-component-axe-correlation-8"
            ],
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "meta": "{{ .metaAlarmRuleID }}",
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
