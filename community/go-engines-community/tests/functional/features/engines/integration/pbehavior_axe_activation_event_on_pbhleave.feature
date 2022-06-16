Feature: send activation event on pbhleave
  I need to be able to trigger rule on alarm activation

  Scenario: given event for new alarm and maintenance pbehavior should not update activation date
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-1",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-1",
      "resource" : "test-resource-pbehavior-axe-activation-event-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-1",
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
              "value": "test-resource-pbehavior-axe-activation-event-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-1",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-1",
      "resource" : "test-resource-pbehavior-axe-activation-event-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-activation-event-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-1",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-1",
            "component" : "test-component-pbehavior-axe-activation-event-1",
            "resource" : "test-resource-pbehavior-axe-activation-event-1",
            "pbehavior_info": {
              "name": "test-pbehavior-axe-activation-event-1"
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    Then the response key "data.0.v.activation_date" should not exist

  Scenario: given event for new alarm and active pbehavior should update activation date
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-2",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-2",
      "resource" : "test-resource-pbehavior-axe-activation-event-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-2",
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
              "value": "test-resource-pbehavior-axe-activation-event-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-2",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-2",
      "resource" : "test-resource-pbehavior-axe-activation-event-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I save response createTimestamp={{ now }}
    When I wait the end of event processing
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-activation-event-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-2",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-2",
            "component" : "test-component-pbehavior-axe-activation-event-2",
            "resource" : "test-resource-pbehavior-axe-activation-event-2",
            "pbehavior_info": {
              "name": "test-pbehavior-axe-activation-event-2"
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate createTimestamp is in range -2,2

  Scenario: given event for new alarm and maintenance pbehavior should update activation date on pbhleave
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-3",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-3",
      "resource" : "test-resource-pbehavior-axe-activation-event-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-axe-activation-event-3",
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
              "value": "test-resource-pbehavior-axe-activation-event-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-axe-activation-event-3",
      "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-axe-activation-event-3",
      "resource" : "test-resource-pbehavior-axe-activation-event-3",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I save response pbhleaveTimestamp={{ now }}
    When I do GET /api/v4/alarms?search=test-resource-pbehavior-axe-activation-event-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "connector" : "test-connector-pbehavior-axe-activation-event-3",
            "connector_name" : "test-connector-name-pbehavior-axe-activation-event-3",
            "component" : "test-component-pbehavior-axe-activation-event-3",
            "resource" : "test-resource-pbehavior-axe-activation-event-3"
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I save response alarmActivationDate={{ (index .lastResponse.data 0).v.activation_date }}
    Then the difference between alarmActivationDate pbhleaveTimestamp is in range -2,2
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
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
              {"_t": "stateinc"},
              {"_t": "statusinc"},
              {"_t": "pbhenter"},
              {"_t": "pbhleave"}
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
