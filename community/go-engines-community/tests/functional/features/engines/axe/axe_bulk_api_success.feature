Feature: update alarm
  I need to be able to update alarm

  @concurrent
  Scenario: given ack should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-1",
        "connector": "test-connector-axe-bulk-api-1",
        "connector_name": "test-connector-name-axe-bulk-api-1",
        "component": "test-component-axe-bulk-api-1",
        "resource": "test-resource-axe-bulk-api-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-1",
        "connector": "test-connector-axe-bulk-api-1",
        "connector_name": "test-connector-name-axe-bulk-api-1",
        "component": "test-component-axe-bulk-api-1",
        "resource": "test-resource-axe-bulk-api-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-1&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/ack:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "comment": "test-comment-axe-bulk-api-1-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "comment": "test-comment-axe-bulk-api-1-2"
      },
      {
        "_id": "test-alarm-not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "comment": "test-comment-axe-bulk-api-1-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "comment": "test-comment-axe-bulk-api-1-2"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-1",
        "resource": "test-resource-axe-bulk-api-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-1",
        "resource": "test-resource-axe-bulk-api-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-1&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-1-1"
            },
            "connector": "test-connector-axe-bulk-api-1",
            "connector_name": "test-connector-name-axe-bulk-api-1",
            "component": "test-component-axe-bulk-api-1",
            "resource": "test-resource-axe-bulk-api-1-1"
          }
        },
        {
          "v": {
            "ack": {
              "_t": "ack",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-1-2"
            },
            "connector": "test-connector-axe-bulk-api-1",
            "connector_name": "test-connector-name-axe-bulk-api-1",
            "component": "test-component-axe-bulk-api-1",
            "resource": "test-resource-axe-bulk-api-1-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "test-comment-axe-bulk-api-1-1"
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
                "m": "test-comment-axe-bulk-api-1-2"
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
  Scenario: given remove ack should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-2",
        "connector": "test-connector-axe-bulk-api-2",
        "connector_name": "test-connector-name-axe-bulk-api-2",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-2",
        "connector": "test-connector-axe-bulk-api-2",
        "connector_name": "test-connector-name-axe-bulk-api-2",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/ack:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}"
      },
      {
        "_id": "{{ .alarmId2 }}"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      },
      {
        "status": 200
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do PUT /api/v4/bulk/alarms/ackremove:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "comment": "test-comment-axe-bulk-api-2-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "comment": "test-comment-axe-bulk-api-2-2"
      },
      {
        "_id": "test-alarm-not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "comment": "test-comment-axe-bulk-api-2-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "comment": "test-comment-axe-bulk-api-2-2"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ackremove",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-1",
        "source_type": "resource"
      },
      {
        "event_type": "ackremove",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-2",
        "resource": "test-resource-axe-bulk-api-2-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-2&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector": "test-connector-axe-bulk-api-2",
            "connector_name": "test-connector-name-axe-bulk-api-2",
            "component": "test-component-axe-bulk-api-2",
            "resource": "test-resource-axe-bulk-api-2-1"
          }
        },
        {
          "v": {
            "connector": "test-connector-axe-bulk-api-2",
            "connector_name": "test-connector-name-axe-bulk-api-2",
            "component": "test-component-axe-bulk-api-2",
            "resource": "test-resource-axe-bulk-api-2-2"
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
    Then the response key "data.0.v.ack" should not exist
    Then the response key "data.1.v.ack" should not exist
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "_t": "ack"
              },
              {
                "_t": "ackremove",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-2-1"
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
                "_t": "stateinc"
              },
              {
                "_t": "statusinc"
              },
              {
                "_t": "ack"
              },
              {
                "_t": "ackremove",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-2-2"
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
  Scenario: given cancel should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-3",
        "connector": "test-connector-axe-bulk-api-3",
        "connector_name": "test-connector-name-axe-bulk-api-3",
        "component": "test-component-axe-bulk-api-3",
        "resource": "test-resource-axe-bulk-api-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-3",
        "connector": "test-connector-axe-bulk-api-3",
        "connector_name": "test-connector-name-axe-bulk-api-3",
        "component": "test-component-axe-bulk-api-3",
        "resource": "test-resource-axe-bulk-api-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-3&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/cancel:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "comment": "test-comment-axe-bulk-api-3-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "comment": "test-comment-axe-bulk-api-3-2"
      },
      {
        "_id": "test-alarm-not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "comment": "test-comment-axe-bulk-api-3-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "comment": "test-comment-axe-bulk-api-3-2"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "cancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-3",
        "resource": "test-resource-axe-bulk-api-3-1",
        "source_type": "resource"
      },
      {
        "event_type": "cancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-3",
        "resource": "test-resource-axe-bulk-api-3-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-3&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-3-1"
            },
            "status": {
              "_t": "statusinc",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-3-1",
              "val": 4
            },
            "connector": "test-connector-axe-bulk-api-3",
            "connector_name": "test-connector-name-axe-bulk-api-3",
            "component": "test-component-axe-bulk-api-3",
            "resource": "test-resource-axe-bulk-api-3-1"
          }
        },
        {
          "v": {
            "canceled": {
              "_t": "cancel",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-3-2"
            },
            "status": {
              "_t": "statusinc",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-3-2",
              "val": 4
            },
            "connector": "test-connector-axe-bulk-api-3",
            "connector_name": "test-connector-name-axe-bulk-api-3",
            "component": "test-component-axe-bulk-api-3",
            "resource": "test-resource-axe-bulk-api-3-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "test-comment-axe-bulk-api-3-1"
              },
              {
                "_t": "statusinc",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-3-1",
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
      },
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
                "m": "test-comment-axe-bulk-api-3-2"
              },
              {
                "_t": "statusinc",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-3-2",
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
  Scenario: given uncancel should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-4",
        "connector": "test-connector-axe-bulk-api-4",
        "connector_name": "test-connector-name-axe-bulk-api-4",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-4",
        "connector": "test-connector-axe-bulk-api-4",
        "connector_name": "test-connector-name-axe-bulk-api-4",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-4&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/cancel:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}"
      },
      {
        "_id": "{{ .alarmId2 }}"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200
      },
      {
        "status": 200
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "cancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "cancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do PUT /api/v4/bulk/alarms/uncancel:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "comment": "test-comment-axe-bulk-api-4-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "comment": "test-comment-axe-bulk-api-4-2"
      },
      {
        "_id": "test-alarm-not-exist"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "comment": "test-comment-axe-bulk-api-4-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "comment": "test-comment-axe-bulk-api-4-2"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "uncancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-1",
        "source_type": "resource"
      },
      {
        "event_type": "uncancel",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-4",
        "resource": "test-resource-axe-bulk-api-4-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-4&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-4-1",
              "val": 1
            },
            "connector": "test-connector-axe-bulk-api-4",
            "connector_name": "test-connector-name-axe-bulk-api-4",
            "component": "test-component-axe-bulk-api-4",
            "resource": "test-resource-axe-bulk-api-4-1"
          }
        },
        {
          "v": {
            "status": {
              "_t": "statusdec",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-4-2",
              "val": 1
            },
            "connector": "test-connector-axe-bulk-api-4",
            "connector_name": "test-connector-name-axe-bulk-api-4",
            "component": "test-component-axe-bulk-api-4",
            "resource": "test-resource-axe-bulk-api-4-2"
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
    Then the response key "data.0.v.canceled" should not exist
    Then the response key "data.1.v.canceled" should not exist
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "_t": "cancel"
              },
              {
                "_t": "statusinc",
                "val": 4
              },
              {
                "_t": "uncancel",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-4-1"
              },
              {
                "_t": "statusdec",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-4-1",
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
      },
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
                "_t": "cancel"
              },
              {
                "_t": "statusinc",
                "val": 4
              },
              {
                "_t": "uncancel",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-4-2"
              },
              {
                "_t": "statusdec",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "test-comment-axe-bulk-api-4-2",
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
  Scenario: given comment should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-5",
        "connector": "test-connector-axe-bulk-api-5",
        "connector_name": "test-connector-name-axe-bulk-api-5",
        "component": "test-component-axe-bulk-api-5",
        "resource": "test-resource-axe-bulk-api-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-5",
        "connector": "test-connector-axe-bulk-api-5",
        "connector_name": "test-connector-name-axe-bulk-api-5",
        "component": "test-component-axe-bulk-api-5",
        "resource": "test-resource-axe-bulk-api-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-5&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/comment:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "comment": "test-comment-axe-bulk-api-5-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "comment": "test-comment-axe-bulk-api-5-2"
      },
      {
        "_id": "test-alarm-not-exist"
      },
      {
        "_id": "test-alarm-not-exist",
        "comment": "test-comment"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "comment": "test-comment-axe-bulk-api-5-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "comment": "test-comment-axe-bulk-api-5-2"
        }
      },
      {
        "status": 400,
        "errors": {
          "comment": "Comment is missing."
        },
        "item": {
          "_id": "test-alarm-not-exist"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist",
          "comment": "test-comment"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-5",
        "resource": "test-resource-axe-bulk-api-5-1",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-5",
        "resource": "test-resource-axe-bulk-api-5-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-5&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-5-1"
            },
            "connector": "test-connector-axe-bulk-api-5",
            "connector_name": "test-connector-name-axe-bulk-api-5",
            "component": "test-component-axe-bulk-api-5",
            "resource": "test-resource-axe-bulk-api-5-1"
          }
        },
        {
          "v": {
            "last_comment": {
              "_t": "comment",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-5-2"
            },
            "connector": "test-connector-axe-bulk-api-5",
            "connector_name": "test-connector-name-axe-bulk-api-5",
            "component": "test-component-axe-bulk-api-5",
            "resource": "test-resource-axe-bulk-api-5-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "test-comment-axe-bulk-api-5-1"
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
                "m": "test-comment-axe-bulk-api-5-2"
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
  Scenario: given assoc ticket should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-6",
        "connector": "test-connector-axe-bulk-api-6",
        "connector_name": "test-connector-name-axe-bulk-api-6",
        "component": "test-component-axe-bulk-api-6",
        "resource": "test-resource-axe-bulk-api-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-6",
        "connector": "test-connector-axe-bulk-api-6",
        "connector_name": "test-connector-name-axe-bulk-api-6",
        "component": "test-component-axe-bulk-api-6",
        "resource": "test-resource-axe-bulk-api-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-6&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/assocticket:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "ticket": "test-ticket-axe-bulk-api-6-1",
        "url": "test-ticket-url-axe-bulk-api-6-1",
        "system_name": "test-ticket-system-name-axe-bulk-api-6-1",
        "data": {
          "test-ticket-param-axe-bulk-api-6-1": "test-ticket-param-val-axe-bulk-api-6-1"
        },
        "comment": "test-comment-axe-bulk-api-6-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "ticket": "test-ticket-axe-bulk-api-6-2",
        "url": "test-ticket-url-axe-bulk-api-6-2",
        "system_name": "test-ticket-system-name-axe-bulk-api-6-2",
        "data": {
          "test-ticket-param-axe-bulk-api-6-2": "test-ticket-param-val-axe-bulk-api-6-2"
        },
        "comment": "test-comment-axe-bulk-api-6-2"
      },
      {
        "_id": "test-alarm-not-exist"
      },
      {
        "_id": "test-alarm-not-exist",
        "ticket": "test-ticket"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "ticket": "test-ticket-axe-bulk-api-6-1",
          "url": "test-ticket-url-axe-bulk-api-6-1",
          "system_name": "test-ticket-system-name-axe-bulk-api-6-1",
          "data": {
            "test-ticket-param-axe-bulk-api-6-1": "test-ticket-param-val-axe-bulk-api-6-1"
          },
          "comment": "test-comment-axe-bulk-api-6-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "ticket": "test-ticket-axe-bulk-api-6-2",
          "url": "test-ticket-url-axe-bulk-api-6-2",
          "system_name": "test-ticket-system-name-axe-bulk-api-6-2",
          "data": {
            "test-ticket-param-axe-bulk-api-6-2": "test-ticket-param-val-axe-bulk-api-6-2"
          },
          "comment": "test-comment-axe-bulk-api-6-2"
        }
      },
      {
        "status": 400,
        "errors": {
          "ticket": "Ticket is missing."
        },
        "item": {
          "_id": "test-alarm-not-exist"
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist",
          "ticket": "test-ticket"
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-6",
        "resource": "test-resource-axe-bulk-api-6-1",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-6",
        "resource": "test-resource-axe-bulk-api-6-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-6&sort_by=v.resource&sort=asc
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
                "m": "Ticket ID: test-ticket-axe-bulk-api-6-1. Ticket URL: test-ticket-url-axe-bulk-api-6-1. Ticket test-ticket-param-axe-bulk-api-6-1: test-ticket-param-val-axe-bulk-api-6-1.",
                "ticket": "test-ticket-axe-bulk-api-6-1",
                "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-1",
                "ticket_url": "test-ticket-url-axe-bulk-api-6-1",
                "ticket_data": {
                  "test-ticket-param-axe-bulk-api-6-1": "test-ticket-param-val-axe-bulk-api-6-1"
                },
                "ticket_comment": "test-comment-axe-bulk-api-6-1"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-bulk-api-6-1. Ticket URL: test-ticket-url-axe-bulk-api-6-1. Ticket test-ticket-param-axe-bulk-api-6-1: test-ticket-param-val-axe-bulk-api-6-1.",
              "ticket": "test-ticket-axe-bulk-api-6-1",
              "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-1",
              "ticket_url": "test-ticket-url-axe-bulk-api-6-1",
              "ticket_data": {
                "test-ticket-param-axe-bulk-api-6-1": "test-ticket-param-val-axe-bulk-api-6-1"
              },
              "ticket_comment": "test-comment-axe-bulk-api-6-1"
            },
            "connector": "test-connector-axe-bulk-api-6",
            "connector_name": "test-connector-name-axe-bulk-api-6",
            "component": "test-component-axe-bulk-api-6",
            "resource": "test-resource-axe-bulk-api-6-1"
          }
        },
        {
          "v": {
            "tickets": [
              {
                "_t": "assocticket",
                "a": "root John Doe admin@canopsis.net",
                "user_id": "root",
                "initiator": "user",
                "m": "Ticket ID: test-ticket-axe-bulk-api-6-2. Ticket URL: test-ticket-url-axe-bulk-api-6-2. Ticket test-ticket-param-axe-bulk-api-6-2: test-ticket-param-val-axe-bulk-api-6-2.",
                "ticket": "test-ticket-axe-bulk-api-6-2",
                "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-2",
                "ticket_url": "test-ticket-url-axe-bulk-api-6-2",
                "ticket_data": {
                  "test-ticket-param-axe-bulk-api-6-2": "test-ticket-param-val-axe-bulk-api-6-2"
                },
                "ticket_comment": "test-comment-axe-bulk-api-6-2"
              }
            ],
            "ticket": {
              "_t": "assocticket",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "Ticket ID: test-ticket-axe-bulk-api-6-2. Ticket URL: test-ticket-url-axe-bulk-api-6-2. Ticket test-ticket-param-axe-bulk-api-6-2: test-ticket-param-val-axe-bulk-api-6-2.",
              "ticket": "test-ticket-axe-bulk-api-6-2",
              "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-2",
              "ticket_url": "test-ticket-url-axe-bulk-api-6-2",
              "ticket_data": {
                "test-ticket-param-axe-bulk-api-6-2": "test-ticket-param-val-axe-bulk-api-6-2"
              },
              "ticket_comment": "test-comment-axe-bulk-api-6-2"
            },
            "connector": "test-connector-axe-bulk-api-6",
            "connector_name": "test-connector-name-axe-bulk-api-6",
            "component": "test-component-axe-bulk-api-6",
            "resource": "test-resource-axe-bulk-api-6-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "Ticket ID: test-ticket-axe-bulk-api-6-1. Ticket URL: test-ticket-url-axe-bulk-api-6-1. Ticket test-ticket-param-axe-bulk-api-6-1: test-ticket-param-val-axe-bulk-api-6-1.",
                "ticket": "test-ticket-axe-bulk-api-6-1",
                "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-1",
                "ticket_url": "test-ticket-url-axe-bulk-api-6-1",
                "ticket_data": {
                  "test-ticket-param-axe-bulk-api-6-1": "test-ticket-param-val-axe-bulk-api-6-1"
                },
                "ticket_comment": "test-comment-axe-bulk-api-6-1"
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
                "m": "Ticket ID: test-ticket-axe-bulk-api-6-2. Ticket URL: test-ticket-url-axe-bulk-api-6-2. Ticket test-ticket-param-axe-bulk-api-6-2: test-ticket-param-val-axe-bulk-api-6-2.",
                "ticket": "test-ticket-axe-bulk-api-6-2",
                "ticket_system_name": "test-ticket-system-name-axe-bulk-api-6-2",
                "ticket_url": "test-ticket-url-axe-bulk-api-6-2",
                "ticket_data": {
                  "test-ticket-param-axe-bulk-api-6-2": "test-ticket-param-val-axe-bulk-api-6-2"
                },
                "ticket_comment": "test-comment-axe-bulk-api-6-2"
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
  Scenario: given change state should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-7",
        "connector": "test-connector-axe-bulk-api-7",
        "connector_name": "test-connector-name-axe-bulk-api-7",
        "component": "test-component-axe-bulk-api-7",
        "resource": "test-resource-axe-bulk-api-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-7",
        "connector": "test-connector-axe-bulk-api-7",
        "connector_name": "test-connector-name-axe-bulk-api-7",
        "component": "test-component-axe-bulk-api-7",
        "resource": "test-resource-axe-bulk-api-7-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-7&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/changestate:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "state": 3,
        "comment": "test-comment-axe-bulk-api-7-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "state": 3,
        "comment": "test-comment-axe-bulk-api-7-2"
      },
      {
        "_id": "test-alarm-not-exist"
      },
      {
        "_id": "test-alarm-not-exist",
        "state": 5
      },
      {
        "_id": "test-alarm-not-exist",
        "state": 0
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "state": 3,
          "comment": "test-comment-axe-bulk-api-7-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "state": 3,
          "comment": "test-comment-axe-bulk-api-7-2"
        }
      },
      {
        "status": 400,
        "errors": {
          "state": "State is missing."
        },
        "item": {
          "_id": "test-alarm-not-exist"
        }
      },
      {
        "status": 400,
        "errors": {
          "state": "State must be one of [0 1 2 3]."
        },
        "item": {
          "_id": "test-alarm-not-exist",
          "state": 5
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist",
          "state": 0
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "changestate",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-7",
        "resource": "test-resource-axe-bulk-api-7-1",
        "source_type": "resource"
      },
      {
        "event_type": "changestate",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-7",
        "resource": "test-resource-axe-bulk-api-7-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-7&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-7-1",
              "val": 3
            },
            "connector": "test-connector-axe-bulk-api-7",
            "connector_name": "test-connector-name-axe-bulk-api-7",
            "component": "test-component-axe-bulk-api-7",
            "resource": "test-resource-axe-bulk-api-7-1"
          }
        },
        {
          "v": {
            "state": {
              "_t": "changestate",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-7-2",
              "val": 3
            },
            "connector": "test-connector-axe-bulk-api-7",
            "connector_name": "test-connector-name-axe-bulk-api-7",
            "component": "test-component-axe-bulk-api-7",
            "resource": "test-resource-axe-bulk-api-7-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "test-comment-axe-bulk-api-7-1",
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
      },
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
                "m": "test-comment-axe-bulk-api-7-2",
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
  Scenario: given snooze should update alarm
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-8",
        "connector": "test-connector-axe-bulk-api-8",
        "connector_name": "test-connector-name-axe-bulk-api-8",
        "component": "test-component-axe-bulk-api-8",
        "resource": "test-resource-axe-bulk-api-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 2,
        "output": "test-output-axe-bulk-api-8",
        "connector": "test-connector-axe-bulk-api-8",
        "connector_name": "test-connector-name-axe-bulk-api-8",
        "component": "test-component-axe-bulk-api-8",
        "resource": "test-resource-axe-bulk-api-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-8&sort_by=v.resource&sort=asc
    Then the response code should be 200
    Then I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    Then I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I do PUT /api/v4/bulk/alarms/snooze:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "duration": {
          "value": 1,
          "unit": "h"
        },
        "comment": "test-comment-axe-bulk-api-8-1"
      },
      {
        "_id": "{{ .alarmId2 }}",
        "duration": {
          "value": 1,
          "unit": "h"
        },
        "comment": "test-comment-axe-bulk-api-8-2"
      },
      {
        "_id": "test-alarm-not-exist"
      },
      {
        "_id": "test-alarm-not-exist",
        "duration": {
          "value": 1,
          "unit": "y"
        }
      },
      {
        "_id": "test-alarm-not-exist",
        "duration": {
          "value": 1,
          "unit": "h"
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "status": 200,
        "id": "{{ .alarmId1 }}",
        "item": {
          "_id": "{{ .alarmId1 }}",
          "duration": {
            "value": 1,
            "unit": "h"
          },
          "comment": "test-comment-axe-bulk-api-8-1"
        }
      },
      {
        "status": 200,
        "id": "{{ .alarmId2 }}",
        "item": {
          "_id": "{{ .alarmId2 }}",
          "duration": {
            "value": 1,
            "unit": "h"
          },
          "comment": "test-comment-axe-bulk-api-8-2"
        }
      },
      {
        "status": 400,
        "errors": {
          "duration.value": "Value is missing.",
          "duration.unit": "Unit is missing."
        },
        "item": {
          "_id": "test-alarm-not-exist"
        }
      },
      {
        "status": 400,
        "errors": {
          "duration": "Duration is invalid."
        },
        "item": {
          "_id": "test-alarm-not-exist",
          "duration": {
            "value": 1,
            "unit": "y"
          }
        }
      },
      {
        "status": 404,
        "error": "Not found",
        "item": {
          "_id": "test-alarm-not-exist",
          "duration": {
            "value": 1,
            "unit": "h"
          }
        }
      }
    ]
    """
    When I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "snooze",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-8",
        "resource": "test-resource-axe-bulk-api-8-1",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "api",
        "connector_name": "api",
        "component": "test-component-axe-bulk-api-8",
        "resource": "test-resource-axe-bulk-api-8-2",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-bulk-api-8&sort_by=v.resource&sort=asc
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
              "m": "test-comment-axe-bulk-api-8-1"
            },
            "connector": "test-connector-axe-bulk-api-8",
            "connector_name": "test-connector-name-axe-bulk-api-8",
            "component": "test-component-axe-bulk-api-8",
            "resource": "test-resource-axe-bulk-api-8-1"
          }
        },
        {
          "v": {
            "snooze": {
              "_t": "snooze",
              "a": "root John Doe admin@canopsis.net",
              "user_id": "root",
              "initiator": "user",
              "m": "test-comment-axe-bulk-api-8-2"
            },
            "connector": "test-connector-axe-bulk-api-8",
            "connector_name": "test-connector-name-axe-bulk-api-8",
            "component": "test-component-axe-bulk-api-8",
            "resource": "test-resource-axe-bulk-api-8-2"
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
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ .alarmId1 }}",
        "steps": {
          "page": 1
        }
      },
      {
        "_id": "{{ .alarmId2 }}",
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
                "m": "test-comment-axe-bulk-api-8-1"
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
                "m": "test-comment-axe-bulk-api-8-2"
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
