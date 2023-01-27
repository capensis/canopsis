Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given check event should create alarm
    Given I am admin
    When I save response createTimestamp={{ now }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-first",
      "connector_name": "test-connector-name-axe-first",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-first",
      "resource": "test-resource-axe-first",
      "state": 2,
      "output": "test-output-axe-first",
      "long_output": "test-long-output-axe-first",
      "author": "test-author-axe-first",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response eventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-first
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-first/test-component-axe-first"
          },
          "infos": {},
          "tags": [],
          "v": {
            "children": [],
            "component": "test-component-axe-first",
            "connector": "test-connector-axe-first",
            "connector_name": "test-connector-name-axe-first",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "test-long-output-axe-first",
            "initial_output": "test-output-axe-first",
            "last_update_date": {{ .eventTimestamp }},
            "long_output": "test-long-output-axe-first",
            "long_output_history": ["test-long-output-axe-first"],
            "output": "test-output-axe-first",
            "parents": [],
            "resource": "test-resource-axe-first",
            "state": {
              "_t": "stateinc",
              "a": "test-connector-axe-first.test-connector-name-axe-first",
              "m": "test-output-axe-first",
              "t": {{ .eventTimestamp }},
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-axe-first.test-connector-name-axe-first",
              "m": "test-output-axe-first",
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
                "a": "test-connector-axe-first.test-connector-name-axe-first",
                "m": "test-output-axe-first",
                "t": {{ .eventTimestamp }},
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-first.test-connector-name-axe-first",
                "m": "test-output-axe-first",
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
      "connector": "test-connector-axe-second",
      "connector_name": "test-connector-name-axe-second",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-second",
      "resource": "test-resource-axe-second",
      "state": 2,
      "output": "test-output-axe-second",
      "long_output": "test-long-output-axe-second-1",
      "author": "test-author-axe-second",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response firstEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-second",
      "connector_name": "test-connector-name-axe-second",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-axe-second",
      "resource": "test-resource-axe-second",
      "state": 3,
      "output": "test-output-axe-second",
      "long_output": "test-long-output-axe-second-2",
      "author": "test-author-axe-second",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response secondEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-second
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "entity": {
            "_id": "test-resource-axe-second/test-component-axe-second"
          },
          "infos": {},
          "v": {
            "children": [],
            "component": "test-component-axe-second",
            "connector": "test-connector-axe-second",
            "connector_name": "test-connector-name-axe-second",
            "infos": {},
            "infos_rule_version": {},
            "initial_long_output": "test-long-output-axe-second-1",
            "initial_output": "test-output-axe-second",
            "last_update_date": {{ .secondEventTimestamp }},
            "long_output": "test-long-output-axe-second-2",
            "long_output_history": ["test-long-output-axe-second-1", "test-long-output-axe-second-2"],
            "output": "test-output-axe-second",
            "parents": [],
            "resource": "test-resource-axe-second",
            "state": {
              "_t": "stateinc",
              "a": "test-connector-axe-second.test-connector-name-axe-second",
              "m": "test-output-axe-second",
              "t": {{ .secondEventTimestamp }},
              "val": 3
            },
            "status": {
              "_t": "statusinc",
              "a": "test-connector-axe-second.test-connector-name-axe-second",
              "m": "test-output-axe-second",
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
                "a": "test-connector-axe-second.test-connector-name-axe-second",
                "m": "test-output-axe-second",
                "t": {{ .firstEventTimestamp }},
                "val": 2
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-second.test-connector-name-axe-second",
                "m": "test-output-axe-second",
                "t": {{ .firstEventTimestamp }},
                "val": 1
              },
              {
                "_t": "stateinc",
                "a": "test-connector-axe-second.test-connector-name-axe-second",
                "m": "test-output-axe-second",
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
      "user_id": "test-author-id-3",
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
              "user_id": "test-author-id-3",
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
      "user_id": "test-author-id-4",
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
      "user_id": "test-author-id-4",
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
                "user_id": "test-author-id-4",
                "m": "test-output-axe-4",
                "t": {{ .ackEventTimestamp }},
                "val": 0
              },
              {
                "_t": "ackremove",
                "a": "test-author-axe-4",
                "user_id": "test-author-id-4",
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
      "user_id": "test-author-id-5",
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
              "user_id": "test-author-id-5",
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
              "a": "test-connector-axe-5.test-connector-name-axe-5",
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
                "a": "test-connector-axe-5.test-connector-name-axe-5",
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
      "user_id": "test-author-id-6",
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
      "user_id": "test-author-id-6",
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
              "a": "test-connector-axe-6.test-connector-name-axe-6",
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
                "user_id": "test-author-id-6",
                "m": "test-output-axe-6",
                "t": {{ .cancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-6.test-connector-name-axe-6",
                "m": "test-output-axe-6",
                "t": {{ .cancelEventTimestamp }},
                "val": 4
              },
              {
                "_t": "uncancel",
                "a": "test-author-axe-6",
                "user_id": "test-author-id-6",
                "m": "test-output-axe-6",
                "t": {{ .uncancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-axe-6.test-connector-name-axe-6",
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
      "user_id": "test-author-id-7",
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
              "user_id": "test-author-id-7",
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
      "user_id": "test-author-id-7",
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
              "user_id": "test-author-id-7",
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
      "event_type": "assocticket",
      "ticket": "test-ticket",
      "ticket_system_name": "test-system-name",
      "ticket_url": "test-ticket-url",
      "ticket_data": {
        "ticket_param_1": "ticket_value_1"
      },
      "ticket_comment": "test-ticket-comment",
      "connector": "test-connector-axe-9",
      "connector_name": "test-connector-name-axe-9",
      "source_type": "resource",
      "component":  "test-component-axe-9",
      "resource": "test-resource-axe-9",
      "author": "test-author-axe-9",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ticketEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-9
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
                "a": "test-author-axe-9",
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
              "a": "test-author-axe-9",
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
            "component": "test-component-axe-9",
            "connector": "test-connector-axe-9",
            "connector_name": "test-connector-name-axe-9",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-9",
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
                "a": "test-author-axe-9",
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
      "connector": "test-connector-axe-10",
      "connector_name": "test-connector-name-axe-10",
      "source_type": "resource",
      "component":  "test-component-axe-10",
      "resource": "test-resource-axe-10",
      "state": 2,
      "output": "test-output-axe-10",
      "long_output": "test-long-output-axe-10",
      "author": "test-author-axe-10",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 3,
      "connector": "test-connector-axe-10",
      "connector_name": "test-connector-name-axe-10",
      "source_type": "resource",
      "component":  "test-component-axe-10",
      "resource": "test-resource-axe-10",
      "output": "test-output-axe-10",
      "long_output": "test-long-output-axe-10",
      "author": "test-author-axe-10",
      "user_id": "test-author-id-10",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response changeStateEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-10
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-10",
            "connector": "test-connector-axe-10",
            "connector_name": "test-connector-name-axe-10",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-10",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-10",
              "user_id": "test-author-id-10",
              "m": "test-output-axe-10",
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
                "a": "test-author-axe-10",
                "user_id": "test-author-id-10",
                "m": "test-output-axe-10",
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
      "connector": "test-connector-axe-10",
      "connector_name": "test-connector-name-axe-10",
      "source_type": "resource",
      "component":  "test-component-axe-10",
      "resource": "test-resource-axe-10",
      "state": 1,
      "output": "test-output-axe-10",
      "long_output": "test-long-output-axe-10",
      "author": "test-author-axe-10"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-10
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-10",
            "connector": "test-connector-axe-10",
            "connector_name": "test-connector-name-axe-10",
            "resource": "test-resource-axe-10",
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

  @concurrent
  Scenario: given snooze event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-11",
      "connector_name": "test-connector-name-axe-11",
      "source_type": "resource",
      "component":  "test-component-axe-11",
      "resource": "test-resource-axe-11",
      "state": 2,
      "output": "test-output-axe-11",
      "long_output": "test-long-output-axe-11",
      "author": "test-author-axe-11",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 3600,
      "connector": "test-connector-axe-11",
      "connector_name": "test-connector-name-axe-11",
      "source_type": "resource",
      "component":  "test-component-axe-11",
      "resource": "test-resource-axe-11",
      "output": "test-output-axe-11",
      "long_output": "test-long-output-axe-11",
      "author": "test-author-axe-11",
      "user_id": "test-author-id-11",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-11
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "test-author-axe-11",
              "user_id": "test-author-id-11",
              "m": "test-output-axe-11",
              "t": {{ .snoozeEventTimestamp }},
              "val": {{ .snoozeEventTimestamp | sumTime 3600 }}
            },
            "component": "test-component-axe-11",
            "connector": "test-connector-axe-11",
            "last_update_date": {{ .checkEventTimestamp }},
            "connector_name": "test-connector-name-axe-11",
            "resource": "test-resource-axe-11",
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
                "_t": "snooze",
                "a": "test-author-axe-11",
                "user_id": "test-author-id-11",
                "m": "test-output-axe-11",
                "t": {{ .snoozeEventTimestamp }},
                "val": {{ .snoozeEventTimestamp | sumTime 3600 }}
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
  Scenario: given unsnooze event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-12",
      "connector_name": "test-connector-name-axe-12",
      "source_type": "resource",
      "component":  "test-component-axe-12",
      "resource": "test-resource-axe-12",
      "state": 2,
      "output": "test-output-axe-12",
      "long_output": "test-long-output-axe-12",
      "author": "test-author-axe-12",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 3600,
      "connector": "test-connector-axe-12",
      "connector_name": "test-connector-name-axe-12",
      "source_type": "resource",
      "component":  "test-component-axe-12",
      "resource": "test-resource-axe-12",
      "output": "test-output-axe-12",
      "long_output": "test-long-output-axe-12",
      "author": "test-author-axe-12",
      "user_id": "test-author-id-12",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "unsnooze",
      "connector": "test-connector-axe-12",
      "connector_name": "test-connector-name-axe-12",
      "source_type": "resource",
      "component":  "test-component-axe-12",
      "resource": "test-resource-axe-12",
      "output": "test-output-axe-12",
      "long_output": "test-long-output-axe-12",
      "author": "test-author-axe-12",
      "user_id": "test-author-id-12",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-12
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-12",
            "connector": "test-connector-axe-12",
            "connector_name": "test-connector-name-axe-12",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-12",
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
    Then the response key "data.0.v.snooze" should not exist
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
                "_t": "snooze",
                "a": "test-author-axe-12",
                "user_id": "test-author-id-12",
                "m": "test-output-axe-12",
                "t": {{ .snoozeEventTimestamp }},
                "val": {{ .snoozeEventTimestamp | sumTime 3600 }}
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
  Scenario: given resolve cancel event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-14",
      "connector_name": "test-connector-name-axe-14",
      "source_type": "resource",
      "component":  "test-component-axe-14",
      "resource": "test-resource-axe-14",
      "state": 2,
      "output": "test-output-axe-14",
      "long_output": "test-long-output-axe-14",
      "author": "test-author-axe-14",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-14",
      "connector_name": "test-connector-name-axe-14",
      "source_type": "resource",
      "component":  "test-component-axe-14",
      "resource": "test-resource-axe-14",
      "output": "test-output-axe-14",
      "long_output": "test-long-output-axe-14",
      "author": "test-author-axe-14",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "connector": "test-connector-axe-14",
      "connector_name": "test-connector-name-axe-14",
      "source_type": "resource",
      "component":  "test-component-axe-14",
      "resource": "test-resource-axe-14",
      "output": "test-output-axe-14",
      "long_output": "test-long-output-axe-14",
      "author": "test-author-axe-14",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-14
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-14",
            "connector": "test-connector-axe-14",
            "connector_name": "test-connector-name-axe-14",
            "last_update_date": {{ .cancelEventTimestamp }},
            "resource": "test-resource-axe-14",
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
    When I save response alarmResolve={{ (index .lastResponse.data 0).v.resolved }}
    Then the difference between alarmResolve resolveTimestamp is in range -2,2
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
                "a": "test-author-axe-14",
                "m": "test-output-axe-14",
                "t": {{ .cancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-connector-axe-14.test-connector-name-axe-14",
                "m": "test-output-axe-14",
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
  Scenario: given resolve close event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-15",
      "connector_name": "test-connector-name-axe-15",
      "source_type": "resource",
      "component":  "test-component-axe-15",
      "resource": "test-resource-axe-15",
      "state": 2,
      "output": "test-output-axe-15",
      "long_output": "test-long-output-axe-15",
      "author": "test-author-axe-15",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "connector": "test-connector-axe-15",
      "connector_name": "test-connector-name-axe-15",
      "source_type": "resource",
      "component":  "test-component-axe-15",
      "resource": "test-resource-axe-15",
      "output": "test-output-axe-15",
      "long_output": "test-long-output-axe-15",
      "author": "test-author-axe-15",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response closeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_close",
      "connector": "test-connector-axe-15",
      "connector_name": "test-connector-name-axe-15",
      "source_type": "resource",
      "component":  "test-component-axe-15",
      "resource": "test-resource-axe-15",
      "output": "test-output-axe-15",
      "long_output": "test-long-output-axe-15",
      "author": "test-author-axe-15",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-15
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-15",
            "connector": "test-connector-axe-15",
            "connector_name": "test-connector-name-axe-15",
            "last_update_date": {{ .closeEventTimestamp }},
            "resource": "test-resource-axe-15",
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
    When I save response alarmResolve={{ (index .lastResponse.data 0).v.resolved }}
    Then the difference between alarmResolve resolveTimestamp is in range -2,2
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
                "a": "test-connector-axe-15.test-connector-name-axe-15",
                "m": "test-output-axe-15",
                "t": {{ .closeEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-axe-15.test-connector-name-axe-15",
                "m": "test-output-axe-15",
                "t": {{ .closeEventTimestamp }},
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
  Scenario: given ack resources event should update resource alarms
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-17",
      "connector_name": "test-connector-name-axe-17",
      "source_type": "resource",
      "component":  "test-component-axe-17",
      "resource": "test-resource-axe-17",
      "state": 2,
      "output": "test-output-axe-17",
      "long_output": "test-long-output-axe-17",
      "author": "test-author-axe-17",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-17",
      "connector_name": "test-connector-name-axe-17",
      "source_type": "component",
      "component":  "test-component-axe-17",
      "state": 2,
      "output": "test-output-axe-17",
      "long_output": "test-long-output-axe-17",
      "author": "test-author-axe-17",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-17",
      "connector_name": "test-connector-name-axe-17",
      "source_type": "component",
      "component":  "test-component-axe-17",
      "ack_resources": true,
      "output": "test-output-axe-17",
      "long_output": "test-long-output-axe-17",
      "author": "test-author-axe-17",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-17",
      "connector_name": "test-connector-name-axe-17",
      "source_type": "resource",
      "component":  "test-component-axe-17",
      "resource":  "test-resource-axe-17"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-17
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-17",
              "m": "test-output-axe-17",
              "t": {{ .ackEventTimestamp }},
              "val": 0
            },
            "component": "test-component-axe-17",
            "connector": "test-connector-axe-17",
            "connector_name": "test-connector-name-axe-17",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-17",
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
                "a": "test-author-axe-17",
                "m": "test-output-axe-17",
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
  Scenario: given change state event should not update alarm state anymore
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-18",
      "connector_name": "test-connector-name-axe-18",
      "source_type": "resource",
      "component":  "test-component-axe-18",
      "resource": "test-resource-axe-18",
      "state": 1,
      "output": "test-output-axe-18",
      "long_output": "test-long-output-axe-18",
      "author": "test-author-axe-18",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-18",
      "connector_name": "test-connector-name-axe-18",
      "source_type": "resource",
      "component":  "test-component-axe-18",
      "resource": "test-resource-axe-18",
      "output": "test-output-axe-18",
      "long_output": "test-long-output-axe-18",
      "author": "test-author-axe-18",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-18",
      "connector_name": "test-connector-name-axe-18",
      "source_type": "resource",
      "component":  "test-component-axe-18",
      "resource": "test-resource-axe-18",
      "state": 3,
      "output": "test-output-axe-18",
      "long_output": "test-long-output-axe-18",
      "author": "test-author-axe-18"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-18
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-18",
            "connector": "test-connector-axe-18",
            "connector_name": "test-connector-name-axe-18",
            "resource": "test-resource-axe-18",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-18",
              "m": "test-output-axe-18",
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
                "_t": "changestate",
                "a": "test-author-axe-18",
                "m": "test-output-axe-18",
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
  Scenario: given changestate with same state as already existed one should not update alarm state anymore
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-20",
      "connector_name": "test-connector-name-axe-20",
      "source_type": "resource",
      "component":  "test-component-axe-20",
      "resource": "test-resource-axe-20",
      "state": 2,
      "output": "test-output-axe-20",
      "long_output": "test-long-output-axe-20",
      "author": "test-author-axe-20",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-20",
      "connector_name": "test-connector-name-axe-20",
      "source_type": "resource",
      "component":  "test-component-axe-20",
      "resource": "test-resource-axe-20",
      "output": "test-output-axe-20",
      "long_output": "test-long-output-axe-20",
      "author": "test-author-axe-20",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-20
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-20",
            "connector": "test-connector-axe-20",
            "connector_name": "test-connector-name-axe-20",
            "resource": "test-resource-axe-20",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-20",
              "m": "test-output-axe-20",
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
                "_t": "changestate",
                "a": "test-author-axe-20",
                "m": "test-output-axe-20",
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
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-20",
      "connector_name": "test-connector-name-axe-20",
      "source_type": "resource",
      "component":  "test-component-axe-20",
      "resource": "test-resource-axe-20",
      "state": 3,
      "output": "test-output-axe-20",
      "long_output": "test-long-output-axe-20",
      "author": "test-author-axe-20"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-20
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-20",
            "connector": "test-connector-axe-20",
            "connector_name": "test-connector-name-axe-20",
            "resource": "test-resource-axe-20",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-20",
              "m": "test-output-axe-20",
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
                "_t": "changestate",
                "a": "test-author-axe-20",
                "m": "test-output-axe-20",
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
  Scenario: given change state event should resolve alarm anyway
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-19",
      "connector_name": "test-connector-name-axe-19",
      "source_type": "resource",
      "component":  "test-component-axe-19",
      "resource": "test-resource-axe-19",
      "state": 1,
      "output": "test-output-axe-19",
      "long_output": "test-long-output-axe-19",
      "author": "test-author-axe-19",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-19",
      "connector_name": "test-connector-name-axe-19",
      "source_type": "resource",
      "component":  "test-component-axe-19",
      "resource": "test-resource-axe-19",
      "output": "test-output-axe-19",
      "long_output": "test-long-output-axe-19",
      "author": "test-author-axe-19",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-19",
      "connector_name": "test-connector-name-axe-19",
      "source_type": "resource",
      "component":  "test-component-axe-19",
      "resource": "test-resource-axe-19",
      "state": 0,
      "output": "test-output-axe-19",
      "long_output": "test-long-output-axe-19",
      "author": "test-author-axe-19"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-19
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-19",
            "connector": "test-connector-axe-19",
            "connector_name": "test-connector-name-axe-19",
            "resource": "test-resource-axe-19",
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
                "_t": "changestate",
                "a": "test-author-axe-19",
                "m": "test-output-axe-19",
                "val": 2
              },
              {
                "_t": "statedec",
                "val": 0
              },
              {
                "_t": "statusdec",
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

  @concurrent
  Scenario: given double ack events should update alarm with double ack
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-21",
      "connector_name": "test-connector-name-axe-21",
      "source_type": "resource",
      "component":  "test-component-axe-21",
      "resource": "test-resource-axe-21",
      "state": 2,
      "output": "test-output-axe-21",
      "long_output": "test-long-output-axe-21",
      "author": "test-author-axe-21"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-21",
      "connector_name": "test-connector-name-axe-21",
      "source_type": "resource",
      "component":  "test-component-axe-21",
      "resource": "test-resource-axe-21",
      "output": "test-output-axe-21",
      "long_output": "test-long-output-axe-21",
      "author": "test-author-axe-21",
      "user_id": "test-author-id-21"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-21
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-21",
              "m": "test-output-axe-21",
              "user_id": "test-author-id-21",
              "val": 0
            },
            "component": "test-component-axe-21",
            "connector": "test-connector-axe-21",
            "connector_name": "test-connector-name-axe-21",
            "resource": "test-resource-axe-21",
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
                "a": "test-author-axe-21",
                "m": "test-output-axe-21",
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
      "event_type": "ack",
      "connector": "test-connector-axe-21",
      "connector_name": "test-connector-name-axe-21",
      "source_type": "resource",
      "component":  "test-component-axe-21",
      "resource": "test-resource-axe-21",
      "output": "new-test-output-axe-21",
      "long_output": "test-long-output-axe-21",
      "author": "test-author-axe-21",
      "user_id": "test-author-id-21"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-21
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-21",
              "m": "new-test-output-axe-21",
              "user_id": "test-author-id-21",
              "val": 0
            },
            "component": "test-component-axe-21",
            "connector": "test-connector-axe-21",
            "connector_name": "test-connector-name-axe-21",
            "resource": "test-resource-axe-21",
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
                "a": "test-author-axe-21",
                "m": "test-output-axe-21",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-21",
                "m": "new-test-output-axe-21",
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
  Scenario: given ticket resources event should update resource alarms
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-22",
      "long_output": "test-long-output-axe-22",
      "author": "test-author-axe-22",
      "timestamp": {{ nowAdd "-10s" }},
      "connector": "test-connector-axe-22",
      "connector_name": "test-connector-name-axe-22",
      "component":  "test-component-axe-22",
      "resource": "test-resource-axe-22",
      "source_type": "resource"
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-22",
      "long_output": "test-long-output-axe-22",
      "author": "test-author-axe-22",
      "timestamp": {{ nowAdd "-10s" }},
      "connector": "test-connector-axe-22",
      "connector_name": "test-connector-name-axe-22",
      "component":  "test-component-axe-22",
      "source_type": "component"
    }
    """
    When I send an event:
    """json
    {
      "event_type": "assocticket",
      "ticket_resources": true,
      "ticket": "test-ticket-axe-22",
      "ticket_system_name": "test-system-name-axe-22",
      "ticket_url": "test-ticket-url-axe-22",
      "ticket_data": {
        "param1": "test-value-param-1-axe-22"
      },
      "ticket_comment": "test-ticket-comment-axe-22",
      "author": "test-author-axe-22",
      "timestamp": {{ nowAdd "-5s" }},
      "connector": "test-connector-axe-22",
      "connector_name": "test-connector-name-axe-22",
      "component":  "test-component-axe-22",
      "source_type": "component"
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-22",
        "connector_name": "test-connector-name-axe-22",
        "component":  "test-component-axe-22",
        "source_type": "component"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-22",
        "connector_name": "test-connector-name-axe-22",
        "component":  "test-component-axe-22",
        "resource":  "test-resource-axe-22",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-22
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "test-author-axe-22",
              "m": "Ticket ID: test-ticket-axe-22. Ticket URL: test-ticket-url-axe-22. Ticket param1: test-value-param-1-axe-22.",
              "ticket": "test-ticket-axe-22",
              "ticket_system_name": "test-system-name-axe-22",
              "ticket_url": "test-ticket-url-axe-22",
              "ticket_data": {
                "param1": "test-value-param-1-axe-22"
              },
              "ticket_comment": "test-ticket-comment-axe-22",
              "t": {{ .ackEventTimestamp }},
              "val": 0
            },
            "tickets": [
              {
                "_t": "assocticket",
                "a": "test-author-axe-22",
                "m": "Ticket ID: test-ticket-axe-22. Ticket URL: test-ticket-url-axe-22. Ticket param1: test-value-param-1-axe-22.",
                "ticket": "test-ticket-axe-22",
                "ticket_system_name": "test-system-name-axe-22",
                "ticket_url": "test-ticket-url-axe-22",
                "ticket_data": {
                  "param1": "test-value-param-1-axe-22"
                },
                "ticket_comment": "test-ticket-comment-axe-22",
                "t": {{ .ackEventTimestamp }},
                "val": 0
              }
            ],
            "component": "test-component-axe-22",
            "connector": "test-connector-axe-22",
            "connector_name": "test-connector-name-axe-22",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-22",
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
                "a": "test-author-axe-22",
                "m": "Ticket ID: test-ticket-axe-22. Ticket URL: test-ticket-url-axe-22. Ticket param1: test-value-param-1-axe-22.",
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
