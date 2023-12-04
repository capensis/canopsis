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
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
              "val": 2
            },
            "status": {
              "_t": "statusinc",
              "a": "engine.correlation",
              "user_id": "",
              "initiator": "system",
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
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/ack:
    """json
    {
      "comment": "test-ack-comment-axe-correlation-third-1"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "ack",
        "connector": "api",
        "connector_name": "api",
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
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/assocticket:
    """json
    {
      "ticket": "test-ticket-axe-correlation-third-1"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "assocticket",
        "connector": "api",
        "connector_name": "api",
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
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/comment:
    """json
    {
      "comment": "test-comment-axe-correlation-third-1"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "comment",
        "connector": "api",
        "connector_name": "api",
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
    When I do PUT /api/v4/alarms/{{ .metaAlarmID }}/snooze:
    """json
    {
      "duration": {
        "value": 1,
        "unit": "h"
      },
      "comment": "test-snooze-comment-axe-correlation-third-1"
    }
    """
    Then the response code should be 204
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "snooze",
        "connector": "api",
        "connector_name": "api",
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
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-ack-comment-axe-correlation-third-1"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "Ticket ID: test-ticket-axe-correlation-third-1."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-comment-axe-correlation-third-1"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-snooze-comment-axe-correlation-third-1"
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
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-ack-comment-axe-correlation-third-1"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "Ticket ID: test-ticket-axe-correlation-third-1."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-comment-axe-correlation-third-1"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-snooze-comment-axe-correlation-third-1"
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
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-ack-comment-axe-correlation-third-1"
                  },
                  "ticket": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "Ticket ID: test-ticket-axe-correlation-third-1."
                  },
                  "last_comment": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-comment-axe-correlation-third-1"
                  },
                  "snooze": {
                    "a": "root John Doe admin@canopsis.net",
                    "user_id": "root",
                    "initiator": "user",
                    "m": "test-snooze-comment-axe-correlation-third-1"
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

  @concurrent
  Scenario: given manual meta alarm with auto_resolve=true should resolve manual metaalarm when the last child is resolved
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-third-2",
        "connector_name": "test-connector-name-axe-correlation-third-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-third-2",
        "resource": "test-resource-axe-correlation-third-2-1",
        "state": 1,
        "output": "test-output-axe-correlation-third-2"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-2-1
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaAlarm-axe-correlation-third-2",
      "comment": "test-metaAlarm-axe-correlation-third-2-comment",
      "alarms": ["{{ .alarmId1 }}"],
      "auto_resolve": true
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaAlarm-axe-correlation-third-2 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaAlarm-axe-correlation-third-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-third-2",
      "connector_name": "test-connector-name-axe-correlation-third-2",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-axe-correlation-third-2",
      "resource": "test-resource-axe-correlation-third-2-1",
      "output": "test-output-axe-correlation-third-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-third-2",
      "connector_name": "test-connector-name-axe-correlation-third-2",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-axe-correlation-third-2",
      "resource": "test-resource-axe-correlation-third-2-1",
      "output": "test-output-axe-correlation-third-2"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-2-1&correlation=true until response code is 200 and response key "data.0.v.resolved" is greater or equal than 1

  @concurrent
  Scenario: given manual meta alarm with auto_resolve=false shouldn't resolve manual metaalarm when the last child is resolved
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-correlation-third-3",
        "connector_name": "test-connector-name-axe-correlation-third-3",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-axe-correlation-third-3",
        "resource": "test-resource-axe-correlation-third-3-1",
        "state": 1,
        "output": "test-output-axe-correlation-third-3"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-3-1
    When I save response alarmId1={{ (index .lastResponse.data 0)._id }}
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaAlarm-axe-correlation-third-3",
      "comment": "test-metaAlarm-axe-correlation-third-3-comment",
      "alarms": ["{{ .alarmId1 }}"],
      "auto_resolve": false
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaAlarm-axe-correlation-third-3 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaAlarm-axe-correlation-third-3"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-third-3",
      "connector_name": "test-connector-name-axe-correlation-third-3",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-axe-correlation-third-3",
      "resource": "test-resource-axe-correlation-third-3-1",
      "output": "test-output-axe-correlation-third-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-correlation-third-3",
      "connector_name": "test-connector-name-axe-correlation-third-3",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-axe-correlation-third-3",
      "resource": "test-resource-axe-correlation-third-3-1",
      "output": "test-output-axe-correlation-third-3"
    }
    """
    When I wait 3s
    Then the response key "data.0.v.resolved" should not exist

  @concurrent
  Scenario: given meta alarm child and check event should update parent last event date
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-axe-correlation-third-4",
      "connector": "test-connector-axe-correlation-third-4",
      "connector_name": "test-connector-name-axe-correlation-third-4",
      "component": "test-component-axe-correlation-third-4",
      "resource": "test-resource-axe-correlation-third-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-4&correlation=false
    Then the response code should be 200
    When I save response alarmLastEventDate1={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "last_event_date": {{ .alarmLastEventDate1 }}
          }
        }
      ]
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-axe-correlation-third-4",
      "connector": "test-connector-axe-correlation-third-4",
      "connector_name": "test-connector-name-axe-correlation-third-4",
      "component": "test-component-axe-correlation-third-4",
      "resource": "test-resource-axe-correlation-third-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-4&correlation=false
    Then the response code should be 200
    When I save response alarmLastEventDate2={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-4&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "v": {
            "last_event_date": {{ .alarmLastEventDate2 }}
          }
        }
      ]
    }
    """

  @concurrent
  Scenario: given meta alarm child and removed children should update parent last event date
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-axe-correlation-third-5",
      "connector": "test-connector-axe-correlation-third-5",
      "connector_name": "test-connector-name-axe-correlation-third-5",
      "component": "test-component-axe-correlation-third-5",
      "resource": "test-resource-axe-correlation-third-5-1",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=false
    Then the response code should be 200
    When I save response alarmLastEventDate1={{ (index .lastResponse.data 0).v.last_event_date }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "v": {
            "last_event_date": {{ .alarmLastEventDate1 }}
          }
        }
      ]
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-axe-correlation-third-5",
        "connector": "test-connector-axe-correlation-third-5",
        "connector_name": "test-connector-name-axe-correlation-third-5",
        "component": "test-component-axe-correlation-third-5",
        "resource": "test-resource-axe-correlation-third-5-2",
        "source_type": "resource"
      },
      {
        "event_type": "check",
        "state": 1,
        "output": "test-output-axe-correlation-third-5",
        "connector": "test-connector-axe-correlation-third-5",
        "connector_name": "test-connector-name-axe-correlation-third-5",
        "component": "test-component-axe-correlation-third-5",
        "resource": "test-resource-axe-correlation-third-5-3",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=false&sort_by=v.resource&sort=asc
    Then the response code should be 200
    When I save response alarmId2={{ (index .lastResponse.data 1)._id }}
    When I save response alarmLastEventDate2={{ (index .lastResponse.data 1).v.last_event_date }}
    When I save response alarmId3={{ (index .lastResponse.data 2)._id }}
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=true until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 3,
          "v": {
            "last_event_date": {{ .alarmLastEventDate2 }}
          }
        }
      ]
    }
    """
    When I save response metaAlarmId={{ (index .lastResponse.data 0)._id }}
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "alarms": ["{{ .alarmId2 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=true&sort_by=v.meta&sort=desc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 2,
          "v": {
            "last_event_date": {{ .alarmLastEventDate2 }}
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-third-5-2"
          }
        }
      ]
    }
    """
    When I do PUT /api/v4/cat/meta-alarms/{{ .metaAlarmId }}/remove:
    """json
    {
      "alarms": ["{{ .alarmId3 }}"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-axe-correlation-third-5&correlation=true&multi_sort[]=v.meta,desc&multi_sort[]=v.resource,asc until response code is 200 and body contains:
    """json
    {
      "data": [
        {
          "children": 1,
          "v": {
            "last_event_date": {{ .alarmLastEventDate1 }}
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-third-5-2"
          }
        },
        {
          "v": {
            "resource": "test-resource-axe-correlation-third-5-3"
          }
        }
      ]
    }
    """
