Feature: update alarm
  I need to be able to update alarm

  @concurrent
  Scenario: given ack should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-1",
      "connector": "test-connector-axe-api-1",
      "connector_name": "test-connector-name-axe-api-1",
      "component": "test-component-axe-api-1",
      "resource": "test-resource-axe-api-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-1
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/ack:
    """json
    {
      "comment": "test-comment-axe-api-1"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "ack",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-1",
      "resource": "test-resource-axe-api-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-1
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
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-1"
            },
            "connector": "test-connector-axe-api-1",
            "connector_name": "test-connector-name-axe-api-1",
            "component": "test-component-axe-api-1",
            "resource": "test-resource-axe-api-1"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-1"
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
      "state": 2,
      "output": "test-output-axe-api-2",
      "connector": "test-connector-axe-api-2",
      "connector_name": "test-connector-name-axe-api-2",
      "component": "test-component-axe-api-2",
      "resource": "test-resource-axe-api-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-2
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/ack:
    """json
    {
      "comment": "test-comment-axe-api-2-1"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "ack",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-2",
      "resource": "test-resource-axe-api-2",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/alarms/{{ .alarmId }}/ackremove:
    """json
    {
      "comment": "test-comment-axe-api-2-2"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "ackremove",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-2",
      "resource": "test-resource-axe-api-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-api-2",
            "connector_name": "test-connector-name-axe-api-2",
            "component": "test-component-axe-api-2",
            "resource": "test-resource-axe-api-2"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-2-1"
              },
              {
                "_t": "ackremove",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-2-2"
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
      "state": 2,
      "output": "test-output-axe-api-3",
      "connector": "test-connector-axe-api-3",
      "connector_name": "test-connector-name-axe-api-3",
      "component": "test-component-axe-api-3",
      "resource": "test-resource-axe-api-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-3
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/cancel:
    """json
    {
      "comment": "test-comment-axe-api-3"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "cancel",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-3",
      "resource": "test-resource-axe-api-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "canceled": {
              "_t": "cancel",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-3"
            },
            "status": {
              "_t": "statusinc",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-3",
              "val": 4
            },
            "connector": "test-connector-axe-api-3",
            "connector_name": "test-connector-name-axe-api-3",
            "component": "test-component-axe-api-3",
            "resource": "test-resource-axe-api-3"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "cancel",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-3"
              },
              {
                "_t": "statusinc",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-3",
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
      "state": 2,
      "output": "test-output-axe-api-4",
      "connector": "test-connector-axe-api-4",
      "connector_name": "test-connector-name-axe-api-4",
      "component": "test-component-axe-api-4",
      "resource": "test-resource-axe-api-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-4
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/cancel:
    """json
    {
      "comment": "test-comment-axe-api-4-1"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "cancel",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-4",
      "resource": "test-resource-axe-api-4",
      "source_type": "resource"
    }
    """
    When I do PUT /api/v4/alarms/{{ .alarmId }}/uncancel:
    """json
    {
      "comment": "test-comment-axe-api-4-2"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "uncancel",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-4",
      "resource": "test-resource-axe-api-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "status": {
              "_t": "statusdec",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-4-2",
              "val": 1
            },
            "connector": "test-connector-axe-api-4",
            "connector_name": "test-connector-name-axe-api-4",
            "component": "test-component-axe-api-4",
            "resource": "test-resource-axe-api-4"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "cancel",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-4-1"
              },
              {
                "_t": "statusinc",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-4-1",
                "val": 4
              },
              {
                "_t": "uncancel",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-4-2"
              },
              {
                "_t": "statusdec",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-4-2",
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
      "state": 2,
      "output": "test-output-axe-api-5",
      "connector": "test-connector-axe-api-5",
      "connector_name": "test-connector-name-axe-api-5",
      "component": "test-component-axe-api-5",
      "resource": "test-resource-axe-api-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-5
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/comment:
    """json
    {
      "comment": "test-comment-axe-api-5"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "comment",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-5",
      "resource": "test-resource-axe-api-5",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "last_comment": {
              "_t": "comment",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-5"
            },
            "connector": "test-connector-axe-api-5",
            "connector_name": "test-connector-name-axe-api-5",
            "component": "test-component-axe-api-5",
            "resource": "test-resource-axe-api-5"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "comment",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-5"
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
  Scenario: given assoc ticket event should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-6",
      "connector": "test-connector-axe-api-6",
      "connector_name": "test-connector-name-axe-api-6",
      "component": "test-component-axe-api-6",
      "resource": "test-resource-axe-api-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-6
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/assocticket:
    """json
    {
      "ticket": "test-ticket-axe-api-3",
      "url": "test-ticket-url-axe-api-3",
      "system_name": "test-ticket-system-name-axe-api-3",
      "data": {
        "test-ticket-param-axe-api-3": "test-ticket-param-val-axe-api-3"
      },
      "comment": "test-comment-axe-api-3"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "assocticket",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-6",
      "resource": "test-resource-axe-api-6",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-6
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
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-api-3. Ticket URL: test-ticket-url-axe-api-3. Ticket test-ticket-param-axe-api-3: test-ticket-param-val-axe-api-3.",
                "ticket": "test-ticket-axe-api-3",
                "ticket_system_name": "test-ticket-system-name-axe-api-3",
                "ticket_url": "test-ticket-url-axe-api-3",
                "ticket_data": {
                  "test-ticket-param-axe-api-3": "test-ticket-param-val-axe-api-3"
                },
                "ticket_comment": "test-comment-axe-api-3"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-api-3. Ticket URL: test-ticket-url-axe-api-3. Ticket test-ticket-param-axe-api-3: test-ticket-param-val-axe-api-3.",
              "ticket": "test-ticket-axe-api-3",
              "ticket_system_name": "test-ticket-system-name-axe-api-3",
              "ticket_url": "test-ticket-url-axe-api-3",
              "ticket_data": {
                "test-ticket-param-axe-api-3": "test-ticket-param-val-axe-api-3"
              },
              "ticket_comment": "test-comment-axe-api-3"
            },
            "connector": "test-connector-axe-api-6",
            "connector_name": "test-connector-name-axe-api-6",
            "component": "test-component-axe-api-6",
            "resource": "test-resource-axe-api-6"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-api-3. Ticket URL: test-ticket-url-axe-api-3. Ticket test-ticket-param-axe-api-3: test-ticket-param-val-axe-api-3.",
                "ticket": "test-ticket-axe-api-3",
                "ticket_system_name": "test-ticket-system-name-axe-api-3",
                "ticket_url": "test-ticket-url-axe-api-3",
                "ticket_data": {
                  "test-ticket-param-axe-api-3": "test-ticket-param-val-axe-api-3"
                },
                "ticket_comment": "test-comment-axe-api-3"
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
      "state": 2,
      "output": "test-output-axe-api-7",
      "connector": "test-connector-axe-api-7",
      "connector_name": "test-connector-name-axe-api-7",
      "component": "test-component-axe-api-7",
      "resource": "test-resource-axe-api-7",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-7
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/changestate:
    """json
    {
      "state": 3,
      "comment": "test-comment-axe-api-7"
    }
    """
    Then the response code should be 204
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "changestate",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-7",
      "resource": "test-resource-axe-api-7",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "state": {
              "_t": "changestate",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-7",
              "val": 3
            },
            "connector": "test-connector-axe-api-7",
            "connector_name": "test-connector-name-axe-api-7",
            "component": "test-component-axe-api-7",
            "resource": "test-resource-axe-api-7"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "changestate",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-7",
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
      "state": 2,
      "output": "test-output-axe-api-8",
      "connector": "test-connector-axe-api-8",
      "connector_name": "test-connector-name-axe-api-8",
      "component": "test-component-axe-api-8",
      "resource": "test-resource-axe-api-8",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-8
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/snooze:
    """json
    {
      "duration": {
        "value": 1,
        "unit": "h"
      },
      "comment": "test-comment-axe-api-8"
    }
    """
    Then the response code should be 204
    When I save response snoozeEventTimestamp={{ now }}
    When I wait the end of event processing which contains:
    """json
    {
      "event_type": "snooze",
      "connector": "api",
      "connector_name": "api",
      "component": "test-component-axe-api-8",
      "resource": "test-resource-axe-api-8",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-8",
              "t": {{ .snoozeEventTimestamp }},
              "val": {{ .snoozeEventTimestamp | sumTime 3600 }}
            },
            "connector": "test-connector-axe-api-8",
            "connector_name": "test-connector-name-axe-api-8",
            "component": "test-component-axe-api-8",
            "resource": "test-resource-axe-api-8"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "snooze",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-8",
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
  Scenario: given ack resources event should update resource alarms
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-9",
      "connector": "test-connector-axe-api-9",
      "connector_name": "test-connector-name-axe-api-9",
      "component": "test-component-axe-api-9",
      "resource": "test-resource-axe-api-9",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 2,
      "output": "test-output-axe-api-9",
      "connector": "test-connector-axe-api-9",
      "connector_name": "test-connector-name-axe-api-9",
      "component": "test-component-axe-api-9",
      "source_type": "component"
    }
    """
    When I do GET /api/v4/alarms?search=test-component-axe-api-9&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/alarms/{{ .alarmId }}/ack:
    """json
    {
      "ack_resources": true,
      "comment": "test-comment-axe-api-9"
    }
    """
    Then the response code should be 204
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-api-9",
        "source_type": "component"
      },
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-api-9",
        "resource": "test-resource-axe-api-9",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-api-9
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
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-api-9"
            },
            "connector": "test-connector-axe-api-9",
            "connector_name": "test-connector-name-axe-api-9",
            "component": "test-component-axe-api-9",
            "resource": "test-resource-axe-api-9"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-api-9"
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
