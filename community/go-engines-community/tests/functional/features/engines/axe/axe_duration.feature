Feature: Durations of an alarm
  I need to be able to get durations of an alarm

  @concurrent
  Scenario: given new alarm should get alarm durations
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 2,
      "output": "test-output-axe-duration-1",
      "timestamp": {{ nowAdd "-10s" }},
      "connector": "test-connector-axe-duration-1",
      "connector_name": "test-connector-name-axe-duration-1",
      "component": "test-component-axe-duration-1",
      "resource": "test-resource-axe-duration-1",
      "source_type": "resource"
    }
    """
    When I wait 2s
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-1
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response minExpectedDuration=2
    When I save response maxExpectedDuration=4
    Then "duration" >= "minExpectedDuration"
    Then "duration" <= "maxExpectedDuration"
    Then "currentStateDuration" >= "minExpectedDuration"
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 3,
      "output": "test-output-axe-duration-1",
      "timestamp": {{ nowAdd "-8s" }},
      "connector": "test-connector-axe-duration-1",
      "connector_name": "test-connector-name-axe-duration-1",
      "component": "test-component-axe-duration-1",
      "resource": "test-resource-axe-duration-1",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-1
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response minExpectedDuration=3
    When I save response maxExpectedDuration=5
    When I save response expectedCurrentStateDuration=1
    Then "duration" >= "minExpectedDuration"
    Then "duration" <= "maxExpectedDuration"
    Then "currentStateDuration" >= "expectedCurrentStateDuration"
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "active_duration": {{ .duration }}
          }
        }
      ]
    }
    """

  @concurrent
  Scenario: given resolved alarm should get alarm durations
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 2,
      "output": "test-output-axe-duration-2",
      "timestamp": {{ nowAdd "-7s" }},
      "connector": "test-connector-axe-duration-2",
      "connector_name": "test-connector-name-axe-duration-2",
      "component": "test-component-axe-duration-2",
      "resource": "test-resource-axe-duration-2",
      "source_type": "resource"
    }
    """
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 3,
      "output": "test-output-axe-duration-2",
      "timestamp": {{ nowAdd "-5s" }},
      "connector": "test-connector-axe-duration-2",
      "connector_name": "test-connector-name-axe-duration-2",
      "component": "test-component-axe-duration-2",
      "resource": "test-resource-axe-duration-2",
      "source_type": "resource"
    }
    """
    When I wait 1s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "output": "test-output-axe-duration-2",
      "connector": "test-connector-axe-duration-2",
      "connector_name": "test-connector-name-axe-duration-2",
      "component": "test-component-axe-duration-2",
      "resource": "test-resource-axe-duration-2",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "output": "test-output-axe-duration-2",
      "connector": "test-connector-axe-duration-2",
      "connector_name": "test-connector-name-axe-duration-2",
      "component": "test-component-axe-duration-2",
      "resource": "test-resource-axe-duration-2",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-2
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response minExpectedDuration=3
    When I save response maxExpectedDuration=5
    When I save response expectedCurrentStateDuration=1
    Then "duration" >= "minExpectedDuration"
    Then "duration" <= "maxExpectedDuration"
    Then "currentStateDuration" >= "expectedCurrentStateDuration"
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "active_duration": {{ .duration }}
          }
        }
      ]
    }
    """
    When I wait 2s
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "duration": {{ .duration }},
            "active_duration": {{ .duration }},
            "current_state_duration": {{ .currentStateDuration }}
          }
        }
      ]
    }
    """

  @concurrent
  Scenario: given unsnooze event should update alarm snooze duration
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 2,
      "output": "test-output-axe-duration-3",
      "timestamp": {{ nowAdd "-10s" }},
      "connector": "test-connector-axe-duration-3",
      "connector_name": "test-connector-name-axe-duration-3",
      "component": "test-component-axe-duration-3",
      "resource": "test-resource-axe-duration-3",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 2,
      "output": "test-output-axe-duration-3",
      "connector": "test-connector-axe-duration-3",
      "connector_name": "test-connector-name-axe-duration-3",
      "component": "test-component-axe-duration-3",
      "resource": "test-resource-axe-duration-3",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "snooze_duration": 0
          }
        }
      ]
    }
    """
    Then I wait the end of event processing which contains:
    """json
    {
      "event_type": "unsnooze",
      "connector": "test-connector-axe-duration-3",
      "connector_name": "test-connector-name-axe-duration-3",
      "source_type": "resource",
      "component": "test-component-axe-duration-3",
      "resource": "test-resource-axe-duration-3"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-3
    Then the response code should be 200
    When I save response snoozeDuration={{ ( index .lastResponse.data 0 ).v.snooze_duration }}
    When I save response minExpectedSnoozeDuration=2
    When I save response maxExpectedSnoozeDuration=5
    Then "snoozeDuration" >= "minExpectedSnoozeDuration"
    Then "snoozeDuration" <= "maxExpectedSnoozeDuration"

  @concurrent
  Scenario: given snooze event should update alarm snooze duration on resolve
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "check",
      "state" : 2,
      "output": "test-output-axe-duration-4",
      "timestamp": {{ nowAdd "-5s" }},
      "connector": "test-connector-axe-duration-4",
      "connector_name": "test-connector-name-axe-duration-4",
      "component": "test-component-axe-duration-4",
      "resource": "test-resource-axe-duration-4",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "snooze",
      "duration": 3000,
      "output": "test-output-axe-duration-4",
      "connector": "test-connector-axe-duration-4",
      "connector_name": "test-connector-name-axe-duration-4",
      "component": "test-component-axe-duration-4",
      "resource": "test-resource-axe-duration-4",
      "source_type": "resource"
    }
    """
    When I wait 2s
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "cancel",
      "output": "test-output-axe-duration-4",
      "connector": "test-connector-axe-duration-4",
      "connector_name": "test-connector-name-axe-duration-4",
      "component": "test-component-axe-duration-4",
      "resource": "test-resource-axe-duration-4",
      "source_type": "resource"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "event_type": "resolve_cancel",
      "output": "test-output-axe-duration-4",
      "connector": "test-connector-axe-duration-4",
      "connector_name": "test-connector-name-axe-duration-4",
      "component": "test-component-axe-duration-4",
      "resource": "test-resource-axe-duration-4",
      "source_type": "resource"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-4
    Then the response code should be 200
    When I save response snoozeDuration={{ ( index .lastResponse.data 0 ).v.snooze_duration }}
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response minExpectedSnoozeDuration=2
    When I save response maxExpectedSnoozeDuration=5
    When I save response minExpectedActiveDuration=0
    When I save response maxExpectedActiveDuration=1
    Then "snoozeDuration" >= "minExpectedSnoozeDuration"
    Then "snoozeDuration" <= "maxExpectedSnoozeDuration"
    Then "activeDuration" >= "minExpectedActiveDuration"
    Then "activeDuration" <= "maxExpectedActiveDuration"
