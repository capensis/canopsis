Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given check event should create alarm
    Given I am admin
    When I save response createTimestamp={{ now }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-1",
      "connector_name": "test-connector-name-axe-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-1",
      "resource": "test-resource-axe-1",
      "state": 2,
      "output": "test-output-axe-1",
      "long_output": "test-long-output-axe-1",
      "author": "test-author-axe-1",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response eventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-1/test-component-axe-1"
          },
          "infos": {},
          "tags": [],
          "v": {
            "children": [],
            "component": "test-component-axe-1",
            "connector": "test-connector-axe-1",
            "connector_name": "test-connector-name-axe-1",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "test-long-output-axe-1",
            "initial_output": "test-output-axe-1",
            "last_update_date": {{ .eventTimestamp }},
            "long_output": "test-long-output-axe-1",
            "long_output_history": ["test-long-output-axe-1"],
            "output": "test-output-axe-1",
            "parents": [],
            "resource": "test-resource-axe-1",
            "state": {
              "_t": "stateinc",
              "a": "test-connector-axe-1.test-connector-name-axe-1",
              "m": "test-output-axe-1",
              "t": {{ .eventTimestamp }},
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-axe-1.test-connector-name-axe-1",
              "m": "test-output-axe-1",
              "t": {{ .eventTimestamp }},
              "val": 1
            },
            "total_state_changes": 1
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
    When I save response alarmTimestamp={{ (index .lastResponse.data 0).t }}
    When I save response alarmCreationDate={{ (index .lastResponse.data 0).v.creation_date }}
    When I save response alarmLastEventDate={{ (index .lastResponse.data 0).v.last_event_date }}
    Then the difference between alarmTimestamp createTimestamp is in range -2,2
    Then the difference between alarmCreationDate createTimestamp is in range -2,2
    Then the difference between alarmLastEventDate createTimestamp is in range -2,2
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
                "a": "test-connector-axe-1.test-connector-name-axe-1",
                "m": "test-output-axe-1",
                "t": {{ .eventTimestamp }},
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-1.test-connector-name-axe-1",
                "m": "test-output-axe-1",
                "t": {{ .eventTimestamp }},
                "val": 1
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
  Scenario: given check event should update alarm
    Given I am admin
    When I save response createTimestamp={{ now }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-2",
      "connector_name": "test-connector-name-axe-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-2",
      "resource": "test-resource-axe-2",
      "state": 2,
      "output": "test-output-axe-2",
      "long_output": "test-long-output-axe-2-1",
      "author": "test-author-axe-2",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response firstEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-2",
      "connector_name": "test-connector-name-axe-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-2",
      "resource": "test-resource-axe-2",
      "state": 3,
      "output": "test-output-axe-2",
      "long_output": "test-long-output-axe-2-2",
      "author": "test-author-axe-2",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response secondEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-2/test-component-axe-2"
          },
          "infos": {},
          "v": {
            "children": [],
            "component": "test-component-axe-2",
            "connector": "test-connector-axe-2",
            "connector_name": "test-connector-name-axe-2",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "test-long-output-axe-2-1",
            "initial_output": "test-output-axe-2",
            "last_update_date": {{ .secondEventTimestamp }},
            "long_output": "test-long-output-axe-2-2",
            "long_output_history": ["test-long-output-axe-2-1", "test-long-output-axe-2-2"],
            "output": "test-output-axe-2",
            "parents": [],
            "resource": "test-resource-axe-2",
            "state": {
              "_t": "stateinc",
              "a": "test-connector-axe-2.test-connector-name-axe-2",
              "m": "test-output-axe-2",
              "t": {{ .secondEventTimestamp }},
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-axe-2.test-connector-name-axe-2",
              "m": "test-output-axe-2",
              "t": {{ .firstEventTimestamp }},
              "val": 1
            },
            "state_changes_since_status_update": 1,
            "total_state_changes": 2
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
    When I save response alarmTimestamp={{ (index .lastResponse.data 0).t }}
    When I save response alarmCreationDate={{ (index .lastResponse.data 0).v.creation_date }}
    Then the difference between alarmTimestamp createTimestamp is in range -2,2
    Then the difference between alarmCreationDate createTimestamp is in range -2,2
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
                "a": "test-connector-axe-2.test-connector-name-axe-2",
                "m": "test-output-axe-2",
                "t": {{ .firstEventTimestamp }},
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-2.test-connector-name-axe-2",
                "m": "test-output-axe-2",
                "t": {{ .firstEventTimestamp }},
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "test-connector-axe-2.test-connector-name-axe-2",
                "m": "test-output-axe-2",
                "t": {{ .secondEventTimestamp }},
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

  @concurrent
  Scenario: given ack event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-3",
      "connector_name": "test-connector-name-axe-3",
      "source_type": "resource",
      "component":  "test-component-axe-3",
      "resource": "test-resource-axe-3",
      "state": 2,
      "output": "test-output-axe-3",
      "long_output": "test-long-output-axe-3",
      "author": "test-author-axe-3",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-3",
      "connector_name": "test-connector-name-axe-3",
      "source_type": "resource",
      "component":  "test-component-axe-3",
      "resource": "test-resource-axe-3",
      "output": "test-output-axe-3",
      "long_output": "test-long-output-axe-3",
      "author": "test-author-axe-3",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-3",
              "m": "test-output-axe-3",
              "t": {{ .ackEventTimestamp }},
              "val": 0
            },
            "component": "test-component-axe-3",
            "connector": "test-connector-axe-3",
            "connector_name": "test-connector-name-axe-3",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-3",
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
                "_t": "ack",
                "a": "test-author-axe-3",
                "m": "test-output-axe-3",
                "t": {{ .ackEventTimestamp }},
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

  @concurrent
  Scenario: given remove ack event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-4",
      "connector_name": "test-connector-name-axe-4",
      "source_type": "resource",
      "component":  "test-component-axe-4",
      "resource": "test-resource-axe-4",
      "state": 2,
      "output": "test-output-axe-4",
      "long_output": "test-long-output-axe-4",
      "author": "test-author-axe-4",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-4",
      "connector_name": "test-connector-name-axe-4",
      "source_type": "resource",
      "component":  "test-component-axe-4",
      "resource": "test-resource-axe-4",
      "output": "test-output-axe-4",
      "long_output": "test-long-output-axe-4",
      "author": "test-author-axe-4",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ackremove",
      "connector": "test-connector-axe-4",
      "connector_name": "test-connector-name-axe-4",
      "source_type": "resource",
      "component":  "test-component-axe-4",
      "resource": "test-resource-axe-4",
      "output": "test-output-axe-4",
      "long_output": "test-long-output-axe-4",
      "author": "test-author-axe-4",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackRemoveEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-4",
            "connector": "test-connector-axe-4",
            "connector_name": "test-connector-name-axe-4",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-4",
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
    Then the response key "data.0.v.ack" should not exist
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
                "_t": "ack",
                "a": "test-author-axe-4",
                "m": "test-output-axe-4",
                "t": {{ .ackEventTimestamp }},
                "val": 0
              },
              {
                "_t": "ackremove",
                "a": "test-author-axe-4",
                "m": "test-output-axe-4",
                "t": {{ .ackRemoveEventTimestamp }},
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
      }
    ]
    """

  @concurrent
  Scenario: given cancel event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-5",
      "connector_name": "test-connector-name-axe-5",
      "source_type": "resource",
      "component":  "test-component-axe-5",
      "resource": "test-resource-axe-5",
      "state": 2,
      "output": "test-output-axe-5",
      "long_output": "test-long-output-axe-5",
      "author": "test-author-axe-5",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-5",
      "connector_name": "test-connector-name-axe-5",
      "source_type": "resource",
      "component":  "test-component-axe-5",
      "resource": "test-resource-axe-5",
      "output": "test-output-axe-5",
      "long_output": "test-long-output-axe-5",
      "author": "test-author-axe-5",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "canceled": {
              "_t": "cancel",
              "a": "test-author-axe-5",
              "m": "test-output-axe-5",
              "t": {{ .cancelEventTimestamp }},
              "val": 0
            },
            "component": "test-component-axe-5",
            "connector": "test-connector-axe-5",
            "connector_name": "test-connector-name-axe-5",
            "last_update_date": {{ .cancelEventTimestamp }},
            "resource": "test-resource-axe-5",
            "state": {
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "test-author-axe-5",
              "m": "test-output-axe-5",
              "t": {{ .cancelEventTimestamp }},
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
                "_t": "cancel",
                "a": "test-author-axe-5",
                "m": "test-output-axe-5",
                "t": {{ .cancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-author-axe-5",
                "m": "test-output-axe-5",
                "t": {{ .cancelEventTimestamp }},
                "val": 4
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
  Scenario: given uncancel event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-6",
      "connector_name": "test-connector-name-axe-6",
      "source_type": "resource",
      "component":  "test-component-axe-6",
      "resource": "test-resource-axe-6",
      "state": 2,
      "output": "test-output-axe-6",
      "long_output": "test-long-output-axe-6",
      "author": "test-author-axe-6",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-6",
      "connector_name": "test-connector-name-axe-6",
      "source_type": "resource",
      "component":  "test-component-axe-6",
      "resource": "test-resource-axe-6",
      "output": "test-output-axe-6",
      "long_output": "test-long-output-axe-6",
      "author": "test-author-axe-6",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "uncancel",
      "connector": "test-connector-axe-6",
      "connector_name": "test-connector-name-axe-6",
      "source_type": "resource",
      "component":  "test-component-axe-6",
      "resource": "test-resource-axe-6",
      "output": "test-output-axe-6",
      "long_output": "test-long-output-axe-6",
      "author": "test-author-axe-6",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response uncancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-6",
            "connector": "test-connector-axe-6",
            "connector_name": "test-connector-name-axe-6",
            "last_update_date": {{ .uncancelEventTimestamp }},
            "resource": "test-resource-axe-6",
            "state": {
              "val": 2
            },
            "status": {
              "_t": "statusdec",
              "a": "test-author-axe-6",
              "m": "test-output-axe-6",
              "t": {{ .uncancelEventTimestamp }},
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
    Then the response key "data.0.v.canceled" should not exist
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
                "_t": "cancel",
                "a": "test-author-axe-6",
                "m": "test-output-axe-6",
                "t": {{ .cancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-author-axe-6",
                "m": "test-output-axe-6",
                "t": {{ .cancelEventTimestamp }},
                "val": 4
              },
              {
                "_t": "uncancel",
                "a": "test-author-axe-6",
                "m": "test-output-axe-6",
                "t": {{ .uncancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-author-axe-6",
                "m": "test-output-axe-6",
                "t": {{ .uncancelEventTimestamp }},
                "val": 1
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 10,
              "total_count": 6
            }
          }
        }
      }
    ]
    """

  @concurrent
  Scenario: given comment event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-7",
      "connector_name": "test-connector-name-axe-7",
      "source_type": "resource",
      "component":  "test-component-axe-7",
      "resource": "test-resource-axe-7",
      "state": 2,
      "output": "test-output-axe-7",
      "long_output": "test-long-output-axe-7",
      "author": "test-author-axe-7",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "comment",
      "connector": "test-connector-axe-7",
      "connector_name": "test-connector-name-axe-7",
      "source_type": "resource",
      "component":  "test-component-axe-7",
      "resource": "test-resource-axe-7",
      "output": "test-output-axe-7-1",
      "long_output": "test-long-output-axe-7",
      "author": "test-author-axe-7",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response commentEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-7",
            "connector": "test-connector-axe-7",
            "connector_name": "test-connector-name-axe-7",
            "last_update_date": {{ .checkEventTimestamp }},
            "last_comment": {
              "_t": "comment",
              "a": "test-author-axe-7",
              "m": "test-output-axe-7-1",
              "t": {{ .commentEventTimestamp }},
              "val": 0
            },
            "resource": "test-resource-axe-7",
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
                "_t": "comment",
                "a": "test-author-axe-7",
                "m": "test-output-axe-7-1",
                "t": {{ .commentEventTimestamp }},
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "comment",
      "connector": "test-connector-axe-7",
      "connector_name": "test-connector-name-axe-7",
      "source_type": "resource",
      "component":  "test-component-axe-7",
      "resource": "test-resource-axe-7",
      "output": "test-output-axe-7-2",
      "long_output": "test-long-output-axe-7",
      "author": "test-author-axe-7",
      "timestamp": {{ nowAdd "-2s" }}
    }
    """
    When I save response commentEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-7",
            "connector": "test-connector-axe-7",
            "connector_name": "test-connector-name-axe-7",
            "last_comment": {
              "_t": "comment",
              "a": "test-author-axe-7",
              "m": "test-output-axe-7-2",
              "t": {{ .commentEventTimestamp }},
              "val": 0
            },
            "resource": "test-resource-axe-7",
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
                "_t": "comment",
                "a": "test-author-axe-7",
                "m": "test-output-axe-7-1",
                "val": 0
              },
              {
                "_t": "comment",
                "a": "test-author-axe-7",
                "m": "test-output-axe-7-2",
                "t": {{ .commentEventTimestamp }},
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
      }
    ]
    """

  @concurrent
  Scenario: given assoc ticket event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-8",
      "connector_name": "test-connector-name-axe-8",
      "source_type": "resource",
      "component":  "test-component-axe-8",
      "resource": "test-resource-axe-8",
      "state": 2,
      "output": "test-output-axe-8",
      "long_output": "test-long-output-axe-8",
      "author": "test-author-axe-8",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "assocticket",
      "ticket": "test-ticket",
      "ticket_system_name": "test-system-name",
      "ticket_url": "test-ticket-url",
      "ticket_data": {
        "ticket_param_1": "ticket_value_1"
      },
      "ticket_comment": "test-ticket-comment",
      "connector": "test-connector-axe-8",
      "connector_name": "test-connector-name-axe-8",
      "source_type": "resource",
      "component":  "test-component-axe-8",
      "resource": "test-resource-axe-8",
      "author": "test-author-axe-8",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ticketEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-8
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
                "a": "test-author-axe-8",
                "m": "Ticket ID: test-ticket. Ticket URL: test-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "t": {{ .ticketEventTimestamp }},
                "ticket": "test-ticket",
                "ticket_system_name": "test-system-name",
                "ticket_url": "test-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-ticket-comment"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "test-author-axe-8",
              "m": "Ticket ID: test-ticket. Ticket URL: test-ticket-url. Ticket ticket_param_1: ticket_value_1.",
              "t": {{ .ticketEventTimestamp }},
              "ticket": "test-ticket",
              "ticket_system_name": "test-system-name",
              "ticket_url": "test-ticket-url",
              "ticket_data": {
                "ticket_param_1": "ticket_value_1"
              },
              "ticket_comment": "test-ticket-comment"
            },
            "component": "test-component-axe-8",
            "connector": "test-connector-axe-8",
            "connector_name": "test-connector-name-axe-8",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-8",
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
                "_t": "assocticket",
                "a": "test-author-axe-8",
                "m": "Ticket ID: test-ticket. Ticket URL: test-ticket-url. Ticket ticket_param_1: ticket_value_1.",
                "t": {{ .ticketEventTimestamp }},
                "ticket": "test-ticket",
                "ticket_system_name": "test-system-name",
                "ticket_url": "test-ticket-url",
                "ticket_data": {
                  "ticket_param_1": "ticket_value_1"
                },
                "ticket_comment": "test-ticket-comment"
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
  Scenario: given change state event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-9",
      "connector_name": "test-connector-name-axe-9",
      "source_type": "resource",
      "component":  "test-component-axe-9",
      "resource": "test-resource-axe-9",
      "state": 2,
      "output": "test-output-axe-9",
      "long_output": "test-long-output-axe-9",
      "author": "test-author-axe-9",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 3,
      "connector": "test-connector-axe-9",
      "connector_name": "test-connector-name-axe-9",
      "source_type": "resource",
      "component":  "test-component-axe-9",
      "resource": "test-resource-axe-9",
      "output": "test-output-axe-9",
      "long_output": "test-long-output-axe-9",
      "author": "test-author-axe-9",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response changeStateEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-9",
            "connector": "test-connector-axe-9",
            "connector_name": "test-connector-name-axe-9",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-9",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-9",
              "m": "test-output-axe-9",
              "t": {{ .changeStateEventTimestamp }},
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
                "_t": "changestate",
                "a": "test-author-axe-9",
                "m": "test-output-axe-9",
                "t": {{ .changeStateEventTimestamp }},
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-9",
      "connector_name": "test-connector-name-axe-9",
      "source_type": "resource",
      "component":  "test-component-axe-9",
      "resource": "test-resource-axe-9",
      "state": 1,
      "output": "test-output-axe-9",
      "long_output": "test-long-output-axe-9",
      "author": "test-author-axe-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-9",
            "connector": "test-connector-axe-9",
            "connector_name": "test-connector-name-axe-9",
            "resource": "test-resource-axe-9",
            "state": {
              "_t": "changestate",
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
                "_t": "changestate",
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
