Feature: Get alarms counters
  I need to be able to get a alarms counters

  @concurrent
  Scenario: given alarms should return alarms counters
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-1",
        "connector_name": "test-connector-name-axe-alarm-count-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-1",
        "resource": "test-resource-axe-alarm-count-get-1-1",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-1"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-1",
        "connector_name": "test-connector-name-axe-alarm-count-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-1",
        "resource": "test-resource-axe-alarm-count-get-1-2",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-1"
      }
    ]
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "total": 2,
      "total_active": 2,
      "snooze": 0,
      "ack": 0,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  @concurrent
  Scenario: given acked alarms should return ack alarms counter
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-2",
        "connector_name": "test-connector-name-axe-alarm-count-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-2",
        "resource": "test-resource-axe-alarm-count-get-2-1",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-2"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-2",
        "connector_name": "test-connector-name-axe-alarm-count-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-2",
        "resource": "test-resource-axe-alarm-count-get-2-2",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-2"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-alarm-count-get-2",
      "connector_name": "test-connector-name-axe-alarm-count-get-2",
      "source_type": "resource",
      "event_type": "ack",
      "component":  "test-component-axe-alarm-count-get-2",
      "resource": "test-resource-axe-alarm-count-get-2-2",
      "output": "test-output-axe-alarm-count-get-2"
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "total": 2,
      "total_active": 2,
      "snooze": 0,
      "ack": 1,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  @concurrent
  Scenario: given alarms with ticket should return ticket alarms counter
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-3",
        "connector_name": "test-connector-name-axe-alarm-count-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-3",
        "resource": "test-resource-axe-alarm-count-get-3-1",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-3"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-3",
        "connector_name": "test-connector-name-axe-alarm-count-get-3",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-3",
        "resource": "test-resource-axe-alarm-count-get-3-2",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-3"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-alarm-count-get-3",
      "connector_name": "test-connector-name-axe-alarm-count-get-3",
      "source_type": "resource",
      "event_type": "assocticket",
      "component":  "test-component-axe-alarm-count-get-3",
      "resource": "test-resource-axe-alarm-count-get-3-2",
      "output": "test-output-axe-alarm-count-get-3",
      "ticket": "test-ticket-axe-alarm-count-get-3"
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-3
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "total": 2,
      "total_active": 2,
      "snooze": 0,
      "ack": 0,
      "ticket": 1,
      "pbehavior_active": 0
    }
    """

  @concurrent
  Scenario: given snoozed alarms should return snooze alarms counter
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-4",
        "connector_name": "test-connector-name-axe-alarm-count-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-4",
        "resource": "test-resource-axe-alarm-count-get-4-1",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-4"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-4",
        "connector_name": "test-connector-name-axe-alarm-count-get-4",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-4",
        "resource": "test-resource-axe-alarm-count-get-4-2",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-4"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-alarm-count-get-4",
      "connector_name": "test-connector-name-axe-alarm-count-get-4",
      "source_type": "resource",
      "event_type": "snooze",
      "component":  "test-component-axe-alarm-count-get-4",
      "resource": "test-resource-axe-alarm-count-get-4-2",
      "output": "test-output-axe-alarm-count-get-4",
      "duration": 3600
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-4
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "total": 2,
      "total_active": 1,
      "snooze": 1,
      "ack": 0,
      "ticket": 0,
      "pbehavior_active": 0
    }
    """

  @concurrent
  Scenario: given get request with opened filter should return alarms by open status
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-1",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-5"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-2",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-5"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-3",
        "state": 2,
        "output": "test-output-axe-alarm-count-get-5"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    [
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "snooze",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-1",
        "output": "test-output-axe-alarm-count-get-5",
        "duration": 3600
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "snooze",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-2",
        "output": "test-output-axe-alarm-count-get-5",
        "duration": 3600
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "assocticket",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-1",
        "output": "test-output-axe-alarm-count-get-5",
        "ticket": "test-ticket-axe-alarm-count-get-5"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "assocticket",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-2",
        "output": "test-output-axe-alarm-count-get-5",
        "ticket": "test-ticket-axe-alarm-count-get-5"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-1",
        "output": "test-output-axe-alarm-count-get-5"
      },
      {
        "connector": "test-connector-axe-alarm-count-get-5",
        "connector_name": "test-connector-name-axe-alarm-count-get-5",
        "source_type": "resource",
        "event_type": "ack",
        "component":  "test-component-axe-alarm-count-get-5",
        "resource": "test-resource-axe-alarm-count-get-5-2",
        "output": "test-output-axe-alarm-count-get-5"
      }
    ]
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-alarm-count-get-5",
      "connector_name": "test-connector-name-axe-alarm-count-get-5",
      "source_type": "resource",
      "event_type": "cancel",
      "component":  "test-component-axe-alarm-count-get-5",
      "resource": "test-resource-axe-alarm-count-get-5-2",
      "output": "test-output-axe-alarm-count-get-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-axe-alarm-count-get-5",
      "connector_name": "test-connector-name-axe-alarm-count-get-5",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component":  "test-component-axe-alarm-count-get-5",
      "resource": "test-resource-axe-alarm-count-get-5-2",
      "output": "test-output-axe-alarm-count-get-5"
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "total": 3,
      "total_active": 1,
      "snooze": 2,
      "ack": 2,
      "ticket": 2,
      "pbehavior_active": 0
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-5&opened=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "total": 2,
      "total_active": 1,
      "snooze": 1,
      "ack": 1,
      "ticket": 1,
      "pbehavior_active": 0
    }
    """
    When I do GET /api/v4/alarm-counters?search=test-resource-axe-alarm-count-get-5&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "total": 1,
      "total_active": 0,
      "snooze": 1,
      "ack": 1,
      "ticket": 1,
      "pbehavior_active": 0
    }
    """

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/alarm-counters
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarm-counters
    Then the response code should be 403
