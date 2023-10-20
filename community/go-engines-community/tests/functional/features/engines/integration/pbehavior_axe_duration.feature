Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  @concurrent
  Scenario: given pbehavior should update alarm inactive duration
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-pbehavior-alarm-duration-1",
      "timestamp": {{ nowAdd "-5s" }},
      "connector": "test-connector-pbehavior-alarm-duration-1",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
      "component": "test-component-pbehavior-alarm-duration-1",
      "resource": "test-resource-pbehavior-alarm-duration-1",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 2,
      "output": "test-output-pbehavior-alarm-duration-1",
      "connector": "test-connector-pbehavior-alarm-duration-1",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
      "component": "test-component-pbehavior-alarm-duration-1",
      "resource": "test-resource-pbehavior-alarm-duration-1",
      "source_type": "resource"
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-duration-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "3s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-duration-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "connector": "test-connector-pbehavior-alarm-duration-1",
        "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
        "component": "test-component-pbehavior-alarm-duration-1",
        "resource": "test-resource-pbehavior-alarm-duration-1",
        "source_type": "resource"
      },
      {
        "event_type": "unsnooze",
        "connector": "test-connector-pbehavior-alarm-duration-1",
        "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
        "component": "test-component-pbehavior-alarm-duration-1",
        "resource": "test-resource-pbehavior-alarm-duration-1",
        "source_type": "resource"
      },
      {
        "event_type": "pbhleave",
        "connector": "test-connector-pbehavior-alarm-duration-1",
        "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
        "component": "test-component-pbehavior-alarm-duration-1",
        "resource": "test-resource-pbehavior-alarm-duration-1",
        "source_type": "resource"
      }
    ]
    """
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-1
    Then the response code should be 200
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response minExpectedActiveDuration=0
    When I save response maxExpectedActiveDuration=1
    Then "activeDuration" >= "minExpectedActiveDuration"
    Then "activeDuration" <= "maxExpectedActiveDuration"
