Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  Scenario: given pbehavior should update alarm pbehavior duration
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
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-1",
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
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-1
    Then the response code should be 200
    When I save response pbhInactiveDuration={{ ( index .lastResponse.data 0 ).v.pbh_inactive_duration }}
    When I save response expectedPbhInactiveDuration=2
    Then "pbhInactiveDuration" >= "expectedPbhInactiveDuration"

  Scenario: given active pbehavior should not update alarm pbehavior duration
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-2",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-2",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-alarm-duration-2",
      "resource": "test-resource-pbehavior-alarm-duration-2",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2s" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-duration-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbh_inactive_duration": 0
          }
        }
      ],
      "meta": {
        "total_count": 1
      }
    }
    """

  Scenario: given pbehavior should update alarm pbehavior duration on resolve
    Given I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-3",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-3",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-pbehavior-alarm-duration-3",
      "resource": "test-resource-pbehavior-alarm-duration-3",
      "state": 1,
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-duration-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-3",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-3",
      "source_type": "resource",
      "event_type": "cancel",
      "component": "test-component-pbehavior-alarm-duration-3",
      "resource": "test-resource-pbehavior-alarm-duration-3",
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector": "test-connector-pbehavior-alarm-duration-3",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-3",
      "source_type": "resource",
      "event_type": "resolve_cancel",
      "component": "test-component-pbehavior-alarm-duration-3",
      "resource": "test-resource-pbehavior-alarm-duration-3",
      "output": "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-3
    Then the response code should be 200
    When I save response pbhInactiveDuration={{ ( index .lastResponse.data 0 ).v.pbh_inactive_duration }}
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response expectedPbhInactiveDuration=2
    When I save response expectedActiveDuration=1
    Then "pbhInactiveDuration" >= "expectedPbhInactiveDuration"
    Then "activeDuration" <= "expectedActiveDuration"

  Scenario: given entity with pbehavior info and new alarm should update alarm pbehavior duration on resolve
    Given I am admin
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-pbehavior-alarm-duration-4",
      "connector": "test-connector-pbehavior-alarm-duration-4",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-4",
      "component": "test-component-pbehavior-alarm-duration-4",
      "resource": "test-resource-pbehavior-alarm-duration-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-duration-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-pbehavior-alarm-duration-4",
      "connector": "test-connector-pbehavior-alarm-duration-4",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-4",
      "component": "test-component-pbehavior-alarm-duration-4",
      "resource": "test-resource-pbehavior-alarm-duration-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I wait 3s
    When I send an event:
    """json
    {
      "event_type": "cancel",
      "output": "test-output-pbehavior-alarm-duration-4",
      "connector": "test-connector-pbehavior-alarm-duration-4",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-4",
      "component": "test-component-pbehavior-alarm-duration-4",
      "resource": "test-resource-pbehavior-alarm-duration-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "resolve_cancel",
      "output": "test-output-pbehavior-alarm-duration-4",
      "connector": "test-connector-pbehavior-alarm-duration-4",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-4",
      "component": "test-component-pbehavior-alarm-duration-4",
      "resource": "test-resource-pbehavior-alarm-duration-4",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-4
    Then the response code should be 200
    When I save response duration={{ ( index .lastResponse.data 0 ).v.duration }}
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "active_duration": 0,
            "pbh_inactive_duration": {{ .duration }}
          }
        }
      ]
    }
    """
    When I save response pbhInactiveDuration={{ ( index .lastResponse.data 0 ).v.pbh_inactive_duration }}
    When I save response expectedPbhInactiveDuration=2
    Then "pbhInactiveDuration" >= "expectedPbhInactiveDuration"

  Scenario: given entity with pbehavior info and new alarm should update alarm pbehavior duration on pbhleave
    Given I am admin
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 0,
      "output": "test-output-pbehavior-alarm-duration-5",
      "connector": "test-connector-pbehavior-alarm-duration-5",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-5",
      "component": "test-component-pbehavior-alarm-duration-5",
      "resource": "test-resource-pbehavior-alarm-duration-5",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-alarm-duration-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "4s" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-alarm-duration-5"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I wait 2s
    When I send an event:
    """json
    {
      "event_type": "check",
      "state": 1,
      "output": "test-output-pbehavior-alarm-duration-5",
      "connector": "test-connector-pbehavior-alarm-duration-5",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-5",
      "component": "test-component-pbehavior-alarm-duration-5",
      "resource": "test-resource-pbehavior-alarm-duration-5",
      "source_type": "resource"
    }
    """
    When I wait the end of 2 events processing
    When I wait 3s
    When I send an event:
    """json
    {
      "event_type": "cancel",
      "output": "test-output-pbehavior-alarm-duration-5",
      "connector": "test-connector-pbehavior-alarm-duration-5",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-5",
      "component": "test-component-pbehavior-alarm-duration-5",
      "resource": "test-resource-pbehavior-alarm-duration-5",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "event_type": "resolve_cancel",
      "output": "test-output-pbehavior-alarm-duration-5",
      "connector": "test-connector-pbehavior-alarm-duration-5",
      "connector_name": "test-connector-name-pbehavior-alarm-duration-5",
      "component": "test-component-pbehavior-alarm-duration-5",
      "resource": "test-resource-pbehavior-alarm-duration-5",
      "source_type": "resource"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-alarm-duration-5
    Then the response code should be 200
    When I save response pbhInactiveDuration={{ ( index .lastResponse.data 0 ).v.pbh_inactive_duration }}
    When I save response activeDuration={{ ( index .lastResponse.data 0 ).v.active_duration }}
    When I save response expectedPbhInactiveDuration=1
    When I save response expectedActiveDuration=2
    Then "pbhInactiveDuration" >= "expectedPbhInactiveDuration"
    Then "activeDuration" >= "expectedActiveDuration"
