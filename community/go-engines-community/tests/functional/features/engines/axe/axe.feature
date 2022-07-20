Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  Scenario: given check event should create alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-axe-1",
      "connector_name" : "test-connector-name-axe-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-1",
      "resource" : "test-resource-axe-1",
      "state" : 2,
      "output" : "test-output-axe-1",
      "long_output" : "test-long-output-axe-1",
      "author" : "test-author-axe-1",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response eventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
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
            "tags": [],
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

  Scenario: given check event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-axe-2",
      "connector_name" : "test-connector-name-axe-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-2",
      "resource" : "test-resource-axe-2",
      "state" : 2,
      "output" : "test-output-axe-2",
      "long_output" : "test-long-output-axe-2-1",
      "author" : "test-author-axe-2",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response firstEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-axe-2",
      "connector_name" : "test-connector-name-axe-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-2",
      "resource" : "test-resource-axe-2",
      "state" : 3,
      "output" : "test-output-axe-2",
      "long_output" : "test-long-output-axe-2-2",
      "author" : "test-author-axe-2",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response secondEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
            "tags": [],
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

  Scenario: given ack event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-3",
      "connector_name" : "test-connector-name-axe-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-3",
      "resource" : "test-resource-axe-3",
      "state" : 2,
      "output" : "test-output-axe-3",
      "long_output" : "test-long-output-axe-3",
      "author" : "test-author-axe-3",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-3",
      "connector_name" : "test-connector-name-axe-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-3",
      "resource" : "test-resource-axe-3",
      "output" : "test-output-axe-3",
      "long_output" : "test-long-output-axe-3",
      "author" : "test-author-axe-3",
      "user_id": "test-author-id-3",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given remove ack event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-4",
      "connector_name" : "test-connector-name-axe-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-4",
      "resource" : "test-resource-axe-4",
      "state" : 2,
      "output" : "test-output-axe-4",
      "long_output" : "test-long-output-axe-4",
      "author" : "test-author-axe-4",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-4",
      "connector_name" : "test-connector-name-axe-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-4",
      "resource" : "test-resource-axe-4",
      "user_id": "test-author-id-4",
      "output" : "test-output-axe-4",
      "long_output" : "test-long-output-axe-4",
      "author" : "test-author-axe-4",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ackremove",
      "connector" : "test-connector-axe-4",
      "connector_name" : "test-connector-name-axe-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-4",
      "resource" : "test-resource-axe-4",
      "user_id": "test-author-id-4",
      "output" : "test-output-axe-4",
      "long_output" : "test-long-output-axe-4",
      "author" : "test-author-axe-4",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackRemoveEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given cancel event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-5",
      "connector_name" : "test-connector-name-axe-5",
      "source_type" : "resource",
      "component" :  "test-component-axe-5",
      "resource" : "test-resource-axe-5",
      "state" : 2,
      "output" : "test-output-axe-5",
      "long_output" : "test-long-output-axe-5",
      "author" : "test-author-axe-5",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-axe-5",
      "connector_name" : "test-connector-name-axe-5",
      "source_type" : "resource",
      "component" :  "test-component-axe-5",
      "resource" : "test-resource-axe-5",
      "output" : "test-output-axe-5",
      "long_output" : "test-long-output-axe-5",
      "author" : "test-author-axe-5",
      "user_id": "test-author-id-5",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given uncancel event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-6",
      "connector_name" : "test-connector-name-axe-6",
      "source_type" : "resource",
      "component" :  "test-component-axe-6",
      "resource" : "test-resource-axe-6",
      "state" : 2,
      "output" : "test-output-axe-6",
      "long_output" : "test-long-output-axe-6",
      "author" : "test-author-axe-6",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-axe-6",
      "connector_name" : "test-connector-name-axe-6",
      "source_type" : "resource",
      "component" :  "test-component-axe-6",
      "resource" : "test-resource-axe-6",
      "output" : "test-output-axe-6",
      "long_output" : "test-long-output-axe-6",
      "author" : "test-author-axe-6",
      "user_id": "test-author-id-6",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "uncancel",
      "connector" : "test-connector-axe-6",
      "connector_name" : "test-connector-name-axe-6",
      "source_type" : "resource",
      "component" :  "test-component-axe-6",
      "resource" : "test-resource-axe-6",
      "output" : "test-output-axe-6",
      "long_output" : "test-long-output-axe-6",
      "author" : "test-author-axe-6",
      "user_id": "test-author-id-6",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response uncancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given comment event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-7",
      "connector_name" : "test-connector-name-axe-7",
      "source_type" : "resource",
      "component" :  "test-component-axe-7",
      "resource" : "test-resource-axe-7",
      "state" : 2,
      "output" : "test-output-axe-7",
      "long_output" : "test-long-output-axe-7",
      "author" : "test-author-axe-7",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "comment",
      "connector" : "test-connector-axe-7",
      "connector_name" : "test-connector-name-axe-7",
      "source_type" : "resource",
      "component" :  "test-component-axe-7",
      "resource" : "test-resource-axe-7",
      "output" : "test-output-axe-7-1",
      "long_output" : "test-long-output-axe-7",
      "author" : "test-author-axe-7",
      "user_id": "test-author-id-7",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response commentEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
    When I send an event:
    """json
    {
      "event_type" : "comment",
      "connector" : "test-connector-axe-7",
      "connector_name" : "test-connector-name-axe-7",
      "source_type" : "resource",
      "component" :  "test-component-axe-7",
      "resource" : "test-resource-axe-7",
      "output" : "test-output-axe-7-2",
      "long_output" : "test-long-output-axe-7",
      "author" : "test-author-axe-7",
      "user_id": "test-author-id-7",
      "timestamp": {{ nowAdd "-2s" }}
    }
    """
    When I save response commentEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given done event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-8",
      "connector_name" : "test-connector-name-axe-8",
      "source_type" : "resource",
      "component" :  "test-component-axe-8",
      "resource" : "test-resource-axe-8",
      "state" : 2,
      "output" : "test-output-axe-8",
      "long_output" : "test-long-output-axe-8",
      "author" : "test-author-axe-8",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "done",
      "connector" : "test-connector-axe-8",
      "connector_name" : "test-connector-name-axe-8",
      "source_type" : "resource",
      "component" :  "test-component-axe-8",
      "resource" : "test-resource-axe-8",
      "output" : "test-output-axe-8",
      "long_output" : "test-long-output-axe-8",
      "author" : "test-author-axe-8",
      "user_id": "test-author-id-8",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response doneEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "done": {
              "_t": "done",
              "a": "test-author-axe-8",
              "user_id": "test-author-id-8",
              "m": "test-output-axe-8",
              "t": {{ .doneEventTimestamp }},
              "val": 0
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
                "_t": "done",
                "a": "test-author-axe-8",
                "user_id": "test-author-id-8",
                "m": "test-output-axe-8",
                "t": {{ .doneEventTimestamp }},
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

  Scenario: given assoc ticket event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-9",
      "connector_name" : "test-connector-name-axe-9",
      "source_type" : "resource",
      "component" :  "test-component-axe-9",
      "resource" : "test-resource-axe-9",
      "state" : 2,
      "output" : "test-output-axe-9",
      "long_output" : "test-long-output-axe-9",
      "author" : "test-author-axe-9",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "assocticket",
      "ticket": "testticket",
      "connector" : "test-connector-axe-9",
      "connector_name" : "test-connector-name-axe-9",
      "source_type" : "resource",
      "component" :  "test-component-axe-9",
      "resource" : "test-resource-axe-9",
      "output" : "test-output-axe-9",
      "long_output" : "test-long-output-axe-9",
      "author" : "test-author-axe-9",
      "user_id": "test-author-id-9",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ticketEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ticket": {
              "_t": "assocticket",
              "a": "test-author-axe-9",
              "m": "testticket",
              "t": {{ .ticketEventTimestamp }},
              "val": "testticket"
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
                "user_id": "test-author-id-9",
                "m": "testticket",
                "t": {{ .ticketEventTimestamp }},
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

  Scenario: given change state event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-10",
      "connector_name" : "test-connector-name-axe-10",
      "source_type" : "resource",
      "component" :  "test-component-axe-10",
      "resource" : "test-resource-axe-10",
      "state" : 2,
      "output" : "test-output-axe-10",
      "long_output" : "test-long-output-axe-10",
      "author" : "test-author-axe-10",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "changestate",
      "state": 3,
      "connector" : "test-connector-axe-10",
      "connector_name" : "test-connector-name-axe-10",
      "source_type" : "resource",
      "component" :  "test-component-axe-10",
      "resource" : "test-resource-axe-10",
      "output" : "test-output-axe-10",
      "long_output" : "test-long-output-axe-10",
      "author" : "test-author-axe-10",
      "user_id": "test-author-id-10",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response changeStateEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-10",
      "connector_name" : "test-connector-name-axe-10",
      "source_type" : "resource",
      "component" :  "test-component-axe-10",
      "resource" : "test-resource-axe-10",
      "state" : 1,
      "output" : "test-output-axe-10",
      "long_output" : "test-long-output-axe-10",
      "author" : "test-author-axe-10"
    }
    """
    When I wait the end of event processing
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

  Scenario: given snooze event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-11",
      "connector_name" : "test-connector-name-axe-11",
      "source_type" : "resource",
      "component" :  "test-component-axe-11",
      "resource" : "test-resource-axe-11",
      "state" : 2,
      "output" : "test-output-axe-11",
      "long_output" : "test-long-output-axe-11",
      "author" : "test-author-axe-11",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 3600,
      "connector" : "test-connector-axe-11",
      "connector_name" : "test-connector-name-axe-11",
      "source_type" : "resource",
      "component" :  "test-component-axe-11",
      "resource" : "test-resource-axe-11",
      "output" : "test-output-axe-11",
      "long_output" : "test-long-output-axe-11",
      "author" : "test-author-axe-11",
      "user_id": "test-author-id-11",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
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

  Scenario: given unsnooze event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-12",
      "connector_name" : "test-connector-name-axe-12",
      "source_type" : "resource",
      "component" :  "test-component-axe-12",
      "resource" : "test-resource-axe-12",
      "state" : 2,
      "output" : "test-output-axe-12",
      "long_output" : "test-long-output-axe-12",
      "author" : "test-author-axe-12",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 3600,
      "connector" : "test-connector-axe-12",
      "connector_name" : "test-connector-name-axe-12",
      "source_type" : "resource",
      "component" :  "test-component-axe-12",
      "resource" : "test-resource-axe-12",
      "output" : "test-output-axe-12",
      "long_output" : "test-long-output-axe-12",
      "author" : "test-author-axe-12",
      "user_id": "test-author-id-12",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "unsnooze",
      "connector" : "test-connector-axe-12",
      "connector_name" : "test-connector-name-axe-12",
      "source_type" : "resource",
      "component" :  "test-component-axe-12",
      "resource" : "test-resource-axe-12",
      "output" : "test-output-axe-12",
      "long_output" : "test-long-output-axe-12",
      "author" : "test-author-axe-12",
      "user_id": "test-author-id-12",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
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

  Scenario: given resolve done event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-13",
      "connector_name" : "test-connector-name-axe-13",
      "source_type" : "resource",
      "component" :  "test-component-axe-13",
      "resource" : "test-resource-axe-13",
      "state" : 2,
      "output" : "test-output-axe-13",
      "long_output" : "test-long-output-axe-13",
      "author" : "test-author-axe-13",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "done",
      "connector" : "test-connector-axe-13",
      "connector_name" : "test-connector-name-axe-13",
      "source_type" : "resource",
      "component" :  "test-component-axe-13",
      "resource" : "test-resource-axe-13",
      "output" : "test-output-axe-13",
      "long_output" : "test-long-output-axe-13",
      "author" : "test-author-axe-13",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response doneEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "resolve_done",
      "connector" : "test-connector-axe-13",
      "connector_name" : "test-connector-name-axe-13",
      "source_type" : "resource",
      "component" :  "test-component-axe-13",
      "resource" : "test-resource-axe-13",
      "output" : "test-output-axe-13",
      "long_output" : "test-long-output-axe-13",
      "author" : "test-author-axe-13",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-13
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-13",
            "connector": "test-connector-axe-13",
            "connector_name": "test-connector-name-axe-13",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-13",
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
                "_t": "done",
                "a": "test-author-axe-13",
                "m": "test-output-axe-13",
                "t": {{ .doneEventTimestamp }},
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

  Scenario: given resolve cancel event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-14",
      "connector_name" : "test-connector-name-axe-14",
      "source_type" : "resource",
      "component" :  "test-component-axe-14",
      "resource" : "test-resource-axe-14",
      "state" : 2,
      "output" : "test-output-axe-14",
      "long_output" : "test-long-output-axe-14",
      "author" : "test-author-axe-14",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-axe-14",
      "connector_name" : "test-connector-name-axe-14",
      "source_type" : "resource",
      "component" :  "test-component-axe-14",
      "resource" : "test-resource-axe-14",
      "output" : "test-output-axe-14",
      "long_output" : "test-long-output-axe-14",
      "author" : "test-author-axe-14",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector" : "test-connector-axe-14",
      "connector_name" : "test-connector-name-axe-14",
      "source_type" : "resource",
      "component" :  "test-component-axe-14",
      "resource" : "test-resource-axe-14",
      "output" : "test-output-axe-14",
      "long_output" : "test-long-output-axe-14",
      "author" : "test-author-axe-14",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I wait the end of event processing
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

  Scenario: given resolve close event should update alarm
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-15",
      "connector_name" : "test-connector-name-axe-15",
      "source_type" : "resource",
      "component" :  "test-component-axe-15",
      "resource" : "test-resource-axe-15",
      "state" : 2,
      "output" : "test-output-axe-15",
      "long_output" : "test-long-output-axe-15",
      "author" : "test-author-axe-15",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "check",
      "state" : 0,
      "connector" : "test-connector-axe-15",
      "connector_name" : "test-connector-name-axe-15",
      "source_type" : "resource",
      "component" :  "test-component-axe-15",
      "resource" : "test-resource-axe-15",
      "output" : "test-output-axe-15",
      "long_output" : "test-long-output-axe-15",
      "author" : "test-author-axe-15",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response closeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "resolve_close",
      "connector" : "test-connector-axe-15",
      "connector_name" : "test-connector-name-axe-15",
      "source_type" : "resource",
      "component" :  "test-component-axe-15",
      "resource" : "test-resource-axe-15",
      "output" : "test-output-axe-15",
      "long_output" : "test-long-output-axe-15",
      "author" : "test-author-axe-15",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I wait the end of event processing
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

  Scenario: given ack resources event should update resource alarms
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-17",
      "connector_name" : "test-connector-name-axe-17",
      "source_type" : "resource",
      "component" :  "test-component-axe-17",
      "resource" : "test-resource-axe-17",
      "state" : 2,
      "output" : "test-output-axe-17",
      "long_output" : "test-long-output-axe-17",
      "author" : "test-author-axe-17",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-17",
      "connector_name" : "test-connector-name-axe-17",
      "source_type" : "component",
      "component" :  "test-component-axe-17",
      "state" : 2,
      "output" : "test-output-axe-17",
      "long_output" : "test-long-output-axe-17",
      "author" : "test-author-axe-17",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-17",
      "connector_name" : "test-connector-name-axe-17",
      "source_type" : "component",
      "component" :  "test-component-axe-17",
      "ack_resources": true,
      "output" : "test-output-axe-17",
      "long_output" : "test-long-output-axe-17",
      "author" : "test-author-axe-17",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response ackEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of 2 events processing
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

  Scenario: given change state event should not update alarm state anymore
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-18",
      "connector_name" : "test-connector-name-axe-18",
      "source_type" : "resource",
      "component" :  "test-component-axe-18",
      "resource" : "test-resource-axe-18",
      "state" : 1,
      "output" : "test-output-axe-18",
      "long_output" : "test-long-output-axe-18",
      "author" : "test-author-axe-18",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "changestate",
      "state": 2,
      "connector" : "test-connector-axe-18",
      "connector_name" : "test-connector-name-axe-18",
      "source_type" : "resource",
      "component" :  "test-component-axe-18",
      "resource" : "test-resource-axe-18",
      "output" : "test-output-axe-18",
      "long_output" : "test-long-output-axe-18",
      "author" : "test-author-axe-18",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-18",
      "connector_name" : "test-connector-name-axe-18",
      "source_type" : "resource",
      "component" :  "test-component-axe-18",
      "resource" : "test-resource-axe-18",
      "state" : 3,
      "output" : "test-output-axe-18",
      "long_output" : "test-long-output-axe-18",
      "author" : "test-author-axe-18"
    }
    """
    When I wait the end of event processing
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

  Scenario: given changestate with same state as already existed one should not update alarm state anymore
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "state" : 2,
      "output" : "test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "changestate",
      "state": 2,
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "output" : "test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
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
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-20",
      "connector_name" : "test-connector-name-axe-20",
      "source_type" : "resource",
      "component" :  "test-component-axe-20",
      "resource" : "test-resource-axe-20",
      "state" : 3,
      "output" : "test-output-axe-20",
      "long_output" : "test-long-output-axe-20",
      "author" : "test-author-axe-20"
    }
    """
    When I wait the end of event processing
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

  Scenario: given change state event should resolve alarm anyway
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-19",
      "connector_name" : "test-connector-name-axe-19",
      "source_type" : "resource",
      "component" :  "test-component-axe-19",
      "resource" : "test-resource-axe-19",
      "state" : 1,
      "output" : "test-output-axe-19",
      "long_output" : "test-long-output-axe-19",
      "author" : "test-author-axe-19",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "changestate",
      "state": 2,
      "connector" : "test-connector-axe-19",
      "connector_name" : "test-connector-name-axe-19",
      "source_type" : "resource",
      "component" :  "test-component-axe-19",
      "resource" : "test-resource-axe-19",
      "output" : "test-output-axe-19",
      "long_output" : "test-long-output-axe-19",
      "author" : "test-author-axe-19",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-19",
      "connector_name" : "test-connector-name-axe-19",
      "source_type" : "resource",
      "component" :  "test-component-axe-19",
      "resource" : "test-resource-axe-19",
      "state" : 0,
      "output" : "test-output-axe-19",
      "long_output" : "test-long-output-axe-19",
      "author" : "test-author-axe-19"
    }
    """
    When I wait the end of event processing
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

  Scenario: given double ack events should update alarm with double ack
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-21",
      "connector_name" : "test-connector-name-axe-21",
      "source_type" : "resource",
      "component" :  "test-component-axe-21",
      "resource" : "test-resource-axe-21",
      "state" : 2,
      "output" : "test-output-axe-21",
      "long_output" : "test-long-output-axe-21",
      "author" : "test-author-axe-21"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-21",
      "connector_name" : "test-connector-name-axe-21",
      "source_type" : "resource",
      "component" :  "test-component-axe-21",
      "resource" : "test-resource-axe-21",
      "output" : "test-output-axe-21",
      "long_output" : "test-long-output-axe-21",
      "author" : "test-author-axe-21",
      "user_id": "test-author-id-21"
    }
    """
    When I wait the end of event processing
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
    When I send an event:
    """json
    {
      "event_type" : "ack",
      "connector" : "test-connector-axe-21",
      "connector_name" : "test-connector-name-axe-21",
      "source_type" : "resource",
      "component" :  "test-component-axe-21",
      "resource" : "test-resource-axe-21",
      "output" : "new-test-output-axe-21",
      "long_output" : "test-long-output-axe-21",
      "author" : "test-author-axe-21",
      "user_id": "test-author-id-21"
    }
    """
    When I wait the end of event processing
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
