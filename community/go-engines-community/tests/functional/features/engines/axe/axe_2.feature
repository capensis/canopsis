Feature: create and update alarm by main event stream
  I need to be able to create and update alarm on event

  @concurrent
  Scenario: given snooze event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-1",
      "connector_name": "test-connector-name-axe-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "state": 2,
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1",
      "author": "test-author-axe-second-1",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 3600,
      "connector": "test-connector-axe-second-1",
      "connector_name": "test-connector-name-axe-second-1",
      "source_type": "resource",
      "component":  "test-component-axe-second-1",
      "resource": "test-resource-axe-second-1",
      "output": "test-output-axe-second-1",
      "long_output": "test-long-output-axe-second-1",
      "author": "test-author-axe-second-1",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I do GET /api/v4/alarms?search=test-resource-axe-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "test-author-axe-second-1",
              "initiator": "external",
              "m": "test-output-axe-second-1",
              "t": {{ .snoozeEventTimestamp }},
              "val": {{ .snoozeEventTimestamp | sumTime 3600 }}
            },
            "component": "test-component-axe-second-1",
            "connector": "test-connector-axe-second-1",
            "last_update_date": {{ .checkEventTimestamp }},
            "connector_name": "test-connector-name-axe-second-1",
            "resource": "test-resource-axe-second-1",
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
                "a": "test-author-axe-second-1",
                "initiator": "external",
                "m": "test-output-axe-second-1",
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
      "connector": "test-connector-axe-second-2",
      "connector_name": "test-connector-name-axe-second-2",
      "source_type": "resource",
      "component":  "test-component-axe-second-2",
      "resource": "test-resource-axe-second-2",
      "state": 2,
      "output": "test-output-axe-second-2",
      "long_output": "test-long-output-axe-second-2",
      "author": "test-author-axe-second-2",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 3600,
      "connector": "test-connector-axe-second-2",
      "connector_name": "test-connector-name-axe-second-2",
      "source_type": "resource",
      "component":  "test-component-axe-second-2",
      "resource": "test-resource-axe-second-2",
      "output": "test-output-axe-second-2",
      "long_output": "test-long-output-axe-second-2",
      "author": "test-author-axe-second-2",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response snoozeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "unsnooze",
      "connector": "test-connector-axe-second-2",
      "connector_name": "test-connector-name-axe-second-2",
      "source_type": "resource",
      "component":  "test-component-axe-second-2",
      "resource": "test-resource-axe-second-2",
      "output": "test-output-axe-second-2",
      "long_output": "test-long-output-axe-second-2",
      "author": "test-author-axe-second-2",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-2",
            "connector": "test-connector-axe-second-2",
            "connector_name": "test-connector-name-axe-second-2",
            "last_update_date": {{ .checkEventTimestamp }},
            "resource": "test-resource-axe-second-2",
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
                "a": "test-author-axe-second-2",
                "initiator": "external",
                "m": "test-output-axe-second-2",
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
      "connector": "test-connector-axe-second-3",
      "connector_name": "test-connector-name-axe-second-3",
      "source_type": "resource",
      "component":  "test-component-axe-second-3",
      "resource": "test-resource-axe-second-3",
      "state": 2,
      "output": "test-output-axe-second-3",
      "long_output": "test-long-output-axe-second-3",
      "author": "test-author-axe-second-3",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "connector": "test-connector-axe-second-3",
      "connector_name": "test-connector-name-axe-second-3",
      "source_type": "resource",
      "component":  "test-component-axe-second-3",
      "resource": "test-resource-axe-second-3",
      "output": "test-output-axe-second-3",
      "long_output": "test-long-output-axe-second-3",
      "author": "test-author-axe-second-3",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response cancelEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "connector": "test-connector-axe-second-3",
      "connector_name": "test-connector-name-axe-second-3",
      "source_type": "resource",
      "component":  "test-component-axe-second-3",
      "resource": "test-resource-axe-second-3",
      "output": "test-output-axe-second-3",
      "long_output": "test-long-output-axe-second-3",
      "author": "test-author-axe-second-3",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-3",
            "connector": "test-connector-axe-second-3",
            "connector_name": "test-connector-name-axe-second-3",
            "last_update_date": {{ .cancelEventTimestamp }},
            "resource": "test-resource-axe-second-3",
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
                "a": "test-author-axe-second-3",
                "initiator": "external",
                "m": "test-output-axe-second-3",
                "t": {{ .cancelEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusinc",
                "a": "test-author-axe-second-3",
                "initiator": "external",
                "m": "test-output-axe-second-3",
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
      "connector": "test-connector-axe-second-4",
      "connector_name": "test-connector-name-axe-second-4",
      "source_type": "resource",
      "component":  "test-component-axe-second-4",
      "resource": "test-resource-axe-second-4",
      "state": 2,
      "output": "test-output-axe-second-4",
      "long_output": "test-long-output-axe-second-4",
      "author": "test-author-axe-second-4",
      "timestamp": {{ nowAdd "-10s" }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 0,
      "connector": "test-connector-axe-second-4",
      "connector_name": "test-connector-name-axe-second-4",
      "source_type": "resource",
      "component":  "test-component-axe-second-4",
      "resource": "test-resource-axe-second-4",
      "output": "test-output-axe-second-4",
      "long_output": "test-long-output-axe-second-4",
      "author": "test-author-axe-second-4",
      "timestamp": {{ nowAdd "-7s" }}
    }
    """
    When I save response closeEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_close",
      "connector": "test-connector-axe-second-4",
      "connector_name": "test-connector-name-axe-second-4",
      "source_type": "resource",
      "component":  "test-component-axe-second-4",
      "resource": "test-resource-axe-second-4",
      "output": "test-output-axe-second-4",
      "long_output": "test-long-output-axe-second-4",
      "author": "test-author-axe-second-4",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I save response resolveTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-axe-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-4",
            "connector": "test-connector-axe-second-4",
            "connector_name": "test-connector-name-axe-second-4",
            "last_update_date": {{ .closeEventTimestamp }},
            "resource": "test-resource-axe-second-4",
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
                "a": "test-connector-axe-second-4.test-connector-name-axe-second-4",
                "initiator": "external",
                "m": "test-output-axe-second-4",
                "t": {{ .closeEventTimestamp }},
                "val": 0
              },
              {
                "_t": "statusdec",
                "a": "test-connector-axe-second-4.test-connector-name-axe-second-4",
                "initiator": "external",
                "m": "test-output-axe-second-4",
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
  Scenario: given change state event should not update alarm state anymore
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-5",
      "connector_name": "test-connector-name-axe-second-5",
      "source_type": "resource",
      "component":  "test-component-axe-second-5",
      "resource": "test-resource-axe-second-5",
      "state": 1,
      "output": "test-output-axe-second-5",
      "long_output": "test-long-output-axe-second-5",
      "author": "test-author-axe-second-5",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-second-5",
      "connector_name": "test-connector-name-axe-second-5",
      "source_type": "resource",
      "component":  "test-component-axe-second-5",
      "resource": "test-resource-axe-second-5",
      "output": "test-output-axe-second-5",
      "long_output": "test-long-output-axe-second-5",
      "author": "test-author-axe-second-5",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-5",
      "connector_name": "test-connector-name-axe-second-5",
      "source_type": "resource",
      "component":  "test-component-axe-second-5",
      "resource": "test-resource-axe-second-5",
      "state": 3,
      "output": "test-output-axe-second-5",
      "long_output": "test-long-output-axe-second-5",
      "author": "test-author-axe-second-5"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-5",
            "connector": "test-connector-axe-second-5",
            "connector_name": "test-connector-name-axe-second-5",
            "resource": "test-resource-axe-second-5",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-second-5",
              "initiator": "external",
              "m": "test-output-axe-second-5",
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
                "a": "test-author-axe-second-5",
                "initiator": "external",
                "m": "test-output-axe-second-5",
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
      "connector": "test-connector-axe-second-6",
      "connector_name": "test-connector-name-axe-second-6",
      "source_type": "resource",
      "component":  "test-component-axe-second-6",
      "resource": "test-resource-axe-second-6",
      "state": 2,
      "output": "test-output-axe-second-6",
      "long_output": "test-long-output-axe-second-6",
      "author": "test-author-axe-second-6",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-second-6",
      "connector_name": "test-connector-name-axe-second-6",
      "source_type": "resource",
      "component":  "test-component-axe-second-6",
      "resource": "test-resource-axe-second-6",
      "output": "test-output-axe-second-6",
      "long_output": "test-long-output-axe-second-6",
      "author": "test-author-axe-second-6",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-6",
            "connector": "test-connector-axe-second-6",
            "connector_name": "test-connector-name-axe-second-6",
            "resource": "test-resource-axe-second-6",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-second-6",
              "initiator": "external",
              "m": "test-output-axe-second-6",
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
                "a": "test-author-axe-second-6",
                "initiator": "external",
                "m": "test-output-axe-second-6",
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
      "connector": "test-connector-axe-second-6",
      "connector_name": "test-connector-name-axe-second-6",
      "source_type": "resource",
      "component":  "test-component-axe-second-6",
      "resource": "test-resource-axe-second-6",
      "state": 3,
      "output": "test-output-axe-second-6",
      "long_output": "test-long-output-axe-second-6",
      "author": "test-author-axe-second-6"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-6",
            "connector": "test-connector-axe-second-6",
            "connector_name": "test-connector-name-axe-second-6",
            "resource": "test-resource-axe-second-6",
            "state": {
              "_t": "changestate",
              "a": "test-author-axe-second-6",
              "initiator": "external",
              "m": "test-output-axe-second-6",
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
                "a": "test-author-axe-second-6",
                "initiator": "external",
                "m": "test-output-axe-second-6",
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
      "connector": "test-connector-axe-second-7",
      "connector_name": "test-connector-name-axe-second-7",
      "source_type": "resource",
      "component":  "test-component-axe-second-7",
      "resource": "test-resource-axe-second-7",
      "state": 1,
      "output": "test-output-axe-second-7",
      "long_output": "test-long-output-axe-second-7",
      "author": "test-author-axe-second-7",
      "timestamp": {{ nowAdd "-19s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "changestate",
      "state": 2,
      "connector": "test-connector-axe-second-7",
      "connector_name": "test-connector-name-axe-second-7",
      "source_type": "resource",
      "component":  "test-component-axe-second-7",
      "resource": "test-resource-axe-second-7",
      "output": "test-output-axe-second-7",
      "long_output": "test-long-output-axe-second-7",
      "author": "test-author-axe-second-7",
      "timestamp": {{ nowAdd "-5s" }}
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-7",
      "connector_name": "test-connector-name-axe-second-7",
      "source_type": "resource",
      "component":  "test-component-axe-second-7",
      "resource": "test-resource-axe-second-7",
      "state": 0,
      "output": "test-output-axe-second-7",
      "long_output": "test-long-output-axe-second-7",
      "author": "test-author-axe-second-7"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-7",
            "connector": "test-connector-axe-second-7",
            "connector_name": "test-connector-name-axe-second-7",
            "resource": "test-resource-axe-second-7",
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
                "a": "test-author-axe-second-7",
                "initiator": "external",
                "m": "test-output-axe-second-7",
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
      "connector": "test-connector-axe-second-8",
      "connector_name": "test-connector-name-axe-second-8",
      "source_type": "resource",
      "component":  "test-component-axe-second-8",
      "resource": "test-resource-axe-second-8",
      "state": 2,
      "output": "test-output-axe-second-8",
      "long_output": "test-long-output-axe-second-8",
      "author": "test-author-axe-second-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "ack",
      "connector": "test-connector-axe-second-8",
      "connector_name": "test-connector-name-axe-second-8",
      "source_type": "resource",
      "component":  "test-component-axe-second-8",
      "resource": "test-resource-axe-second-8",
      "output": "test-output-axe-second-8",
      "long_output": "test-long-output-axe-second-8",
      "author": "test-author-axe-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-second-8",
              "initiator": "external",
              "m": "test-output-axe-second-8",
              "val": 0
            },
            "component": "test-component-axe-second-8",
            "connector": "test-connector-axe-second-8",
            "connector_name": "test-connector-name-axe-second-8",
            "resource": "test-resource-axe-second-8",
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
                "a": "test-author-axe-second-8",
                "initiator": "external",
                "m": "test-output-axe-second-8",
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
      "connector": "test-connector-axe-second-8",
      "connector_name": "test-connector-name-axe-second-8",
      "source_type": "resource",
      "component":  "test-component-axe-second-8",
      "resource": "test-resource-axe-second-8",
      "output": "new-test-output-axe-second-8",
      "long_output": "test-long-output-axe-second-8",
      "author": "test-author-axe-second-8"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "test-author-axe-second-8",
              "initiator": "external",
              "m": "new-test-output-axe-second-8",
              "val": 0
            },
            "component": "test-component-axe-second-8",
            "connector": "test-connector-axe-second-8",
            "connector_name": "test-connector-name-axe-second-8",
            "resource": "test-resource-axe-second-8",
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
                "a": "test-author-axe-second-8",
                "initiator": "external",
                "m": "test-output-axe-second-8",
                "val": 0
              },
              {
                "_t": "ack",
                "a": "test-author-axe-second-8",
                "initiator": "external",
                "m": "new-test-output-axe-second-8",
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
  Scenario: given events with different connector for the same resource should use right connectors in alarm steps
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-9",
      "connector_name": "test-connector-name-axe-second-9",
      "source_type": "resource",
      "component":  "test-component-axe-second-9",
      "resource": "test-resource-axe-second-9",
      "state": 1,
      "output": "test-output-axe-second-9",
      "long_output": "test-long-output-axe-second-9"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-9",
      "connector_name": "test-connector-name-axe-second-9-new",
      "source_type": "resource",
      "component":  "test-component-axe-second-9",
      "resource": "test-resource-axe-second-9",
      "state": 2,
      "output": "test-output-axe-second-9",
      "long_output": "test-long-output-axe-second-9"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-9-new",
      "connector_name": "test-connector-name-axe-second-9",
      "source_type": "resource",
      "component":  "test-component-axe-second-9",
      "resource": "test-resource-axe-second-9",
      "state": 3,
      "output": "test-output-axe-second-9",
      "long_output": "test-long-output-axe-second-9"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "connector": "test-connector-axe-second-9-new",
      "connector_name": "test-connector-name-axe-second-9-new",
      "source_type": "resource",
      "component":  "test-component-axe-second-9",
      "resource": "test-resource-axe-second-9",
      "state": 0,
      "output": "test-output-axe-second-9",
      "long_output": "test-long-output-axe-second-9"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-second-9",
            "connector": "test-connector-axe-second-9",
            "connector_name": "test-connector-name-axe-second-9",
            "resource": "test-resource-axe-second-9"
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
                "val": 1,
                "a": "test-connector-axe-second-9.test-connector-name-axe-second-9",
                "initiator": "external"
              },
              {
                "_t": "statusinc",
                "val": 1,
                "a": "test-connector-axe-second-9.test-connector-name-axe-second-9",
                "initiator": "external"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "a": "test-connector-axe-second-9.test-connector-name-axe-second-9-new",
                "initiator": "external"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "a": "test-connector-axe-second-9-new.test-connector-name-axe-second-9",
                "initiator": "external"
              },
              {
                "_t": "statedec",
                "val": 0,
                "a": "test-connector-axe-second-9-new.test-connector-name-axe-second-9-new",
                "initiator": "external"
              },
              {
                "_t": "statusdec",
                "val": 0,
                "a": "test-connector-axe-second-9-new.test-connector-name-axe-second-9-new",
                "initiator": "external"
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
