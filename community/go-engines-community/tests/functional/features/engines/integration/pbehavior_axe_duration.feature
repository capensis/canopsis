Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  Scenario: given pbehavior should update alarm inactive duration
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-1",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-alarm-duration-1",
      "resource": "test-resource-pbehavior-alarm-duration-1",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-1",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-1",
      "source_type": "resource",
      "event_type": "snooze",
      "duration": 2,
      "component": "test-component-pbehavior-alarm-duration-1",
      "resource": "test-resource-pbehavior-alarm-duration-1",
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
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
    When I wait the end of 3 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-1
    Then the response code should be 200
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response expectedActiveDuration=0
    Then the difference between activeDuration expectedActiveDuration is in range 0,1
