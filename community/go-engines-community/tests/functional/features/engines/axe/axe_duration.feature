Feature: Durations of an alarm
  I need to be able to get durations of an alarm

  Scenario: given new alarm should get alarm durations
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-1",
      "connector_name" : "test-connector-name-axe-duration-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-1",
      "resource" : "test-resource-axe-duration-1",
      "state" : 2,
      "output" : "test-output-axe-duration-1"
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-1
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response expectedDuration=2
    Then "duration" >= "expectedDuration"
    Then "currentStateDuration" >= "expectedDuration"
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-1",
      "connector_name" : "test-connector-name-axe-duration-1",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-1",
      "resource" : "test-resource-axe-duration-1",
      "state" : 3,
      "output" : "test-output-axe-duration-1"
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-1
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response expectedDuration=3
    When I save response expectedCurrentStateDuration=1
    Then "duration" >= "expectedDuration"
    Then "currentStateDuration" >= "expectedCurrentStateDuration"
    Then "currentStateDuration" < "expectedDuration"
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

  Scenario: given resolved alarm should get alarm durations
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-2",
      "connector_name" : "test-connector-name-axe-duration-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-2",
      "resource" : "test-resource-axe-duration-2",
      "state" : 2,
      "output" : "test-output-axe-duration-2"
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-2",
      "connector_name" : "test-connector-name-axe-duration-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-2",
      "resource" : "test-resource-axe-duration-2",
      "state" : 3,
      "output" : "test-output-axe-duration-2"
    }
    """
    When I wait the end of event processing
    When I wait 1s
    When I send an event:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-axe-duration-2",
      "connector_name" : "test-connector-name-axe-duration-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-2",
      "resource" : "test-resource-axe-duration-2",
      "output" : "test-output-axe-duration-2"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector" : "test-connector-axe-duration-2",
      "connector_name" : "test-connector-name-axe-duration-2",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-2",
      "resource" : "test-resource-axe-duration-2",
      "output" : "test-output-axe-duration-2"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-2
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    When I save response currentStateDuration={{ ( index .lastResponse.data 0 ).v.current_state_duration }}
    When I save response expectedDuration=3
    When I save response expectedCurrentStateDuration=1
    Then "duration" >= "expectedDuration"
    Then "currentStateDuration" >= "expectedCurrentStateDuration"
    Then "currentStateDuration" < "expectedDuration"
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

  Scenario: given unsnooze event should update alarm snooze duration
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-3",
      "connector_name" : "test-connector-name-axe-duration-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-3",
      "resource" : "test-resource-axe-duration-3",
      "state" : 2,
      "output" : "test-output-axe-duration-3"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 2,
      "connector" : "test-connector-axe-duration-3",
      "connector_name" : "test-connector-name-axe-duration-3",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-3",
      "resource" : "test-resource-axe-duration-3",
      "output" : "test-output-axe-duration-3"
    }
    """
    When I wait the end of event processing
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
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-3
    Then the response code should be 200
    When I save response snoozeDuration={{ ( index .lastResponse.data 0 ).v.snooze_duration }}
    When I save response expectedSnoozeDuration=2
    Then "snoozeDuration" >= "expectedSnoozeDuration"

  Scenario: given snooze event should update alarm snooze duration on resolve
    Given I am admin
    When I send an event:
    """json
    {
      "event_type" : "check",
      "connector" : "test-connector-axe-duration-4",
      "connector_name" : "test-connector-name-axe-duration-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-4",
      "resource" : "test-resource-axe-duration-4",
      "state" : 2,
      "output" : "test-output-axe-duration-4"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "snooze",
      "duration": 3000,
      "connector" : "test-connector-axe-duration-4",
      "connector_name" : "test-connector-name-axe-duration-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-4",
      "resource" : "test-resource-axe-duration-4",
      "output" : "test-output-axe-duration-4"
    }
    """
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "event_type" : "cancel",
      "connector" : "test-connector-axe-duration-4",
      "connector_name" : "test-connector-name-axe-duration-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-4",
      "resource" : "test-resource-axe-duration-4",
      "output" : "test-output-axe-duration-4"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type" : "resolve_cancel",
      "connector" : "test-connector-axe-duration-4",
      "connector_name" : "test-connector-name-axe-duration-4",
      "source_type" : "resource",
      "component" :  "test-component-axe-duration-4",
      "resource" : "test-resource-axe-duration-4",
      "output" : "test-output-axe-duration-4"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-axe-duration-4
    Then the response code should be 200
    When I save response snoozeDuration={{ ( index .lastResponse.data 0 ).v.snooze_duration }}
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response expectedSnoozeDuration=2
    When I save response expectedActiveDuration=1
    Then "snoozeDuration" >= "expectedSnoozeDuration"
    Then "activeDuration" <= "expectedActiveDuration"
