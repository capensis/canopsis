Feature: create and update meta alarm
  I need to be able to create and update meta alarm

  @concurrent
  Scenario: given manual meta alarm and added child should inherit metaalarm actions
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-1",
        "state": 1,
        "output": "test-output-axe-correlation-third-1"
      },
      {
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-2",
        "state": 2,
        "output": "test-output-axe-correlation-third-1"
      },
      {
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-3",
        "state": 3,
        "output": "test-output-axe-correlation-third-1"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-1&sort_by=v.resource&sort=asc
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I save response alarmId3={{ (index .lastResponse.data 2)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaAlarm-axe-correlation-third-1",
      "comment": "test-metaAlarm-axe-correlation-third-1-comment",
      "alarms": ["{{ .alarmId1 }}", "{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-1-1&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "display_name": "test-metaAlarm-axe-correlation-third-1"
          }
        }
      ]
    }
    """
    When I save response metaAlarmID={{ (index .lastResponse.data 0)._id }}
    When I save response metaAlarmEntityID={{ (index .lastResponse.data 0).entity._id }}
    When I save response metaAlarmConnector={{ (index .lastResponse.data 0).v.connector }}
    When I save response metaAlarmConnectorName={{ (index .lastResponse.data 0).v.connector_name }}
    When I save response metaAlarmComponent={{ (index .lastResponse.data 0).v.component }}
    When I save response metaAlarmResource={{ (index .lastResponse.data 0).v.resource }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-1&correlation=true&sort_by=v.meta&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "output": "test-metaAlarm-axe-correlation-third-1-comment",
            "component": "metaalarm",
            "connector": "engine",
            "connector_name": "correlation",
            "display_name": "test-metaAlarm-axe-correlation-third-1",
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
            "resource": "test-resource-axe-correlation-third-1-3"
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
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "test-ack-19"
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
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "assocticket",
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "ticket": "ticket-19"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "comment",
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "output": "comment-19"
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "{{ .metaAlarmConnector }}",
        "connector_name": "{{ .metaAlarmConnectorName }}",
        "component": "{{ .metaAlarmComponent }}",
        "resource": "{{ .metaAlarmResource }}",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I send an event:
    """json
    {
      "connector": "{{ .metaAlarmConnector }}",
      "connector_name": "{{ .metaAlarmConnectorName }}",
      "source_type": "resource",
      "event_type": "snooze",
      "component": "{{ .metaAlarmComponent }}",
      "resource": "{{ .metaAlarmResource }}",
      "duration": 3600
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
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-1",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-2",
        "source_type": "resource"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ .metaAlarmID }}/add:
    """json
    {
      "comment": "test-metaAlarm-axe-correlation-third-1-comment",
      "alarms": ["{{ .alarmId3 }}"]
    }
    """
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "assocticket",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "comment",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-3",
        "source_type": "resource"
      },
      {
        "event_type": "snooze",
        "connector": "test-connector-axe-correlation-third-1",
        "connector_name": "test-connector-name-axe-correlation-third-1",
        "component": "test-component-axe-correlation-third-1",
        "resource": "test-resource-axe-correlation-third-1-3",
        "source_type": "resource"
      }
    ]
    """
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
                  "component": "test-component-axe-correlation-third-1",
                  "connector": "test-connector-axe-correlation-third-1",
                  "connector_name": "test-connector-name-axe-correlation-third-1",
                  "resource": "test-resource-axe-correlation-third-1-1",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net"
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-third-1",
                  "connector": "test-connector-axe-correlation-third-1",
                  "connector_name": "test-connector-name-axe-correlation-third-1",
                  "resource": "test-resource-axe-correlation-third-1-2",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net"
                  }
                }
              },
              {
                "v": {
                  "component": "test-component-axe-correlation-third-1",
                  "connector": "test-connector-axe-correlation-third-1",
                  "connector_name": "test-connector-name-axe-correlation-third-1",
                  "resource": "test-resource-axe-correlation-third-1-3",
                  "parents": ["{{ .metaAlarmEntityID }}"],
                  "ack": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "test-ack-19"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "Ticket ID: ticket-19."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "m": "comment-19"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net"
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
