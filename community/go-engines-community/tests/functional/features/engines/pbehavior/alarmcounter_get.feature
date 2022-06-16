Feature: Get alarms counters
  I need to be able to get a alarms counters

  Scenario: given alarms in pbehavior should return pbehavior alarms counters
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-pbehavior-alarm-count-get-1",
        "connector_name": "test-connector-name-pbehavior-alarm-count-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-count-get-1",
        "resource": "test-resource-pbehavior-alarm-count-get-1-1",
        "state": 2,
        "output": "test-output-pbehavior-alarm-count-get-1"
      },
      {
        "connector": "test-connector-pbehavior-alarm-count-get-1",
        "connector_name": "test-connector-name-pbehavior-alarm-count-get-1",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-count-get-1",
        "resource": "test-resource-pbehavior-alarm-count-get-1-2",
        "state": 2,
        "output": "test-output-pbehavior-alarm-count-get-1"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pbehavior-alarm-count-get-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-count-get-1-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/alarm-counters?search=test-resource-pbehavior-alarm-count-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "total": 2,
      "total_active": 1,
      "snooze": 0,
      "ack": 0,
      "ticket": 0,
      "pbehavior_active": 1
    }
    """

  Scenario: given alarms in active pbehavior should not return pbehavior alarms counters
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-pbehavior-alarm-count-get-2",
        "connector_name": "test-connector-name-pbehavior-alarm-count-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-count-get-2",
        "resource": "test-resource-pbehavior-alarm-count-get-2-1",
        "state": 2,
        "output": "test-output-pbehavior-alarm-count-get-2"
      },
      {
        "connector": "test-connector-pbehavior-alarm-count-get-2",
        "connector_name": "test-connector-name-pbehavior-alarm-count-get-2",
        "source_type": "resource",
        "event_type": "check",
        "component":  "test-component-pbehavior-alarm-count-get-2",
        "resource": "test-resource-pbehavior-alarm-count-get-2-2",
        "state": 2,
        "output": "test-output-pbehavior-alarm-count-get-2"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pbehavior-alarm-count-get-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-count-get-2-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/alarm-counters?search=test-resource-pbehavior-alarm-count-get-2
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
